package matching



type Proof struct {
	/*
	RS: save the matching results
	CSG: save the search space
	CSGRe: save the search space that removed duplicate vertices
	*/
	RS []map[int]int
	CSG map[int][]int
	ExpandID int
	CSGRe map[int][]int
}


func (g *Graph) AuthMatching(query QueryGraph, proof *Proof) {
	/*
	Obtaining all matching results and their verification objects
	*/

	// obtain CSG
	proof.CSG = make(map[int][]int)
	for _, cl := range query.CandidateSets {
		for _, c := range cl {
			proof.CSG[c] = g.adj[c]
		}
	}

	// obtain RS
	expandId := g.GetExpandQueryVertex(query)
	proof.ExpandID = expandId
	for _, candid := range query.CandidateSets[expandId] {
		g.authEE(candid, expandId, query, proof)
	}
}

func (g *Graph) authEE(candidateId, expandQId int, query QueryGraph, proof *Proof) {
	/*
	Expanding the data graph from the given candidate vertex to enumerate matching results and collect verification objects
	*/
	expL := 1
	preMatched := make(map[int]int)
	preMatched[expandQId] = candidateId
	g.authMatch(expL, expandQId, query, preMatched, proof)
}

func (g *Graph) authMatch(expL int, expQId int, query QueryGraph, preMatched map[int]int, proof *Proof){
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
	visited := make(map[int]bool)
	for _, v := range preMatched {
		visited[v] = true
	}
	var gVer []int // the graph vertices need to be expanded in current layer
	classes := make(map[int][]int)
	if expL == 1 {
		cVer := preMatched[expQId]
		for _, n := range g.adj[cVer] {
			if !visited[n] { // get one unvisited graph vertex n of the current layer
				for _, c := range qPresentVer { // check current graph vertex n belong to which query vertex's candidate set
					if query.CandidateSetsB[c][n] { // graph vertex n may belong to the candidate set of query vertex c
						//oneVer.CSG[n] = g.adj[n]
						classes[c] = append(classes[c], n)
					}
				}
			}
		}
	} else {
		for _, q := range query.QVList[expQId].PendingExpand[expL-1] {
			gVer = append(gVer, preMatched[q])
		}
		repeat := make(map[int]bool)  // avoid visited repeat vertex in current layer
		for _, v := range gVer { // expand each graph vertex of current layer
			//oneVer.CSG[v] = g.adj[v]
			for _, n := range g.adj[v] {
				if !visited[n] && !repeat[n] { // get one unvisited graph vertex n of the current layer
					repeat[n] = true
					for _, c := range qPresentVer { // check current graph vertex n belong to which query vertex's candidate set
						fg := true
						if query.CandidateSetsB[c][n] { // graph vertex n may belong to the candidate set of query vertex c
							//oneVer.CSG[n] = g.adj[n]
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
	}
	// if one of query vertices' candidate set is empty then return
	if len(classes) < len(qPresentVer) {
		return
	}

	// 3. obtain current layer's matched results (Enumeration)
	var curRes []map[int]int
	curRes = g.ObtainCurRes(classes, query, qPresentVer)

	// if present layer has no media result then return
	if len(curRes) == 0 {
		return
	}
	// 4. combine current layer's result with pre result
	for _, cur := range curRes {
		for k, v := range preMatched {
			cur[k] = v
		}
	}

	// 5. if present layer is the last layer then add the filterMedia into res
	if expL == len(query.QVList[expQId].ExpandLayer) {
		proof.RS = append(proof.RS, curRes...)
		return
	} else {
		// else continue matching
		for _, eachM := range curRes {
			g.authMatch(expL+1, expQId, query, eachM, proof)
		}
	}
}
