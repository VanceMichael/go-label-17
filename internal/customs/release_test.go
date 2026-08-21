package customs

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/VanceMichael/go-base-airbridge/internal/domain"
)

type releaseSignerFunc func(context.Context, Case) error

func (f releaseSignerFunc) SignRelease(ctx context.Context, c Case) error { return f(ctx, c) }

func TestFailedReleaseSigningRestoresReviewableCase(t *testing.T) {
	now := time.Date(2026, 8, 21, 7, 0, 0, 0, time.UTC)
	workflow := NewWorkflow()
	if err := workflow.Open(context.Background(), Case{ID: "case-17", ShipmentID: "shipment-17", TenantID: "cargo-east"}); err != nil {
		t.Fatalf("open case: %v", err)
	}
	if err := workflow.Attach(context.Background(), "shipment-17", Document{Number: "doc-17", IssuedAt: now.Add(-time.Hour), ExpiresAt: now.Add(time.Hour)}); err != nil {
		t.Fatalf("attach document: %v", err)
	}
	signErr := errors.New("customs signing service unavailable")
	_, err := workflow.Release(context.Background(), "shipment-17", "broker-17", now, releaseSignerFunc(func(context.Context, Case) error { return signErr }))
	if !errors.Is(err, signErr) {
		t.Fatalf("release error = %v, want signing error", err)
	}
	current, err := workflow.Get("shipment-17")
	if err != nil {
		t.Fatalf("get after failed release: %v", err)
	}
	if current.Status != domain.CustomsPending {
		t.Fatalf("status after failed signing = %s, want pending", current.Status)
	}

	released, err := workflow.Release(context.Background(), "shipment-17", "broker-17", now.Add(time.Minute), releaseSignerFunc(func(context.Context, Case) error { return nil }))
	if err != nil {
		t.Fatalf("retry release: %v", err)
	}
	if released.Status != domain.CustomsReleased || len(released.Notes) != 1 {
		t.Fatalf("released case = %+v", released)
	}
}
