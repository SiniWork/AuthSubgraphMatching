package main

import (
	"Corgi/src/matching"
	"Corgi/src/mpt"
	"fmt"
)

func main(){

	fmt.Println("----------------Loading Graph----------------")
	g := new(matching.Graph)
	dataset := "wn"
	switch dataset {
	case "ex":
		g.LoadUnGraphFromTxt("./data/example1.txt")
		g.AssignLabel("./data/example1_label.txt")
		g.ObtainPathFeature("./data/pf/JExample.json")
	case "ye":
		g.LoadUnGraphFromTxt("./data/yeast.txt")
		g.AssignLabel("./data/yeast_label.txt")
		g.ObtainPathFeature("./data/pf/JYeast.json")
	case "hu":
		g.LoadUnGraphFromTxt("./data/human.txt")
		g.AssignLabel("./data/human_label.txt")
		g.ObtainPathFeature("./data/pf/JHuman.json")
	case "wn":
		g.LoadUnGraphFromTxt("./data/wordnet.txt")
		g.AssignLabel("./data/wordnet_label.txt")
		g.ObtainPathFeature("./data/pf/JWordnet.json")
	case "db":
		g.LoadUnGraphFromTxt("./data/dblp.txt")
		g.AssignLabel("./data/dblp_label.txt")
		g.ObtainPathFeature("./data/pf/JDblp.json")
	case "am":
		g.LoadUnGraphFromTxt("./data/amazon.txt")
		g.AssignLabel("./data/amazon_label.txt")
		g.ObtainPathFeature("./data/pf/JAmazon.json")
	case "yt":
		g.LoadUnGraphFromTxt("./data/youtube.txt")
		g.AssignLabel("./data/youtube_label.txt")
		g.ObtainPathFeature("./data/pf/JYoutube.json")
	case "lj":
		g.LoadUnGraphFromTxt("./data/livejournal.txt")
		g.AssignLabel("./data/livejournal_label.txt")
		g.ObtainPathFeature("./data/pf/JLivejournal.json")
	}

	fmt.Println("----------------Building MVPTree----------------")
	trie := mpt.NewTrie()
	for k, v := range g.NeiStr {
		byteKey := []byte(k)
		for _, e := range v {
			trie.Insert(byteKey, e, g.NeiHashes[e], g.Vertices[e].Content)
		}
	}
	//RD := trie.HashRoot()

	fmt.Println("----------------Loading Query----------------")
	var q matching.QueryGraph
	workload := 1
	switch workload {
	case 1:
		q = matching.LoadProcessing("./data/query1.txt", "./data/query1_label.txt")
	case 2:
		q = matching.LoadProcessing("./data/query2.txt", "./data/query2_label.txt")
	case 3:
		q = matching.LoadProcessing("./data/query3.txt", "./data/query3_label.txt")
	case 4:
		q = matching.LoadProcessing("./data/query4.txt", "./data/query4_label.txt")
	case 5:
		q = matching.LoadProcessing("./data/query5.txt", "./data/query5_label.txt")
	case 6:
		q = matching.LoadProcessing("./data/query6.txt", "./data/query6_label.txt")
	}

	fmt.Println("----------------Test for candidate set filtering----------------")
	CS := make(map[int][]int)
	CSF := make(map[int][]int)

	for str, ul := range q.NeiStr {
		fmt.Println("present key: ", str)
		C := trie.GetCandidate([]byte(str))
		for _, u := range ul {
			CS[u] = C
		}
	}
	total := 0
	for k, lis := range CS {
		fmt.Println(k, len(lis))
		total = total + len(lis)
	}
	fmt.Println(total)
	fmt.Println("after filter: ")
	for k, c := range CS {
		for _, v := range c {
			flag := true
			for pa, num := range q.PathFeature[k] {
				if _, ok := g.PathFeature[v][pa]; !ok {
					flag = false
					break
				} else if len(g.PathFeature[v][pa]) < num {
					flag = false
					break
				}
			}
			if flag {
				CSF[k] = append(CSF[k], v)
			}
		}
	}
	RealSum := 0
	for k, lis := range CSF {
		fmt.Println(k, len(lis))
		RealSum = RealSum + len(lis)
	}
	fmt.Println(RealSum)


	//fmt.Println("----------------Authenticated Filtering----------------")
	//VO := verification.VO{}
	//VO.NodeList = trie.AuthFilter(&q)
	//fmt.Println(q.CandidateSets)
	//
	//fmt.Println("----------------Authenticated Matching----------------")
	//VO2 := g.AuthMatching(q)
	//VO.CSG = VO2.CSG
	//VO.FP = VO2.FP
	//VO.RS = VO2.RS
	//
	//fmt.Println("----------------Verification----------------")
	//F, _ := VO.Authentication(q, RD)
	//fmt.Println(F)

}