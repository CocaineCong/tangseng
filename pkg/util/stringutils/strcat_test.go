package stringutils

import (
	"fmt"
	"testing"
)

func TestStrConcat(t *testing.T) {
	a := []string{"123", "321"}
	s := StrConcat(a)
	fmt.Println(s)
}
