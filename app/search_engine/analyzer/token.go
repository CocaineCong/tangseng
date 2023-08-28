package analyzer

import (
	"fmt"
	"strings"
)

// Tokenization 分词返回结构
type Tokenization struct {
	Token    string // 词条
	Position int64  // 词条在文本的位置
	Offset   int64  // 偏移量
}

// Ngram 分词
func Ngram(content string, n int64) ([]Tokenization, error) {
	if n < 1 {
		return nil, fmt.Errorf("ngram n must >= 1")
	}
	content = ignoredChar(content)
	var token []Tokenization
	if n >= int64(len([]rune(content))) {
		token = append(token, Tokenization{content, 0, 0})
		return token, nil
	}

	i := int64(0)
	num := len([]rune(content))
	for i = 0; i < int64(num); i++ {
		t := []rune{}
		if i+n > int64(num) {
			break
		} else {
			t = []rune(content)[i : i+n]
		}
		token = append(token, Tokenization{
			Token:    string(t),
			Position: i,
		})
	}

	return token, nil
}

// GseCut 分词 IK
func GseCut(content string) ([]Tokenization, error) {
	content = ignoredChar(content)
	c := GobalSeg.Segment([]byte(content))
	token := make([]Tokenization, 0)
	for _, v := range c {
		if v.Token().Text() == " " { // 这个空格去掉，英文就断在一起了
			continue
		}
		token = append(token, Tokenization{
			Token:    v.Token().Text(),
			Position: int64(v.Start()),
			Offset:   int64(v.End()),
		})
	}

	return token, nil
}

func ignoredChar(str string) string {
	for _, c := range str {
		switch c {
		case '\f', '\n', '\r', '\t', '\v', '!', '"', '#', '$', '%', '&',
			'\'', '(', ')', '*', '+', ',', '-', '.', '/', ':', ';', '<', '=', '>',
			'?', '@', '[', '\\', '【', '】', ']', '“', '”', '「', '」', '★', '^', '·', '_', '`', '{', '|', '}', '~', '《', '》', '：',
			'（', '）', 0x3000, 0x3001, 0x3002, 0xFF01, 0xFF0C, 0xFF1B, 0xFF1F:
			str = strings.ReplaceAll(str, string(c), "")
		}
	}
	return str
}
