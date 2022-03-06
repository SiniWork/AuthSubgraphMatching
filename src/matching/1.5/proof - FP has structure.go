package matching

import (
	"fmt"
)

type OneProof struct {
	MS []map[int]int
	CSG map[int][]int
}

type Proof struct {
	/*
	RS: save the matching results
	CSG: save the search space
	FP: save the false positive candidate vertices
	*/
	RS []map[int]int
	CSG map[int][]int
	FP map[int]map[int][]int
}

func (g *Graph) AuthMatching(query QueryGraph) Proof {
	/*
	Obtaining all matching results and their verification objects
	*/
	var proof Proof
	proof.FP = make(map[int]map[int][]int)
	proof.CSG = make(map[int][]int)

	// obtain CSG
	AllCSG := make(map[int][]int)
	for _, u := range query.QVList {
		for _, v := range u.Candidates {
			AllCSG[v] = g.adj[v]
		}
	}
	fmt.Println("total adj: ", len(AllCSG))

	// obtain FP
	CSAdj := make(map[int]map[int]bool)
	ResMap := make(map[int]map[int]bool)
	for _, u := range query.QVList {
		CSAdj[u.Id] = make(map[int]bool)
		ResMap[u.Id] = make(map[int]bool)
		for _, v := range u.Candidates {
			for _, vn := range g.adj[v] {
				CSAdj[u.Id][vn] = true
			}
		}
	}
	var flag bool
	newCS := make(map[int][]int)
	for _, u := range query.QVList {
		proof.FP[u.Id] = make(map[int][]int)
		for _, v := range u.Candidates {
			flag = true
			for _, un := range query.Adj[u.Id] {
				if yes, _ := CSAdj[un][v]; !yes {
					proof.FP[u.Id][un] = append(proof.FP[u.Id][un], v)
					flag = false
					break
				}
			}
			if flag {
				newCS[u.Id] = append(newCS[u.Id], v)
			}
		}
	}
	for k, u := range newCS {
		AddCandidate(u, &query, k)
	}

	//obtain RS & CSG
	expandId := GetExpandQueryVertex(query.QVList)
	pendingVertex := query.QVList[expandId]
	for _, candid := range pendingVertex.Candidates {
		oneRes := g.authEE(candid, expandId, query)
		proof.RS = append(proof.RS, oneRes.MS...)
		for _, m := range oneRes.MS {
			for k, v := range m {
				ResMap[k][v] = true
			}
		}
		for k, c := range oneRes.CSG {
			proof.CSG[k] = c
		}
	}

	// complete FP
	index := len(query.QVList)
	for _, u := range query.QVList {
		for _, v := range u.Candidates {
			if yes, _ := ResMap[u.Id][v]; !yes {
				proof.FP[u.Id][index] = append(proof.FP[u.Id][index], v)
			}
		}
	}

	return proof
}

func (g *Graph) authEE(candidateId, expandQId int, query QueryGraph) OneProof {
	/*
	Expanding the data graph from the given candidate vertex to enumerate matching results and collect verification objects
	*/
	var oneProof OneProof

	oneProof.CSG = make(map[int][]int)
	expL := 1
	preMatched := make(map[int]int)
	preMatched[expandQId] = candidateId
	g.authMatch(expL, expandQId, query, preMatched, &oneProof)
	return oneProof
}

