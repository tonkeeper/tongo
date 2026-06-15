package boc

import (
	"testing"
	"time"
)

func TestResetCountersDAGNoExponentialBlowup(t *testing.T) {
	const depth = 60
	cells := make([]*Cell, depth+1)
	cells[depth] = NewCell()
	for i := depth - 1; i >= 0; i-- {
		c := NewCell()
		if err := c.AddRef(cells[i+1]); err != nil {
			t.Fatalf("AddRef: %v", err)
		}
		if err := c.AddRef(cells[i+1]); err != nil {
			t.Fatalf("AddRef: %v", err)
		}
		cells[i] = c
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		cells[0].ResetCounters()
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("ResetCounters did not terminate — exponential DAG blowup regressed")
	}
}
