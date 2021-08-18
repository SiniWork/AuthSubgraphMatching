package main

import (
	"Corgi/src/mpt"
)

func main() {
	t := new(mpt.MerklePartialTree)
	var mkt mpt.Operator
	mkt = t
	mkt.Insert("data")
}

