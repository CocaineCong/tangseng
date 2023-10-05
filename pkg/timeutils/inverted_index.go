package timeutils

import (
	"fmt"

	"github.com/golang-module/carbon"
)

// GetTodayDate return 2023-10-01
func GetTodayDate() string {
	return carbon.Now().ToDateString()
}

// GetMonthDate return 2023-10
func GetMonthDate() string {
	year := carbon.Now().Year()
	month := carbon.Now().Month()
	return fmt.Sprintf("%d-%d", year, month)
}

// GetSeasonDate return 2023-Autumn
func GetSeasonDate() string {
	year := carbon.Now().Year()
	season := carbon.Now().Season()
	return fmt.Sprintf("%d-%s", year, season)
}
