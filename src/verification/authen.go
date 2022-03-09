package verification

import (
	"Corgi/src/matching"
	"Corgi/src/mpt"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"reflect"
	"sort"
)

type VO struct {
	/*
	NodeList: save the visited nodes of the MVPTree and the digest of the siblings of the visited nodes
	CSG: save the search space
	FP: save the false positive vertices
	 */
	NodeList map[string]mpt.Proof
	CSG map[int][]int
	FP map[int][]int
	RS []map[int]int
	CSGMatrix map[int]map[int]bool
}

type NodeIndex struct {
	/*
	node: the visited node of MVPTree
	index: the index of the visited node in the list "Proof.Nodes"
	 */
	no mpt.Node
	index int
}

type pendingPath struct {
	key []byte
	node NodeIndex
}

func (vo *VO) Authentication(query matching.QueryGraph, RD []byte) (bool, []map[int]int) {
	/*
	Verifying whether the matching result satisfy correctness and completeness
	*/

	// search MVPTree and recompute the root digest based on VO.N and VO.CSG
	CS := make(map[int][]int)
	for str, ul := range query.NeiStr {
		fmt.Println("verify key: ", str)
		C, RRD := recomputeRD([]byte(str), vo.NodeList[str], vo.CSG)
		if string(RRD) != string(RD) {
			return false, nil
		}
		for _, u := range ul {
			CS[u] = C
		}
	}
	vo.CSGMatrix = make(map[int]map[int]bool)
	for k, j := range vo.CSG {
		vo.CSGMatrix[k] = make(map[int]bool)
		for _, n := range j {
			vo.CSGMatrix[k][n] = true
		}
	}

	// check whether the combination of vertices of VO.FP and RS is the same as the CS
	combinCS := make(map[int][]int)
	combinCSM := make(map[int]map[int]bool)
	for _, u := range query.QVList {
		combinCSM[u.Id] = make(map[int]bool)
	}
	for _, m := range vo.RS {
		for u, v := range m {
			if !combinCSM[u][v] {
				combinCS[u] = append(combinCS[u], v)
			}
		}
	}
	for u, _ := range combinCS {
		combinCS[u] = append(combinCS[u], vo.FP[u]...)
		sort.Ints(combinCS[u])
		sort.Ints(CS[u])
		if !reflect.DeepEqual(combinCS[u], CS[u]) {
			return false, nil
		}
	}

	// check whether the false positive vertices in VO.FP are real false positive vertices
	// reconstruct the result graph RG'
	expandId := matching.GetExpandQueryVertex(query)
	var RS []map[int]int
	for _, candid := range query.CandidateSets[expandId] {
		oneRes := vo.ExEnum(candid, expandId, query)
		RS = append(RS, oneRes.MS...)
	}
	// recount the number of matching results
	if len(RS) != len(vo.RS) {
		return false, nil
	}

	return true, RS
}

