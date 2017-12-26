package run

import (
	"context"
	"os"
	"os/signal"
	"errors"
)

func WithCancelBy(f interface{}, sig ...os.Signal) error {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, sig...)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		select {
		case <-c: cancel()
		case <-ctx.Done():
		}
	}()
	if fn, ok := f.(func(ctx context.Context)error); ok {
		return fn(ctx)
	} else {
		f.(func(ctx context.Context))(ctx)
		return nil
	}
}

func WithCancelByInterruptErr(f func(ctx context.Context)error) error {
	return WithCancelBy(f, os.Interrupt)
}

func WithCancelByInterrupt(f func(ctx context.Context)) {
	WithCancelBy(f, os.Interrupt)
}

func Interrupted(ctx context.Context) bool {
	select {
	case <- ctx.Done():
		return true
	default:
		return false
	}
}

var InterruptedError = errors.New("interrupted")

func InterruptedErr(ctx context.Context) error {
	if Interrupted(ctx) {
		return InterruptedError
	}
	return nil
}
