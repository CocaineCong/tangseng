package segment

import (
	"fmt"
	"testing"

	"github.com/CocaineCong/tangseng/app/search-engine/internal/query"
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

func TestMergeInvertedIndex(t *testing.T) {
	base := make(InvertedIndexHash)
	token := query.Tokenization{
		Token:    "测试文本",
		Position: 10,
		Offset:   100,
	}
	err := Token2PostingsLists(base, token, 2)
	if err != nil {
		fmt.Println("Token2PostingsLists", err)
	}
	fmt.Println("base", base)

	addDoc := make(InvertedIndexHash)
	token2 := query.Tokenization{
		Token:    "测试文本2",
		Position: 101,
		Offset:   1002,
	}
	err = Token2PostingsLists(addDoc, token2, 3)
	if err != nil {
		fmt.Println("Token2PostingsLists2", err)
	}
	fmt.Println("docDoc", addDoc)
	MergeInvertedIndex(base, addDoc)
	fmt.Println(base)
}
