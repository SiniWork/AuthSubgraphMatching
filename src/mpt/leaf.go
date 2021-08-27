package mpt

import "github.com/ethereum/go-ethereum/crypto"


type LeafNode struct {
	Path []byte
	Value []byte
	flags nodeFlag
}

func NewLeafNode(key, value string) *LeafNode {
	return &LeafNode{
		Path: []byte(key),
		Value: []byte(value),
	}
}

func (l LeafNode) Hash() []byte {
	return crypto.Keccak256(l.Serialize())
}

func (l LeafNode) Raw() []interface{} {
	path := l.Path
	raw := []interface{}{path, l.Value}
	return raw
}

func (l LeafNode) Serialize() []byte {
	return Serialize(l)
}
