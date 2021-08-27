package mpt

import "github.com/ethereum/go-ethereum/crypto"

type ExtensionNode struct {
		Path []byte
		Next Node
		flags nodeFlag
}

func NewExtensionNode(path []byte, next Node) *ExtensionNode {
	return &ExtensionNode{
		Path: path,
		Next: next,
	}
}

func (e ExtensionNode) Hash() []byte {
	return crypto.Keccak256(e.Serialize())
}

func (e ExtensionNode) Raw() []interface{} {
	hashes := make([]interface{}, 2)
	hashes[0] = e.Path
	if len(Serialize(e.Next)) >= 32 {
		hashes[1] = e.Next.Hash()
	} else {
		hashes[1] = e.Next.Raw()
	}
	return hashes
}

func (e ExtensionNode) Serialize() []byte {
	return Serialize(e)
}