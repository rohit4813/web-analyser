package ctx_test

import (
	"context"
	"github.com/google/uuid"
	"testing"
	iCtx "web-analyser/internal/utils/ctx"
)

func TestRequestID(t *testing.T) {
	t.Run("Test getting request ID from the context", func(t *testing.T) {
		u := uuid.New().String()
		ctx := context.WithValue(context.Background(), iCtx.KeyRequestID, u)
		got := iCtx.RequestID(ctx)
		if got != u {
			t.Fatalf("Expected:%v, Got:%v", u, got)

		}
	})
}

func TestSetRequestID(t *testing.T) {
	t.Run("Test setting request ID in the context", func(t *testing.T) {
		u := uuid.New().String()
		ctx := iCtx.SetRequestID(context.Background(), u)
		got := iCtx.RequestID(ctx)
		if iCtx.RequestID(ctx) != u {
			t.Fatalf("Expected:%v, Got:%v", u, got)
		}
	})
}
