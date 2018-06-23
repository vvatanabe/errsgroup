package errsgroup

import (
	"context"
	"sync"
)

type Option func(*Conf)

type Conf struct {
	limitSize   int
	errChanSize int
}

func LimitSize(size int) Option {
	return func(c *Conf) {
		c.limitSize = size
	}
}

func ErrorChanelSize(size int) Option {
	return func(c *Conf) {
		c.errChanSize = size
	}
}

func NewGroup(opts ...Option) *Group {
	conf := &Conf{
		limitSize:   1,
		errChanSize: 0,
	}
	for _, opt := range opts {
		opt(conf)
	}
	return &Group{
		limit:   make(chan struct{}, conf.limitSize),
		errChan: make(chan error, conf.errChanSize),
		wg:      sync.WaitGroup{},
	}
}

func WithContext(ctx context.Context, opts ...Option) (*Group, context.Context) {
	g := NewGroup(opts...)
	ctx, g.cancel = context.WithCancel(ctx)
	return g, ctx
}

type Group struct {
	limit   chan struct{}
	errChan chan error
	wg      sync.WaitGroup
	cancel  func()
}

func (g *Group) Go(f func() error) {
	g.wg.Add(1)
	go func() {
		g.limit <- struct{}{}
		defer func() {
			<-g.limit
			g.wg.Done()
		}()
		if err := f(); err != nil {
			g.errChan <- err
			if g.cancel != nil {
				g.cancel()
			}
		}
	}()
}

func (g *Group) Wait() []error {
	go func() {
		g.wg.Wait()
		if g.cancel != nil {
			g.cancel()
		}
		close(g.errChan)
	}()
	var errs []error
	for r := range g.errChan {
		errs = append(errs, r)
	}
	return errs
}
