package mpt

import (
	"errors"
	"fmt"
)

type potentialPath struct {
	key []byte
	node Node
}

type Trie struct {
	root Node
}
func (t *Trie) GetRoot() Node {
	return t.root
}
func NewTrie() *Trie {
	return &Trie{}
}

func (t *Trie) Insert(key []byte, value int) error {
	/*
	Inserting (key, value) into trie
	key: the key to be inserted
	value: the value to be inserted
	*/

	if len(key) == 0 {
		return errors.New("the key is empty")
	}
	node := &t.root
	var pre = node
	var recordB byte
	for {
		if IsEmptyNode(*node) {
			leaf := NewLeafNode(key, value)
			*node = leaf
			return nil
		}

		if leaf, ok := (*node).(*LeafNode); ok {
			matched := PrefixMatchedLen(leaf.Path, key)
			// first case: full matched
			if matched == len(key) && matched == len(leaf.Path) {
				leaf.Value = append(leaf.Value, value)
				return nil
			}
			// second case: no matched
			branch := NewBranchNode()
			if matched == 0 {
				if preBranch, yes := (*pre).(*BranchNode); yes {
					preBranch.SetBranch(recordB, branch)
				}
				*node = branch
				if len(key) == 0 {
					branch.SetValue(value)
					oldLeaf := NewLeafNode(leaf.Path[1:], leaf.Value)
					branch.SetBranch(leaf.Path[0],oldLeaf)
					return nil
				}
				if len(leaf.Path) == 0 {
					branch.SetValue(leaf.Value)
					newLeaf := NewLeafNode(key[1:], value)
					branch.SetBranch(key[0], newLeaf)
					return nil
				}
				oldLeaf := NewLeafNode(leaf.Path[1:], leaf.Value)
				branch.SetBranch(leaf.Path[0],oldLeaf)
				newLeaf := NewLeafNode(key[1:], value)
				branch.SetBranch(key[0], newLeaf)
				return nil
			}
			// third case: part matched
			ext := NewExtensionNode(leaf.Path[:matched], branch)
			*node = ext
			if preBranch, yes := (*pre).(*BranchNode); yes {
				preBranch.SetBranch(recordB, ext)
			}
			if matched == len(leaf.Path) {
				branch.SetValue(leaf.Value)
				branchKey, leafKey := key[matched], key[matched+1:]
				newLeaf := NewLeafNode(leafKey, value)
				branch.SetBranch(branchKey, newLeaf)
			} else if matched == len(key) {
				branch.SetValue(value)
				oldBranchKey, oldLeafKey := leaf.Path[matched], leaf.Path[matched+1:]
				oldLeaf := NewLeafNode(oldLeafKey, leaf.Value)
				branch.SetBranch(oldBranchKey, oldLeaf)
			} else {
				oldBranchKey, oldLeafKey := leaf.Path[matched], leaf.Path[matched+1:]
				oldLeaf := NewLeafNode(oldLeafKey, leaf.Value)
				branch.SetBranch(oldBranchKey, oldLeaf)
				branchKey, leafKey := key[matched], key[matched+1:]
				newLeaf := NewLeafNode(leafKey, value)
				branch.SetBranch(branchKey, newLeaf)
			}
			return nil
		}

		if branch, ok := (*node).(*BranchNode); ok {
			if len(key) == 0 {
				if branch.Value != nil{
					branch.Value = append(branch.Value, value)
				} else {
					branch.SetValue(value)
				}
				return nil
			}
			pre = node
			recordB = key[0]
			b, remaining := key[0], key[1:]
			key = remaining
			tmp := branch.GetBranch(b)
			if tmp == nil {
				leaf := NewLeafNode(key, value)
				branch.SetBranch(b, leaf)
				return nil
			} else {
				node = &tmp
				continue
			}
		}

		if ext, ok := (*node).(*ExtensionNode); ok {
			matched := PrefixMatchedLen(ext.Path, key)
			// first case: full matched
			if  matched == len(ext.Path) {
				key = key[matched:]
				node = &ext.Next
				continue
			}
			// second case: no matched
			branch := NewBranchNode()
			if matched == 0 {
				if preBranch, ok := (*pre).(*BranchNode); ok {
					preBranch.SetBranch(recordB, branch)
				}
				extBranchKey, extRemainingKey := ext.Path[0], ext.Path[1:]
				if len(extRemainingKey) == 0 {
					branch.SetBranch(extBranchKey, ext.Next)
				} else {
					newExt := NewExtensionNode(extRemainingKey, ext.Next)
					branch.SetBranch(extBranchKey, newExt)
				}
				leafBranchKey, leafRemainingKey := key[0], key[1:]
				newLeaf := NewLeafNode(leafRemainingKey, value)
				branch.SetBranch(leafBranchKey, newLeaf)
				*node = branch
				return nil
			}
			// third case: part matched
			commonKey, branchKey, extRemainingKey := ext.Path[:matched], ext.Path[matched], ext.Path[matched+1:]
			oldExt := NewExtensionNode(commonKey, branch)
			if preBranch, ok := (*pre).(*BranchNode); ok {
				preBranch.SetBranch(recordB, oldExt)
			}
			if len(extRemainingKey) == 0 {
				branch.SetBranch(branchKey, ext.Next)
			} else {
				newExt := NewExtensionNode(extRemainingKey, ext.Next)
				branch.SetBranch(branchKey, newExt)
			}
			if len(commonKey) == len(key) {
				branch.SetValue(value)
			} else {
				leafBranchKey, leafRemainingKey := key[matched], key[matched+1:]
				newLeaf := NewLeafNode(leafRemainingKey, value)
				branch.SetBranch(leafBranchKey, newLeaf)
			}
			*node = oldExt
			return nil
		}
		panic("unknown type")
	}
}

