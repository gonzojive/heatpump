// Package queueclient provides a rich client library wrapping the gRPC client for CommandQueueService.
package queueclient

import (
	"context"
	"fmt"
	"sync"

	"github.com/golang/glog"
	"github.com/gonzojive/heatpump/util/retry"

	qpb "github.com/gonzojive/heatpump/proto/command_queue"
)

const mainTopic = "thermostat-commands"

var retryConfig = &retry.Config{
	IsRetriable: func(err error) bool { return true },
}

// New returns a wrapped version of the command queue service client.
func New(rawProtoClient qpb.CommandQueueServiceClient) *Client {
	return &Client{rawProtoClient}
}

// Client wraps a raw proto client for CommandQueueService in an easier to use
// object.
type Client struct {
	raw qpb.CommandQueueServiceClient
}

func (c *Client) Listen(ctx context.Context, handler func(command *Command)) error {
	listenAndProcessMessages := func(ctx context.Context) error {
		stream, err := c.raw.Listen(ctx)
		if err != nil {
			return fmt.Errorf("error initiating RPC: %w", err)
		}

		if err := stream.Send(&qpb.ListenRequest{
			Request: &qpb.ListenRequest_SubscribeRequest_{
				SubscribeRequest: &qpb.ListenRequest_SubscribeRequest{
					Topics: []string{mainTopic},
				},
			},
		}); err != nil {
			return err
		}

		for {
			glog.Infof("waiting for message from queue server...")
			got, err := stream.Recv()
			if err != nil {
				return fmt.Errorf("error while waiting for commands: %w", err)
			}

			if msg := got.GetMessageResponse(); msg != nil {
				glog.Infof("got message %q", string(msg.GetPayload()))

				cmd := &Command{
					stream:    stream,
					messageID: msg.GetId(),
					payload:   msg.GetPayload(),
				}
				go func() {
					defer cmd.Nack()
					handler(cmd)
				}()
			}
		}
	}

	return retryConfig.Start(ctx, func(ctx context.Context) error {
		if err := listenAndProcessMessages(ctx); err != nil {
			glog.Errorf("error while processing messages, retrying: %v", err)
			return err
		}
		return nil
	})
}

type Command struct {
	messageID string
	payload   []byte
	stream    qpb.CommandQueueService_ListenClient
	ackOnce   sync.Once

	ackErr error
}

func (c *Command) ack(ack bool) {
	c.ackOnce.Do(func() {
		c.ackErr = c.stream.Send(&qpb.ListenRequest{
			Request: &qpb.ListenRequest_AckRequest_{
				AckRequest: &qpb.ListenRequest_AckRequest{
					MessageId: c.messageID,
					Nack:      !ack,
				},
			},
		})
	})
}

func (c *Command) Ack() {
	c.ack(true)
}

func (c *Command) Nack() {
	c.ack(false)
}

func (c *Command) Payload() []byte {
	return c.payload
}
