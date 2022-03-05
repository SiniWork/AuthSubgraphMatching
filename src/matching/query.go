package matching

import (
	"sort"
)

type QVertex struct {
	/*
	OneHopStr: the dictionary label sequence of 1-hop neighbor
	ExpandLayer: save each layer's vertices, the layer index start from 1
	*/
	Id int
	Label byte
	OneHopStr []byte
	ExpandLayer map[int][]int
}

type CandiQVertex struct {
	/*
	Base: the original query vertex structure
	Candidates: current vertex's candidate set
	CandidateB: the 'map' format of candidate set, used as bloom filter
	*/
	Base QVertex
	Candidates []int
	CandidateB map[int]bool
}

type QueryGraph struct {
	/*
	CQVList: the original query vertex structure
	Adj: the adjacency list
	Matrix: the 'map' format of Adj
	*/
	CQVList []CandiQVertex
	Adj map[int][]int
	Matrix map[int]map[int]bool
}

func QueryPreProcessing(queryFile, queryLabelFile string) QueryGraph {
	/*
	Preprocessing the query graph
	*/
	var queryG QueryGraph
	var query Graph
	query.LoadUnGraphFromTxt(queryFile)
	query.AssignLabel(queryLabelFile)
	queryG.Adj = query.adj
	queryG.Matrix = query.matrix

	var temp []int
	for k, _ := range query.vertices {
		temp = append(temp, k)
	}
	sort.Ints(temp)
	for _, i := range temp {
		v := query.vertices[i]
		qV := QVertex{Id: v.id, Label: v.label}
		qV.OneHopStr = append(qV.OneHopStr, v.label)
		for _, nei := range query.adj[v.id] {
			qV.OneHopStr = append(qV.OneHopStr, query.vertices[nei].label)
		}
		qV.ExpandLayer = expandGraph(v.id, query.adj)
		cQV := CandiQVertex{Base: qV}
		queryG.CQVList = append(queryG.CQVList, cQV)
	}
	return queryG
}

func AttachCandidate(candiList [][]int, qG *QueryGraph) {
	for k, candiL := range candiList {
		qG.CQVList[k].Candidates = candiL
		qG.CQVList[k].CandidateB = make(map[int]bool)
		for _, v := range candiL {
			qG.CQVList[k].CandidateB[v] = true
		}
	}
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