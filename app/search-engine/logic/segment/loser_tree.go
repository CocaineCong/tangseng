package segment

import (
	"fmt"
	"strconv"
	"time"

	"github.com/CocaineCong/tangseng/app/search-engine/logic/storage"
	"github.com/CocaineCong/tangseng/app/search-engine/logic/types"
	"github.com/CocaineCong/tangseng/consts"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

// https://www.cnblogs.com/qianye/archive/2012/11/25/2787923.html 胜者树和败者树

type TermNode struct {
	*storage.KvInfo
	Seg *Segment
}

// LoserTree 败者数
type LoserTree struct {
	tree     []int // 索引表示顺序，0表示最小值，value表示对应的leaves的index
	leaves   []*TermNode
	levelsCh []chan storage.KvInfo
}

func NewSegLoserTree(leaves []*TermNode, leavesCh []chan storage.KvInfo) *LoserTree {
	k := len(leaves)
	lt := &LoserTree{
		tree:     make([]int, k),
		leaves:   leaves,
		levelsCh: leavesCh,
	}
	if k > 0 {
		lt.initWinner(0)
	}

	return lt
}

// 整体逻辑 输的留下来，赢的向上比
func (lt *LoserTree) initWinner(idx int) int {
	// 根节点有一个父节点，存储最小值
	if idx == 0 {
		lt.tree[0] = lt.initWinner(1)
		return lt.tree[0]
	}
	if idx >= len(lt.tree) {
		return idx - len(lt.tree)
	}

	left := lt.initWinner(idx * 2)
	right := lt.initWinner(idx*2 + 1)
	log.LogrusObj.Infof("left:%d, right:%d ", left, right)

	// 左不为空，右为空，则记录右边
	if lt.leaves[left] != nil && lt.leaves[right] == nil {
		left, right = right, left
	}

	if lt.leaves[left] != nil && lt.leaves[right] != nil {
		leftCh := <-lt.levelsCh[left]
		rightCh := <-lt.levelsCh[right]

		lt.leaves[left].KvInfo = &storage.KvInfo{
			Key:   leftCh.Key,
			Value: leftCh.Value,
		}

		lt.leaves[right].KvInfo = &storage.KvInfo{
			Key:   rightCh.Key,
			Value: rightCh.Value,
		}

		log.LogrusObj.Infof("left Ch:%s,rightCh:%s", leftCh.Key, rightCh.Key)
		if string(leftCh.Key) < string(rightCh.Key) {
			left, right = right, left
		}

	}
	// 左边的节点比右边的节点大
	// 记录败者 即 记录较大的节点索引 较小的继续向上比较
	lt.tree[idx] = left

	return right
}

// Pop 弹出最小值
func (lt *LoserTree) Pop() (res *TermNode) {
	if len(lt.tree) == 0 {
		return nil
	}
	// 取出最小的索引
	leafWinIdx := lt.tree[0]
	// 找到对应的叶子节点
	winner := lt.leaves[leafWinIdx]
	winnerCh := lt.levelsCh[leafWinIdx]

	// 更新对应index里节点的值
	// 如果是最后一个节点，标识为nil
	if winner == nil {
		log.LogrusObj.Infof("数据已读取完毕，winner.key == nil")
		lt.leaves[leafWinIdx] = nil
		res = nil
	} else {
		var target TermNode
		target = *winner
		res = &target
		log.LogrusObj.Infof("winner:%s", winner.Key)

		// 获取下一轮的key和value
		termCh, isOpen := <-winnerCh
		// channel 已关闭
		if !isOpen {
			log.LogrusObj.Infof("channel 已关闭")
			lt.leaves[leafWinIdx] = nil
		} else {
			// 重新赋值
			lt.leaves[leafWinIdx].KvInfo = &storage.KvInfo{
				Key:   termCh.Key,
				Value: termCh.Value,
			}
		}
	}

	// 读取父节点
	treeIdx := (leafWinIdx + len(lt.tree)) / 2

	for treeIdx != 0 {
		loserLeafIdx := lt.tree[treeIdx]
		if lt.leaves[loserLeafIdx] == nil {
			// 如果为nil ，则父节点的idx设置为该索引，不为空的继续向上比较
			lt.tree[treeIdx] = loserLeafIdx
		} else {
			// 如果该叶子节点已经读取完毕，则将父节点的idx设置为该索引
			if lt.leaves[leafWinIdx] == nil {
				loserLeafIdx, leafWinIdx = leafWinIdx, loserLeafIdx
			} else if string(lt.leaves[loserLeafIdx].Key) < string(lt.leaves[leafWinIdx].Key) {
				loserLeafIdx, leafWinIdx = leafWinIdx, loserLeafIdx
			}

			// 更新
			lt.tree[treeIdx] = loserLeafIdx
		}
		treeIdx /= 2
	}
	lt.tree[0] = leafWinIdx

	time.Sleep(1e8)

	return res
}

// MergeKTermSegments 多路归并，合并term数据，合并后需要一起处理合并倒排表数据
func MergeKTermSegments(list []*TermNode, chList []chan storage.KvInfo) (res InvertedIndexHash, err error) {
	// 初始化
	lt := NewSegLoserTree(list, chList)
	res = make(InvertedIndexHash)

	for {
		node := lt.Pop()
		if node == nil {
			break
		}
		log.LogrusObj.Infof("pop node key:%+v,value:%v", string(node.Key), node.Value)
		val, err := storage.Bytes2TermVal(node.Value)
		if err != nil {
			return
		}
		log.LogrusObj.Infof("val:%+v", val)
		c, err := node.Seg.GetInvertedDoc(val.Offset, val.Size)
		if err != nil {
			log.LogrusObj.Errorf("FetchPostings getDocInfo err: %v", err)
			return
		}
		pos, err := DecodePostings(c)
		if err != nil {
			log.LogrusObj.Errorf("FetchPostings DecodePostings err: %v", err)
			return
		}
		log.LogrusObj.Infof("pop node key:%+v,value:%v,count:%d", string(node.Key), val, pos.DocCount)
		if p, ok := res[string(node.Key)]; ok {
			p.DocCount += pos.DocCount
			p.PostingsList = MergePostings(p.PostingsList, pos.PostingsList)
			continue
		}
		res[string(node.Key)] = &types.InvertedIndexValue{
			Token:        string(node.Key),
			DocCount:     pos.DocCount,
			PostingsList: pos.PostingsList,
		}
	}

	return
}

// MergeKForwardSegments 合并正排
func MergeKForwardSegments(seg *Segment, list []*TermNode, chList []chan storage.KvInfo) error {
	// 初始化
	lt := NewSegLoserTree(list, chList)
	count := int64(0)
	for {
		node := lt.Pop()
		if node == nil {
			break
		}
		// 正排中的总数，需要单独操作 TODO 正排总数字段考虑下其他存储or实现方式
		if string(node.Key) == consts.ForwardCountKey {
			c, err := strconv.Atoi(string(node.Value))
			if err != nil {
				return fmt.Errorf("strconv.Atoi err:%s", err)
			}
			count += int64(c)
			continue
		}
		err := seg.PutForwardByKV(node.Key, node.Value)
		if err != nil {
			return fmt.Errorf("Put Error:%v", err)
		}
		log.LogrusObj.Infof("pop node key:%s,value:%s", node.Key, node.Value)

	}
	// 更新count
	return seg.UpdateForwardCount(count)
}
