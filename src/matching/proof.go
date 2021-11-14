package matching

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"reflect"
)


type AuxiliaryInfo struct {
	vertexList map[int]Vertex
	adj map[int][]int
}

type OneVertexProof struct {
	MatchedRes []map[int]int
	Aux AuxiliaryInfo
}

type Proof struct {
	Evidence map[int]OneVertexProof
	RemainGHash []byte
}

func (g *Graph) Prove(query QueryGraph) Proof {
	/*
	Obtaining all sub graphs that matched the given query graph in the data graph and their auxiliary information
	*/
	var VO Proof
	VO.RemainGHash = g.ComputingGHash()
	expandId := GetExpandQueryVertex(query.CQVList)
	pendingVertex := query.CQVList[expandId]
	VO.Evidence = make(map[int]OneVertexProof)
	for _, candid := range pendingVertex.Candidates {
		oneVerProof, subGHash := g.authExpandOneVertex(candid, expandId, query)
		VO.Evidence[candid] = oneVerProof
		VO.RemainGHash = xor(VO.RemainGHash, subGHash)
	}
	return VO
}

func (g *Graph) authExpandOneVertex(candidateId, expandQId int, query QueryGraph) (OneVertexProof, []byte) {
	/*
		Expanding the data graph from the given candidate vertex to obtain matched results and verification objects
	*/
	var oneVerProof OneVertexProof
	expL := 1
	preMatched := make(map[int]int)
	preMatched[expandQId] = candidateId
	oneVerProof.Aux.adj = make(map[int][]int)
	oneVerProof.Aux.adj[candidateId] = g.adj[candidateId]
	oneVerProof.Aux.vertexList = make(map[int]Vertex)
	g.authMatching(expL, expandQId, query, preMatched, &oneVerProof)

	subGHash := g.computingHashVal(g.vertices[candidateId])
	for v, l := range oneVerProof.Aux.adj {
		if v != candidateId {
			subGHash = xor(subGHash, g.computingHashVal(g.vertices[v]))
		}
		for _, e := range l {
			oneVerProof.Aux.vertexList[e] = g.vertices[e]
		}
	}
	return oneVerProof, subGHash
}

func (g *Graph) authMatching(expL int, expQId int, query QueryGraph, preMatched map[int]int, oneVer *OneVertexProof){
	/*
		expT: still need expanding times
		gVer: the set of vertices that need to be expanded
		expQId: the starting expansion query vertex
		preMatched: already matched part
		res: save the result
	*/
	if expL > len(query.CQVList[expQId].Base.ExpandLayer) {
		return
	}
	// 1. get the query vertices of the current layer as well as each vertex's candidate set
	qPresentVer := query.CQVList[expQId].Base.ExpandLayer[expL]

	// 2. get the graph vertices of the current layer and classify them
	classes := make(map[int][]int)
	visited := make(map[int]bool)
	for _, v := range preMatched {
		visited[v] = true
	}
	var gVer []int // the graph vertices need to be expanded in current layer
	if expL == 1 {
		gVer = append(gVer, preMatched[expQId])
	} else {
		for _, q := range query.CQVList[expQId].Base.ExpandLayer[expL-1] {
			gVer = append(gVer, preMatched[q])
		}
	}
	repeat := make(map[int]bool)  // avoid visited repeat vertex in current layer
	for _, v := range gVer { // expand each graph vertex of current layer
		for _, n := range g.adj[v] {
			if !visited[n] && !repeat[n] { // get one unvisited graph vertex n of the current layer
				oneVer.Aux.adj[n] = g.adj[n]
				repeat[n] = true
				for _, c := range qPresentVer { // check current graph vertex n belong to which query vertex's candidate set
					flag := true
					if query.CQVList[c].CandidateB[n] { // graph vertex n may belong to the candidate set of query vertex c
						for pre, _ := range preMatched { // check whether the connectivity of query vertex c with its pre vertices and the connectivity of graph vertex n with its correspond pre vertices are consistent
							if query.Matrix[c][pre] && !g.matrix[n][preMatched[pre]] { // not consist
								flag = false
								break
							}
						}
						if flag { // graph vertex n indeed belong to the candidate set of query vertex c
							classes[c] = append(classes[c], n)
						}
					}
				}
			}
		}
	}
	// if one of query vertices' candidate set is empty then return
	if len(classes) < len(qPresentVer) {
		return
	}

	// 3. obtain current layer's matched results
	curRes := g.ObtainCurRes(classes, query, qPresentVer)
	// if present layer has no media result then return
	if len(curRes) == 0 {
		return
	}

	// 4. combine current layer's result with pre result
	totalRes := curRes
	for _, cur := range totalRes {
		for k, v := range preMatched {
			cur[k] = v
		}
	}

	// 5. if present layer is the last layer then add the filterMedia into res
	if expL == len(query.CQVList[expQId].Base.ExpandLayer) {
		oneVer.MatchedRes = append(oneVer.MatchedRes, totalRes...)
		return
	} else {
		// else continue matching
		for _, eachM := range totalRes {
			g.authMatching(expL+1, expQId, query, eachM, oneVer)
		}
	}
}

