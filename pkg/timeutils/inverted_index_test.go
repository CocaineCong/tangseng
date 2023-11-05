package timeutils

import (
	"fmt"
	"testing"
)

func TestGetTodayFormatDaily(t *testing.T) {
	a := GetTodayDate()
	fmt.Println(a)
}

func TestGetMonthDate(t *testing.T) {
	a := GetMonthDate()
	fmt.Println(a)
}

func TestGetSeasonDate(t *testing.T) {
	a := GetSeasonDate()
	fmt.Println(a)
}

func TestGetNowTime(t *testing.T) {
	b := GetNowTime()
	fmt.Println(b)
}
