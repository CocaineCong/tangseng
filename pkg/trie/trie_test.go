package trie

import (
	"bytes"
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
	// str := "啊啊睡觉滴哦叫爱哦大家爱哦大家哦i"
	// var chars = []rune(str) // 将str转换为rune数组（[]rune）
	// for i := range chars {
	// 	fmt.Println(chars[i], string(chars[i]))
	// }
	t2 := NewTrie()
	t2.Insert("啊啊啊")
	t2.Insert("则美")
	t2.Insert("成型")
	t2.Traverse()
	rootByte, err := t2.Root.Children.MarshalJSON()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(rootByte))

	t3 := NewTrie()
	t3.Root = NewTrieNode()
	err = t3.Root.Children.UnmarshalJSON(rootByte)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("starting", t3)
}

func TestJsonMarshall(t *testing.T) {
	a := "{\"则\":{\"is_end\":false,\"children\":{\"美\":{\"is_end\":true,\"children\":{}}}},\"啊\":{\"is_end\":false,\"children\":{\"啊\":{\"is_end\":false,\"children\":{\"啊\":{\"is_end\":true,\"children\":{}}}}}},\"成\":{\"is_end\":false,\"children\":{\"型\":{\"is_end\":true,\"children\":{}}}}}\n"
	b := []byte(a)
	replaced := bytes.Replace(b, []byte("children"), []byte("children_recall"), -1)
	fmt.Println(string(replaced))
	node, err := ParseTrieNode(string(replaced))
	if err != nil {
		fmt.Println(err)
	}

	// 使用转换后的 TrieNode 结构体
	fmt.Println(node)
	trie := NewTrie()
	trie.Root = node
	fmt.Println("Traverse")
	trie.TraverseForRecall()
	alist := trie.FindAllByPrefixForRecall("成")
	fmt.Println(alist)
}