func (g *Graph) authMatch(expL int, expQId int, query QueryGraph, preMatched map[int]int, oneVer *OneProof){
	/*
	Authenticated recursively enumerating each layer's matched results then generating final results
	expT: still need expanding times
	expQId: the starting expansion query vertex
	preMatched: already matched part
	oneVer: save the result and auxiliary information
	lastLayerR: the vertex list of result that exist in last layer
	*/

	if expL > len(query.QVList[expQId].ExpandLayer) {
		return
	}
	// 1. get the query vertices of the current layer as well as each vertex's candidate set
	qPresentVer := query.QVList[expQId].ExpandLayer[expL]

	// 2. get the graph vertices of the current layer and classify them (Exploration)
	classes := make(map[int][]int)
	visited := make(map[int]bool)
	for _, v := range preMatched {
		visited[v] = true
	}
	var gVer []int // the graph vertices need to be expanded in current layer
	if expL == 1 {
		gVer = append(gVer, preMatched[expQId])
	} else {
		for _, q := range query.QVList[expQId].ExpandLayer[expL-1] {
			gVer = append(gVer, preMatched[q])
		}
	}
	repeat := make(map[int]bool)  // avoid visited repeat vertex in current layer
	for _, v := range gVer { // expand each graph vertex of current layer
		oneVer.CSG[v] = g.adj[v]
		for _, n := range g.adj[v] {
			if !visited[n] && !repeat[n] { // get one unvisited graph vertex n of the current layer
				repeat[n] = true
				for _, c := range qPresentVer { // check current graph vertex n belong to which query vertex's candidate set
					flag := true
					if query.QVList[c].CandidateB[n] { // graph vertex n may belong to the candidate set of query vertex c
						oneVer.CSG[n] = g.adj[n]
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

	// 3. obtain current layer's matched results (Enumeration)
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
	if expL == len(query.QVList[expQId].ExpandLayer) {
		oneVer.MS = append(oneVer.MS, totalRes...)
		return
	} else {
		// else continue matching
		for _, eachM := range totalRes {
			g.authMatch(expL+1, expQId, query, eachM, oneVer)
		}
	}
}


//
//func Verify(proof Proof, gHash []byte, query QueryGraph) bool {
//	/*
//	Verifying the results whether are correctness and completeness
//	 */
//	var newGHashVal []byte
//	for candi, oneVertex := range proof.Evidence {
//		if !oneVertex.checkingRes(candi, query) {
//			fmt.Println("result check fail")
//			return false
//		}
//	}
//	newGHash := crypto.Keccak256(newGHashVal)
//	if string(gHash) != string(newGHash) {
//		fmt.Println("gHash recompute fail")
//		return false
//	}
//	return true
//}
//
//func (g *OneVertexProof) checkingRes(candidate int, query QueryGraph) bool {
//	/*
//	Checking whether the reMatching results are the same as received results
//	 */
//	var result []map[int]int
//	expandQId := GetExpandQueryVertex(query.QVList)
//	expL := 1
//	preMatched := make(map[int]int)
//	preMatched[expandQId] = candidate
//	//matrix := make(map[int]map[int]bool)
//	//for k, l := range g.Aux.adj {
//	//	matrix[k] = make(map[int]bool)
//	//	for _, v := range l {
//	//		matrix[k][v] = true
//	//	}
//	//}
//	g.reMatching(expL, expandQId, query, preMatched, &result)
//
//	if len(result) != len(g.MatchedRes) {
//		return false
//	} else {
//		for i:=0; i<len(result); i++ {
//			if !reflect.DeepEqual(result[i], g.MatchedRes[i]){
//				return false
//			}
//		}
//	}
//	return true
//}
//
//func (g *OneVertexProof) reMatching(expL int, expQId int, query QueryGraph, preMatched map[int]int, res *[]map[int]int){
//	/*
//	ReComputing results in corresponded subgraph
//	expT: still need expanding times
//	expQId: the starting expansion query vertex
//	preMatched: already matched part
//	res: save the result
//	*/
//	if expL > len(query.QVList[expQId].ExpandLayer) {
//		return
//	}
//	// 1. get the query vertices of the current layer as well as each vertex's candidate set
//	qPresentVer := query.QVList[expQId].ExpandLayer[expL]
//
//	// 2. get the graph vertices of the current layer and classify them
//	classes := make(map[int][]int)
//	visited := make(map[int]bool)
//	for _, v := range preMatched {
//		visited[v] = true
//	}
//	var gVer []int // the graph vertices need to be expanded in current layer
//	if expL == 1 {
//		gVer = append(gVer, preMatched[expQId])
//	} else {
//		for _, q := range query.QVList[expQId].ExpandLayer[expL-1] {
//			gVer = append(gVer, preMatched[q])
//		}
//	}
//	repeat := make(map[int]bool)  // avoid visited repeat vertex in current layer
//	for _, v := range gVer { // expand each graph vertex of current layer
//		for _, n := range g.Aux.adj[v] {
//			if !visited[n] && !repeat[n] { // get one unvisited graph vertex n of the current layer
//				repeat[n] = true
//				for _, c := range qPresentVer { // check current graph vertex n belong to which query vertex's candidate set
//					flag := true
//					if query.QVList[c].CandidateB[n] { // graph vertex n may belong to the candidate set of query vertex c
//						for pre, _ := range preMatched { // check whether the connectivity of query vertex c with its pre vertices and the connectivity of graph vertex n with its correspond pre vertices are consistent
//							if query.Matrix[c][pre] && !g.Aux.matrix[n][preMatched[pre]] { // not consist
//								flag = false
//								break
//							}
//						}
//						if flag { // graph vertex n indeed belong to the candidate set of query vertex c
//							classes[c] = append(classes[c], n)
//						}
//					}
//				}
//			}
//		}
//	}
//	// if one of query vertices' candidate set is empty then return
//	if len(classes) < len(qPresentVer) {
//		return
//	}
//
//	// 3. obtain current layer's matched results
//	curRes := g.reObtainCurRes(classes, query, qPresentVer)
//	// if present layer has no media result then return
//	if len(curRes) == 0 {
//		return
//	}
//
//	// 4. combine current layer's result with pre result
//	totalRes := curRes
//	for _, cur := range totalRes {
//		for k, v := range preMatched {
//			cur[k] = v
//		}
//	}
//
//	// 5. if present layer is the last layer then add the filterMedia into res
//	if expL == len(query.QVList[expQId].ExpandLayer) {
//		*res = append(*res, totalRes...)
//		return
//	} else {
//		// else continue matching
//		for _, eachM := range totalRes {
//			g.reMatching(expL+1, expQId, query, eachM, res)
//		}
//	}
//}
//
//func (g *OneVertexProof) reObtainCurRes(classes map[int][]int, query QueryGraph, qVer []int) []map[int]int {
//	/*
//	Obtain current layer's matched results
//	*/
//
//	var matchedRes []map[int]int
//
//	// find all edges between query vertices in current layer
//	qVerCurAdj := make(map[int][]int)
//	for i:=0; i<len(qVer); i++ {
//		qVerCurAdj[qVer[i]] = []int{}
//		for j:=0; j<len(qVer); j++ {
//			if query.Matrix[qVer[i]][qVer[j]] {
//				qVerCurAdj[qVer[i]] = append(qVerCurAdj[qVer[i]], qVer[j])
//			}
//		}
//	}
//
//	// using BFS find all connected part, meanwhile generating part results
//	visited := make(map[int]bool)
//	var queue []int
//	var partResults []map[int][]int
//	for _, k:= range qVer {
//		if !visited[k] {
//			visited[k] = true
//			queue = append(queue, k)
//			onePartRes := make(map[int][]int)
//			onePartRes[k] = classes[k]
//			//sort.Ints(onePartRes[k])
//			for len(queue) != 0 {
//				v := queue[0]
//				queue = queue[1:]
//				for _, n := range qVerCurAdj[v] {
//					if !visited[n] {
//						visited[n] = true
//						queue = append(queue, n)
//						onePartRes = g.join(onePartRes, v, n, classes[n], qVerCurAdj[n])
//					}
//				}
//			}
//			if len(onePartRes) != 0 {
//				partResults = append(partResults, onePartRes)
//			}
//		}
//	}
//	if len(partResults) == 0 {
//		return matchedRes
//	}
//	// combine all part results
//	var agent []int
//	for _, par := range partResults {
//		for k, _ := range par {
//			agent = append(agent, k)
//			break
//		}
//	}
//	oneRes := make(map[int]int)
//	ProductPlus(partResults, &matchedRes, agent, 0, oneRes)
//	return matchedRes
//}
//
//func (g *OneVertexProof) join(curRes map[int][]int, v1, v2 int, v2Candi, v2Nei []int) map[int][]int {
//	/*
//	Join the vertex v2 to current results
//	*/
//	newCurRes := make(map[int][]int)
//	for i, c1 := range curRes[v1] {
//		for _, c2 := range v2Candi {
//			flag := false
//			if g.Aux.matrix[c1][c2] {
//				flag = true
//				// judge the connectivity with other matching vertices
//				for _, n := range v2Nei { // check each neighbor of v2 whether in matched res or not
//					if _, ok := curRes[n]; ok { // neighbor belong to res
//						if !g.Aux.matrix[curRes[n][i]][c2]{ // the connectivity is not satisfied
//							flag = false
//							break
//						}
//					}
//				}
//				// satisfy the demand so that produce a new match
//				if flag {
//					for k, _ := range curRes {
//						newCurRes[k] = append(newCurRes[k], curRes[k][i])
//					}
//					newCurRes[v2] = append(newCurRes[v2], c2)
//				}
//			}
//		}
//	}
//	return newCurRes
//}
//
//func (p *Proof) Size() (int, int) {
//	/*
//	Counting the size of the Proof and result
//	 */
//	totalSize := 0
//	resultSize := 0
//	adjSize := 0
//	resultNum := make(map[int]bool)
//	const verSize = 18
//	num := make(map[int]bool)
//
//	// total size
//	for _, v := range p.Evidence {
//		for _, m := range v.MatchedRes {
//			for _, r := range m {
//				resultNum[r] = true
//			}
//		}
//		//totalSize = totalSize + len(v.Aux.vertexList) * verSize
//		for i, _ := range v.Aux.vertexList {
//			num[i] = true
//		}
//		for _, lis := range v.Aux.adj {
//			adjSize = adjSize + len(lis) * 8
//		}
//	}
//
//	totalSize = totalSize + len(num) * verSize + adjSize
//	fmt.Println("the number of vertices of ours: ", len(num))
//	// result size
//	resultSize = resultSize + len(resultNum) * verSize
//	fmt.Println("the number of vertices of result: ", len(resultNum))
//	for _, v := range p.Evidence {
//		for k, lis := range v.Aux.adj {
//			if resultNum[k] {
//				resultSize = resultSize + len(lis) * 8
//			}
//		}
//	}
//	return totalSize, resultSize
//}
