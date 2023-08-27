package trie

import (
	"fmt"
	"testing"
)

func TestTrieTree(t *testing.T) {
	// example
	t1 := NewTrie()
	t1.Insert("hello")
	t1.Insert("world")
	fmt.Println("t1")
	t1.Traverse()
	t2 := NewTrie()
	t2.Insert("hello")
	t2.Insert("golang")
	t2.Insert("programming")
	fmt.Println("t2")
	t2.Traverse()

	t1.Merge(t2)
	fmt.Println("t1 merge")
	t1.Traverse()

	r := t1.FindAllByPrefix("he")
	fmt.Println(r)
}

func TestBinaryTree(t *testing.T) {
	str := "啊啊睡觉滴哦叫爱哦大家爱哦大家哦i"
	var chars = []rune(str) // 将str转换为rune数组（[]rune）
	for i := range chars {
		fmt.Println(chars[i], string(chars[i]))
	}
	t2 := NewTrie()
	t2.Insert("啊啊啊")
	t2.Insert("则美")
	t2.Insert("成型")
	t2.Traverse()
}
