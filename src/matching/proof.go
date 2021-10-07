package matching

import "fmt"


type AuxiliaryInfo struct {
	vertexList map[int]Vertex
	adj map[int][]int
}

type Proof struct {
	MatchedRes []map[int]int
	Aux AuxiliaryInfo
}

func (g *Graph) Prove(query QueryGraph) []Proof {
	/*
	Obtaining all sub graphs that matched the given query graph in the data graph and their auxiliary information
	*/
	var VO []Proof
	expandId := getExpandQueryVertex(query.CQVList)
	pendingVertex := query.CQVList[expandId]
	layers := len(pendingVertex.Base.ExpandLayer)
	for _, candid := range pendingVertex.Candidates {
		fmt.Println(candid)
		VO = append(VO, Proof{MatchedRes: g.matchingV1(candid, expandId, query), Aux: g.getAuxiliaryInfo(candid, layers)})
	}
	return VO
}

func (g *Graph) getAuxiliaryInfo(start, layers int) AuxiliaryInfo {

	aux := AuxiliaryInfo{adj: make(map[int][]int), vertexList: make(map[int]Vertex)}
	visited := make(map[int]bool)
	visited[start] = true
	aux.vertexList[start] = g.vertices[start]
	aux.adj[start] = g.adj[start]

	hopVertices := make(map[int][]int)
	hopVertices[0] = append(hopVertices[0], start)
	for hop:=0; hop<layers; hop++ {
		for _, k := range hopVertices[hop]{
			for _, j := range g.adj[k] {
				if !visited[j] {
					visited[j] = true
					aux.vertexList[j] = g.vertices[j]
					aux.adj[j] = g.adj[j]
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
	return aux
}


func Verify(proof []Proof, gHash []byte) {

}
