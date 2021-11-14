package main

import (
	"Corgi/src/matching"
	"Corgi/src/mpt"
	"fmt"
	"time"
)

/*
part test
 */
//time1: phase 1 time
//time2: phase 1 verify time
//time3: phase 2 time
//time4: phase 2 verify time

func main(){

	fmt.Println("----------------loading graph--------------")
	g := new(matching.Graph)

	g.LoadUnGraphFromTxt("./data/yeast.txt")
	g.AssignLabel("./data/yeast_label.txt")

	//g.LoadUnGraphFromTxt("./data/human.txt")
	//g.AssignLabel("./data/human_label.txt")

	//g.LoadUnGraphFromTxt("./data/wordnet.txt")
	//g.AssignLabel("./data/wordnet_label.txt")

	//g.LoadUnGraphFromTxt("./data/dblp.txt")
	//g.AssignLabel("./data/dblp_label.txt")

	//g.LoadUnGraphFromTxt("./data/amazon.txt")
	//g.AssignLabel("./data/amazon_label.txt")

	//g.LoadUnGraphFromTxt("./data/youtube.txt")
	//g.AssignLabel("./data/youtube_label.txt")

	//g.LoadUnGraphFromTxt("./data/livejournal.txt")
	//g.AssignLabel("./data/livejournal_label.txt")

	g.StatisticNeiStr()

	fmt.Println("----------------building mpt--------------")
	trie := mpt.NewTrie()
	for k, v := range g.NeiStr {
		byteKey := []byte(k)
		for _, n := range v {
			trie.Insert(byteKey, n)
		}
	}

	fmt.Println("----------------loading query--------------")
	qG := matching.QueryPreProcessing("./data/query1.txt", "./data/query1_label.txt")

	//fmt.Println("----------------generating VO1 then verifying it--------------") // get time3 and VO1 size
	//qExId := matching.GetExpandQueryVertex(qG.CQVList)
	//fmt.Println("query vertex string: ", string(qG.CQVList[qExId].Base.OneHopStr))
	//_, VO1, _ := trie.Prove(qG.CQVList[qExId].Base.OneHopStr)
	//rootH, _ := trie.HashRoot()
	//startT3 := time.Now()
	//fmt.Println(mpt.Verify(rootH, qG.CQVList[qExId].Base.OneHopStr, VO1))
	//time3 := time.Since(startT3)
	//fmt.Println("the number of nodes in VO1: ", len(VO1.Nodes))
	//fmt.Println("the size of VO1: ", VO1.Size(), "Byte")
	//fmt.Println("the time of verifying VO1 is: ", time3)


	fmt.Println("----------------generating candidates for each query vertex--------------") // get time1
	var candiList [][]int
	startT1 := time.Now()
	for k, each := range qG.CQVList {
		each.Candidates = trie.GetCandidate(each.Base.OneHopStr)
		fmt.Println("present string: ", string(qG.CQVList[k].Base.OneHopStr), ", its candidates: ", len(each.Candidates))
		candiList = append(candiList, each.Candidates)
	}
	time1 := time.Since(startT1)
	fmt.Println("the time of phase1 is: ", time1)
	matching.AttachCandidate(candiList, &qG)

	//fmt.Println("----------------optimizing test----------------")
	//startT := time.Now()
	//fmt.Println("the number of total results: ", len(g.ObtainMatchedGraphs(qG)))
	////fmt.Println("the number of results: ", len(g.ConObtainMatchedGraphs(qG)))
	//tm := time.Since(startT)
	//fmt.Println("time: ", tm)

	fmt.Println("----------------generating matched graphs for query graph then verifying VO2--------------") // get time2 and time4 as well as VO2 size
	g.ComputingGHash()
	startT2 := time.Now()
	VO2 := g.Prove(qG)
	time2 := time.Since(startT2)
	fmt.Println("the time of phase2 is: ", time2)
	fmt.Println("the number of evidence: ", len(VO2.Evidence))
	fmt.Println("the size of VO2: ", VO2.Size(), "Byte")
	startT4 := time.Now()
	fmt.Println(matching.Verify(VO2, g.GHash, qG))
	time4 := time.Since(startT4)
	fmt.Println("the time of verifying VO2 is: ", time4)

}


/*
other test
*/
//func main() {
//	//tool.ConfigLabelForG("./data/livejournal.txt", "./data/livejournal_label.txt")
//	//fmt.Println(tool.CheckGraphLabel("./data/livejournal.txt", "./data/livejournal_label.txt"))
//}