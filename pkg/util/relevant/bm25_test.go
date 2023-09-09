package relevant

import (
	"fmt"
	"sort"
	"testing"

	"github.com/CocaineCong/tangseng/app/search_engine/analyzer"
	"github.com/CocaineCong/tangseng/config"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

var bodyRecallReason = []string{
	"呀哈哈，怎么说啊太可恶了王国之泪我还没玩够就要写搜索引擎！！写完这个项目马上去海拉鲁大陆拯救塞尔达！！",
	"我们当然可以使用有序数组，二叉搜索树，哈希表等等来存储所有的用户id。但是无论是有序数组还是二叉搜索树，这两种数据结构都是基于二分查找的思想从中间元素开始查起的。",
	"如果我们使用bit类型来存储，就是原来的32倍了，非常亏贼！而这种以bit为单位构建数组的方案就叫做bitmap，也就是位图。",
	"虽然位图相对于原始数组来说，在元素存储上已经有了很大的优化，但如果我们还想进一步优化存储空间，要怎么做呢？",
	"数组的每个成员是一个链表。该数据结构所容纳的所有元素均包含一个指针，用于元素间的链接。我们根据元素的自身特征把元素分配到不同的链表中去，也是根据这些特征，找到正确的链表，再从链表中找出这个元素",
	"当发生哈希冲突时，重新找到空闲的位置，然后插入元素。寻址方式有多种，常用的有线性寻址、二次方寻址、双重哈希寻址等等",
	"而布隆过滤器和位图最大的区别就是我们不再使用一位来表示一个对象，而是使用N位来表示一个对象。这样两个对象的N位都相同的概率就会大大降低了，就能大大缓解哈希冲突了。",
}

func TestMain(m *testing.M) {
	// 这个文件相对于config.yaml的位置
	re := config.ConfigReader{FileName: "../../../config/config.yaml"}
	config.InitConfigForTest(&re)
	analyzer.InitSeg()
	log.InitLog()
	fmt.Println("Write tests on values: ", config.Conf)
	m.Run()
}

func TestBM25(t *testing.T) {
	corpus, _ := MakeCorpus(bodyRecallReason)
	docs := MakeDocuments(bodyRecallReason, corpus)
	tf := New()

	for _, doc := range docs {
		tf.Add(doc)
	}
	tf.CalculateIDF()
	token := Doc{corpus["王国"]}
	tokenScores := BM25(tf, token, docs, 1.5, 0.75)
	sort.Sort(sort.Reverse(tokenScores))
	fmt.Printf("Top 3 Relevant Docs to \"whenever i find\":\n")
	for _, d := range tokenScores[:3] {
		if d.Score == 0.0 {
			continue
		}
		fmt.Printf("\tID   : %d\n\tScore: %1.3f\n\tDoc  : %q\n", d.ID, d.Score, bodyRecallReason[d.ID])
	}
}
