package hook

import (
	"context"
	"errors"
	"testing"
)

func TestHookManager_Action(t *testing.T) {
	hm := NewHookManager()
	ctx := context.Background()

	var order []int
	hm.AddAction("test.event", 20, func(ctx context.Context, payload any) error {
		order = append(order, 20)
		return nil
	})
	hm.AddAction("test.event", 10, func(ctx context.Context, payload any) error {
		order = append(order, 10)
		return nil
	})
	hm.AddAction("test.event", 15, func(ctx context.Context, payload any) error {
		order = append(order, 15)
		return nil
	})

	err := hm.DoAction(ctx, "test.event", "data")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(order) != 3 || order[0] != 10 || order[1] != 15 || order[2] != 20 {
		t.Fatalf("unexpected action execution order: %v", order)
	}
}

func TestHookManager_ActionError(t *testing.T) {
	hm := NewHookManager()
	ctx := context.Background()

	expectedErr := errors.New("boom")
	hm.AddAction("test.fail", 10, func(ctx context.Context, payload any) error {
		return expectedErr
	})

	err := hm.DoAction(ctx, "test.fail", nil)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}

func TestHookManager_Filter(t *testing.T) {
	hm := NewHookManager()
	ctx := context.Background()

	// Filter: calculate score
	hm.AddFilter("calculate.score", 10, func(ctx context.Context, data any) (any, error) {
		score := data.(int)
		return score + 5, nil
	})
	hm.AddFilter("calculate.score", 20, func(ctx context.Context, data any) (any, error) {
		score := data.(int)
		return score * 2, nil
	})

	res, err := hm.ApplyFilter(ctx, "calculate.score", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// (10 + 5) * 2 = 30
	if res.(int) != 30 {
		t.Fatalf("expected 30, got %v", res)
	}
}
