package codec

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"os"
	"testing"

	"github.com/CocaineCong/tangseng/app/search_engine/types"
	"github.com/CocaineCong/tangseng/config"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

func TestMain(m *testing.M) {
	re := config.ConfigReader{FileName: "../../../config/config.yaml"}
	config.InitConfigForTest(&re)
	log.InitLog()
	fmt.Println("Write tests on values: ", config.Conf)
	m.Run()
}

func TestIntToBytes(t *testing.T) {
	p := new(types.PostingsList)
	buf := bytes.NewBuffer([]byte{})
	p.DocId = 100
	err := BinaryWrite(buf, p.DocId)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(buf)
	var c int64
	err = binary.Read(buf, binary.LittleEndian, &c)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(c)
}

func TestBinarySize(t *testing.T) {
	v := []int{1, 2, 3}
	a := binary.Size(v)
	fmt.Println(a)
}

type Person struct {
	Name string
	Age  uint8
}

func TestBinaryByteRead(t *testing.T) {
	// 创建一个字节流，表示一个 Person 结构体的二进制数据
	buf := bytes.NewBuffer([]byte{
		0x41, 0x6c, 0x69, 0x63, 0x65, // Name: "Alice"
		0x1E, // Age: 30
	})

	var p Person

	// 使用 binary.Read 解析二进制数据到结构体 p 中
	err := binary.Read(buf, binary.BigEndian, &p)
	if err != nil {
		fmt.Println("解析错误:", err)
		return
	}

	fmt.Printf("Name: %s, Age: %d\n", p.Name, p.Age)
}

func TestBinaryByteWrite(t *testing.T) {
	// 创建一个字节流，表示一个 Person 结构体的二进制数据
	buf := bytes.NewBuffer([]byte{})

	var p = Person{
		Name: "FanOne",
		Age:  22,
	}

	err := binary.Write(buf, binary.LittleEndian, p)
	if err != nil {
		fmt.Println("解析错误:", err)
		return
	}
	var p2 = Person{}
	err = binary.Write(buf, binary.LittleEndian, p2)
	if err != nil {
		fmt.Println("err", err)
	}
	fmt.Printf("Name: %s, Age: %d\n", p2.Name, p2.Age)
}

func TestGobEncoding(t *testing.T) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	s := &Person{
		Name: "FanOne",
		Age:  22,
	}
	if err := enc.Encode(s); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v\n", buf.Bytes())
	dec := gob.NewDecoder(buf)
	var s2 *Person
	if err := dec.Decode(&s2); err != nil {
		fmt.Println(err)
	}
	fmt.Println(s2)
}

func TestGobDecoding(t *testing.T) {
	buf := new(bytes.Buffer)
	dec := gob.NewDecoder(buf)
	var s2 *Person
	if err := dec.Decode(&s2); err != nil {
		fmt.Println(err)
	}
	fmt.Println(s2)
}

type TermValue struct {
	DocCount int64
	Offset   int64
	Size     int64
}

func TestGobByte(t *testing.T) {
	a := []byte{12, 255, 129, 2, 1, 2, 255, 130, 0, 1, 4, 0, 0, 8, 255, 130, 0, 2, 254, 27, 86, 46}
	buffer := bytes.NewBuffer(a)
	fmt.Println(buffer)
	p := new(TermValue)
	err := gob.NewDecoder(buffer).Decode(&p)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p)
}

func TestGobFile(t *testing.T) {
	filePath := "../../../app/search-engine/data/db/0.term"
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
	}
	buffer := []byte{}
	_, err = f.Read(buffer)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(buffer))
}

func TestBinaryStruct(t *testing.T) {
	buf := new(bytes.Buffer)
	p := &TermValue{
		DocCount: 100,
		Offset:   5,
		Size:     15,
	}
	err := BinaryEncoding(buf, p)
	if err != nil {
		fmt.Println(err)
	}

	p2 := &TermValue{
		DocCount: 200,
		Offset:   25,
		Size:     25,
	}
	err = BinaryEncoding(buf, p2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(buf.Bytes())

	p3 := new(TermValue)
	err = BinaryDecoding(buf, p3)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p3)

	p4 := new(TermValue)
	err = BinaryDecoding(buf, p4)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p4)

	p5 := new(TermValue)
	err = BinaryDecoding(buf, p5)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p5)
}

func TestBinarySlice(t *testing.T) {
	buf := new(bytes.Buffer)
	p := make([]int, 0)
	for i := 0; i < 5; i++ {
		p = append(p, i)
	}
	fmt.Println(p)
	err := BinaryEncoding(buf, p)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(buf)
	p2 := make([]int, 0)

	err = BinaryDecoding(buf, p2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p2)
}
