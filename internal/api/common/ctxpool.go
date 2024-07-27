package common

import (
	"context"
	"fmt"
	"sync"
)

type CreateContextFunc func() (context.Context, error)

type CtxPool struct {
	ctxCh chan context.Context
	size  int
	fn    CreateContextFunc
	mu    sync.RWMutex
}

func NewCtxPool(size int, fn CreateContextFunc) *CtxPool {
	ch := make(chan context.Context, size)
	return &CtxPool{
		ctxCh: ch,
		size:  0,
		fn:    fn,
		mu:    sync.RWMutex{},
	}
}

func (p *CtxPool) Put(ctx context.Context) {
	p.ctxCh <- ctx
}

func (p *CtxPool) Get() (context.Context, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.ctxCh) == 0 {
		if p.size < cap(p.ctxCh) {
			if ctx, err := p.fn(); err != nil {
				return nil, fmt.Errorf("call create context error: %w", err)
			} else {
				p.size++
				p.ctxCh <- ctx
			}
		}
	}
	return <-p.ctxCh, nil
}
