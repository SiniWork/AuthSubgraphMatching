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

	//fmt.Println("--------------mpt proving test----------------------------------")
	//rootHash, _ := trie.HashRoot()
	//fmt.Println(rootHash)
	//key := "CBC"
	//_, VO1, _ := trie.Prove([]byte(key))
	//fmt.Println(mpt.Verify(rootHash, []byte(key), VO1))

	fmt.Println("----------------matching proving test------------------------------")
	fmt.Println("----------------getting candidates for each query vertex-----------")
	var candiList [][]int
	for _, each := range qG.CQVList {
		each.Candidates = trie.GetCandidate(each.Base.OneHopStr)
		//fmt.Println(each.Base.Id, each.Candidates)
		candiList = append(candiList, each.Candidates)
	}
	matching.AttachCandidate(candiList, &qG)

	//fmt.Println(g.ComputingGHash())
	VO2 := g.Prove(qG)
	fmt.Println(matching.Verify(VO2, g.GHash, qG))



	/*
	matching test
	*/
	//fmt.Println("----------------matching below--------------")
	//fmt.Println("----------------getting matched graph for query graph--------------")
	//g.ObtainMatchedGraphs(qG)
	//matched := g.ObtainMatchedGraphs(qG)
	//fmt.Println(matched)






	//fmt.Println("----------------------finding test---------------------------------")
	//key := "CBC"
	//v2, _ := trie.GetExactOne([]byte(key))
	////v2 := trie.GetCandidate([]byte(key))
	//fmt.Println(v2)

	//fmt.Println("----------------------insert and find test-------------------------")
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
//}
