package notebook

import (
	"reflect"
	"testing"

	"github.com/msp301/graph"
	"github.com/msp301/zb"
)

func TestRecentListsNewestNotesFirst(t *testing.T) {
	g := graph.New()
	g.Add(202401010001, "note", zb.Note{Id: 202401010001, Start: 2})
	g.Add(202401010003, "note", zb.Note{Id: 202401010003, Start: 2})
	g.Add(202401010002, "note", zb.Note{Id: 202401010002, Start: 2})
	g.Add(202401010004, "note", nil)
	g.Add(1, "tag", "#test")

	book := &Notebook{Notes: g}

	var got []uint64
	for _, result := range book.Recent(2) {
		got = append(got, result.Value.(graph.Vertex).Id)
	}

	want := []uint64{202401010003, 202401010002}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected: %v\nGot: %v", want, got)
	}
}
