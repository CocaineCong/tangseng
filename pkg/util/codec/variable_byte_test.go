package codec

import (
	"fmt"
	"testing"
)

func TestVBDecode(t *testing.T) {
	// 测试VB编码和解码
	var num uint32 = 333
	fmt.Println("原始整数: ", num)

	encoded := VBEncodeNumber(num)
	fmt.Println("VB编码: ", encoded)

}
