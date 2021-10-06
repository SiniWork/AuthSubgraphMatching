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

func (g *Graph) ObtainMatchedGraphs(query QueryGraph) []map[int]int {
	/*
	Obtaining all sub graphs that matched the given query graph in the data graph
	 */
	var result []map[int]int
	expandId := getExpandQueryVertex(query.CQVList)
	pendingVertex := query.CQVList[expandId]
	for _, candid := range pendingVertex.Candidates {
		res := g.matchingV1(candid, expandId, query)
		result = append(result, res...)
	}
	return result
}

func (g *Graph) matchingV1(candidateId, expandQId int, query QueryGraph) []map[int]int {
	/*
	Expanding the data graph from the given candidate vertex to obtain matched results
	*/
	var result []map[int]int
	visited := g.setVisited(candidateId, len(query.CQVList[expandQId].Base.ExpandLayer))

	expL := 1
	var gVer []int
	gVer = append(gVer, candidateId)
	preMatched := make(map[int]int)
	preMatched[expandQId] = candidateId
	g.matchingV2(expL, gVer, expandQId, query, visited, preMatched, &result)

	return result
}

func (g *Graph) matchingV2(expL int, gVer []int, expQId int, query QueryGraph, visited []map[int]bool, preMatched map[int]int, res *[]map[int]int) {
	/*
	expT: still need expanding times
	gVer: the set of vertices that need to be expanded
	expQId: the starting expansion query vertex
	queryList: all the query vertices and their related information
	visited: for checking whether present vertex has been visited in last layer
	preMatched: already matched part
	res: save the result
	 */
	if expL > len(query.CQVList[expQId].Base.ExpandLayer) {
		return
	}

	// get the vertices of the current layer of the data graph
	var gPresentVer []int
	for _, k := range gVer{
		for _, j := range g.adj[k] {
			if !visited[expL-1][j] {
				gPresentVer = append(gPresentVer, j)
			}
		}
	}

	// get the vertices of the current layer of the query graph and the candidates of the vertices
	qPresentVer := query.CQVList[expQId].Base.ExpandLayer[expL]
	qVerCandi := make(map[int]map[int]bool)
	for _, qV := range qPresentVer {
		candi := make(map[int]bool) // play the role of bloom filter
		for _, c := range query.CQVList[qV].Candidates {
			candi[c] = true
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
			return
		}
	}

	// obtain media results and filter these results and present layer vertices
	var media []map[int]int
	sort.Ints(qPresentVer)
	oneMap := make(map[int]int)
	Product(matched, &media, qPresentVer, 0, oneMap)
	var filterMedia []map[int]int
	filterVer := g.Filter(preMatched, media, &filterMedia, query.Adj)
	// if present layer has no media result then return
	if len(filterMedia) == 0 {
		return
	}

	// if present layer is the last layer then add the filterMedia into res
	if expL == len(query.CQVList[expQId].Base.ExpandLayer) {
		*res = append(*res, filterMedia...)
	} else {
		// else continue matching
		for i, eachM := range filterMedia {
			g.matchingV2(expL+1, filterVer[i], expQId, query, visited, eachM, res)
		}
	}
	return
}

func (g *Graph) Filter(preMatched map[int]int, raw []map[int]int, fine *[]map[int]int, qAdj map[int][]int) [][]int{
	/*
	Filter the raw media results rely on the connectivity of query graph and the vertices of present layer
	 */
	var verList [][]int
	var flag = true
	for _, r := range raw {
		if checkDuplicateVal(r) {
			continue
		}
		presentMatched := append2Map(preMatched, r)
		flag = true
		var verL []int
		I:
		for k1, v1 := range presentMatched {
			k1Nei := make(map[int]bool)
			for _, kn := range qAdj[k1]{
				k1Nei[kn] = true
			}
			v1Nei := make(map[int]bool)
			for _, vn := range g.adj[v1]{
				v1Nei[vn] = true
			}
			for k2, v2 := range preMatched {
				if k1 == k2 {
					continue
				} else if !connected(k1Nei, k2, v1Nei, v2){
					flag = false
					break I
				}
			}
		}
		if flag {
			for _, v := range r {
				verL = append(verL, v)
			}
			verList = append(verList, verL)
			*fine = append(*fine, presentMatched)
		}
	}
	return verList
}

func (g *Graph) setVisited(candidateId, layers int) []map[int]bool {
	/*
	Expanding 'layer' times from the given start vertex 'candidateID', and setting the visited status for the vertices of layer
	 */
	var res []map[int]bool
	visi := make(map[int]bool)
	visi[candidateId] = true
	res = append(res, visi)

	hopVertices := make(map[int][]int)
	hopVertices[0] = append(hopVertices[0], candidateId)
	for hop:=0; hop < layers; hop++ {
		visited := make(map[int]bool)
		for _, k := range hopVertices[hop]{
			for _, j := range g.adj[k] {
				if !res[hop][j] {
					visited[j] = true
					hopVertices[hop+1] = append(hopVertices[hop+1], j)
				}
			}
		}
		res = append(res, visited)
	}
	return res
}

func checkDuplicateVal(mp map[int]int) bool {
	/*
	Checking whether the given map has the same value
	 */
	vMp := make(map[int]int)
	for _, v := range mp {
		if _, ok := vMp[v]; ok {
			return true
		} else {
			vMp[v] = 1
		}
	}
	return false
}

func append2Map(mp1, mp2 map[int]int) map[int]int {
	res := make(map[int]int)
	for k, v := range mp2 {
		res[k] = v
	}
	for k, v := range mp1 {
		res[k] = v
	}
	return res
}

func copyMap(orig map[int]int) map[int]int {
	cp := make(map[int]int)
	for k, v := range orig {
		cp[k] = v
	}
	return cp
}

func connected(qNei map[int]bool, qV int, gNei map[int]bool, gV int) bool {
	/*
	Checking whether the connection relationship between the two graph vertices is the same as the two query vertices
	 */
	if qNei[qV] && gNei[gV] {
		return true
	} else if !qNei[qV] && !gNei[gV] {
		return true
	}
	return false
}

func Product(matchedMap map[int][]int, res *[]map[int]int, qV []int, level int, oneMap map[int]int) {
	/*
	Permutation and combination on multiple lists
	 */
	if level < len(matchedMap) {
		for i:= 0; i<len(matchedMap[qV[level]]); i++ {
			oneMap[qV[level]] = matchedMap[qV[level]][i]
			Product(matchedMap, res, qV, level+1, oneMap)
		}
	} else {
		newMp := copyMap(oneMap)
		*res = append(*res, newMp)
	}
}

func getExpandQueryVertex(qList []CandiQVertex) int {
	/*
	Computing the weights for each query vertex and choose the smallest
	 */
	bias := 0.5
	index := 0
	coe := 10000000000.0000
	for i, each := range qList {
		temp := float64(len(each.Candidates))*(1-bias) + float64(len(each.Base.ExpandLayer))*bias
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
	/*
	Computing the XOR result of the given two byte array
	 */
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