func (t *Trie) GetExactOne(key []byte) ([]int, bool){
	/*
	Get the element depends on the given key
	 */

	node := t.root
	for {
		if IsEmptyNode(node) {
			return nil, false
		}

		if leaf, ok := node.(*LeafNode); ok {
			fmt.Println("leaf node") // for test
			matched := PrefixMatchedLen(leaf.Path, key)
			if matched != len(leaf.Path) || matched != len(key) {
				return nil, false
			}
			return leaf.Value, true
		}

		if branch, ok := node.(*BranchNode); ok {
			fmt.Println("branch node") // for test
			if len(key) == 0 {
				return branch.Value, branch.HasValue()
			}
			b, remaining := key[0], key[1:]
			key = remaining
			node = branch.GetBranch(b)
			continue
		}

		if ext, ok := node.(*ExtensionNode); ok {
			fmt.Println("extension node") // for test
			matched := PrefixMatchedLen(ext.Path, key)
			if matched < len(ext.Path) {
				return nil, false
			}
			key = key[matched:]
			node = ext.Next
			continue
		}
		panic("not found")
	}
}

func (t *Trie) GetCandidate(key []byte) []int{
	/*
	get results that include given key
	 */

	var result []int
	if len(key) == 0 {
		return result
	}
	if root, ok := t.root.(*BranchNode); ok {
		node := root.GetBranch(key[0])
		key = key[1:]
		var latence []potentialPath
		for {
			if IsEmptyNode(node) {
				if len(latence) == 0 {
					return result
				}
				key = latence[0].key
				node = latence[0].node
				latence = latence[1:]
			}

			if leaf, ok := node.(*LeafNode); ok {
				fmt.Println("leaf node")
				matched := PrefixMatchedLen(leaf.Path, key)
				if matched == len(key) || IsContain(leaf.Path[matched:], key[matched:]){
					result = append(result, leaf.Value...)
				}
				if len(latence) == 0 {
					return result
				}
				key = latence[0].key
				node = latence[0].node
				latence = latence[1:]
				continue
			}

			if branch, ok := node.(*BranchNode); ok {
				fmt.Println("branch node")
				if len(key) == 0 {
					latence = append(latence, ToBeAdd(key, *branch)...)
					result = append(result, branch.Value...)
					if len(latence) == 0 {
						return result
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
				fmt.Println("extension node")
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
								return result
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
	return result
}

func (t *Trie)PrintTrie() {
	if t.root == nil {
		return
	}
	printNode(t.root)
	return
}

func (t *Trie) HashRoot() ([]byte,error) {
	/*
	computing the root hash
	*/
	if t.root == nil {
		return []byte{}, errors.New("the trie is empty")
	}
	hashed := hash(&t.root)
	return hashed, nil
}

func hash(node *Node) []byte {
	/*
	computing root hash of the subtree corresponding to the given node
	 */
	switch (*node).(type) {
	case *LeafNode:
		leaf, _ := (*node).(*LeafNode)
		hashed := leaf.Hash()
		leaf.flags.hash = hashed
		return hashed
	case *ExtensionNode:
		ext, _ := (*node).(*ExtensionNode)
		hash(&ext.Next)
		ext.flags.hash = ext.Hash()
		return ext.flags.hash
	case *BranchNode:
		branch, _ := (*node).(*BranchNode)
		for i:=0; i < BranchSize; i++ {
			if child := branch.Branches[i]; child != nil {
				hash(&child)
			}
		}
		branch.flags.hash = branch.Hash()
		return branch.flags.hash
	}
	return nil
}

func printNode(node Node) {
	switch (node).(type) {
	case *LeafNode:
		leaf, _ := (node).(*LeafNode)
		fmt.Println("LeafNode hash: ", leaf.flags.hash)
		fmt.Println(leaf.Value)
		return
	case *ExtensionNode:
		ext, _ := (node).(*ExtensionNode)
		fmt.Println("ExtensionNode hash: ", ext.flags.hash)
		printNode(ext.Next)
		return
	case *BranchNode:
		branch, _ := (node).(*BranchNode)
		fmt.Println("BranchNode hash: ", branch.flags.hash)
		for i:=0; i<BranchSize; i++ {
			if child := branch.Branches[i]; child != nil {
				printNode(child)
			}
		}
		return
	}
	return
}

func PrefixMatchedLen(node1, node2 []byte) int {
	matched := 0
	for i := 0; i < len(node1) && i < len(node2); i++ {
		n1, n2 := node1[i], node2[i]
		if n1 == n2 {
			matched++
		} else {
			break
		}
	}
	return matched
}

func IsContain(node1, node2 []byte) bool {
	/*
	Judge the key of node1 whether contain the key of the node2
	*/

	for _, v := range node1 {
		if len(node2) == 0 {
			return true
		} else {
			if v > node2[0] {
				return false
			} else if v == node2[0] {
				node2 = node2[1:]
				continue
			}
		}
	}
	if len(node2) != 0 {
		return false
	} else {
		return true
	}
}

func ContainJudge(node1, node2 []byte) (bool, int) {
	/*
	Judge the key of node1 whether contain the key of the node2, if not, return the position in node2
	 */

	i := 0
	for _, v := range node1 {
		if len(node2) == 0 {
			return true, -1
		} else {
			if v > node2[0] {
				return false, i
			} else if v == node2[0] {
				i = i + 1
				node2 = node2[1:]
				continue
			}
		}
	}
	if len(node2) != 0 {
		return false, i
	} else {
		return true, -1
	}
}

func ToBeAdd(key []byte, node BranchNode) []potentialPath {
	var subBranches []Node
	if len(key) == 0 {
		subBranches = node.Branches[:len(node.Branches)]
	} else {
		subBranches = node.Branches[:key[0]-'A']
	}

	var result []potentialPath
	for _, v := range subBranches{
		if IsEmptyNode(v) {
			continue
		}
		p := potentialPath{key, v}
		result = append(result, p)
	}
	return result
}

func Test(node LeafNode) bool {
	newLeaf := NewLeafNode(node.Path, node.Value)
	newLeaf1 := NewLeafNode(node.Path, node.Value)
	if newLeaf == newLeaf1 {
		return true
	}
	fmt.Println(newLeaf)
	fmt.Println(newLeaf1)
	return false
}