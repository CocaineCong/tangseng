package se

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/pkg/errors"
	"io"
	"os"
	"path"
	"strconv"
	"time"
)

func IntToBytes(n int) []byte {
	data := int64(n)
	byteBuf := bytes.NewBuffer([]byte{})
	err := binary.Write(byteBuf, binary.BigEndian, data)
	if err != nil {
		return nil
	}
	return byteBuf.Bytes()
}

func StrToBytes(s string) []byte {
	byteBuf := bytes.NewBuffer([]byte{})
	err := binary.Write(byteBuf, binary.BigEndian, &s) // nolint:golint,staticcheck
	if err != nil {
		return nil
	}
	return byteBuf.Bytes()
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
	if err != nil { // 文件不存在
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
		return 0, errors.Wrap(err, "os.Stat error")
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, errors.Wrap(errors.Errorf("%s is not a regular file", src), "sourceFile IsRegular error")
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, errors.Wrap(err, "os.open error")
	}
	defer func(source *os.File) {
		err := source.Close()
		if err != nil {
			log.LogrusObj.Errorln(err)
		}
	}(source)

	destination, err := os.Create(dst)
	if err != nil {
		return 0, errors.Wrap(err, "os.Create error")
	}
	defer func(destination *os.File) {
		err := destination.Close()
		if err != nil {
			log.LogrusObj.Errorln(err)
		}
	}(destination)
	nBytes, err := io.Copy(destination, source)
	return nBytes, errors.Wrap(err, "os.Copy error")
}

func ArrayUnique(arr []string) []string {
	size := len(arr)
	result := make([]string, 0, size)
	temp := map[string]struct{}{}
	for i := 0; i < size; i++ {
		if _, ok := temp[arr[i]]; !ok {
			temp[arr[i]] = struct{}{}
			result = append(result, arr[i])
		}
	}
	return result
}
