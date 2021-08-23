package main

import (
	"Corgi/src/matching"
)

func main() {

	// Txt file
	//fileName := "./data/output.txt"
	//t := new(matching.Graph)
	//t.LoadGraphFromTxt(fileName)
	//t.AssignLabel("")
	//t.ObtainNeiStr()
	//t.Get()

	// Excel file
	fileName := "./data/test.xlsx"
	t := new(matching.Graph)
	t.LoadGraphFromExcel(fileName)
	t.AssignLabel("")
	t.ObtainNeiStr()
	t.Get()


	//tools.ExelToTxt("test.xlsx")



}

