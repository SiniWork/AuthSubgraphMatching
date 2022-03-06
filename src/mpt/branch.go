package mpt

import (
	"github.com/ethereum/go-ethereum/crypto"
	"sort"
	"strconv"
)

const BranchSize = 4

type BranchNode struct {
	Branches [BranchSize]Node
	Value map[int][]byte
	Content map[int]string
	flags    nodeFlag
}

func NewBranchNode() *BranchNode {
	return &BranchNode{
		Branches: [BranchSize]Node{},
		flags: newFlag(),
	}
}

func (b *BranchNode) SetBranch(bit byte, node Node) {
	b.Branches[int(bit)-65] = node
}

func (b *BranchNode) GetBranch(bit byte) Node {
	return b.Branches[int(bit)-65]
}

func (b *BranchNode) RemoveBranch(bit byte) {
	b.Branches[int(bit)-65] = nil
}

func (b *BranchNode) SetValue(value map[int][]byte, content map[int]string) {
	b.Value = value
	b.Content = content
}

func (b *BranchNode) RemoveValue() {
	b.Value = nil
}

func (b BranchNode) HasValue() bool {
	return b.Value != nil
}

func (b BranchNode) Hash() []byte {
	return crypto.Keccak256(b.Serialize())
}

func (b BranchNode) Raw() []interface{} {
	hashes := make([]interface{}, BranchSize+1)
	for i := 0; i < BranchSize; i++ {
		if b.Branches[i] == nil {
			hashes[i] = " "
		} else {
			node := b.Branches[i]
			hashes[i] = node.Hash()
		}
	}
	var valueStr []string
	var keys []int
	for k, _ := range b.Value {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, v := range keys {
		valueStr = append(valueStr, strconv.Itoa(v))
		valueStr = append(valueStr, string(b.Value[v]))
		valueStr = append(valueStr, b.Content[v])
	}
	hashes[BranchSize] = valueStr
	return hashes
}

func (b BranchNode) Serialize() []byte {
	return Serialize(b)
}

