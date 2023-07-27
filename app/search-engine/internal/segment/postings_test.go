package segment

import (
	"fmt"
	"testing"
)

func TestMergePostings(t *testing.T) {
	p1 := &PostingsList{
		DocId:         1,
		Positions:     []int64{1, 2, 3, 4},
		PositionCount: 4,
		Next: &PostingsList{
			DocId:         4,
			Positions:     []int64{2, 4, 5},
			PositionCount: 5,
			Next:          nil,
		},
	}
	p2 := &PostingsList{
		DocId:         2,
		Positions:     []int64{12, 22, 32, 42},
		PositionCount: 42,
		Next: &PostingsList{
			DocId:         3,
			Positions:     []int64{22, 42, 52},
			PositionCount: 52,
			Next:          nil,
		},
	}
	res := MergePostings(p1, p2)
	for res != nil {
		fmt.Println(res)
		res = res.Next
	}
}
