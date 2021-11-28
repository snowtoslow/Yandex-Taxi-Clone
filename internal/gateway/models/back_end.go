package models

import (
	grpc_ratelimit "github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	"google.golang.org/grpc"
	"net/url"
	"sync"
)

type Backend struct {
	URL            *url.URL
	Alive          bool
	Mux            sync.RWMutex
	Limiter        grpc_ratelimit.Limiter
	GrpcClientConn *grpc.ClientConn
}

func (b *Backend) SetConn(opts ...grpc.DialOption) error {
	conn, err := grpc.Dial(b.URL.Host, opts...)
	if err != nil {
		return err
	}
	b.GrpcClientConn = conn
	return nil
}

func (b *Backend) SetAlive(alive bool) {
	b.Mux.Lock()
	b.Alive = alive
	b.Mux.Unlock()
}

// IsAlive returns true when backend is alive
func (b *Backend) IsAlive() (alive bool) {
	b.Mux.RLock()
	alive = b.Alive
	b.Mux.RUnlock()
	return
}
