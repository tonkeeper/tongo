package code

import (
	"testing"
	"time"

	"github.com/tonkeeper/tongo/boc"
)

func TestFindLibrariesDAGNoExponentialBlowup(t *testing.T) {
	const depth = 60
	cells := make([]*boc.Cell, depth+1)
	cells[depth] = boc.NewCell()
	for i := depth - 1; i >= 0; i-- {
		c := boc.NewCell()
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
		libs, err := FindLibraries(cells[0])
		if err != nil {
			t.Errorf("FindLibraries: %v", err)
		}
		if len(libs) != 0 {
			t.Errorf("expected no libraries, got %d", len(libs))
		}
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("FindLibraries did not terminate — exponential DAG blowup regressed")
	}
}
