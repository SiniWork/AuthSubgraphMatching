package matching

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"io/ioutil"
	"path"
	"sort"
	"strconv"
	"strings"
)

type Vertex struct {
	id int
	label byte
	content string
	hashVal []byte
}

type Graph struct {
	/*
	vertices: vertex list
	adj: the adjacency list
	neiStr: statistic the one-hop neighborhood string for each vertex
	 */

	vertices []Vertex
	adj map[int][]int
	NeiStr map[string][]int
	GHash []byte
}

type QVertex struct {
	base QueryVertex
	candidates []int
}


func (g *Graph) LoadGraphFromTxt(fileName string) error {
	/*
	loading the graph from txt file and saving it into an adjacency list adj
	 */

	g.adj = make(map[int][]int)
	content, err := readTxtFile(fileName)
	if err != nil {
		fmt.Println("Read file error!", err)
		return err
	}
	for _, line := range content {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		edge := strings.Split(line, " ")
		fr, err := strconv.Atoi(edge[0])
		if err != nil {
			return err
		}
		en, err := strconv.Atoi(edge[1])
		if err!= nil {
			return err
		}
		g.adj[fr] = append(g.adj[fr], en)
	}
	return nil
}

func (g *Graph) LoadGraphFromExcel(fileName string) error {
	/*
	loading the graph from Excel file and saving it into an adjacency list adj
	*/
	g.adj = make(map[int][]int)
	xlsx, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Println("Open file error!", err)
		return err
	}
	rows := xlsx.GetRows("Sheet1")
	for _, row := range rows {
		fr, err := strconv.Atoi(row[0])
		if err != nil {
			return err
		}
		en, err := strconv.Atoi(row[1])
		if err!= nil {
			return err
		}
		fmt.Println(fr, en)
		g.adj[fr] = append(g.adj[fr], en)
	}
	return nil
}

