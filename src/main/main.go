package main

import (
	"Corgi/src/matching"
)

func main() {

	fileName := "./data/output.txt"
	t := new(matching.Graph)
	t.LoadGraphFromTxt(fileName)
	t.AssignLabel()
	t.ObtainNeiStr()
	t.Get()

}

