package se

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func Walk(dirPath string) []string {
	var fileList []string
	filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			fileList = append(fileList, path)
		}
		return nil
	})
	// fmt.Println(len(fileList))

	return fileList
}

func GetMd5(filePath string) string {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	md5h := md5.New()
	io.Copy(md5h, file)
	md5 := hex.EncodeToString(md5h.Sum(nil))
	return md5
}

func DirCHeckAndMk(dir string) {
	_, err := os.Stat(dir)
	if err != nil { // 文件不存在
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			fmt.Printf("permission denied![%v]\n", err)
			panic(err)
		}
	}
}

// ExistFile 判断所给的路径文件/文件夹是否存在
func ExistFile(path string) bool {
	_, err := os.Stat(path) // os.Stat 获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
