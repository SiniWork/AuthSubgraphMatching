package mpt

import (
	"github.com/ethereum/go-ethereum/crypto"
)


type LeafNode struct {
	Path []byte
	Value []int
	flags nodeFlag
}

func NewLeafNode(key []byte, value interface{}) *LeafNode {
	switch value.(type) {
	case int:
		return &LeafNode{
			Path: key,
			Value: []int{value.(int)},
			flags: newFlag(),
		}
	case []int:
		return &LeafNode{
		Path: key,
		Value: value.([]int),
		flags: newFlag(),
		}
	}
	return nil
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