func (g *Graph) getAuxiliaryInfo(start, layers int) (AuxiliaryInfo, []byte) {
	/*
	Obtaining the ground truth subgraph that contains some matched results, these results are expanded and matched from the same vertex 'start'
	 */
	aux := AuxiliaryInfo{adj: make(map[int][]int), vertexList: make(map[int]Vertex)}
	visited := make(map[int]bool)
	visited[start] = true
	aux.vertexList[start] = g.vertices[start]
	aux.adj[start] = g.adj[start]
	groupH := g.computingHashVal(g.vertices[start])

	hopVertices := make(map[int][]int)
	hopVertices[0] = append(hopVertices[0], start)
	for hop:=0; hop<layers; hop++ {
		for _, k := range hopVertices[hop]{
			for _, j := range g.adj[k] {
				if !visited[j] {
					visited[j] = true
					aux.vertexList[j] = g.vertices[j]
					aux.adj[j] = g.adj[j]
					groupH = xor(groupH, g.computingHashVal(g.vertices[j]))
					hopVertices[hop+1] = append(hopVertices[hop+1], j)
				}
			}
		}
	}

	// scan vertices of the last layer to complete the neighbors within one hop of each vertex
	for _, i := range hopVertices[layers] {
		for _, j := range g.adj[i] {
			if !visited[j] {
				aux.vertexList[j] = g.vertices[j]
			}
		}
	}
	return aux, groupH
}

func Verify(proof Proof, gHash []byte, query QueryGraph) bool {
	/*
	Verifying the result is correctness and completeness
	 */
	var newGHashVal []byte
	newGHashVal = proof.RemainGHash
	for candi, oneVertex := range proof.Evidence {
		if !oneVertex.checkingRes(candi, query) {
			return false
		}
		newGHashVal = xor(newGHashVal, oneVertex.Aux.getSubGraphHash())
	}
	newGHash := crypto.Keccak256(newGHashVal)
	if string(gHash) != string(newGHash) {
		fmt.Println("bug in here")
		return false
	}
	return true
}

func (au *AuxiliaryInfo) getSubGraphHash() []byte {
	/*
	Computing the hash of the given subgraph
	 */
	var subHashVal []byte
	i := 0
	for k := range au.adj {
		if i == 0 {
			subHashVal = au.getOneVerHash(k)
			i++
		} else {
			subHashVal = xor(subHashVal, au.getOneVerHash(k))
		}
	}
	return subHashVal
}

func (au *AuxiliaryInfo) getOneVerHash(vId int) []byte {
	var outXor = hash(au.vertexList[vId])
	for _, nei := range au.adj[vId] {
		outXor = xor(outXor, hash(au.vertexList[nei]))
	}
	return crypto.Keccak256(outXor)
}

func (g *OneVertexProof) checkingRes(candidate int, query QueryGraph) bool {
	/*
	Checking whether the reMatching results are the same as received results
	 */
	var result []map[int]int
	expandQId := GetExpandQueryVertex(query.CQVList)
	expL := 1
	preMatched := make(map[int]int)
	preMatched[expandQId] = candidate
	matrix := make(map[int]map[int]bool)
	for k, l := range g.Aux.adj {
		matrix[k] = make(map[int]bool)
		for _, v := range l {
			matrix[k][v] = true
		}
	}
	g.reMatching(expL, expandQId, query, preMatched, matrix, &result)

	if len(result) != len(g.MatchedRes) {
		return false
	} else {
		for i:=0; i<len(result); i++ {
			if !reflect.DeepEqual(result[i], g.MatchedRes[i]){
				return false
			}
		}
	}
	return true
}

