// Package queueserver implements a gRPC service run on a (cloud) server that
// relays Pub/Sub messages to IoT devices.
package queueserver

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"

	"cloud.google.com/go/pubsub"
	"github.com/golang/glog"
	"github.com/gonzojive/heatpump/cloud/acls"
	"github.com/gonzojive/heatpump/util/lockutil"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/gonzojive/heatpump/proto/command_queue"
)

// UserIDAttribute is the name of the pub/sub message attribute with the user id.
const UserIDAttribute = "user-id"

type Service struct {
	pb.UnimplementedCommandQueueServiceServer
	googleCloudProjectID string
	subscriptionName     string
	aclsService          *acls.Service

	// Used to communicate between the pubsub receiver callback and the Listen()
	// implementation.
	listenersByID *lockutil.LockedValue[map[string]*activeListener]
}

type activeListener struct {
	user          *acls.Identity
	handleMessage func(msg *pubsub.Message) error
}

var _ pb.CommandQueueServiceServer = (*Service)(nil)

// New creats a new Google pub/sub backed service.
func New(aclsService *acls.Service, googleCloudProjectID, subscriptionName string) *Service {
	return &Service{
		aclsService:          aclsService,
		googleCloudProjectID: googleCloudProjectID,
		subscriptionName:     subscriptionName,
		listenersByID:        lockutil.NewLockedValue(map[string]*activeListener{}),
	}
}

// Run listens to a pub sub topic until context is cancelled.
func (s *Service) Run(ctx context.Context) error {
	client, err := pubsub.NewClient(ctx, s.googleCloudProjectID)
	if err != nil {
		return fmt.Errorf("failed to initialize pubsub client: %w", err)
	}
	sub := client.Subscription(s.subscriptionName)

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			// listenersByIDLock.Lock()
			// 		listenersByID[listener.user.ID()]
			// 		listenersByIDLock.Unlock()
			// TODO(reddaly): Fix race for access to listener map.
			userID := msg.Attributes[UserIDAttribute]
			if userID == "" {
				glog.Errorf("Invalid Pub/Sub message received with no user id %q attribute; attributes = %v", UserIDAttribute, msg.Attributes)
				msg.Nack()
				return
			}
			listener := lockutil.GetMapValue(s.listenersByID, userID)
			if listener == nil {
				msg.Nack()
				glog.Infof("received message %q with user ID %q, but no listeners exist for that id, so nacking the message", userID, msg.ID)
				return
			}

			if err := listener.handleMessage(msg); err != nil {
				glog.Errorf("error processing message through callback: %v", err)
				msg.Nack()
			}
		})
	})
	return eg.Wait()
}

// Listens for events.
func (s *Service) Listen(stream pb.CommandQueueService_ListenServer) error {
	ident, err := s.aclsService.IdentityFromContext(stream.Context())
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "Client must pass a %q gRPC metadata header to identify the device using a DeviceAccessToken: %v", acls.DeviceAccessTokenMetadataKey, err)
	}
	glog.Infof("authenticated user: %v", ident)

	idCounter := &requestIDCounter{}
	eg, ctx := errgroup.WithContext(stream.Context())

	outstandingMessages := lockutil.NewLockedValue(map[string]*outstandingClientMessage{})

	relayPubSubCommandToClientAndWaitForAck := func(msg *pubsub.Message) error {
		msgID := idCounter.next()
		ocm := &outstandingClientMessage{
			clientMessageID: msgID,
			done:            make(chan struct{}),
			msg:             msg,
			ackStatus:       notFinalized,
		}
		lockutil.SetMapValue(outstandingMessages, msgID, ocm)

		defer lockutil.DeleteMapValue(outstandingMessages, msgID)

		if err := stream.Send(&pb.ListenResponse{
			Response: &pb.ListenResponse_MessageResponse{
				MessageResponse: &pb.MessageResponse{
					Id:      msgID,
					Payload: msg.Data,
				},
			},
		}); err != nil {
			return fmt.Errorf("failed to send MessageResponse to client: %w", err)
		}

		if err := ocm.waitForAck(ctx); err != nil {
			return fmt.Errorf("waiting for ack of message id %q failed: %w", msgID, err)
		}

		return nil
	}

	eg.Go(func() error {
		var listenerForThisRPC *activeListener
		for {
			req, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				// Done processing requests... stop receiving pubsub messages
				// and exit.
				return nil
			}
			if err != nil {
				return fmt.Errorf("error receiving from stream: %w", err)
			}
			if subReq := req.GetSubscribeRequest(); subReq != nil {
				if listenerForThisRPC != nil {
					return status.Errorf(codes.InvalidArgument, "SubscribeRequest can only be sent once")
				}

				listenerForThisRPC = &activeListener{
					user:          ident,
					handleMessage: relayPubSubCommandToClientAndWaitForAck,
				}
				if _, loaded := lockutil.LoadOrStoreMapValue(s.listenersByID, ident.ID(), listenerForThisRPC); loaded {
					return fmt.Errorf("could not subscribe to messages because there is already a registered subscriber for user %q", ident.ID())
				}
				// Once finished processing the stream, remove this listener.
				defer lockutil.DeleteMapValue(s.listenersByID, ident.ID())
			}
			if ackReq := req.GetAckRequest(); ackReq != nil {
				if ackReq.GetMessageId() == "" {
					return status.Errorf(codes.InvalidArgument, "message_id is empty")
				}
				ocm := lockutil.GetMapValue(outstandingMessages, ackReq.GetMessageId())
				if ocm == nil {
					return status.Errorf(codes.InvalidArgument, "message_id %q does not correspond to any message id", ackReq.GetMessageId())
				}
				if ackReq.GetNack() {
					ocm.nack()
				} else {
					ocm.ack()
				}
			}
		}
	})

	return eg.Wait()
}

func (s *Service) ListTopics(ctx context.Context, req *pb.ListTopicsRequest) (*pb.ListTopicsResponse, error) {
	ident, err := s.aclsService.IdentityFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Client must pass a %q gRPC metadata header to identify the device using a DeviceAccessToken: %v", acls.DeviceAccessTokenMetadataKey, err)
	}
	glog.Infof("authenticated user: %v", ident)

	return &pb.ListTopicsResponse{
		Topics: []string{
			"thermostat-commands",
		},
	}, nil
}

type requestIDCounter struct {
	nextValue atomic.Int64
}

func (c *requestIDCounter) next() string {
	return fmt.Sprintf("request-%d", c.nextValue.Add(1))
}

// lazyValue is a lazily-initialized value
type lazyValue[T any] struct {
	once        sync.Once
	value       T
	err         error
	initializer func(ctx context.Context) (T, error)
}

func (lv *lazyValue[T]) get(ctx context.Context) (T, error) {
	lv.once.Do(func() {
		lv.value, lv.err = lv.initializer(ctx)
	})
	return lv.value, lv.err
}

func sendToChannelOrError[T any](ctx context.Context, ch chan T, thing T) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case ch <- thing:
		return nil
	}
}
