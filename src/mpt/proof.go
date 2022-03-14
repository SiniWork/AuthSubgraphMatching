package mpt

import (
	"Corgi/src/matching"
	"fmt"
)

type Proof struct {
	/*
	Nodes: visited node list of one query vertex
	NodeRelationMap: the connection between the nodes stored in Nodes
	*/
	Nodes []Node
	NodeRelationMap map[int]map[int8]int
}

func (t *Trie) AuthFilter(q *matching.QueryGraph) map[string]Proof {
	/*
	obtaining candidate vertex set and merkle proof for all query vertices
	*/
	nodeList := make(map[string]Proof)
	q.CandidateSets = make(map[int][]int)
	q.CandidateSetsB = make(map[int]map[int]bool)

	for str, ul := range q.NeiStr {
		fmt.Println("present key: ", str)
		C, P, _ := t.AuthSearch([]byte(str))
		nodeList[str] = P
		for _, u := range ul {
			q.CandidateSets[u] = C
			q.CandidateSetsB[u] = make(map[int]bool)
			for _, v := range C {
				q.CandidateSetsB[u][v] = true
			}
		}
	}
	return nodeList
}

func (t *Trie) AuthSearch(key []byte) ([]int, Proof, bool) {
	/*
	obtaining candidate vertex set and merkle proof for the given key (one query vertex)
	*/

	var proof Proof
	nodeExist := make(map[Node]int)
	proof.NodeRelationMap = make(map[int]map[int8]int)
	if len(key) == 0 {
		return nil, proof, false
	}
	var result []int
	if root, ok := t.root.(*BranchNode); ok {
		//fmt.Println("branch node")
		node := root.GetBranch(key[0])
		key = key[1:]
		nodeExist[root] = len(proof.Nodes)
		proof.Nodes = append(proof.Nodes, root)
		var latence []potentialPath
		for {
			if IsEmptyNode(node) {
				if len(latence) == 0 {
					break
				}
				key = latence[0].key
				node = latence[0].node
				latence = latence[1:]
			}

			if leaf, ok := node.(*LeafNode); ok {
				//fmt.Println("leaf node")
				nodeExist[leaf] = len(proof.Nodes)
				proof.Nodes = append(proof.Nodes, leaf)
				matched := PrefixMatchedLen(leaf.Path, key)
				if matched == len(key) || IsContain(leaf.Path[matched:], key[matched:]){
					for k, _ := range leaf.Value {
						result = append(result, k)
					}
				}
				if len(latence) == 0 {
					break
				}
				key = latence[0].key
				node = latence[0].node
				latence = latence[1:]
				continue
			}

			if branch, ok := node.(*BranchNode); ok {
				//fmt.Println("branch node")
				nodeExist[branch] = len(proof.Nodes)
				proof.Nodes = append(proof.Nodes, branch)
				if len(key) == 0 {
					latence = append(latence, ToBeAdd(key, *branch)...)
					for k, _ := range branch.Value {
						result = append(result, k)
					}
					if len(latence) == 0 {
						break
					}
					key = latence[0].key
					node = latence[0].node
					latence = latence[1:]
					continue
				} else {
					latence = append(latence, ToBeAdd(key, *branch)...)
					b, remaining := key[0], key[1:]
					key = remaining
					node = branch.GetBranch(b)
					continue
				}
			}

			if ext, ok := node.(*ExtensionNode); ok {
				//fmt.Println("extension node")
				nodeExist[ext] = len(proof.Nodes)
				proof.Nodes = append(proof.Nodes, ext)
				matched := PrefixMatchedLen(ext.Path, key)
				if matched < len(ext.Path) && matched < len(key){
					if ext.Path[len(ext.Path)-1] < key[matched] {
						key = key[matched:]
						node = ext.Next
						continue
					} else {
						containAll, i := ContainJudge(ext.Path[matched:], key[matched:])
						if containAll{
							key = []byte{}
							node = ext.Next
							continue
						} else if ext.Path[len(ext.Path)-1] < key[i] {
							key = key[i:]
							node = ext.Next
							continue
						} else {
							if len(latence) == 0 {
								break
							}
							key = latence[0].key
							node = latence[0].node
							latence = latence[1:]
							continue
						}
					}
				} else {
					key = key[matched:]
					node = ext.Next
					continue
				}
			}
		}
	}
	// make up the relation between nodes
	for _, node := range proof.Nodes {
		switch node.(type) {
		case *LeafNode:
			continue
		case *BranchNode:
			branch, _ := (node).(*BranchNode)
			proof.NodeRelationMap[nodeExist[branch]] = make(map[int8]int)
			var i int8
			for i=0; i<int8(len(branch.Branches)); i++ {
				child := branch.Branches[i]
				if IsEmptyNode(child) {
					continue
				}
				if j, yes := nodeExist[child]; yes {
					proof.NodeRelationMap[nodeExist[branch]][i] = j
				} else {
					var hashNode HashNode
					if leaf, ok := child.(*LeafNode); ok {
						hashNode = NewHashNode(leaf.flags.hash)
					} else if bran, y := child.(*BranchNode); y {
						hashNode = NewHashNode(bran.flags.hash)
					} else {
						hashNode = NewHashNode(child.(*ExtensionNode).flags.hash)
					}
					index := len(proof.Nodes)
					proof.Nodes = append(proof.Nodes, hashNode)
					proof.NodeRelationMap[nodeExist[branch]][i] = index
				}
			}
		case *ExtensionNode:
			ext, _ := (node).(*ExtensionNode)
			proof.NodeRelationMap[nodeExist[ext]] = make(map[int8]int)
			if j, yes := nodeExist[ext.Next]; yes {
				proof.NodeRelationMap[nodeExist[ext]][0] = j
			} else {
				v := NewHashNode(ext.Next.(*BranchNode).flags.hash)
				index := len(proof.Nodes)
				proof.Nodes = append(proof.Nodes, v)
				proof.NodeRelationMap[nodeExist[ext]][0] = index
			}
		}
	}
	return result, proof, true
}

