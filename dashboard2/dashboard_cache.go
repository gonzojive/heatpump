package dashboard

import (
	"context"
	"fmt"
	"io"
	"sort"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/gonzojive/heatpump/cx34"
	"github.com/gonzojive/heatpump/proto/chiltrix"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	backgroundRefreshRate = time.Minute * 10
)

var backgrundStartTime = time.Now().Add(-time.Hour * 1)

type cache struct {
	client               chiltrix.HistorianClient
	m                    map[int64]*cx34.State
	lock                 *sync.RWMutex
	refreshCh            chan *refreshRequest
	validStart, validEnd time.Time
}

type refreshRequest struct {
	span     span
	callback func(err error)
}

func (c *cache) refresh(span span, callback func(error)) {
	if within(span, c.readValidSpan()) {
		callback(nil)
		return
	}
	select {
	case c.refreshCh <- &refreshRequest{
		span:     span,
		callback: callback,
	}:
	default:
	}
}

func (c *cache) readValidSpan() span {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return span{c.validStart, c.validEnd}
}

func (c *cache) doRefresh(ctx context.Context, s span) error {
	validSpan := c.readValidSpan()
	if validSpan.start.IsZero() || validSpan.end.IsZero() {
		return c.doSingleSpanRefresh(ctx, s)
	}
	if within(s, validSpan) {
		return nil
	}
	var spans []span
	if s.start.Before(validSpan.start) {
		if err := c.doSingleSpanRefresh(ctx, span{s.start, validSpan.start}); err != nil {
			return err
		}
	}
	if s.end.After(validSpan.end) {
		spans = append(spans)
		if err := c.doSingleSpanRefresh(ctx, span{validSpan.end, s.end}); err != nil {
			return err
		}
	}
	return nil
}

func (c *cache) doSingleSpanRefresh(ctx context.Context, s span) error {
	putState := func(s *cx34.State) {
		c.lock.Lock()
		defer c.lock.Unlock()
		c.m[s.CollectionTime().UnixNano()] = s
	}
	updateValidityPeriod := func() {
		c.lock.Lock()
		defer c.lock.Unlock()
		if s.end.IsZero() || s.end.After(c.validEnd) {
			c.validEnd = s.end
		}
		if c.validStart.IsZero() || s.start.Before(c.validStart) {
			c.validStart = s.start
		}
	}

	glog.Infof("querying waterpi for %s worth of state (%s-%s)", s.end.Sub(s.start), s.start, s.end)
	queryClient, err := c.client.QueryStream(ctx, &chiltrix.QueryStreamRequest{
		StartTime: timestamppb.New(s.start.Add(time.Second * -10)), // buffer to make up for race condition writing to database vs querying database (this call)
		EndTime:   timestamppb.New(s.end),
	})
	if err != nil {
		return err
	}
	count := 0
	for {
		resp, err := queryClient.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		s, err := cx34.StateFromProto(resp.GetState())
		if err != nil {
			return err
		}
		putState(s)
		count++
	}
	glog.Infof("got %d states from %s query (%s-%s)", count, s.end.Sub(s.start), s.start, s.end)
	updateValidityPeriod()
	return nil
}

func (c *cache) queryStates(ctx context.Context, span span) ([]*cx34.State, error) {
	debugStart := time.Now()
	defer func() {
		glog.Infof("cache query took %s", time.Now().Sub(debugStart))
	}()
	var err error
	wg := sync.WaitGroup{}
	wg.Add(1)
	c.refresh(span, func(e error) {
		err = e
		wg.Done()
	})
	wg.Wait()
	if err != nil {
		return nil, err
	}

	tryCached := func() ([]*cx34.State, error) {
		if !c.containsPeriod(span) {
			return nil, fmt.Errorf("valid period doesn't contain %+v", span)
		}

		c.lock.RLock()
		states := c.cachedStates()
		c.lock.RUnlock()
		return filterStates(states, func(s *cx34.State) bool {
			t := s.CollectionTime()
			return between(t, span.start, span.end)
		}), nil
	}

	got, err := tryCached()
	if err != nil {
		return nil, fmt.Errorf("internal error; cache after refresh still didn't return usable results: %w", err)
	}
	return got, nil
}

func between(t, start, end time.Time) bool {
	if t.Before(start) {
		return false
	}
	if !t.Before(end) {
		return false
	}
	return true
}

type span struct{ start, end time.Time }

// within reports if a is entirely within b.
func within(a, b span) bool {
	if a.start.Before(b.start) {
		return false
	}
	if a.end.After(b.end) {
		return false
	}
	return true
}

// cachedStates returns the states stored in the cache. Should only be called
// when the lock is held.
func (c *cache) cachedStates() []*cx34.State {
	var states []*cx34.State
	for _, v := range c.m {
		states = append(states, v)
	}

	sort.Slice(states, func(i, j int) bool {
		return states[i].CollectionTime().Before(states[j].CollectionTime())
	})
	return states
}

func (c *cache) containsPeriod(sp span) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return within(sp, span{c.validStart, c.validEnd})
}

func newCache(ctx context.Context, client chiltrix.HistorianClient) (*cache, io.Closer, error) {
	c := &cache{
		client:    client,
		m:         make(map[int64]*cx34.State),
		lock:      &sync.RWMutex{},
		refreshCh: make(chan *refreshRequest, 100),
	}
	var finalErr error
	finished := make(chan struct{})
	stop := make(chan struct{})

	go func() {
		backgroundTicker := time.NewTicker(backgroundRefreshRate)
		defer close(finished)
		defer backgroundTicker.Stop()
		for {
			select {
			case <-ctx.Done():
				glog.Errorf("context error, stopping cache: %v", ctx.Err())
				return
			case <-stop:
				return
			case req := <-c.refreshCh:
				if c.containsPeriod(req.span) {
					req.callback(nil)
					continue
				}

				if err := c.doRefresh(ctx, req.span); err != nil {
					glog.Errorf("failed to refresh heat pump state: %v", err)
					req.callback(err)
					continue
				}
				if !c.containsPeriod(req.span) {
					req.callback(fmt.Errorf("internal data refresh error - validity period does not contain requested period"))
					continue
				}
				req.callback(nil)

			case <-backgroundTicker.C:
				if err := c.doRefresh(ctx, span{backgrundStartTime, time.Now()}); err != nil {
					glog.Errorf("failed to refresh heat pump state: %v", err)
				}
			}
		}
	}()

	return c, closerFunc(func() error {
		close(c.refreshCh)
		<-finished
		return finalErr
	}), nil
}

type closerFunc func() error

func (f closerFunc) Close() error {
	return f()
}

func filterStates(states []*cx34.State, pred func(s *cx34.State) bool) []*cx34.State {
	var out []*cx34.State
	for _, s := range states {
		if pred(s) {
			out = append(out, s)
		}
	}
	return out
}
