package matching

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"io/ioutil"
	"path"
	"sort"
	"strconv"
	"strings"
)

type Vertex struct {
	id int
	label byte
	content string
}

type Graph struct {
	/*
	vertices: vertex list
	adj: the adjacency list
	neiStr: statistic the one-hop neighborhood string for each vertex
	 */

	vertices []Vertex
	adj map[int][]int
	NeiStr map[string][]int
}

func (g *Graph) LoadGraphFromTxt(fileName string) error {
	/*
	loading the graph from txt file and saving it into an adjacency list adj
	 */

	g.adj = make(map[int][]int)
	content, err := readTxtFile(fileName)
	if err != nil {
		fmt.Println("Read file error!", err)
		return err
	}
	for _, line := range content {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		edge := strings.Split(line, " ")
		fr, err := strconv.Atoi(edge[0])
		if err != nil {
			return err
		}
		en, err := strconv.Atoi(edge[1])
		if err!= nil {
			return err
		}
		g.adj[fr] = append(g.adj[fr], en)
	}
	return nil
}

func (g *Graph) LoadGraphFromExcel(fileName string) error {
	/*
	loading the graph from Excel file and saving it into an adjacency list adj
	*/
	g.adj = make(map[int][]int)
	xlsx, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Println("Open file error!", err)
		return err
	}
	rows := xlsx.GetRows("Sheet1")
	for _, row := range rows {
		fr, err := strconv.Atoi(row[0])
		if err != nil {
			return err
		}
		en, err := strconv.Atoi(row[1])
		if err!= nil {
			return err
		}
		fmt.Println(fr, en)
		g.adj[fr] = append(g.adj[fr], en)
	}
	return nil
}

func (g *Graph) AssignLabel(fileName string) error {
	/*
	randomly assign a label to each vertex or read the vertex's label from a txt file
	then saving them into a list g.vertices
	 */
	var labelSet []byte
	if fileName == "" {
		labelSet = []byte{'A', 'B', 'C', 'D', 'E', 'F','G', 'H', 'I', 'J', 'K', 'L','M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
	} else {
		content, err := readTxtFile(fileName)
		if err != nil {
			fmt.Println("Read file error!", err)
			return err
		}
		for _, line := range content {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			labelSet = append(labelSet, line[0])
			}
		}
	for i := 0; i < len(g.adj); i++ {
		var v Vertex
		v.id = i
		v.label = labelSet[i]
		v.content = "lsy"
		g.vertices = append(g.vertices, v)
	}
	return nil
}

func (g *Graph) ObtainNeiStr() error {
	/*
	Generating the one-hop neighborhood string for each vertex and saving into a map g.neiStr
	g.neiStr, key is the one-hop neighborhood string, value is the list of vertex
	 */

	g.NeiStr = make(map[string][]int)
	for k, v := range g.adj {
		str := string(g.vertices[k].label)
		var nei []string
		for _, i := range v {
			nei = append(nei, string(g.vertices[i].label))
		}
		sort.Strings(nei)
		for _, t := range nei {
			str = str + t
		}
		g.NeiStr[str] = append(g.NeiStr[str], k)
	}
	return nil
}

func (g *Graph) Print() error {
	for k, v := range g.NeiStr {
		fmt.Println(k, v)
	}
	return nil
}

func readTxtFile(filePath string) ([]string, error) {
	fileSuffix :=  path.Ext(filePath)
	result := []string{}
	if fileSuffix == ".txt" {
		cont, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Println("Open file error!", err)
			return result, err
		}
		s := string(cont)
		result = strings.Split(s, "\n")
		return result, nil
	} else {
		return result, errors.New("file format error")
	}
}





