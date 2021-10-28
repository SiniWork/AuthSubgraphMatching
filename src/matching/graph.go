package matching

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"io/ioutil"
	"path"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Vertex struct {
	id int
	label byte
	content string
}

type Graph struct {
	/*
	vertices: vertex list
	adj: the adjacency list
	neiStr: statistic the one-hop neighborhood string for each vertex
	 */

	vertices map[int]Vertex
	adj map[int][]int
	matrix map[int]map[int]bool
	NeiStr map[string][]int
	GHash []byte
}


func (g *Graph) LoadUnGraphFromTxt(fileName string) error {
	/*
	loading the graph from txt file and saving it into an adjacency list adj, the subscripts start from 0
	 */
	g.adj = make(map[int][]int)
	g.matrix = make(map[int]map[int]bool)
	content, err := readTxtFile(fileName)
	if err != nil {
		fmt.Println("Read file error!", err)
		return err
	}
	splitStr := " "
	if find := strings.Contains(content[0], ","); find {
		splitStr = ","
	} else if find := strings.Contains(content[0], "	"); find {
		splitStr = "	"
	}
	// determine edge is one-way (flag = false) or two-way (flag = true)
	var target string
	flag := true
	for i, line := range content {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if i == 0 {
			edge := strings.Split(line, splitStr)
			target = edge[1] + splitStr + edge[0]
			continue
		}
		if line == target {
			flag = false
			break
		}
	}
	if flag { // case1: two-way
		for _, line := range content {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			edge := strings.Split(line, splitStr)
			fr, err := strconv.Atoi(edge[0])
			if err != nil {
				return err
			}
			en, err := strconv.Atoi(edge[1])
			if err!= nil {
				return err
			}
			g.adj[fr] = append(g.adj[fr], en)
			g.adj[en] = append(g.adj[en], fr)
			// build matrix
			if g.matrix[fr] == nil {
				g.matrix[fr] = make(map[int]bool)
			}
			if g.matrix[en] == nil {
				g.matrix[en] = make(map[int]bool)
			}
			g.matrix[fr][en] = true
			g.matrix[en][fr] = true
		}
	} else { // case2: one-way
		for _, line := range content {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			edge := strings.Split(line, splitStr)
			fr, err := strconv.Atoi(edge[0])
			if err != nil {
				return err
			}
			en, err := strconv.Atoi(edge[1])
			if err!= nil {
				return err
			}
			g.adj[fr] = append(g.adj[fr], en)
			// build matrix
			if g.matrix[fr] == nil {
				g.matrix[fr] = make(map[int]bool)
			}
			g.matrix[fr][en] = true
		}
	}
	return nil
}

