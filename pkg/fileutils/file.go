package fileutils

import (
	"os"
	"path/filepath"
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

// GetFiles 获取文件夹下的所有文件
func GetFiles(folder string) (res []string) {
	files, _ := os.ReadDir(folder)
	folderAbs, _ := filepath.Abs(folder)
	for _, file := range files {
		if file.IsDir() {
			GetFiles(folder + "/" + file.Name())
		} else {
			res = append(res, folderAbs+"/"+file.Name())
		}
	}

	return
}
