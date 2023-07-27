package segment

import (
	"sync"
)

type SkipNodeInt struct {
	key   int64
	value interface{}
	next  []*SkipNodeInt
}

type SkipListInt struct {
	SkipNodeInt
	mutex  sync.RWMutex
	update []*SkipNodeInt
	maxl   int
	skip   int
	level  int
	length int32
}

func NewSkipListInt(skip ...int) *SkipListInt {
	list := &SkipListInt{}
	list.maxl = 32
	list.skip = 4
	list.level = 0
	list.length = 0
	list.SkipNodeInt.next = make([]*SkipNodeInt, list.maxl)
	list.update = make([]*SkipNodeInt, list.maxl)
	if len(skip) == 1 && skip[0] > 1 {
		list.skip = skip[0]
	}
	return list
}
