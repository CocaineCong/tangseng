package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"time"
)

func IntToBytes(n int) []byte {
	data := int64(n)
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, data)
	return bytebuf.Bytes()
}

func StrToBytes(s string) []byte {
	data := string(s)
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, data)
	return bytebuf.Bytes()
}

func Tmd5() string {
	timeInt := time.Now().Unix()
	return StrToMd5(strconv.Itoa(int(timeInt)))
}

func StrToMd5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func GetWd() string {
	cpath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dataDir := ".tversion"
	dataPath := path.Join(cpath, dataDir)
	_, err = os.Stat(dataPath)
	if err != nil { //文件不存在
		err = os.Mkdir(dataPath, os.ModePerm)
		if err != nil {
			fmt.Printf("permission denied![%v]\n", err)
			panic(err)
		}
	}

	return cpath
}

func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func ArrayUnique(arr []string) []string{
	size := len(arr)
	result := make([]string, 0, size)
	temp := map[string]struct{}{}
	for i:=0; i < size; i++ {
		if _,ok := temp[arr[i]]; ok != true {
			temp[arr[i]] = struct{}{}
			result = append(result, arr[i])
		}
	}
	return result
}