func (g *Graph) LoadDireGraphFromTxt(fileName string) error {
	/*
		loading the graph from txt file and saving it into an adjacency list adj, the subscripts start from 0
	*/

	g.adj = make(map[int][]int)
	content, err := readTxtFile(fileName)
	if err != nil {
		fmt.Println("Read file error!", err)
		return err
	}
	splitStr := " "
	if find := strings.Contains(content[0], ","); find {
		splitStr = ","
	} else if find := strings.Contains(content[0], "	"); find {
		splitStr = "	"
	}
	for _, line := range content {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		edge := strings.Split(line, splitStr)
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

func (g *Graph) AssignLabel(labelFile string) error {
	/*
	Assigning a label to each vertex
	 */
	g.vertices = make(map[int]Vertex)
	labelSet := make(map[int]string)
	if labelFile != "" {
		content, err := readTxtFile(labelFile)
		if err != nil {
			fmt.Println("Read file error!", err)
			return err
		}
		for _, line := range content {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			onePair := strings.Split(line, " ")
			key, _ := strconv.Atoi(onePair[0])
			labelSet[key] = onePair[1]
		}
	}
	for k, _ := range g.adj {
		var v Vertex
		v.id = k
		v.label = []byte(labelSet[k])[0]
		g.vertices[k] = v
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
		hashVal := g.computingHashVal(ver)
		if i == 0 {
			accHashVal = hashVal
		} else {
			accHashVal = xor(accHashVal, hashVal)
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
	expandId := GetExpandQueryVertex(query.CQVList)
	pendingVertex := query.CQVList[expandId]

	for _, candid := range pendingVertex.Candidates {
		fmt.Println("processing: ", candid, "degree is: ", len(g.adj[candid]))
		res := g.matchingV1(candid, expandId, query)
		fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!the number of generating results: ", len(res))
		result = append(result, res...)
	}

	//fmt.Println(len(g.matchingV1(643, expandId, query)))

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

func (g *Graph) ObtainMatchedGraphsCon(query QueryGraph) []map[int]int {
	/*
		Obtaining all sub graphs that matched the given query graph in the data graph
	*/
	var result []map[int]int
	expandId := GetExpandQueryVertex(query.CQVList)
	pendingVertex := query.CQVList[expandId]
	lenT := len(pendingVertex.Candidates)

	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)
	chs := make([]chan []map[int]int, cpus)
	start := 0
	interval := lenT / cpus
	for i := 0; i < len(chs); i++ {
		chs[i] = make(chan []map[int]int, 1)
		task := pendingVertex.Candidates[start:start+interval]
		start = start + interval
		qExpand := expandId
		qG := query
		go g.matchingV1Con(task, qExpand, qG, chs[i])
	}
	for _, ch := range chs {
		res := <- ch
		result = append(result, res...)
	}
	return result
}

func (g *Graph) matchingV1Con(candidateIdList []int, expandQId int, query QueryGraph, res chan []map[int]int)  {
	/*
		Expanding the data graph from the given candidate vertex to obtain matched results
	*/
	startT1 := time.Now()
	var resultA []map[int]int
	for _, id := range candidateIdList {
		//fmt.Println(id)
		var resultP []map[int]int
		visited := g.setVisited(id, len(query.CQVList[expandQId].Base.ExpandLayer))

		expL := 1
		var gVer []int
		gVer = append(gVer, id)
		preMatched := make(map[int]int)
		preMatched[expandQId] = id
		g.matchingV2(expL, gVer, expandQId, query, visited, preMatched, &resultP)
		resultA = append(resultA, resultP...)
	}
	res <- resultA
	time1 := time.Since(startT1)
	fmt.Println("the time of phase2 is: ", time1)
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

	// 1. get the vertices of the current layer of the data graph
	ts1 := time.Now()
	var gPresentVer []int
	repeat := make(map[int]bool)
	for _, k := range gVer{
		for _, j := range g.adj[k] {
			if !visited[expL-1][j] && !repeat[j]{
				repeat[j] = true
				gPresentVer = append(gPresentVer, j)
			}
		}
	}
	t1 := time.Since(ts1)
	timeMax := t1
	tFlag := 1

	// 2. get the vertices of the current layer of the query graph and the candidates of the vertices
	ts2 := time.Now()
	qPresentVer := query.CQVList[expQId].Base.ExpandLayer[expL]
	qVerCandi := make(map[int]map[int]bool)
	for _, qV := range qPresentVer {
		candi := make(map[int]bool) // play the role of bloom filter
		for _, c := range query.CQVList[qV].Candidates {
			candi[c] = true
		}
		qVerCandi[qV] = candi
	}
	t2 := time.Since(ts2)
	if t2 > timeMax {
		timeMax = t2
		tFlag = 2
	}

	// 3. classify the vertices of the current layer of the data graph according to query candidates map
	ts3 := time.Now()
	matched := make(map[int][]int)
	//fmt.Println("gPresentVer: ", gPresentVer)
	for _, gV := range gPresentVer {
		for qV, qVC := range qVerCandi {
			if _, ok := qVC[gV]; ok {
				matched[qV] = append(matched[qV], gV)
			}
		}
	}
	t3 := time.Since(ts3)
	if t3 > timeMax {
		timeMax = t3
		tFlag = 3
	}
	// if no matched then return
	//fmt.Println("matched: ", matched)
	if len(matched) < len(qPresentVer) {
		return
	}

	// 4. obtain media results
	ts4 := time.Now()
	var media []map[int]int
	oneMap := make(map[int]int)
	Product(matched, &media, qPresentVer, 0, oneMap)
	fmt.Println("media result: ", len(media))
	t4 := time.Since(ts4)
	if t4 > timeMax {
		timeMax = t4
		tFlag = 4
	}

	// 5. filter these results as well as present layer vertices
	ts5 := time.Now()
	var filterMedia []map[int]int
	filterVer := g.Filter(preMatched, media, &filterMedia, query.Adj)
	//fmt.Println("present result: ", filterMedia)
	t5 := time.Since(ts5)
	if t5 > timeMax {
		timeMax = t5
		tFlag = 5
	}
	fmt.Println("the most time consuming step is: ", tFlag, "the time is: ", timeMax)

	// if present layer has no media result then return
	if len(filterMedia) == 0 {
		return
	}

	// 6. if present layer is the last layer then add the filterMedia into res
	if expL == len(query.CQVList[expQId].Base.ExpandLayer) {
		*res = append(*res, filterMedia...)
		return
	} else {
		// else continue matching
		for i, eachM := range filterMedia {
			g.matchingV2(expL+1, filterVer[i], expQId, query, visited, eachM, res)
		}
	}
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
		presentMatched := append2IntMap(preMatched, r)
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
			for k2, v2 := range presentMatched {
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
		for k, v := range res[hop] {
			visited[k] = v
		}
		res = append(res, visited)
	}
	return res
}

func (g *Graph)ObtainCurRes(matchedMap map[int][]int, query QueryGraph, qVer []int) {
	/*
		Obtain current layer's matched results
	*/

	//var matchedRes []map[int]int

	// find all edges in current layer
	qVerCurAdj := make(map[int][]int)
	qVDegree := make(map[int]int) // sort it depend on degree
	for i:=0; i<len(qVer); i++ {
		for j:=i+1; j<len(qVer); j++ {
			if query.Matrix[qVer[i]][qVer[j]] {
				qVerCurAdj[qVer[i]] = append(qVerCurAdj[qVer[i]], qVer[j])
			}
		}
		qVDegree[qVer[i]] = len(qVerCurAdj[qVer[i]])
	}

	// using BFS find all connected part, meanwhile generating part results
	visited := make(map[int]bool)
	var queue []int
	var partResults []map[int][]int
	for k, _:= range qVDegree {
		if !visited[k] {
			visited[k] = true
			queue = append(queue, k)
			onePartRes := make(map[int][]int)
			onePartRes[k] = matchedMap[k]
			sort.Ints(onePartRes[k])
			for len(queue) != 0 {
				v := queue[0]
				queue = queue[1:]
				for _, n := range qVerCurAdj[v] {
					if !visited[n] {
						visited[n] = true
						queue = append(queue, n)
						g.join(onePartRes, n, matchedMap[n], qVerCurAdj[n])
					}
				}
			}
			partResults = append(partResults, onePartRes)
		}
	}
	// combine all part results
}
func (g *Graph)join(curRes map[int][]int, v2 int, v2Candi, v2Nei []int) {
	/*
	join the vertex v2 to current results
	 */

	for _, n := range v2Nei {
		if _, ok := curRes[n]; ok {
			newCurRes := make(map[int][]int)
			for i, c1 := range curRes[n] {
				for _, c2 := range v2Candi {
					if g.matrix[c1][c2]{
						for k, _ := range curRes {
							newCurRes[k] = append(newCurRes[k], curRes[k][i])
						}
						if _, yes := curRes[v2]; !yes {
							newCurRes[v2] = append(newCurRes[v2], c2)
						} else {
							newCurRes[v2][i] = c2
						}
					}
				}
			}
			curRes = newCurRes
		}
	}
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

func append2IntMap(mp1, mp2 map[int]int) map[int]int {
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

func GetExpandQueryVertex(qList []CandiQVertex) int {
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
			coe = temp
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


