package segment

import (
	"time"

	"github.com/CocaineCong/Go-SearchEngine/app/search-engine/internal/storage"
	log "github.com/CocaineCong/Go-SearchEngine/pkg/logger"
)

// https://www.cnblogs.com/qianye/archive/2012/11/25/2787923.html 胜者树和败者树

type TermNode struct {
	*storage.KvInfo
	Seg *Segment //
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
