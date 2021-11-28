package models

import (
	//"context"
	grpc_ratelimit "github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	"time"

	//"go.uber.org/ratelimit"
	"golang.org/x/time/rate"
)

type limiter struct {
	//ratelimit.Limiter
	Limiter *rate.Limiter
}

func NewLimiter(count int) grpc_ratelimit.Limiter {
	magicRate := rate.Every(time.Second)
	return &limiter{
		Limiter: rate.NewLimiter(magicRate, count),
	}
}

func (l *limiter) Limit() bool {
	return l.Limiter.Allow()
}
