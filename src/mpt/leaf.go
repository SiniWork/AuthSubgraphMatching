package mpt

import (
	"github.com/ethereum/go-ethereum/crypto"
)


type LeafNode struct {
	Path []byte
	Value []string
	flags nodeFlag
}

func NewLeafNode(key []byte, value interface{}) *LeafNode {
	switch value.(type) {
	case string:
		return &LeafNode{
			Path: key,
			Value: []string{value.(string)},
			flags: newFlag(),
		}
	case []string:
		return &LeafNode{
		Path: key,
		Value: value.([]string),
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
