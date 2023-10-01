package timeutils

import (
	"github.com/golang-module/carbon"
)

// GetTodayDate return 2023-10-01
func GetTodayDate() string {
	return carbon.Now().ToDateString()
}
