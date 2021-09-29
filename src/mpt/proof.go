package mpt

type Proof struct {
	Nodes []Node
	NodeRelationMap map[int]map[int]int
}

type NodeID struct {
	no Node
	id int
}

func (t *Trie) Prove(key []byte) (Proof, bool) {
	/*
	obtaining merkle proof of the given key
	*/
	var proof Proof
	nodeExist := make(map[Node]int)
	proof.NodeRelationMap = make(map[int]map[int]int)
	if len(key) == 0 {
		return proof, false
	}
	var result []string
	if root, ok := t.root.(*BranchNode); ok {
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
		case *HashNode:
			continue
		case *BranchNode:
			branch, _ := (node).(*BranchNode)
			for i:=0; i<len(branch.Branches); i++ {
				child := branch.Branches[i]
				if IsEmptyNode(child) {
					continue
				}
				if j, yes := nodeExist[child]; yes {
					proof.NodeRelationMap[nodeExist[branch]] = make(map[int]int)
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
			if j, yes := nodeExist[ext.Next]; yes {
				proof.NodeRelationMap[nodeExist[ext]] = make(map[int]int)
				proof.NodeRelationMap[nodeExist[ext]][0] = j
			} else {
				v := NewHashNode(ext.Next.(*BranchNode).flags.hash)
				index := len(proof.Nodes)
				proof.Nodes = append(proof.Nodes, v)
				proof.NodeRelationMap[nodeExist[ext]] = make(map[int]int)
				proof.NodeRelationMap[nodeExist[ext]][0] = index
			}
		}
	}
	return proof, true
}

func VerifyProof(rootHash []byte, key []byte, proof Proof)  bool {
	/*
	proving the given key indeed exist in trie
	 */
	var nodeList []NodeID
	for i, n := range proof.Nodes {
		nId := NodeID{no: n, id: i}
		nodeList = append(nodeList, nId)
	}
	root := nodeList[0]
	newHash := ComputeHash(key, root, nodeList, proof.NodeRelationMap)
	if string(newHash) == string(rootHash) {
		return true
	} else {
		return false
	}
}

func ComputeHash(key []byte, node NodeID, nodeList []NodeID, relation map[int]map[int]int) []byte {
	switch node.no.(type) {
	case *LeafNode:
		leaf, _ := (node.no).(*LeafNode)
		hashed := leaf.Hash()
		return hashed
	case *ExtensionNode:
		ext, _ := (node.no).(*ExtensionNode)
		index := relation[node.id][0]
		ext.Next = nodeList[index].no
		ComputeHash(key, nodeList[index], nodeList, relation)
		hashed := ext.Hash()
		return hashed
	case *BranchNode:
		branch, _ := (node.no).(*BranchNode)
		for i:=0; i < BranchSize; i++ {
			if branch.Branches[i] != nil {
				index := relation[node.id][i]
				branch.Branches[i] = nodeList[index].no
				ComputeHash(key, nodeList[index], nodeList, relation)
			}
		}
		hashed := branch.Hash()
		return hashed
	case *HashNode:
		hs, _ := (node.no).(*HashNode)
		hashed := hs.hash
		return hashed
	}
	return nil
}