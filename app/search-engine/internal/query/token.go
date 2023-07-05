package query

import (
	"fmt"
	"strings"
)

// Tokenization 分词返回结构
type Tokenization struct {
	Token    string
	Position int64
}

// Ngram 分词 后续看情况换成ik吧..
func Ngram(content string, n int64) ([]Tokenization, error) {
	if n < 1 {
		return nil, fmt.Errorf("ngram n must >= 1")
	}
	content = ignoredChar(content)
	var token []Tokenization
	if n >= int64(len([]rune(content))) {
		token = append(token, Tokenization{content, 0})
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

// GseCut 分词
func GseCut(content string) ([]Tokenization, error) {
	content = ignoredChar(content)
	c := GobalSeg.Pos(content)
	token := make([]Tokenization, 0)
	for _, v := range c {
		token = append(token, Tokenization{
			Token:    v.Text,
			Position: int64(strings.Index(content, v.Text)),
		})
	}

	return token, nil
}

func ignoredChar(str string) string {
	for _, c := range str {
		switch c {
		case ' ', '\f', '\n', '\r', '\t', '\v', '!', '"', '#', '$', '%', '&',
			'\'', '(', ')', '*', '+', ',', '-', '.', '/', ':', ';', '<', '=', '>',
			'?', '@', '[', '\\', ']', '^', '_', '`', '{', '|', '}', '~',
			0x3000, 0x3001, 0x3002, 0xFF08, 0xFF09, 0xFF01, 0xFF0C, 0xFF1A, 0xFF1B, 0xFF1F:
			str = strings.ReplaceAll(str, string(c), "")
		}
	}
	return str
}
