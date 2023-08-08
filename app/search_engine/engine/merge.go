package engine

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	segment2 "github.com/CocaineCong/tangseng/app/search_engine/segment"
	"github.com/CocaineCong/tangseng/app/search_engine/storage"
	"github.com/CocaineCong/tangseng/config"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

// MergeScheduler 合并调度器
type MergeScheduler struct {
	Message chan *MergeMessage
	Meta    *Meta

	sync.WaitGroup
}

// MergeMessage 合并队列
type MergeMessage []*segment2.SegInfo

// NewScheduler 创建调度器
func NewScheduler(meta *Meta) *MergeScheduler {
	ch := make(chan *MergeMessage, config.Conf.SeConfig.MergeChannelSize)
	return &MergeScheduler{
		Message: ch,
		Meta:    meta,
	}
}

// Merge 合并入口
func (m *MergeScheduler) Merge() {
	for {
		select {
		case segs := <-m.Message:
			log.LogrusObj.Infof("Merge Msg:%v", segs)
			// 合并
			err := m.merge(segs)
			if err != nil {
				log.LogrusObj.Infof("merge error:%v", err)
			}
		case <-time.After(1e9):
			log.LogrusObj.Infof("sleep 1s...")
		}
	}
}

// Close 关闭调度器
func (m *MergeScheduler) Close() {
	// 保证所有的merge执行完毕
	m.Wait()
}

// MayMerge 判断是否需要merge 通过 meta 数据中的seg info 来计算
func (m *MergeScheduler) MayMerge() {
	// 已存在超过2个segment，则需要判断seg是否需要merge
	if len(m.Meta.SegMeta.SegInfo) <= 1 {
		log.LogrusObj.Infof("seg count:%v,no need merge", len(m.Meta.SegMeta.SegInfo))
		return
	}
	mess, isNeed := m.calculateSegs()
	if !isNeed {
		return
	}
	m.Add(1)
	m.Message <- mess
	log.LogrusObj.Infof("merge segs:%v", mess)

}

// 计算是否有段需要合并
func (m *MergeScheduler) calculateSegs() (*MergeMessage, bool) {
	segs := m.Meta.SegMeta.SegInfo
	log.LogrusObj.Infof("segs %v", segs)
	_, ok0 := segs[0]
	_, ok2 := segs[1]
	if !ok0 || !ok2 {
		return nil, false
	}
	// 判断是否需要合并

	segList := make([]*segment2.SegInfo, 0)
	segList = append(segList, segs[0])
	segList = append(segList, segs[1])

	mes := MergeMessage(segList)
	return &mes, true
}

// term表需要合并k个升序，以及处理对应的倒排索引，正排表直接merge即可
func (m *MergeScheduler) mergeSegments(segs *MergeMessage) (err error) {
	segmentDBs := make([]*segment2.Segment, 0)
	var docSize int64 = 0
	for _, segInfo := range []*segment2.SegInfo(*segs) {
		docSize += segInfo.SegSize
		s := segment2.NewSegment(segInfo.SegId)
		segmentDBs = append(segmentDBs, s)
	}

	if len(segmentDBs) == 0 {
		log.LogrusObj.Infof("no segment to merge")
		return
	}

	termNodes := make([]*segment2.TermNode, 0)
	termChs := make([]chan storage.KvInfo, 0)

	forNodes := make([]*segment2.TermNode, 0)
	forChs := make([]chan storage.KvInfo, 0)

	for _, seg := range segmentDBs {
		termNode := new(segment2.TermNode)
		termNode.Seg = seg

		// 开启协程遍历读取
		termCh := make(chan storage.KvInfo)
		go func() { // TODO：协程发生panic怎么办
			err = seg.GetInvertedTermCursor(termCh)
			if err != nil {
				log.LogrusObj.Errorln("GetInvertedTermCursor", err)
				return
			}
		}()

		forCh := make(chan storage.KvInfo)
		go func() {
			err = seg.GetForwardCursor(forCh)
			if err != nil {
				log.LogrusObj.Errorln("GetForwardCursor", err)
				return
			}
		}()

		termNodes = append(termNodes, termNode)
		termChs = append(termChs, termCh)

		forNodes = append(forNodes, new(segment2.TermNode))
		forChs = append(forChs, forCh)
	}

	// 合并term和倒排数据，返回合并后的数据
	res, err := segment2.MergeKTermSegments(termNodes, termChs)
	if err != nil {
		log.LogrusObj.Errorf("MergeKTermSegments:%v", err)
		return
	}

	engineTmp := NewEngine(m.Meta, segment2.MergeMode)
	// 落盘
	err = engineTmp.Seg[engineTmp.CurrSegId].Flush(res)
	if err != nil {
		log.LogrusObj.Errorf("NewEngine-Flush:%v", err)
		return
	}
	log.LogrusObj.Infof("start forward:%s", strings.Repeat("-", 20))

	// 合并正排数据
	err = segment2.MergeKForwardSegments(engineTmp.Seg[engineTmp.CurrSegId], forNodes, forChs)
	if err != nil {
		log.LogrusObj.Infof("forward merge error:%v", err)
		return err
	}

	// 更新 meta info
	err = m.Meta.UpdateSegMeta(engineTmp.CurrSegId, docSize)
	if err != nil {
		log.LogrusObj.Infof("update seg meta err:%v", err)
		return err
	}

	// 删除老的segs
	err = m.deleteOldSeg(*segs)
	if err != nil {
		log.LogrusObj.Infof("update seg meta err:%v", err)
		return err
	}

	return nil
}

func (m *MergeScheduler) deleteOldSeg(segInfos []*segment2.SegInfo) error {
	for _, segInfo := range segInfos {
		if s, ok := m.Meta.SegMeta.SegInfo[segInfo.SegId]; ok {
			s.IsMerging = false
			delete(m.Meta.SegMeta.SegInfo, segInfo.SegId)
			err := m.deleteSegFile(segInfo.SegId)
			if err != nil {
				log.LogrusObj.Infof("err:%v", err)
				return err
			}
		} else {
			return fmt.Errorf("delete old seg error:%v", segInfo)
		}
	}

	return nil
}

func (m *MergeScheduler) deleteSegFile(segId segment2.SegId) error {
	term, inverted, forward := segment2.GetDbName(segId)
	log.LogrusObj.Infof("delete seg file forward:%s,invert:%s,term:%s", term, inverted, forward)
	err := os.Remove(inverted)
	if err != nil {
		return err
	}
	os.Remove(term)
	if err != nil {
		return err
	}
	os.Remove(forward)
	if err != nil {
		return err
	}

	return nil
}

// merge 合并segment
func (m *MergeScheduler) merge(segs *MergeMessage) error {
	defer m.Done()
	log.LogrusObj.Infof("merge segs:%v", segs)
	// 恢复 seg is_merging 状态
	defer func() {
		for _, seg := range ([]*segment2.SegInfo)(*segs) {
			// 如果merge失败，没有删除旧seg，需要恢复
			if s, ok := m.Meta.SegMeta.SegInfo[seg.SegId]; ok {
				s.IsMerging = false
			}
		}
	}()

	// 合并
	err := m.mergeSegments(segs)
	if err != nil {
		log.LogrusObj.Errorf("merge err:%v", err)
		return err
	}

	return nil
}
