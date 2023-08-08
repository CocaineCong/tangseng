package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/CocaineCong/tangseng/app/search_engine/segment"
	"github.com/CocaineCong/tangseng/config"
	"github.com/CocaineCong/tangseng/consts"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/pkg/util/se"
)

var metaFile = "segments.json" // 存储的元数据文件，包括各种属性信息

// Meta 元数据
type Meta struct {
	sync.RWMutex
	Version    string           `json:"version"` // 版本号
	IndexCount int64            `json:"index"`
	SegMeta    *segment.SegMeta `json:"seg_meta"`
	Path       string           `json:"path"` // 元数据文件存储路径
}

func ParseMeta() (*Meta, error) {
	// 文件不存在表示没有相关数据，第一次创建
	metaFile = config.Conf.SeConfig.StoragePath + metaFile
	if !se.ExistFile(metaFile) {
		log.LogrusObj.Infof("segMetaFile:%s not exist", metaFile)
		_, err := os.Create(metaFile)
		if err != nil {
			return nil, fmt.Errorf("create segmentsGenFile err:%v", err)
		}
		m := &Meta{
			Version: config.Conf.SeConfig.Version,
			Path:    metaFile,
			SegMeta: &segment.SegMeta{
				NextSeg:  0,
				SegCount: 0,
				SegInfo:  make(map[segment.SegId]*segment.SegInfo, 0),
			},
			// TODO:初始化读取正排索引
			IndexCount: 0,
		}
		err = writeMeta(m)
		if err != nil {
			return nil, fmt.Errorf("write seg err:%v", err)
		}
		return m, nil
	}

	return readMeta(metaFile)
}

func readMeta(metaFile string) (*Meta, error) {
	metaByte, err := os.ReadFile(metaFile)
	if err != nil {
		return nil, fmt.Errorf("read file err:%v", err)
	}
	h := new(Meta)
	err = json.Unmarshal(metaByte, &h)
	if err != nil {
		return nil, fmt.Errorf("ParseHeader err:%v", err)
	}
	h.Path = metaFile
	return h, nil
}

func writeMeta(m *Meta) (err error) {
	f, err := os.OpenFile(m.Path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0766)
	if err != nil {
		return
	}
	defer f.Close()
	b, _ := json.Marshal(m)
	_, err = f.Write(b)
	if err != nil {
		return
	}

	return
}

// SyncMeta 同步元数据到文件
func (m *Meta) SyncMeta() (err error) {
	return writeMeta(m)
}

// SyncByTicker 定时同步元数据
func (m *Meta) SyncByTicker(ticker *time.Ticker) {
	// 清理计时器
	for {
		log.LogrusObj.Infof("ticker start:%s,next seg id:%d", time.Now().Format(consts.LayOutTimeFormat), m.SegMeta.NextSeg)
		err := m.SyncMeta()
		if err != nil {
			log.LogrusObj.Error(err)
			return
		}
		<-ticker.C
	}
}

// UpdateSegMeta 更新 SegMeta 并同步
func (m *Meta) UpdateSegMeta(segId segment.SegId, indexCount int64) (err error) {
	err = m.SegMeta.UpdateSegMeta(segId, indexCount)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	return m.SyncMeta()
}
