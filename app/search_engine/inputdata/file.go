package inputData

import (
	"os"
	"strings"

	log "github.com/CocaineCong/tangseng/pkg/logger"
)

func ReadFiles(fileName []string) []string {
	docList := make([]string, 0)
	for _, sourceName := range fileName {
		docs := readFile(sourceName)
		if docs != nil && len(docs) > 0 {
			docList = append(docList, docs...)
		}
	}
	return docList
}

// 可改用流读取
func readFile(fileName string) []string {
	content, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	docList := strings.Split(string(content), "\n")
	if len(docList) == 0 {
		log.LogrusObj.Infof("readFile err: %v", "docList is empty\n")
		return nil
	}
	return docList
}
