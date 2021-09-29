package mpt

import "github.com/ethereum/go-ethereum/crypto"

const BranchSize = 26

type BranchNode struct {
	Branches [BranchSize]Node
	Value    []string
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

func (b *BranchNode) SetValue(value interface{}) {
	switch value.(type) {
	case string:
		b.Value = []string{value.(string)}
	case []string:
		b.Value = value.([]string)
	}
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
	hashes := make([]interface{}, BranchSize)
	for i := 0; i < BranchSize-1; i++ {
		if b.Branches[i] == nil {
			hashes[i] = " "
		} else {
			node := b.Branches[i]
			hashes[i] = node.Hash()
		}
	}
	hashes[BranchSize-1] = b.Value
	return hashes
}

func (b BranchNode) Serialize() []byte {
	return Serialize(b)
}

