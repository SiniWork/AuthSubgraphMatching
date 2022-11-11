package main

import (
	"Corgi/eth"
	"Corgi/src/matching"
	"Corgi/src/mvp"
	"Corgi/src/verification"
	"fmt"
)

func main(){

	sample := "./data/"
	dataset := "email" // email, wordnet, dblp, youtube, patents
	workload := dataset+"/query"+"1"
	//fmt.Println("----------------workload: ", workload, " dataset: ", dataset)
	//fmt.Println("----------------Loading Graph----------------")
	g := new(matching.Graph)
	g.LoadUnDireGraphFromTxt(sample+dataset+".txt")
	g.AssignLabel(sample+dataset+"_label.txt")

	fmt.Println("----------------Upload Root Digest of MVPTree----------------")
	trie := mvp.NewTrie()
	for k, v := range g.NeiStr {
		byteKey := []byte(k)
		for _, e := range v {
			trie.Insert(byteKey, e, g.NeiHashes[e], g.Vertices[e].Content)
		}
	}
	RD := trie.HashRoot()
	key := "RootDigest"
	eth.CommitEth(key, string(RD))

	fmt.Println("----------------Query Processing----------------")
	var q matching.QueryGraph
	q = matching.LoadProcessing(sample+"query/"+workload+".txt", sample+"query/"+workload+"_label.txt")

	//fmt.Println("----------------Authenticated Filtering----------------")
	VO := verification.VO{}
	VO.NodeList, VO.NodeListB = trie.AuthFilter(&q)

	//fmt.Println("----------------Authenticated Searching----------------")
	VO2 := matching.Proof{}
	g.AuthMatching(q, &VO2)
	VO.CSG = VO2.CSG
	VO.ExpandID = VO2.ExpandID
	VO.RS = VO2.RS
	//fmt.Println("the number of total results is: ", len(VO.RS))

	fmt.Println("----------------Request Root Digest of MVPTree From On-chain----------------")
	onchainRD := []byte(eth.QueryEth("RootDigest").([]interface{})[1].(string))
	fmt.Println("----------------Verification----------------")
	fmt.Println(VO.Authentication(q, onchainRD))

}