func recomputeRD(key []byte, proof mpt.Proof, CSG map[int][]int) ([]int, []byte) {
	/*
	Search the MVPTree to obtain the CS and check the completeness of the CS then recomputing the root digest
	*/

	var nodeList []NodeIndex
	for i, n := range proof.Nodes {
		nId := NodeIndex{no: n, index: i}
		nodeList = append(nodeList, nId)
	}
	rootNode := nodeList[0]

	var result []int
	if len(key) == 0 {
		return nil, nil
	}
	if root, ok := rootNode.no.(*mpt.BranchNode); ok {
		//fmt.Println("branch node")
		rebuildRelation(&rootNode, nodeList, proof.NodeRelationMap)
		nodeInd := nodeList[proof.NodeRelationMap[rootNode.index][int8(key[0] - 'A')]]
		key = key[1:]
		var latence []pendingPath
		for {
			if mpt.IsEmptyNode(nodeInd.no) {
				if len(latence) == 0 {
					return result, root.Hash()
				}
				key = latence[0].key
				nodeInd = latence[0].node
				latence = latence[1:]
			}

			if leaf, ok := nodeInd.no.(*mpt.LeafNode); ok {
				//fmt.Println("leaf node")
				matched := mpt.PrefixMatchedLen(leaf.Path, key)
				if matched == len(key) || mpt.IsContain(leaf.Path[matched:], key[matched:]){
					for k, _ := range leaf.Value {
						result = append(result, k)
						if _, yes := CSG[k]; yes {
							leaf.Value[k] = crypto.Keccak256(Serialize(CSG[k]))
						}
					}

				}
				if len(latence) == 0 {
					return result, root.Hash()
				}
				key = latence[0].key
				nodeInd = latence[0].node
				latence = latence[1:]
				continue
			}

			if branch, ok := nodeInd.no.(*mpt.BranchNode); ok {
				//fmt.Println("branch node")
				rebuildRelation(&nodeInd, nodeList, proof.NodeRelationMap)
				// check the unsatisfied branches are indeed unsatisfied
				if len(key) == 0 {
					//_, pendPaths := checkBranch(key, node, nodeList, relation)
					right, pendPaths := checkBranch(key, nodeInd, nodeList, proof.NodeRelationMap)
					if !right {
						return nil, nil
					}
					latence = append(latence, pendPaths...)
					for k, _ := range branch.Value {
						result = append(result, k)
						if _, yes := CSG[k]; yes {
							branch.Value[k] = crypto.Keccak256(Serialize(CSG[k]))
						}
					}
					if len(latence) == 0 {
						return result, root.Hash()
					}
					key = latence[0].key
					nodeInd = latence[0].node
					latence = latence[1:]
					continue
				} else {
					//_, pendPaths := checkBranch(key, node, nodeList, relation)
					right, pendPaths := checkBranch(key, nodeInd, nodeList, proof.NodeRelationMap)
					if !right {
						return nil, nil
					}
					latence = append(latence, pendPaths...)
					b, remaining := key[0], key[1:]
					key = remaining
					if branch.Branches[b-'A'] != nil {
						nodeInd = nodeList[proof.NodeRelationMap[nodeInd.index][int8(b - 'A')]]
					} else {
						nodeInd = NodeIndex{}
					}
					continue
				}
			}

			if ext, ok := nodeInd.no.(*mpt.ExtensionNode); ok {
				//fmt.Println("extension node")
				rebuildRelation(&nodeInd, nodeList, proof.NodeRelationMap)
				matched := mpt.PrefixMatchedLen(ext.Path, key)
				if matched < len(ext.Path) && matched < len(key){
					if ext.Path[len(ext.Path)-1] < key[matched] {
						key = key[matched:]
						nodeInd = nodeList[proof.NodeRelationMap[nodeInd.index][0]]
						continue
					} else {
						containAll, i := mpt.ContainJudge(ext.Path[matched:], key[matched:])
						if containAll{
							key = []byte{}
							nodeInd = nodeList[proof.NodeRelationMap[nodeInd.index][0]]
							continue
						} else if ext.Path[len(ext.Path)-1] < key[i] {
							key = key[i:]
							nodeInd = nodeList[proof.NodeRelationMap[nodeInd.index][0]]
							continue
						} else {
							if len(latence) == 0 {
								return result, root.Hash()
							}
							key = latence[0].key
							nodeInd = latence[0].node
							latence = latence[1:]
							continue
						}
					}
				} else {
					key = key[matched:]
					nodeInd = nodeList[proof.NodeRelationMap[nodeInd.index][0]]
					continue
				}
			}
		}
	}
	return nil, nil
}

func rebuildRelation(nodeI *NodeIndex, nodeList []NodeIndex, relation map[int]map[int8]int) {
	/*
	Restoring the children of the given node
	*/
	var i int8
	if branch, ok := nodeI.no.(*mpt.BranchNode); ok {
		for i=0; i < mpt.BranchSize; i++ {
			if branch.Branches[i] != nil {
				index := relation[nodeI.index][i]
				branch.Branches[i] = nodeList[index].no
			}
		}
	}
	if ext, ok := (nodeI.no).(*mpt.ExtensionNode); ok {
		index := relation[nodeI.index][0]
		ext.Next = nodeList[index].no
	}
}

func checkBranch(key []byte, b NodeIndex, nodeList []NodeIndex, relation map[int]map[int8]int) (bool, []pendingPath) {
	/*
	Checking whether the branch that needs to be added is empty
	*/
	if branch, yes:= b.no.(*mpt.BranchNode); yes {
		var subBranches int8
		if len(key) == 0 {
			subBranches = mpt.BranchSize
		} else {
			subBranches = int8(len(branch.Branches[:key[0]-'A']))
		}
		var result []pendingPath
		var i int8
		for i=0; i < subBranches; i++ {
			if branch.Branches[i] != nil {
				if _, yes := nodeList[relation[b.index][i]].no.(mpt.HashNode); yes {
					return false, nil
				}
				p := pendingPath{key, nodeList[relation[b.index][i]]}
				result = append(result, p)
			}
		}
		return true, result
	}
	return false, nil
}

