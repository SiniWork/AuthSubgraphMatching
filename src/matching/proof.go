package matching

import (
	"github.com/ethereum/go-ethereum/crypto"
)


type AuxiliaryInfo struct {
	vertexList map[int]Vertex
	adj map[int][]int
}

type OneGroupProof struct {
	MatchedRes []map[int]int
	Aux AuxiliaryInfo
}

type Proof struct {
	Evidence map[int]OneGroupProof
	RemainGHash []byte
}

func (g *Graph) Prove(query QueryGraph) Proof {
	/*
	Obtaining all sub graphs that matched the given query graph in the data graph and their auxiliary information
	*/
	var VO Proof
	VO.RemainGHash = g.ComputingGHash()
	expandId := getExpandQueryVertex(query.CQVList)
	pendingVertex := query.CQVList[expandId]
	VO.Evidence = make(map[int]OneGroupProof)
	layers := len(pendingVertex.Base.ExpandLayer)
	for _, candid := range pendingVertex.Candidates {
		aux, groupHash := g.getAuxiliaryInfo(candid, layers)
		res := g.matchingV1(candid, expandId, query)
		VO.Evidence[candid] = OneGroupProof{MatchedRes: res, Aux: aux}
		VO.RemainGHash = xor(VO.RemainGHash, groupHash)
	}
	return VO
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
	for candi, oneGroup := range proof.Evidence {
		if !oneGroup.checkingRes(candi, query) {
			return false
		}
		newGHashVal = xor(newGHashVal, oneGroup.Aux.getSubGraphHash())
	}
	newGHash := crypto.Keccak256(newGHashVal)
	if string(gHash) != string(newGHash) {
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
	for k := range au.adj {     // bug in here
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

func (g *OneGroupProof) reMatching(expL int, gVer []int, expQId int, query QueryGraph, visited []map[int]bool, preMatched map[int]int, res *[]map[int]int) {
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
			g.reMatching(expL+1, filterVer[i], expQId, query, visited, eachM, res)
		}
	}
}

func (g *OneGroupProof) checkingRes(candidate int, query QueryGraph) bool {
	/*
	Checking whether the reMatching results are the same as received results
	 */
	var result []map[int]int
	expandQId := getExpandQueryVertex(query.CQVList)
	visited := g.setVisited(candidate, len(query.CQVList[expandQId].Base.ExpandLayer))
	expL := 1
	var gVer []int
	gVer = append(gVer, candidate)
	preMatched := make(map[int]int)
	preMatched[expandQId] = candidate
	g.reMatching(expL, gVer, expandQId, query, visited, preMatched, &result)

	if len(result) != len(g.MatchedRes) {
		return false
	} else {
		//for i:=0; i<len(result); i++ {
		//	if !reflect.DeepEqual(result[i], g.MatchedRes[i]){
		//		return false
		//	}
		//}
	}
	return true
}

func (g *OneGroupProof) setVisited(candidate, layers int) []map[int]bool{
	/*
		Expanding 'layer' times from the given start vertex 'candidateID', and setting the visited status for the vertices of layer
	*/
	var res []map[int]bool
	visi := make(map[int]bool)
	visi[candidate] = true
	res = append(res, visi)

	hopVertices := make(map[int][]int)
	hopVertices[0] = append(hopVertices[0], candidate)
	for hop:=0; hop < layers; hop++ {
		visited := make(map[int]bool)
		for _, k := range hopVertices[hop]{
			for _, j := range g.Aux.adj[k] {
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

func (g *OneGroupProof) filter(preMatched map[int]int, raw []map[int]int, fine *[]map[int]int, qAdj map[int][]int) [][]int {
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