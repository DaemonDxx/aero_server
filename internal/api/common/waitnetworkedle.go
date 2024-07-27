package common

import (
	"context"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func WaitForNetworkIdle() chromedp.ActionFunc {
	return func(ctx context.Context) error {
		resCh := make(chan struct{}, 1)
		ctxCn, cancel := context.WithCancel(ctx)
		chromedp.ListenTarget(ctxCn, func(ev interface{}) {
			switch e := ev.(type) {
			case *page.EventLifecycleEvent:
				if e.Name == "networkIdle" {
					cancel()
					resCh <- struct{}{}
				}
			}
		})

		select {
		case <-resCh:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
