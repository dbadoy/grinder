package checkpoint

import (
	"os"
	"testing"
)

func TestCheckpoint(t *testing.T) {
	var (
		kind = "temp"
		n    = uint64(4521)
	)

	cp := New(DefaultBasePath, kind)

	if cp.Checkpoint() != 0 {
		t.Fatalf("invalid checkpoint value, want: %v got: %v", 0, cp.Checkpoint())
	}

	if err := cp.SetCheckpoint(n); err != nil {
		t.Fatal(err)
	}

	if cp.Checkpoint() != n {
		t.Fatalf("Checkpoint.SetCheckpoint failure, want: %v got: %v", n, cp.Checkpoint())
	}

	if err := cp.Increase(); err != nil {
		t.Fatalf("Checkpoint.Increase failure, want: %v got: %v", n+1, cp.Checkpoint())
	}

	if cp.Checkpoint() != n+1 {
		t.Fatalf("Checkpoint.Increase failure, want: %v got: %v", n, cp.Checkpoint())
	}

	// Create new object
	obj := New(DefaultBasePath, kind)
	if obj.Checkpoint() != n+1 {
		t.Fatalf("invalid load checkpoint value, want: %v got: %v", n+1, obj.Checkpoint())
	}

	emptyCP := New(DefaultBasePath, "empty")
	if emptyCP.Checkpoint() != 0 {
		t.Fatalf("invalid load checkpoint value, want: %v got: %v", 0, emptyCP.Checkpoint())
	}

	if err := os.RemoveAll(DefaultBasePath); err != nil {
		t.Fatal(err)
	}
}
