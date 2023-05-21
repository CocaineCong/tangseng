package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-ego/gse"
	"github.com/go-ego/gse/hmm/pos"
	_ "github.com/go-sql-driver/mysql" // 初始化
	"gonum.org/v1/gonum/floats"
	"se/internal/utils"
)

var (
	seg    gse.Segmenter
	posSeg pos.Segmenter

	text = "hello world"
)

type InputData struct {
	Key  string            `json:"key"`
	Data map[string]string `json:"data"`
}

// HMM新词发现模式
func cut() {
	// 搜索引擎模式 提供更多可能
	hmm := seg.CutSearch(text, true)
	fmt.Println(hmm)
}

// 使用最短路径和动态规划分词 效率更高
func segcut() {
	// 搜索模式主要用于给搜索引擎提供尽可能多的关键字
	segs := seg.ModeSegment([]byte(text), true)
	fmt.Println(gse.ToSlice(segs, true))
}

func main() {

	query := "什么时候给我 发货"
	// 加载默认英文词典
	_ = seg.LoadDict()
	// 加载简体中文词典
	_ = seg.LoadDict("zh_s")
	// 加载停用词表
	_ = seg.LoadStop()

	segs := seg.ModeSegment([]byte(query), true)
	slice := gse.ToSlice(segs, true)
	// fmt.Println(slice)

	dsn := "root@(localhost:3306)/url_test"
	// 打开数据库链接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("model err:", err)
		return
	}
	// 关闭数据库链接
	defer db.Close()
	// fmt.Println("数据库链接成功")

	var urlId int
	var urlDesc, purl string

	if len(slice) != 0 {
		docs := [][]string{}
		doc := []string{}
		idmp := make(map[string]string)

		for i := 0; i < len(slice); i++ {
			// slice[i]] =
			url := "http://127.0.0.1:5200/index?index=" + slice[i] + "&value=" + slice[i]
			// fmt.Println(url)
			resp, err := http.Get(url)

			if err != nil {
				log.Println(err)
				return
			}
			defer resp.Body.Close()
			buf := bytes.NewBuffer(make([]byte, 0, 512))
			_, _ = buf.ReadFrom(resp.Body)
			// fmt.Println(string(buf.Bytes()))
			mp := make(map[string][]string)
			err = json.Unmarshal(buf.Bytes(), &mp)
			if err != nil {
				// fmt.Println(mp)
				if len(mp["data"]) != 0 {
					for i := 0; i < len(mp["data"]); i++ {
						// fmt.Println(mp["data"][i])
						parse := "select * from url_info where url_id=" + mp["data"][i]
						rows := db.QueryRow(parse)
						_ = rows.Scan(&urlId, &urlDesc, &purl)
						// fmt.Println(url_id, url_desc, purl)
						idmp[urlDesc] = purl
						doc = append(doc, urlDesc)
						segs := seg.ModeSegment([]byte(urlDesc), true)
						slice := gse.ToSlice(segs, true)
						docs = append(docs, slice)
					}
				}
			} else {
				fmt.Println(mp["message"])
			}
		}
		// fmt.Println(docs)
		utils.FeatureSelect(docs)
		scores := utils.DocsScore(slice)
		score := utils.Dense2slice(scores)
		inds := make([]int, len(score))
		floats.Argsort(score, inds)
		last3 := inds[len(inds)-3:]
		utils.Reverse(&last3)
		// fmt.Println(last3)
		fmt.Println("query:", query)
		fmt.Println("最相似的三个结果")
		for _, idx := range last3 {
			fmt.Println(doc[idx], idmp[doc[idx]])
		}
	}
}