func Serialize(nei []int) []byte {
	raw := []interface{}{}
	for _, n := range nei {
		raw = append(raw, byte(n))
	}
	rlp, err := rlp.EncodeToBytes(raw)
	if err != nil {
		panic(err)
	}
	return rlp
}

func (vo *VO) ExEnum(candidateId, expandQId int, query matching.QueryGraph) matching.OneProof {
	/*
		Expanding the data graph from the given candidate vertex to enumerate matching results and collect verification objects
	*/
	var oneProof matching.OneProof

	oneProof.CSG = make(map[int][]int)
	expL := 1
	preMatched := make(map[int]int)
	preMatched[expandQId] = candidateId
	vo.Match(expL, expandQId, query, preMatched, &oneProof)
	return oneProof
}

func (vo *VO) Match(expL int, expQId int, query matching.QueryGraph, preMatched map[int]int, oneVer *matching.OneProof){
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
		for _, n := range vo.CSG[v] {
			if !visited[n] && !repeat[n] { // get one unvisited graph vertex n of the current layer
				repeat[n] = true
				for _, c := range qPresentVer { // check current graph vertex n belong to which query vertex's candidate set
					fg := true
					if query.CandidateSetsB[c][n] { // graph vertex n may belong to the candidate set of query vertex c
						for pre, _ := range preMatched { // check whether the connectivity of query vertex c with its pre vertices and the connectivity of graph vertex n with its correspond pre vertices are consistent
							if query.Matrix[c][pre] && !vo.CSGMatrix[n][preMatched[pre]] { // not consist
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
	curRes := vo.ObtainCurRes(classes, query, qPresentVer)
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
			vo.Match(expL+1, expQId, query, eachM, oneVer)
		}
	}
}

func (vo *VO) ObtainCurRes(classes map[int][]int, query matching.QueryGraph, qVer []int) []map[int]int {
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
						onePartRes = vo.join(onePartRes, v, n, classes[n], qVerCurAdj[n])
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
	matching.ProductPlus(partResults, &matchedRes, agent, 0, oneRes)
	return matchedRes
}

func (vo *VO) join(curRes map[int][]int, v1, v2 int, v2Candi, v2Nei []int) map[int][]int {
	/*
		Join the vertex v2 to current results
	*/
	newCurRes := make(map[int][]int)
	for i, c1 := range curRes[v1] {
		for _, c2 := range v2Candi {
			fg := false
			if vo.CSGMatrix[c1][c2] {
				fg = true
				// judge the connectivity with other matching vertices
				for _, n := range v2Nei { // check each neighbor of v2 whether in matched res or not
					if _, ok := curRes[n]; ok { // neighbor belong to res
						if !vo.CSGMatrix[curRes[n][i]][c2]{ // the connectivity is not satisfied
							fg = false
							break
						}
					}
				}
				// satisfy the demand so that produce a new match
				if fg {
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

//func (p mpt.Proof) Size() (int, int) {
//	/*
//	Counting the size of the Proof
//	*/
//	var totalSize int
//	var resultSize int
//
//	for _, node := range p.Nodes {
//		if leaf, ok := node.(*LeafNode); ok {
//			leafSize := len(leaf.Path) + len(leaf.Value) * 8
//			totalSize = totalSize + leafSize
//			resultSize = resultSize + leafSize - len(leaf.Value) * 8
//		} else if branch, ok := node.(*BranchNode); ok {
//			branchSize := BranchSize * 8 + len(branch.Value) * 8
//			totalSize = totalSize + branchSize
//			resultSize = resultSize + branchSize - len(branch.Value) * 8
//		} else if ext, ok := node.(*ExtensionNode); ok {
//			extSize := len(ext.Path) + 8
//			totalSize = totalSize + extSize
//		} else if hs, ok := node.(HashNode); ok {
//			hashSize := len(hs.hash)
//			totalSize = totalSize + hashSize
//		}
//	}
//	return totalSize, resultSize
//}

