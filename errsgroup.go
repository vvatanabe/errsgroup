package errsgroup

import (
	"context"
	"sync"
)

type Option func(*Conf)

type Conf struct {
	limitSize int
}

func LimitSize(size int) Option {
	return func(c *Conf) {
		c.limitSize = size
	}
}

const errChanSize = 1

func NewGroup(opts ...Option) *Group {
	conf := &Conf{1}
	for _, opt := range opts {
		opt(conf)
	}
	return &Group{
		limitChan: make(chan struct{}, conf.limitSize),
		errChan:   make(chan error, errChanSize),
	}

}

func WithContext(ctx context.Context, opts ...Option) (*Group, context.Context) {
	g := NewGroup(opts...)
	ctx, g.cancel = context.WithCancel(ctx)
	return g, ctx
}

type Group struct {
	limitChan  chan struct{}
	errChan    chan error
	wg         sync.WaitGroup
	cancelOnce sync.Once
	cancel     func()
}

func (g *Group) Go(f func() error) {
	g.wg.Add(1)
	go func() {
		g.limitChan <- struct{}{}
		if err := f(); err != nil {
			g.errChan <- err
			g.cancelOnceIfExist()
		}
		<-g.limitChan
		g.wg.Done()
	}()
}

func (g *Group) Wait() []error {
	go func() {
		g.wg.Wait()
		g.cancelOnceIfExist()
		close(g.errChan)
	}()
	var errs []error
	for r := range g.errChan {
		errs = append(errs, r)
	}
	return errs
}

func (g *Group) cancelOnceIfExist() {
	if g.cancel != nil {
		g.cancelOnce.Do(func() {
			g.cancel()
		})
	}
}
