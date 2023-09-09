package bloom_filter

import (
	"fmt"
	"testing"
)

func TestBloomFilter(t *testing.T) {
	bf := NewBloomFilter(1000, 0.01)

	bf.Add("呀哈哈哈")
	bf.Add("国王资料")
	bf.Add("国王之泪")

	fmt.Println(bf.Contains("国王"))         // false
	fmt.Println(bf.Contains("哈哈"))         // false
	fmt.Println(bf.Contains("国王之泪"))       // true
	fmt.Println(bf.Contains("durian"))     // false
	fmt.Println(bf.Contains("elderberry")) // false
}
