package mpt

type Proof struct {
	Nodes []Node
	NodeRelationMap map[int]map[int8]int
}

type NodeCom struct {
	no Node
	id int
}

type pendingPath struct {
	key []byte
	node NodeCom
}

var flag = true // early stop when verifying

func (t *Trie) Prove(key []byte) ([]int, Proof, bool) {
	/*
	obtaining merkle proof of the given key
	*/
	var proof Proof
	t.HashRoot()
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
					result = append(result, leaf.Value...)
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
					result = append(result, branch.Value...)
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

func Verify(rootHash []byte, key []byte, proof Proof)  bool {
	/*
	Verifying whether the result satisfy correctness and completeness
	 */
	var nodeList []NodeCom
	for i, n := range proof.Nodes {
		nId := NodeCom{no: n, id: i}
		nodeList = append(nodeList, nId)
	}
	root := nodeList[0]
	newHash := reComputeHash(key, root, nodeList, proof.NodeRelationMap)
	if string(newHash) == string(rootHash) {
		return true
	} else {
		return false
	}
}

func reComputeHash(key []byte, nodeC NodeCom, nodeList []NodeCom, relation map[int]map[int8]int) []byte{
	/*
	Recomputing the root hash of the trie (verifying the correctness) and verifying the completeness
	*/

	var result []int
	if len(key) == 0 {
		return nil
	}
	if root, ok := nodeC.no.(*BranchNode); ok {
		rebuildRelation(&nodeC, nodeList, relation)
		node := nodeList[relation[nodeC.id][int8(key[0] - 'A')]]
		key = key[1:]
		var latence []pendingPath
		for {
			if IsEmptyNode(node.no) {
				if len(latence) == 0 {
					return hashRoot(root)
				}
				key = latence[0].key
				node = latence[0].node
				latence = latence[1:]
			}

			if leaf, ok := node.no.(*LeafNode); ok {
				matched := PrefixMatchedLen(leaf.Path, key)
				if matched == len(key) || IsContain(leaf.Path[matched:], key[matched:]){
					result = append(result, leaf.Value...)
				}
				if len(latence) == 0 {
					return hashRoot(root)
				}
				key = latence[0].key
				node = latence[0].node
				latence = latence[1:]
				continue
			}

			if branch, ok := node.no.(*BranchNode); ok {
				rebuildRelation(&node, nodeList, relation)
				// check the unsatisfied branches are indeed unsatisfied
				if len(key) == 0 {
					//_, pendPaths := checkBranch(key, node, nodeList, relation)
					right, pendPaths := checkBranch(key, node, nodeList, relation)
					if !right {
						return nil
					}
					latence = append(latence, pendPaths...)
					result = append(result, branch.Value...)
					if len(latence) == 0 {
						return hashRoot(root)
					}
					key = latence[0].key
					node = latence[0].node
					latence = latence[1:]
					continue
				} else {
					//_, pendPaths := checkBranch(key, node, nodeList, relation)
					right, pendPaths := checkBranch(key, node, nodeList, relation)
					if !right {
						return nil
					}
					latence = append(latence, pendPaths...)
					b, remaining := key[0], key[1:]
					key = remaining
					if branch.Branches[b-'A'] != nil {
						node = nodeList[relation[node.id][int8(b - 'A')]]
					} else {
						node = NodeCom{}
					}
					continue
				}
			}

			if ext, ok := node.no.(*ExtensionNode); ok {
				rebuildRelation(&node, nodeList, relation)
				matched := PrefixMatchedLen(ext.Path, key)
				if matched < len(ext.Path) && matched < len(key){
					if ext.Path[len(ext.Path)-1] < key[matched] {
						key = key[matched:]
						node = nodeList[relation[node.id][0]]
						continue
					} else {
						containAll, i := ContainJudge(ext.Path[matched:], key[matched:])
						if containAll{
							key = []byte{}
							node = nodeList[relation[node.id][0]]
							continue
						} else if ext.Path[len(ext.Path)-1] < key[i] {
							key = key[i:]
							node = nodeList[relation[node.id][0]]
							continue
						} else {
							if len(latence) == 0 {
								return hashRoot(root)
							}
							key = latence[0].key
							node = latence[0].node
							latence = latence[1:]
							continue
						}
					}
				} else {
					key = key[matched:]
					node = nodeList[relation[node.id][0]]
					continue
				}
			}
		}
	}
	return nil
}

func rebuildRelation(nodeC *NodeCom, nodeList []NodeCom, relation map[int]map[int8]int) {
	/*
	Restoring the children of the given node
	 */
	var i int8
	if branch, ok := nodeC.no.(*BranchNode); ok {
		for i=0; i < BranchSize; i++ {
			if branch.Branches[i] != nil {
				index := relation[nodeC.id][i]
				branch.Branches[i] = nodeList[index].no
			}
		}
	}
	if ext, ok := (nodeC.no).(*ExtensionNode); ok {
		index := relation[nodeC.id][0]
		ext.Next = nodeList[index].no
	}
}

func hashRoot(root Node) []byte{
	return root.Hash()
}

func checkBranch(key []byte, node NodeCom, nodeList []NodeCom, relation map[int]map[int8]int) (bool, []pendingPath) {
	/*
	Checking whether the branch that needs to be added is empty
	 */
	if branch, yes:= node.no.(*BranchNode); yes {
		var subBranches int8
		if len(key) == 0 {
			subBranches = BranchSize
		} else {
			subBranches = int8(len(branch.Branches[:key[0]-'A']))
		}
		var result []pendingPath
		var i int8
		for i=0; i < subBranches; i++ {
			if branch.Branches[i] != nil {
				if _, yes := nodeList[relation[node.id][i]].no.(HashNode); yes {
					return false, nil
				}
				p := pendingPath{key, nodeList[relation[node.id][i]]}
				result = append(result, p)
			}
		}
		return true, result
	}
	return false, nil
}

func (p *Proof) Size() int {
	/*
	Counting the size of the Proof
	*/
	var totalSize int
	for _, node := range p.Nodes {
		if leaf, ok := node.(*LeafNode); ok {
			leafSize := len(leaf.Path) + len(leaf.Value)*8
			totalSize = totalSize + leafSize
		} else if branch, ok := node.(*BranchNode); ok {
			branchSize := BranchSize*8 + len(branch.Value)*8
			totalSize = totalSize + branchSize
		} else if ext, ok := node.(*ExtensionNode); ok {
			extSize := len(ext.Path) + 8
			totalSize = totalSize + extSize
		} else if hs, ok := node.(HashNode); ok {
			hashSize := len(hs.hash)
			totalSize = totalSize + hashSize
		}
	}
	return totalSize
}