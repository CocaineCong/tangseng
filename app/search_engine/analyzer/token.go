package analyzer

import (
	"strings"
)

// GseCutForRecall 分词 召回专用
func GseCutForRecall(content string) (token []string, err error) {
	content = ignoredChar(content)
	c := GlobalSega.CutSearch(content, true)
	token = make([]string, 0)
	for _, v := range c {
		if v == " " {
			continue
		}
		token = append(token, v)
	}

	return
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
