package utils

import (
	"fmt"
	"math"
	"sort"
	"time"
)

type wordTfIdf struct {
	nworld string
	value  float64
}

func main() {
	start := currentTimeMillis()
	FeatureSelect(Load())


	cost := currentTimeMillis() - start
	fmt.Printf("耗时 %d ms ",cost)

}

type wordTfIdfs []wordTfIdf
type Interface interface {
	Len() int
	Less(i, j int) bool
	Swap(i, j int)
}

func (us wordTfIdfs) Len() int {
	return len(us)
}
func (us wordTfIdfs) Less(i, j int) bool {
	return us[i].value > us[j].value
}
func (us wordTfIdfs) Swap(i, j int) {
	us[i], us[j] = us[j], us[i]
}

func currentTimeMillis() int64 {
	return time.Now().UnixNano() / 1000000
}
func FeatureSelect(list_words [][]string) {
	docFrequency := make(map[string]float64, 0)
	sumWorlds := 0;
	for _, wordList := range list_words {
		for _, v := range wordList {
			docFrequency[v] += 1
			sumWorlds++;
		}
	}
	wordTf := make(map[string]float64)
	for k, _ := range docFrequency {
		wordTf[k] = docFrequency[k] / float64(sumWorlds)
	}
	docNum := float64(len(list_words))
	wordIdf := make(map[string]float64)
	wordDoc := make(map[string]float64, 0)
	for k, _ := range docFrequency {
		for _, v := range list_words {
			for _, vs := range v {
				if (k == vs) {
					wordDoc[k] += 1
					break
				}
			}
		}
	}
	for k, _ := range docFrequency {
		wordIdf[k] = math.Log(docNum / (wordDoc[k] + 1))
	}
	var wordifS wordTfIdfs
	for k, _ := range docFrequency {
		var wti wordTfIdf
		wti.nworld = k
		wti.value = wordTf[k] * wordIdf[k]
		wordifS = append(wordifS, wti)
	}
	sort.Sort(wordifS)
	fmt.Println(wordifS)
}

func Load() [][]string {
	slice := [][]string{
		{"my", "dog", "has", "flea", "problems", "help", "please"},
		{"maybe", "not", "take", "him", "to", "dog", "park", "stupid"},
		{"my", "dalmation", "is", "so", "cute", "I", "love", "him"},
		{"stop", "posting", "stupid", "worthless", "garbage"},
		{"mr", "licks", "ate", "my", "steak", "how", "to", "stop", "him"},
		{"quit", "buying", "worthless", "dog", "food", "stupid"},
	}
	return slice
}