package mpt

import (
"fmt"
)

func (t *MerklePartialTree) Insert(data interface{}) error {
	fmt.Println("Insert element: ", data)
	return nil
}
