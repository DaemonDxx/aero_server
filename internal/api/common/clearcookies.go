package common

import (
	"context"
	"github.com/chromedp/cdproto/storage"
	"github.com/chromedp/chromedp"
)

func ClearCookies() chromedp.ActionFunc {
	return func(ctx context.Context) error {
		p := storage.ClearCookies()
		return p.Do(ctx)
	}
}
