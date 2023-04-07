package server

import (
	"fmt"
	"testing"

	"github.com/dbadoy/grinder/pkg/database/memdb"
	"github.com/dbadoy/grinder/server/cft"
	"github.com/dbadoy/grinder/server/dto"
)

func TestJournalRevert(t *testing.T) {
	var (
		mdb      = memdb.New()
		journals = make([]journalObject, 0)

		n = 10
	)

	engine, err := cft.NewSoloEngine(nil, mdb, nil)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < n; i++ {
		key := []byte(fmt.Sprintf("%d", i))

		engine.Insert(key, &dto.Contract{})
		journals = append(journals, &insertContract{key})
	}

	if mdb.Size() != n {
		t.Fatalf("TestJournalRevert, want: %d, got: %d", n, mdb.Size())
	}

	for _, task := range journals {
		task.revert(engine)
	}

	if mdb.Size() != 0 {
		t.Fatalf("TestJournalRevert, want: 0, got: %d", mdb.Size())
	}
}
