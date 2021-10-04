package matching

type QueryVertex struct {
	Id int
	Label byte
	OneHopStr []byte
	ExpandLayer map[int][]int
}

func QueryPreProcessing(queryFile, queryLabelFile string) []QueryVertex {
	/*
		Preprocessing the query graph
	*/
	var query Graph
	var queryVertices []QueryVertex
	query.LoadGraphFromTxt(queryFile)
	query.AssignLabel(queryLabelFile)

	for _, v := range query.vertices {
		qv := QueryVertex{Id: v.id, Label: v.label}
		qv.OneHopStr = append(qv.OneHopStr, v.label)
		for _, nei := range query.adj[v.id] {
			qv.OneHopStr = append(qv.OneHopStr, query.vertices[nei].label)
		}
		qv.ExpandLayer = expandGraph(v.id, query.adj)
		queryVertices = append(queryVertices, qv)
	}
	return queryVertices
}

func expandGraph(v int, adj map[int][]int) map[int][]int {
	/*
	Expanding the given graph one hop at a time and recoding each hop's vertices, the start hop is 1
	*/
	hopVertices := make(map[int][]int)
	hop := 1
	hopVertices[hop] = adj[v]
	visited := initialVisited(len(adj), adj[v])
	visited[v] = true
	for {
		if allVisited(visited) {
			break
		}
		for _, k := range hopVertices[hop]{
			for _, j := range adj[k] {
				if !visited[j] {
					visited[j] = true
					hopVertices[hop+1] = append(hopVertices[hop+1], j)
				}
			}
		}
		hop++
	}
	return hopVertices
}

func initialVisited(length int, ini []int) map[int]bool {
	visited := make(map[int]bool)
	for i:=0; i<length; i++ {
		visited[i] = false
	}
	for _, e := range ini {
		visited[e] = true
	}
	return visited
}

func allVisited(visi map[int]bool) bool {
	for _, f := range visi {
		if !f {
			return false
		}
	}
	return true
}