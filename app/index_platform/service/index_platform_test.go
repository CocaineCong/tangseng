package service

import (
	"fmt"
	"testing"

	"github.com/golang-module/carbon"
)

func TestCarbon(t *testing.T) {
	a := carbon.NewCarbon().Now().String()
	fmt.Println(a)
}
