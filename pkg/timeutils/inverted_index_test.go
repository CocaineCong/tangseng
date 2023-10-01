package timeutils

import (
	"fmt"
	"testing"
)

func TestGetTodayFormatDaily(t *testing.T) {
	a := GetTodayDate()
	fmt.Println(a)
}
