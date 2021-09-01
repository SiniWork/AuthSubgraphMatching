package mpt

type Trie struct {
	root Node
}

func NewTrie() *Trie {
	return &Trie{}
}

func (t *Trie) Insert(key []byte, value string) error {
	/*
		Inserting (key, value) into trie
		key: the key to be inserted
	*/

	node := &t.root
	for {
		if IsEmptyNode(*node) {
			leaf := NewLeafNode(key, value)
			*node = leaf
			return nil
		}

		if leaf, ok := (*node).(*LeafNode); ok {
			matched := PrefixMatchedLen(leaf.Path, key)

			if matched == len(key) && matched == len(leaf.Path) {
				leaf.Value = append(leaf.Value, value)
				return nil
			}

			branch := NewBranchNode()
			if matched == len(leaf.Path) {
				// matched cover the key of the present leaf node
				branch.SetValue(leaf.Value)
			}
			if matched == len(key) {
				// matched cover the key of the insert node
				branch.SetValue(value)
			}
			if matched > 0 {
				// the proportion of matched is a part of the key of the present leaf node and the insert node
				// create an extension node for the shared match
				ext := NewExtensionNode(leaf.Path[:matched], branch)
				*node = ext
			} else {
				// no matched, don't need extension node
				*node = branch
			}

			if matched < len(leaf.Path) {
				// present leaf have dismatched
				branchKey, leafKey := leaf.Path[matched], leaf.Path[matched+1:]
				newLeaf := NewLeafNode(leafKey, leaf.Value)
				branch.SetBranch(branchKey, newLeaf)
			}
			if matched < len(key) {
				// insert key have dismatched
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
			b, remaining := key[0], key[1:]
			key = remaining
			tmp := branch.GetBranch(b)
			node = &tmp
			continue
		}

		if ext, ok := (*node).(*ExtensionNode); ok {
			matched := PrefixMatchedLen(ext.Path, key)
			if matched < len(ext.Path) {
				extKey, branchKey, extRemainingKey := ext.Path[:matched], ext.Path[matched], ext.Path[matched+1:]
				nodeBranchKey, nodeLeafKey := key[matched], key[matched+1:]
				branch := NewBranchNode()
				if len(extRemainingKey) == 0 {
					branch.SetBranch(branchKey, ext.Next)
				} else {
					newExt := NewExtensionNode(extRemainingKey, ext.Next)
					branch.SetBranch(branchKey, newExt)
				}
				remainingLeaf := NewLeafNode(nodeLeafKey, value)
				branch.SetBranch(nodeBranchKey, remainingLeaf)

				if len(extKey) == 0 {
					// there is no shared extension key, then don't need the extension node
					*node = branch
				} else {
					// otherwise create a new extension node
					*node = NewExtensionNode(extKey, branch)
				}
				return nil
			}
			key = key[matched:]
			node = &ext.Next
			continue
		}
		panic("unknown type")
	}
}

func (t *Trie) GetExactOne(key []byte) ([]string, bool){
	/*
	Get the element depends on the given key
	 */

	node := t.root
	for {
		if IsEmptyNode(node) {
			return nil, false
		}

		if leaf, ok := node.(*LeafNode); ok {
			matched := PrefixMatchedLen(leaf.Path, key)
			if matched != len(leaf.Path) || matched != len(key) {
				return nil, false
			}
			return leaf.Value, true
		}

		if branch, ok := node.(*BranchNode); ok {
			if len(key) == 0 {
				return branch.Value, branch.HasValue()
			}
			b, remaining := key[0], key[1:]
			key = remaining
			node = branch.GetBranch(b)
			continue
		}

		if ext, ok := node.(*ExtensionNode); ok {
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

func (t *Trie) GetCandidate(key []byte) (map[string][]string, bool){
	/*
	get results that include given key
	 */

	node := t.root
	result := make(map[string][]string)
	for {
		if IsEmptyNode(node) {
			return result, false
		}

	}
}

func (t *Trie) Print() {

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


