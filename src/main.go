package main

import (
	"Corgi/src/matching"
	"Corgi/src/mpt"
	"fmt"
	"strconv"
)

func main() {

	/*
	graph test
	 */
	fmt.Println("----------------graph below--------------")
	g := new(matching.Graph)
	g.LoadGraphFromTxt("./data/example1.txt")
	g.AssignLabel("./data/example1_label.txt")
	g.ObtainNeiStr()
	g.Print()

	/*
	mpt test
	*/
	fmt.Println("----------------trie below--------------")
	trie := mpt.NewTrie()
	for k, v := range g.NeiStr {
		byteKey := []byte(k)
		for _, n := range v {
			trie.Insert(byteKey, strconv.Itoa(n))
		}
	}
	key := "C"
	//v2, _ := trie.GetExactOne([]byte(key))
	v2 := trie.GetCandidate([]byte(key))
	fmt.Println(v2)


	// the range of key: A-Z
	// test 1
	//trie.Insert([]byte{'A', 'B', 'C', 'D'}, "hello")
	//trie.Insert([]byte{'A', 'B'}, "li")
	//trie.Insert([]byte{'A', 'B', 'C'}, "world")
	//trie.Insert([]byte{'B', 'B', 'C', 'D'}, "si")
	//trie.Insert([]byte{'C', 'B', 'C'}, "yu")
	//v2, _ := trie.GetExactOne([]byte{'B', 'B', 'C', 'D'})
	//v2 := trie.GetCandidate([]byte{'A', 'B'})
	//fmt.Println(v2)

	// test 2
	//trie.Insert([]byte{'A', 'B'}, "hello")
	//trie.Insert([]byte{'B', 'B', 'C'}, "world")
	//trie.Insert([]byte{'B', 'B', 'C', 'D'}, "siyu")
	//v2, _ := trie.GetExactOne([]byte{'B', 'B', 'C'})
	//v2, _ := trie.GetExactOne([]byte{'B', 'B', 'C'})
	//fmt.Println(v2)

	// test 3
	//trie.Insert([]byte{'A', 'B', 'C'}, "hello")
	//trie.Insert([]byte{'B', 'B', 'C'}, "world")
	//trie.Insert([]byte{'C', 'B', 'C'}, "siyu")
	//v2, _ := trie.GetExactOne([]byte{'C', 'B', 'C'})
	//fmt.Println(v2)









}

