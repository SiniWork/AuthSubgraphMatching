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

	trie := mpt.NewTrie()
	// test 1
	trie.Insert([]byte{'A', 'B', 'C', 'D'}, "hello")
	trie.Insert([]byte{'A', 'B'}, "siyu")
	trie.Insert([]byte{'A', 'B', 'C'}, "world")
	trie.Insert([]byte{'B', 'B', 'C', 'D'}, "li")

	fmt.Println("-------------")
	v2, _ := trie.GetExactOne([]byte{'B', 'B', 'C', 'D'})
	fmt.Println(v2)

	// test 2
	//trie.Insert([]byte{'A', 'B'}, "hello")
	//trie.Insert([]byte{'B', 'B', 'C'}, "world")
	//trie.Insert([]byte{'B', 'B', 'C', 'D'}, "siyu")
	//fmt.Println("-------------")
	//v2, _ := trie.GetExactOne([]byte{'B', 'B', 'C'})
	//v2, _ := trie.GetExactOne([]byte{'B', 'B', 'C'})
	//fmt.Println(v2)

	// test 3
	//trie.Insert([]byte{'A', 'B', 'C'}, "hello")
	//trie.Insert([]byte{'B', 'B', 'C'}, "world")
	//trie.Insert([]byte{'C', 'B', 'C'}, "siyu")
	//fmt.Println("-------------")
	//v2, _ := trie.GetExactOne([]byte{'C', 'B', 'C'})
	//fmt.Println(v2)









}

