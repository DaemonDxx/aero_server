package lks

import "context"

type BrowserPool interface {
	Put(ctx context.Context)
	Get() (context.Context, error)
}