func (g *OneVertexProof) reMatching1(expL int, gVer []int, expQId int, query QueryGraph, visited []map[int]bool, preMatched map[int]int, res *[]map[int]int) {
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
	repeat := make(map[int]bool)
	for _, k := range gVer{
		for _, j := range g.Aux.adj[k] {
			if !visited[expL-1][j] && !repeat[j]{
				repeat[j] = true
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

	// classify the vertices of the current layer of the data graph according to query candidates map
	matched := make(map[int][]int)
	//fmt.Println("gPresentVer: ", gPresentVer)
	for _, gV := range gPresentVer {
		for qV, qVC := range qVerCandi {
			if _, ok := qVC[gV]; ok {
				matched[qV] = append(matched[qV], gV)
			}
		}
	}
	// if no matched then return
	//fmt.Println("matched: ", matched)
	if len(matched) < len(qPresentVer) {
		return
	}

	// obtain media results and filter these results and present layer vertices
	var media []map[int]int
	oneMap := make(map[int]int)
	Product(matched, &media, qPresentVer, 0, oneMap)
	//fmt.Println("media result: ", media)
	var filterMedia []map[int]int
	filterVer := g.filter(preMatched, media, &filterMedia, query.Adj)
	//fmt.Println("present result: ", filterMedia)
	// if present layer has no media result then return
	if len(filterMedia) == 0 {
		return
	}

	// if present layer is the last layer then add the filterMedia into res
	if expL == len(query.CQVList[expQId].Base.ExpandLayer) {
		*res = append(*res, filterMedia...)
		return
	} else {
		// else continue matching
		for i, eachM := range filterMedia {
			g.reMatching1(expL+1, filterVer[i], expQId, query, visited, eachM, res)
		}
	}
}

func (g *OneVertexProof) reMatching(expL int, expQId int, query QueryGraph, preMatched map[int]int, matrix map[int]map[int]bool, res *[]map[int]int){
	/*
		expT: still need expanding times
		expQId: the starting expansion query vertex
		preMatched: already matched part
		res: save the result
	*/
	if expL > len(query.CQVList[expQId].Base.ExpandLayer) {
		return
	}
	// 1. get the query vertices of the current layer as well as each vertex's candidate set
	qPresentVer := query.CQVList[expQId].Base.ExpandLayer[expL]

	// 2. get the graph vertices of the current layer and classify them
	classes := make(map[int][]int)
	visited := make(map[int]bool)
	for _, v := range preMatched {
		visited[v] = true
	}
	var gVer []int // the graph vertices need to be expanded in current layer
	if expL == 1 {
		gVer = append(gVer, preMatched[expQId])
	} else {
		for _, q := range query.CQVList[expQId].Base.ExpandLayer[expL-1] {
			gVer = append(gVer, preMatched[q])
		}
	}
	repeat := make(map[int]bool)  // avoid visited repeat vertex in current layer
	for _, v := range gVer { // expand each graph vertex of current layer
		for _, n := range g.Aux.adj[v] {
			if !visited[n] && !repeat[n] { // get one unvisited graph vertex n of the current layer
				repeat[n] = true
				for _, c := range qPresentVer { // check current graph vertex n belong to which query vertex's candidate set
					flag := true
					if query.CQVList[c].CandidateB[n] { // graph vertex n may belong to the candidate set of query vertex c
						for pre, _ := range preMatched { // check whether the connectivity of query vertex c with its pre vertices and the connectivity of graph vertex n with its correspond pre vertices are consistent
							if query.Matrix[c][pre] && !matrix[n][preMatched[pre]] { // not consist
								flag = false
								break
							}
						}
						if flag { // graph vertex n indeed belong to the candidate set of query vertex c
							classes[c] = append(classes[c], n)
						}
					}
				}
			}
		}
	}
	// if one of query vertices' candidate set is empty then return
	if len(classes) < len(qPresentVer) {
		return
	}

	// 3. obtain current layer's matched results
	curRes := g.reObtainCurRes(classes, query, qPresentVer, matrix)
	// if present layer has no media result then return
	if len(curRes) == 0 {
		return
	}

	// 4. combine current layer's result with pre result
	totalRes := curRes
	for _, cur := range totalRes {
		for k, v := range preMatched {
			cur[k] = v
		}
	}

	// 5. if present layer is the last layer then add the filterMedia into res
	if expL == len(query.CQVList[expQId].Base.ExpandLayer) {
		*res = append(*res, totalRes...)
		return
	} else {
		// else continue matching
		for _, eachM := range totalRes {
			g.reMatching(expL+1, expQId, query, eachM, matrix, res)
		}
	}
}

func (g *OneVertexProof) reObtainCurRes(classes map[int][]int, query QueryGraph, qVer []int, matrix map[int]map[int]bool) []map[int]int {
	/*
		Obtain current layer's matched results
	*/

	var matchedRes []map[int]int

	// find all edges between query vertices in current layer
	qVerCurAdj := make(map[int][]int)
	for i:=0; i<len(qVer); i++ {
		qVerCurAdj[qVer[i]] = []int{}
		for j:=0; j<len(qVer); j++ {
			if query.Matrix[qVer[i]][qVer[j]] {
				qVerCurAdj[qVer[i]] = append(qVerCurAdj[qVer[i]], qVer[j])
			}
		}
	}

	// using BFS find all connected part, meanwhile generating part results
	visited := make(map[int]bool)
	var queue []int
	var partResults []map[int][]int
	for _, k:= range qVer {
		if !visited[k] {
			visited[k] = true
			queue = append(queue, k)
			onePartRes := make(map[int][]int)
			onePartRes[k] = classes[k]
			//sort.Ints(onePartRes[k])
			for len(queue) != 0 {
				v := queue[0]
				queue = queue[1:]
				for _, n := range qVerCurAdj[v] {
					if !visited[n] {
						visited[n] = true
						queue = append(queue, n)
						onePartRes = g.join(onePartRes, v, n, classes[n], qVerCurAdj[n], matrix)
					}
				}
			}
			if len(onePartRes) != 0 {
				partResults = append(partResults, onePartRes)
			}
		}
	}
	if len(partResults) == 0 {
		return matchedRes
	}
	// combine all part results
	var agent []int
	for _, par := range partResults {
		for k, _ := range par {
			agent = append(agent, k)
			break
		}
	}
	oneRes := make(map[int]int)
	ProductPlus(partResults, &matchedRes, agent, 0, oneRes)
	return matchedRes
}

func (g *OneVertexProof) join(curRes map[int][]int, v1, v2 int, v2Candi, v2Nei []int, matrix map[int]map[int]bool) map[int][]int {
	/*
		join the vertex v2 to current results
	*/
	newCurRes := make(map[int][]int)
	for i, c1 := range curRes[v1] {
		for _, c2 := range v2Candi {
			flag := false
			if matrix[c1][c2] {
				flag = true
				// judge the connectivity with other matching vertices
				for _, n := range v2Nei { // check each neighbor of v2 whether in matched res or not
					if _, ok := curRes[n]; ok { // neighbor belong to res
						if !matrix[curRes[n][i]][c2]{ // the connectivity is not satisfied
							flag = false
							break
						}
					}
				}
				// satisfy the demand so that produce a new match
				if flag {
					for k, _ := range curRes {
						newCurRes[k] = append(newCurRes[k], curRes[k][i])
					}
					newCurRes[v2] = append(newCurRes[v2], c2)
				}
			}
		}
	}
	return newCurRes
}

func (g *OneVertexProof) filter(preMatched map[int]int, raw []map[int]int, fine *[]map[int]int, qAdj map[int][]int) [][]int {
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
			for _, vn := range g.Aux.adj[v1]{
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

func (p *Proof) Size() int {
	/*
	Counting the size of the Proof
	 */
	totalSize := 0
	for _, v := range p.Evidence {
		totalSize = totalSize + len(v.Aux.vertexList)*9
		for _, lis := range v.Aux.adj{
			totalSize = totalSize + len(lis)*8
		}
	}
	return totalSize
}