package mpt

import (
	"github.com/ethereum/go-ethereum/rlp"
)


type Node interface {
	Hash() []byte
	Raw() []interface{}
}

type nodeFlag struct {
	/*
		hash: hash of the node
		dirty: whether the node has changes
	*/
	hash []byte
	dirty bool
}

func Hash(node Node) []byte {
	return node.Hash()
}

func Serialize(node Node) []byte {
	var raw interface{}
	raw = node.Raw()
	rlp, err := rlp.EncodeToBytes(raw)
	if err != nil {
		panic(err)
	}
	return rlp
}



