package main

import (
	"Corgi/src/mpt"
	"fmt"
)

func main() {

	/*
	graph test
	 */
	//g := new(matching.Graph)
	//g.LoadGraphFromTxt("./data/example1.txt")
	//g.AssignLabel("./data/example1_label.txt")
	//g.ObtainNeiStr()
	//g.Print()

	/*
	mpt test
	 */
	// the range of key: A-Z

	//trie := mpt.NewTrie()
	//trie.Insert([]byte{'A', 'B', 'C'}, "hello")
	//trie.Insert([]byte{'A', 'B', 'C'}, "world")
	//v1, f1 := trie.GetExactOne([]byte{'A', 'B', 'C', 'D'})
	//v2, f2 := trie.GetExactOne([]byte{'A', 'B', 'C'})
	//fmt.Println(v1, f1)
	//fmt.Println(v2, f2)

	bo, i := mpt.ContainJudge([]byte("ABBCDEEFFGH"), []byte("ZZZZZ"))
	fmt.Println(bo, i)





}

