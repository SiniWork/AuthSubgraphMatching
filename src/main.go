package main

import (
	"Corgi/src/matching"
	"Corgi/src/mpt"
	"fmt"
)

func main() {

	/*
	graph test
	 */
	fmt.Println("----------------graph below--------------")
	g := new(matching.Graph)
	g.LoadGraphFromTxt("./data/example1.txt")
	g.AssignLabel("./data/example1_label.txt")
	g.StatisticNeiStr()
	//g.Print()
	qG := matching.QueryPreProcessing("./data/query1.txt", "./data/query1_label.txt")

	/*
	mpt test
	*/
	fmt.Println("----------------trie below--------------")
	fmt.Println("----------------building trie according to above graph--------------")
	trie := mpt.NewTrie()
	for k, v := range g.NeiStr {
		byteKey := []byte(k)
		for _, n := range v {
			trie.Insert(byteKey, n)
		}
	}
	fmt.Println("----------------getting candidates for each query vertex--------------")
	var candiList [][]int
	for _, each := range qG.CQVList {
		each.Candidates = trie.GetCandidate(each.Base.OneHopStr)
		//fmt.Println(each.Base.Id, each.Candidates)
		candiList = append(candiList, each.Candidates)
	}
	matching.AttachCandidate(candiList, &qG)

	/*
	matching test
	*/
	fmt.Println("----------------matching below--------------")
	fmt.Println("----------------getting matched graph for query graph--------------")
	//g.ObtainMatchedGraphs(qG)
	matched := g.ObtainMatchedGraphs(qG)
	fmt.Println(matched)


	//fmt.Println("----------------------computing hash-------------------------------")
	//rootHash, _ := trie.HashRoot()
	//fmt.Println(rootHash)

	// fmt.Println("----------------------proving test---------------------------------")
	//key := "CBC"
	//trie.Prove([]byte(key))
	//VO1, _ := trie.Prove([]byte(key))
	//fmt.Println(mpt.VerifyProof(rootHash, []byte(key), VO1))

	//fmt.Println("----------------------finding test---------------------------------")
	//key := "CBC"
	//v2, _ := trie.GetExactOne([]byte(key))
	////v2 := trie.GetCandidate([]byte(key))
	//fmt.Println(v2)

	//fmt.Println("----------------------insert and find test---------------------------------")
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



// other test
//func main() {
//	test := make(map[int][]int)
//	test[0] = []int{3, 4}
//	test[1] = []int{1, 2}
//	var res []map[int]int
//	one := make(map[int]int)
//	matching.Product(test, &res, []int{0, 1}, 0, one)
//	fmt.Println(res)
//}

