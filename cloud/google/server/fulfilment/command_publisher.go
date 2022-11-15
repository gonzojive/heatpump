package fulfilment

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/gonzojive/heatpump/cloud/google/cloudconfig"
	"github.com/gonzojive/heatpump/proto/fancoil"
	"google.golang.org/protobuf/proto"

	cpb "github.com/gonzojive/heatpump/proto/controller"
)

type commandPublisher struct {
	topic *pubsub.Topic
}

func newCommandPublisher(c *pubsub.Client) *commandPublisher {
	return &commandPublisher{c.Topic(cloudconfig.CommandsTopic)}
}

func (cp *commandPublisher) executeFanCoilCommand(ctx context.Context, cmd *fancoil.SetStateRequest) error {
	data, err := proto.Marshal(&cpb.Command{
		Command: &cpb.Command_SetStateRequest{
			SetStateRequest: cmd,
		},
	})
	if err != nil {
		return err
	}

	if _, err := cp.topic.Publish(ctx, &pubsub.Message{
		Data: data,
		Attributes: map[string]string{
			"user-id": "redshouse",
		},
	}).Get(ctx); err != nil {
		return fmt.Errorf("failed to publish pub/sub command: %w", err)
	}
	return nil
}

func (cp *commandPublisher) Stop() {
	cp.topic.Stop()
}
