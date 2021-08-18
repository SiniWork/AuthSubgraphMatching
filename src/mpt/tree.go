package mpt

type MerklePartialTree struct {

}

type Operator interface {
	Insert(data interface{}) error
}