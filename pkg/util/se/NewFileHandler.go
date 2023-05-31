package se

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type FileHandler struct {
	filePath string
	file     *os.File
}

const FILEDIR = "/data/index/"

// 一个字段一个文件
func NewFileHandler(field string) *FileHandler {
	root := GetPath()
	filePath := root + FILEDIR + field + ".bin"
	var fp *os.File
	var err error
	if FileExist(filePath) {
		fp, err = os.OpenFile(filePath, os.O_RDWR, 0666)
		if err != nil {
			fmt.Println("open file:", err)
		}
	} else {
		fp, err = os.Create(filePath)
		if err != nil {
			fmt.Println("create file:", err)
		}
	}
	fileHandler := new(FileHandler)
	fileHandler.filePath = filePath
	fileHandler.file = fp
	return fileHandler
}

// 从指定的位置读取一个int64
func (fh *FileHandler) ReadInt64(start int64) int64 {
	buf := make([]byte, 8)
	_, err := fh.file.ReadAt(buf, start)
	if err != nil {
		if err == io.EOF {
			return -1
		}
	}
	return bytetoint(buf) // 把读取的字节转为int64
}

// 指定的地方写入int64,不传就获取文件最后的下标
func (fh *FileHandler) WriteInt64(value, start int64) int64 {
	if start < 1 {
		start, _ = fh.file.Seek(0, io.SeekEnd) // 表示0到文件end的偏移量
	}
	b := inttobyte(value)
	_, err := fh.file.WriteAt(b, start) // n表示写入的字节数，data是int64,所以n=8, 使用writeAt不能使用追加模式
	if err != nil {
		fmt.Println(err)
	}
	return start
}

// 从start下标读取len个int64
func (fh *FileHandler) ReadDocIdsArry(start, len int64) []int64 {
	var i int64 = 0
	res := make([]int64, 0, len)
	for ; i < len; i++ {
		start = start + i*8
		num := fh.ReadInt64(start)
		if num <= 0 { // 越界了就直接返回
			break
		}
		res = append(res, num)
	}
	return res
}

func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// []byte 转化 int64
func bytetoint(by []byte) int64 {
	var num int64
	b_buf := bytes.NewBuffer(by)
	binary.Read(b_buf, binary.BigEndian, &num)
	return num
}

// int64 转 []byte
func inttobyte(num int64) []byte {
	b_buf := new(bytes.Buffer)
	binary.Write(b_buf, binary.BigEndian, &num) // num类型不能是int
	return b_buf.Bytes()
}

// 获取当前程序目录
func GetPath() string {
	path, _ := os.Getwd()
	return path
}
