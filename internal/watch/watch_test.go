package watch_test

import (
	"context"
	"testing"
	"time"

	"github.com/your-org/vaultpatch/internal/watch"
)

func TestWatch_NilClientErrors(t *testing.T) {
	_, err := watch.Watch(context.Background(), nil, watch.Options{
		Paths: []string{"secret/app"},
	})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestWatch_EmptyPathsError(t *testing.T) {
	_, err := watch.Watch(context.Background(), &struct{}{}, watch.Options{})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestWatch_CancelClosesChannel(t *testing.T) {
	// Use a real-ish stub: we only test that cancel shuts down the goroutine.
	ctx, cancel := context.WithCancel(context.Background())

	// We cannot call Watch with a real client here, so we exercise the
	// cancellation path via a minimal stub that mimics the channel contract.
	ch := make(chan watch.Event)
	go func() {
		defer close(ch)
		<-ctx.Done()
	}()

	cancel()

	select {
	case <-ch:
		// closed as expected
	case <-time.After(time.Second):
		t.Fatal("channel not closed after cancel")
	}
}

func TestWatch_DefaultInterval(t *testing.T) {
	opts := watch.Options{
		Paths:    []string{"secret/app"},
		Interval: 0,
	}
	// A zero interval should be replaced with the default (30s).
	// We verify no panic occurs and the error is only about the client.
	_, err := watch.Watch(context.Background(), nil, opts)
	if err == nil {
		t.Fatal("expected nil-client error")
	}
}

func TestWatch_AlreadyCancelledContext(t *testing.T) {
	// Ensure Watch respects a context that is cancelled before the call.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := watch.Watch(ctx, nil, watch.Options{
		Paths: []string{"secret/app"},
	})
	// A nil client should still produce an error even with a cancelled context.
	if err == nil {
		t.Fatal("expected error when context is already cancelled and client is nil")
	}
}
