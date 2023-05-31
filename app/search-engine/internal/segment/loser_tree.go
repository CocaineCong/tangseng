package segment

import (
	"github.com/CocaineCong/Go-SearchEngine/app/search-engine/internal/storage"
)

// https://www.cnblogs.com/qianye/archive/2012/11/25/2787923.html 胜者树和败者树

type TermNode struct {
	*storage.KvInfo
	Seg *Segment //
}

// LoserTree 败者数
type LoserTree struct {
	tree   []int // 索引表示顺序，0表示最小值，value表示对应的leaves的index
	leaves []*TermNode
}
