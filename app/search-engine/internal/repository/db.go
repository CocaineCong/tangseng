package repository

import (
	"errors"
	"os"
	"path"
	"reflect"
	"se/internal/index"
	"se/internal/inputdata"
	"se/internal/utils"
	"strings"
	"sync"
)

type Table struct {
	Name     string
	IndexSet *index.Set
	RevIndex map[string]*index.RevIndex
	mu       sync.RWMutex
}

const (
	BasePathDir string = "./data"
)

//select * from "abc" where viisible = 1
func NewTable(name string) *Table {
	utils.DirCHeckAndMk(path.Join(BasePathDir, name))
	return &Table{
		Name: name,
	}
}

func GetTable(name string) *Table {
	var table *Table = &Table{
		Name: name,
		IndexSet: &index.Set{
			Sets: []string{},
			Path: "",
		},
	}

	tableDir := path.Join(BasePathDir, name)
	_, err := os.Stat(tableDir)
	if err != nil { //数据表不存在
		table = NewTable(name)
	}

	//索引集合
	indexSetPath := path.Join(tableDir, "indexSet") // ./data/abc/indexSet
	indexSet := index.GetIndexSet(indexSetPath)
	table.IndexSet = indexSet //将索引集合加载至内存, 后续的程序应该都来这里取，以节省io操作次数

	//初始化倒排map
	table.RevIndex = map[string]*index.RevIndex{}

	//倒排文件
	indexDir := path.Join(tableDir, "index")
	utils.DirCHeckAndMk(indexDir)
	files := utils.Walk(indexDir) //倒排索引存放位置  // ./data/abc/index/rev_def
	for _, revfilepath := range files {
		if strings.Index(revfilepath, "rev_") != -1 { //倒排索引文件
			revIndex := index.GetRevIndex(revfilepath)
			table.RevIndex[revIndex.Name] = revIndex
		}
	}

	return table
}

// 新增倒排索引
func (t *Table) AddRevIndex(name string) *index.RevIndex {
	revIndexPath := path.Join(BasePathDir, t.Name, "index", "rev_"+name)
	revIndex := index.GetRevIndex(revIndexPath)
	t.RevIndex[name] = revIndex
	return revIndex
}

// Insert 插入数据
func (t *Table) Insert(inData *inputData.InputData) (bool, error) {
	key := inData.Key
	if key == "" {
		return false, errors.New("插入数据的Key不能为空字符串")
	}

	iData, ok := inData.Data.(map[string]interface{})
	if !ok { //断言失败应该报错
		return false, errors.New("插入数据的类型必须为字符串类型的键值对，键值都必须为字符串，不能有复杂层级,收到：" + reflect.TypeOf(inData.Data).String())
	}

	//整理本次会出现的所有索引
	indexSetArr := []string{}
	for i, _ := range iData {
		t.IndexSet.Sets = append(t.IndexSet.Sets, i) //带上上次的结果
		indexSetArr = append(indexSetArr, i)         //本次处理的所有索引，后面处理倒排使用
	}

	//构建倒排索引
	for _, revIndexName := range indexSetArr {
		var revIndex *index.RevIndex
		if _, ok := t.RevIndex[revIndexName]; !ok {
			revIndex = t.AddRevIndex(revIndexName) //新增倒排索引
		} else {
			revIndex = t.RevIndex[revIndexName] //现有倒排索引
		}

		value, ok := iData[revIndexName].(string)
		if !ok { //断言失败
			return false, errors.New("插入数据的类型必须为字符串类型的键值对，键值都必须为字符串，不能有复杂层级,收到：" + reflect.TypeOf(iData[revIndexName]).String())
		}

		//不存在
		if _, ok := revIndex.Data[value]; !ok {
			revIndex.Data[value] = []string{}
		}

		revIndex.Data[value] = append(revIndex.Data[value], key)
	}

	return true, nil
}

// Save 数据保存
func (t *Table) Save() {
	t.mu.Lock()
	defer t.mu.Unlock()

	//索引去重
	t.IndexSet.Sets = utils.ArrayUnique(t.IndexSet.Sets)
	t.IndexSet.Save()

	//倒排去重
	for index, revIndex := range t.RevIndex {
		for val, _ := range revIndex.Data {
			t.RevIndex[index].Data[val] = utils.ArrayUnique(t.RevIndex[index].Data[val])
			t.RevIndex[index].Save()
		}
	}
}

// CheckIndexExist 检查索引是否存在
func (t *Table) CheckIndexExist(indexName string) bool {
	for _, index := range t.IndexSet.Sets {
		if index == indexName {
			return true
		}
	}

	return false
}

// Search 单值查询
func (t *Table) Search(indexName string, value string) ([]string, error) {
	ret := []string{}
	indexFound := t.CheckIndexExist(indexName)
	if !indexFound {
		return ret, errors.New("未找到相关索引【" + indexName + "】")
	}

	return t._search(indexName, value)
}

// Search 单值查询私有方法，用于减少不必要的查询
func (t *Table) _search(indexName string, value string) ([]string, error) {
	ret := []string{}
	revIndex := t.RevIndex[indexName]
	if _, ok := revIndex.Data[value]; !ok { //未找到indexName 下值为value的所有key
		return ret, nil
	}

	return revIndex.Data[value], nil
}

// MultiSearch 单索引复合查询
func (t *Table) MultiSearch(indexName string, value []string) ([]string, error) {
	ret := []string{}
	indexFound := t.CheckIndexExist(indexName)
	if !indexFound {
		return ret, errors.New("未找到相关索引【" + indexName + "】")
	}

	for _, val := range value {
		res, err := t._search(indexName, val)
		if err != nil {
			return []string{}, err
		}

		ret = append(ret, res...)
	}

	return utils.ArrayUnique(ret), nil
}

// AllIndex 获取所有index的倒排索引 todo::考虑limit
func (t *Table) AllIndex(limit int) map[string]map[string][]string {
	ret := map[string]map[string][]string{}
	for _, index := range t.IndexSet.Sets {
		revIndex := t.RevIndex[index]
		//ret[index] = revIndex.Data
		for val, keys := range revIndex.Data {
			if limit <= len(keys) {
				ret[index][val] = keys[0:limit]
			} else {
				ret[index][val] = keys
			}
		}
	}
	return ret
}

// AllIndexCount 获取所有索引的统计
func (t *Table) AllIndexCount() map[string]map[string]int {
	var ret map[string]map[string]int
	ret = make(map[string]map[string]int)
	for _, index := range t.IndexSet.Sets {
		revIndex := t.RevIndex[index]
		ret[index] = map[string]int{}
		for val, keys := range revIndex.Data {
			ret[index][val] = len(keys)
		}
	}

	return ret
}