func (g *Graph) AssignLabel(fileName string) error {
	/*
	randomly assign a label to each vertex or read the vertex's label from a txt file
	then saving them into a list g.vertices
	 */
	var labelSet []byte
	if fileName == "" {
		labelSet = []byte{'A', 'B', 'C', 'D', 'E', 'F','G', 'H', 'I', 'J', 'K', 'L','M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
	} else {
		content, err := readTxtFile(fileName)
		if err != nil {
			fmt.Println("Read file error!", err)
			return err
		}
		for _, line := range content {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			labelSet = append(labelSet, line[0])
			}
		}
	for i := 0; i < len(g.adj); i++ {
		var v Vertex
		v.id = i
		v.label = labelSet[i]
		v.content = "lsy"
		g.vertices = append(g.vertices, v)
	}
	return nil
}

func (g *Graph) StatisticNeiStr() error {
	/*
	Generating the one-hop neighborhood string for each vertex and saving into a map g.neiStr
	g.neiStr, key is the one-hop neighborhood string, value is the list of vertex that have the same neighborhood string
	 */

	g.NeiStr = make(map[string][]int)
	for k, v := range g.adj {
		str := string(g.vertices[k].label)
		var nei []string
		for _, i := range v {
			nei = append(nei, string(g.vertices[i].label))
		}
		sort.Strings(nei)
		for _, t := range nei {
			str = str + t
		}
		g.NeiStr[str] = append(g.NeiStr[str], k)
	}
	return nil
}

func (g *Graph) Print() error {
	for k, v := range g.NeiStr {
		fmt.Println(k, v)
	}
	return nil
}

func (g *Graph) ComputingGHash() []byte {
	/*
	Computing the GHash
	 */
	var accHashVal []byte
	for i, ver := range g.vertices {
		ver.hashVal = g.computingHashVal(ver)
		if i == 0 {
			accHashVal = ver.hashVal
		} else {
			accHashVal = xor(accHashVal, ver.hashVal)
		}

	}
	g.GHash = crypto.Keccak256(accHashVal)
	return accHashVal
}

func (g *Graph) computingHashVal(v Vertex) []byte {
	/*
	Computing hashVal for vertex v
	 */
	index := v.id
	var outXor = hash(v)
	if len(g.adj[index]) == 0 {
		return nil
	} else {
		for _, nei := range g.adj[index] {
			outXor = xor(outXor, hash(g.vertices[nei]))
		}
	}
	return crypto.Keccak256(outXor)
}

func (g *Graph) ObtainMatchedGraph(query []QVertex) []map[int]int {
	/*
	Obtaining all sub graphs that matched the given query graph in the data graph
	 */
	var result []map[int]int
	expandId := getExpandQueryVertex(query)
	pendingVertex := query[expandId]
	for _, candid := range pendingVertex.candidates {
		res := g.matchingV1(candid, expandId, query)
		result = append(result, res...)
	}
	return result
}

//func (g *Graph) matchingV1(candidateId, expandQId int, queryList []QVertex) []map[int]int {
//	/*
//	Expanding the data graph from the given candidate vertex to obtain matched results
//	 */
//	var result []map[int]int
//	exSet := make(map[int][]int)
//	exSet[1] = g.adj[candidateId]
//	visited := initialVisited(len(g.adj), g.adj[candidateId])
//	visited[candidateId] = true
//	for i:=1; i<=len(queryList[expandQId].base.ExpandLayer); i++ {
//		// get the vertices of the current layer of the data graph
//		var gPresentVer []int
//		if i == 1 {
//			gPresentVer = exSet[i]
//			res := make(map[int]int)
//			res[expandQId] = candidateId
//			result = append(result, res)
//		} else {
//
//		}
//
//		// get the vertices of the current layer of the query graph and the candidates of the vertices
//		qPresentVer := queryList[expandQId].base.ExpandLayer[i]
//		qVerCandi := make(map[int]map[int]int)
//		for _, qV := range qPresentVer {
//			candi := make(map[int]int) // play the role of bloom filter
//			for _, c := range queryList[qV].candidates {
//				candi[c] = c
//			}
//			qVerCandi[qV] = candi
//		}
//
//		// classic the vertices of the current layer of the data graph according to query candidates map
//		matched := make(map[int][]int)
//		for _, gV := range gPresentVer {
//			for qV, qVC := range qVerCandi {
//				if _, ok := qVC[gV]; ok {
//					matched[qV] = append(matched[qV], gV)
//				}
//			}
//		}
//
//		// obtain media result
//
//
//	}
//	return nil
//}

func (g *Graph) matchingV1(candidateId, expandQId int, queryList []QVertex) []map[int]int {
	/*
	Expanding the data graph from the given candidate vertex to obtain matched results
	*/
	var result []map[int]int
	exSet := make(map[int][]int)
	exSet[1] = g.adj[candidateId]
	visited := initialVisited(len(g.adj), g.adj[candidateId])
	visited[candidateId] = true
	for i:=1; i<=len(queryList[expandQId].base.ExpandLayer); i++ {
		// get the vertices of the current layer of the data graph
		var gPresentVer []int
		if i == 1 {
			gPresentVer = exSet[i]
			res := make(map[int]int)
			res[expandQId] = candidateId
			result = append(result, res)
		} else {

		}

		// get the vertices of the current layer of the query graph and the candidates of the vertices
		qPresentVer := queryList[expandQId].base.ExpandLayer[i]
		qVerCandi := make(map[int]map[int]int)
		for _, qV := range qPresentVer {
			candi := make(map[int]int) // play the role of bloom filter
			for _, c := range queryList[qV].candidates {
				candi[c] = c
			}
			qVerCandi[qV] = candi
		}

		// classic the vertices of the current layer of the data graph according to query candidates map
		matched := make(map[int][]int)
		for _, gV := range gPresentVer {
			for qV, qVC := range qVerCandi {
				if _, ok := qVC[gV]; ok {
					matched[qV] = append(matched[qV], gV)
				}
			}
		}
		// obtain media result


	}
	return nil
}

func (g *Graph) matchingV2(expT int, gVer []int, expQId int, queryList []QVertex, visited map[int]bool, preMatched map[int]int) []map[int]int{
	/*
	expT: still needed expanding times
	gVer: the set of vertices that need to be expanded
	expQId: the starting expansion query vertex
	queryList: all the query vertices and their related information
	visited: show whether the vertex has been visited
	preMatched: already matched part
	 */
	if expT > len(queryList[expQId].base.ExpandLayer) {
		return nil
	}
	// get the vertices of the current layer of the data graph
	var gPresentVer []int
	for _, k := range gVer{
		for _, j := range g.adj[k] {
			if !visited[j] {
				visited[j] = true
				gPresentVer = append(gPresentVer, j)
			}
		}
	}
	// get the vertices of the current layer of the query graph and the candidates of the vertices
	qPresentVer := queryList[expQId].base.ExpandLayer[expT]
	qVerCandi := make(map[int]map[int]int)
	for _, qV := range qPresentVer {
		candi := make(map[int]int) // play the role of bloom filter
		for _, c := range queryList[qV].candidates {
			candi[c] = c
		}
		qVerCandi[qV] = candi
	}
	// classic the vertices of the current layer of the data graph according to query candidates map
	matched := make(map[int][]int)
	for _, gV := range gPresentVer {
		for qV, qVC := range qVerCandi {
			if _, ok := qVC[gV]; ok {
				matched[qV] = append(matched[qV], gV)
			}
		}
	}
	// if no matched then return
	for _, v := range matched {
		if len(v) == 0 {
			return nil
		}
	}
	// obtain media results and filter these results
	var media []map[int]int
	sort.Ints(qPresentVer)
	oneMap := make(map[int]int)
	Product(matched, &media, qPresentVer, 0, oneMap)


	return nil
}

// func (g* Graph) Prove() 2


func Product(matchedMap map[int][]int, res *[]map[int]int, qV []int, level int, oneMap map[int]int) {
	if level < len(matchedMap) {
		for i:= 0; i<len(matchedMap[qV[level]]); i++ {
			oneMap[qV[level]] = matchedMap[qV[level]][i]
			Product(matchedMap, res, qV, level+1, oneMap)
		}
	} else {
		*res = append(*res, oneMap)
	}
}


func getExpandQueryVertex(qList []QVertex) int {
	bias := 0.5
	index := 0
	coe := 10000000000.0000
	for i, each := range qList {
		temp := float64(len(each.candidates))*(1-bias) + float64(len(each.base.ExpandLayer))*bias
		if temp < coe {
			index = i
		}
	}
	return index
}

func readTxtFile(filePath string) ([]string, error) {
	fileSuffix :=  path.Ext(filePath)
	result := []string{}
	if fileSuffix == ".txt" {
		cont, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Println("Open file error!", err)
			return result, err
		}
		s := string(cont)
		result = strings.Split(s, "\n")
		return result, nil
	} else {
		return result, errors.New("file format error")
	}
}

func (v *Vertex) Serialize() []byte {
	raw := []interface{}{byte(v.id), v.label, v.content}
	rlp, err := rlp.EncodeToBytes(raw)
	if err != nil {
		panic(err)
	}
	return rlp
}

func hash(v Vertex) []byte {
	return crypto.Keccak256(v.Serialize())
}

func xor(str1, str2 []byte) []byte {
	var res []byte
	if len(str1) != len(str2) {
		return res
	} else {
		for i:=0; i<len(str1); i++ {
			res = append(res, str1[i] ^ str2[i])
		}
	}
	return res
}


