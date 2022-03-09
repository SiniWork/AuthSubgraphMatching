package matching

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
	FP map[int][]int
}

func (g *Graph) AuthMatching(query QueryGraph) Proof {
	/*
	Obtaining all matching results and their verification objects
	*/
	var proof Proof
	proof.FP = make(map[int][]int)
	proof.CSG = make(map[int][]int)

	//obtain RS & CSG
	ResMap := make(map[int]map[int]bool)
	for k, _ := range query.QVList {
		ResMap[k] = make(map[int]bool)
	}
	expandId := GetExpandQueryVertex(query)
	for _, candid := range query.CandidateSets[expandId] {
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
	for u, c := range query.CandidateSets {
		for _, v := range c {
			if yes, _ := ResMap[u][v]; !yes {
				proof.FP[u] = append(proof.FP[u], v)
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
					fg := true
					if query.CandidateSetsB[c][n] { // graph vertex n may belong to the candidate set of query vertex c
						oneVer.CSG[n] = g.adj[n]
						for pre, _ := range preMatched { // check whether the connectivity of query vertex c with its pre vertices and the connectivity of graph vertex n with its correspond pre vertices are consistent
							if query.Matrix[c][pre] && !g.matrix[n][preMatched[pre]] { // not consist
								fg = false
								break
							}
						}
						if fg { // graph vertex n indeed belong to the candidate set of query vertex c
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