func (t *Trie) AuthenFilterPlus(q *matching.QueryGraph, g matching.Graph) map[string]Proof {
	nodeList := make(map[string]Proof)
	q.CandidateSets = make(map[int][]int)
	q.CandidateSetsB = make(map[int]map[int]bool)

	for str, ul := range q.NeiStr {
		fmt.Println("present key: ", str)
		pf := q.PathFeature[q.NeiStr[str][0]]
		C, P, _ := t.AuthSearchPlus([]byte(str), g, pf)
		nodeList[str] = P
		for _, u := range ul {
			q.CandidateSets[u] = C
			q.CandidateSetsB[u] = make(map[int]bool)
			for _, v := range C {
				q.CandidateSetsB[u][v] = true
			}
		}
	}
	return nodeList
}

func (t *Trie) AuthSearchPlus(key []byte, g matching.Graph, pf map[string]int) ([]int, Proof, bool) {
	/*
		obtaining candidate vertex set and merkle proof for the given key (one query vertex)
	*/

	var proof Proof
	nodeExist := make(map[Node]int)
	proof.NodeRelationMap = make(map[int]map[int8]int)
	if len(key) == 0 {
		return nil, proof, false
	}
	var result []int
	if root, ok := t.root.(*BranchNode); ok {
		//fmt.Println("branch node")
		node := root.GetBranch(key[0])
		key = key[1:]
		nodeExist[root] = len(proof.Nodes)
		proof.Nodes = append(proof.Nodes, root)
		var latence []potentialPath
		for {
			if IsEmptyNode(node) {
				if len(latence) == 0 {
					break
				}
				key = latence[0].key
				node = latence[0].node
				latence = latence[1:]
			}

			if leaf, ok := node.(*LeafNode); ok {
				//fmt.Println("leaf node")
				nodeExist[leaf] = len(proof.Nodes)
				proof.Nodes = append(proof.Nodes, leaf)
				matched := PrefixMatchedLen(leaf.Path, key)
				if matched == len(key) || IsContain(leaf.Path[matched:], key[matched:]){
					for k, _ := range leaf.Value {
						flag := true
						for pa, num := range pf {
							if _, yes := g.PathFeature[k][pa]; !yes {
								flag = false
								break
							} else if len(g.PathFeature[k][pa]) < num {
								flag = false
								break
							}
						}
						if flag {
							result = append(result, k)
						}
					}
				}
				if len(latence) == 0 {
					break
				}
				key = latence[0].key
				node = latence[0].node
				latence = latence[1:]
				continue
			}

			if branch, ok := node.(*BranchNode); ok {
				//fmt.Println("branch node")
				nodeExist[branch] = len(proof.Nodes)
				proof.Nodes = append(proof.Nodes, branch)
				if len(key) == 0 {
					latence = append(latence, ToBeAdd(key, *branch)...)
					for k, _ := range branch.Value {
						flag := true
						for pa, num := range pf {
							if _, yes := g.PathFeature[k][pa]; !yes {
								flag = false
								break
							} else if len(g.PathFeature[k][pa]) < num {
								flag = false
								break
							}
						}
						if flag {
							result = append(result, k)
						}
					}
					if len(latence) == 0 {
						break
					}
					key = latence[0].key
					node = latence[0].node
					latence = latence[1:]
					continue
				} else {
					latence = append(latence, ToBeAdd(key, *branch)...)
					b, remaining := key[0], key[1:]
					key = remaining
					node = branch.GetBranch(b)
					continue
				}
			}

			if ext, ok := node.(*ExtensionNode); ok {
				//fmt.Println("extension node")
				nodeExist[ext] = len(proof.Nodes)
				proof.Nodes = append(proof.Nodes, ext)
				matched := PrefixMatchedLen(ext.Path, key)
				if matched < len(ext.Path) && matched < len(key){
					if ext.Path[len(ext.Path)-1] < key[matched] {
						key = key[matched:]
						node = ext.Next
						continue
					} else {
						containAll, i := ContainJudge(ext.Path[matched:], key[matched:])
						if containAll{
							key = []byte{}
							node = ext.Next
							continue
						} else if ext.Path[len(ext.Path)-1] < key[i] {
							key = key[i:]
							node = ext.Next
							continue
						} else {
							if len(latence) == 0 {
								break
							}
							key = latence[0].key
							node = latence[0].node
							latence = latence[1:]
							continue
						}
					}
				} else {
					key = key[matched:]
					node = ext.Next
					continue
				}
			}
		}
	}
	// make up the relation between nodes
	for _, node := range proof.Nodes {
		switch node.(type) {
		case *LeafNode:
			continue
		case *BranchNode:
			branch, _ := (node).(*BranchNode)
			proof.NodeRelationMap[nodeExist[branch]] = make(map[int8]int)
			var i int8
			for i=0; i<int8(len(branch.Branches)); i++ {
				child := branch.Branches[i]
				if IsEmptyNode(child) {
					continue
				}
				if j, yes := nodeExist[child]; yes {
					proof.NodeRelationMap[nodeExist[branch]][i] = j
				} else {
					var hashNode HashNode
					if leaf, ok := child.(*LeafNode); ok {
						hashNode = NewHashNode(leaf.flags.hash)
					} else if bran, y := child.(*BranchNode); y {
						hashNode = NewHashNode(bran.flags.hash)
					} else {
						hashNode = NewHashNode(child.(*ExtensionNode).flags.hash)
					}
					index := len(proof.Nodes)
					proof.Nodes = append(proof.Nodes, hashNode)
					proof.NodeRelationMap[nodeExist[branch]][i] = index
				}
			}
		case *ExtensionNode:
			ext, _ := (node).(*ExtensionNode)
			proof.NodeRelationMap[nodeExist[ext]] = make(map[int8]int)
			if j, yes := nodeExist[ext.Next]; yes {
				proof.NodeRelationMap[nodeExist[ext]][0] = j
			} else {
				v := NewHashNode(ext.Next.(*BranchNode).flags.hash)
				index := len(proof.Nodes)
				proof.Nodes = append(proof.Nodes, v)
				proof.NodeRelationMap[nodeExist[ext]][0] = index
			}
		}
	}
	return result, proof, true
}