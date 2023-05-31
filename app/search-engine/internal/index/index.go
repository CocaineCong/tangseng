package index

import (
	"os"
	"strings"

	"github.com/CocaineCong/Go-SearchEngine/pkg/util/se"
)

// 倒排索引
type RevIndex struct {
	Name string              `json:"name"`
	Data map[string][]string `json:"data"`
	Path string              `json:"path"`
}

// 索引集合
type Set struct {
	Sets []string `json:"sets"`
	Path string   `json:"path"`
}

func NewSet() *Set {
	return &Set{
		Sets: []string{},
	}
}

// GetIndexSet 查找索引集合，不存在会创建
func GetIndexSet(setFilePath string) *Set {
	_, err := os.Stat(setFilePath)
	if err != nil { // 文件不存在则新建索引集合
		indexSet := NewSet()
		indexSet.Path = setFilePath
		_, err := se.DumpJson(setFilePath, indexSet)
		if err != nil {
			panic(err)
		}

		return indexSet
	}

	indexSet := &Set{}
	se.LoadJson(setFilePath, indexSet)
	return indexSet
}

func (indexSet *Set) Save() {
	_, err := se.DumpJson(indexSet.Path, indexSet)
	if err != nil {
		panic(err)
	}
}

func NewRevIndex(name string) *RevIndex {
	return &RevIndex{
		Name: name,
		Data: map[string][]string{},
	}
}

// GetRevIndex 查找倒排索引文件，不存在会创建
func GetRevIndex(revIndexFilePath string) *RevIndex {
	indexInfo := strings.Split(revIndexFilePath, "rev_")
	indexName := indexInfo[1]
	_, err := os.Stat(revIndexFilePath)
	if err != nil { // 文件不存在则新建索引集合
		revIndex := NewRevIndex(indexName)
		revIndex.Path = revIndexFilePath
		_, err = se.DumpJson(revIndexFilePath, revIndex)
		if err != nil {
			panic(err)
		}

		return revIndex
	}

	revIndex := &RevIndex{}
	se.LoadJson(revIndexFilePath, revIndex)
	return revIndex
}

func (revIndex *RevIndex) Save() {
	_, err := se.DumpJson(revIndex.Path, revIndex)
	if err != nil {
		panic(err)
	}
}
