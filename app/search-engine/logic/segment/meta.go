package segment

import (
	"fmt"
	"sync"
)

// Mode 查询 or 索引模式
type Mode int64

const (
	SearchMode Mode = 1 // 查询模式
	IndexMode  Mode = 2 // 索引模式
	MergeMode  Mode = 3 // seg merge 模式
)

type SegId int64

// SegMeta 元数据
type SegMeta struct {
	NextSeg  SegId              `json:"next_seg"`
	SegCount int64              `json:"seg_count"`
	SegInfo  map[SegId]*SegInfo `json:"seg_info"` // TODO replace sync.map

	sync.Mutex
}

// SegInfo 段信息
type SegInfo struct {
	SegId            SegId `json:"seg_name"`           // 段前缀名
	SegSize          int64 `json:"seg_size"`           // 写入doc数量
	InvertedFileSize int64 `json:"inverted_file_size"` // 写入inverted文件大小
	ForwardFileSize  int64 `json:"forward_file_size"`  // 写入forward文件大小
	DelSize          int64 `json:"del_size"`           // 删除文档数量
	DelFileSize      int64 `json:"del_file_size"`      // 删除文档文件大小
	TermSize         int64 `json:"term_size"`          // term文档文件大小
	TermFileSize     int64 `json:"term_file_size"`     // term文件大小
	ReferenceCount   int64 `json:"reference_count"`    // 引入计数
	IsReading        bool  `json:"is_reading"`         // 是否正在被读取
	IsMerging        bool  `json:"is_merging"`         // 是否正在参与合并
}

// newSegment 创建新的segment只创建，更新next seg ，不更新current seg
func newSegmentInfo(segId SegId) *SegInfo {
	return &SegInfo{
		SegId:   segId,
		SegSize: 0,
	}
}

// UpdateSegMeta 更新段信息
func (m *SegMeta) UpdateSegMeta(segId SegId, indexCount int64) error {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.SegInfo[segId]; !ok {
		return fmt.Errorf("segId:%d is not exist", segId)
	}
	m.SegInfo[segId].SegSize = indexCount

	return nil
}

func (m *SegMeta) NewSegmentItem() error {
	m.Lock()
	defer m.Unlock()
	seg := newSegmentInfo(m.NextSeg)
	if _, ok := m.SegInfo[seg.SegId]; ok {
		return fmt.Errorf("seg:%d is exist", seg.SegId)
	}
	m.SegInfo[seg.SegId] = seg
	m.SegCount++
	m.NextSeg++

	return nil
}
