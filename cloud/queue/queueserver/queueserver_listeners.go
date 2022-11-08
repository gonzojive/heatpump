package queueserver

import (
	"context"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/golang/glog"
	"github.com/gonzojive/heatpump/util/lockutil"
)

type listenerSet struct {
	idToListener *lockutil.LockedValue[map[string]*activeListener]
}

func newListenerSet() *listenerSet {
	return &listenerSet{
		idToListener: lockutil.NewLockedValue(map[string]*activeListener{}),
	}
}

func (set *listenerSet) get(id string) *activeListener {
	var got *activeListener
	set.idToListener.Observe(func(old map[string]*activeListener) {
		got = old[id]
	})
	return got
}

func (set *listenerSet) remove(id string) *activeListener {
	var got *activeListener
	set.idToListener.Observe(func(old map[string]*activeListener) {
		got = old[id]
	})
	return got
}

// outstandingClientMessage is used to keep track of messages that may not have
// been relayed to the client that have not yet been acked or nacked.
type outstandingClientMessage struct {
	clientMessageID string
	msg             *pubsub.Message
	done            chan struct{}

	ackOnce   sync.Once
	ackStatus ackStatus
}

func (m *outstandingClientMessage) ack() {
	m.ackOnce.Do(func() {
		m.msg.Ack()
		m.ackStatus = acked
		close(m.done)
		glog.Infof("acked message %q", m.msg.ID)
	})
}

func (m *outstandingClientMessage) nack() {
	m.ackOnce.Do(func() {
		m.msg.Nack()
		m.ackStatus = nacked
		close(m.done)
	})
}

func (m *outstandingClientMessage) waitForAck(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-m.done:
		return nil
	}
}

// ackStatus differentiates between states of the acknowledgement of a message.
type ackStatus byte

const (
	notFinalized ackStatus = 0
	acked        ackStatus = 1
	nacked       ackStatus = 2
)
