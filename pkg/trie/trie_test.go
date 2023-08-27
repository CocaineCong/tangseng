package trie

import (
	"bytes"
	"encoding/gob"
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
	t2 := NewTrie()
	t3 := NewTrie()
	t2.Insert("hello")
	t2.Insert("golang")
	t2.Insert("programming")
	buf := new(bytes.Buffer)
	err := gob.NewEncoder(buf).Encode(t2)
	fmt.Println(err)
	err = gob.NewDecoder(buf).Decode(t3)
	fmt.Println(err)

	t3.Traverse()
}
