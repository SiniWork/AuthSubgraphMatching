package tool

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"
)

func CheckVerEdge(fileName string) (map[string]bool, int){
	/*
	Return the vertex list and the number of edge
	*/
	fileSuffix :=  path.Ext(fileName)
	vertices := make(map[string]bool)
	edgesN := 0
	if fileSuffix == ".txt" {
		cont, err := ioutil.ReadFile(fileName)
		if err != nil {
			fmt.Println("Open file error!", err)
			return nil, 0
		}
		content := strings.Split(string(cont), "\n")
		splitStr := " "
		if find := strings.Contains(content[0], ","); find {
			splitStr = ","
		} else if find := strings.Contains(content[0], "	"); find {
			splitStr = "	"
		}
		for _, line := range content {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			edge := strings.Split(line, splitStr)
			vertices[edge[0]] = true
			vertices[edge[1]] = true
			edgesN++
		}
	}
	fmt.Println("the number of vertices: ", len(vertices))
	return vertices, edgesN
}

func RandomGenerateLabel(verL map[string]bool, filePath string) {
	/*
	Randomly choosing a label from the set {'A', 'B', 'C', 'D'} to a vertex
	 */
	labelSet := map[int]string {0:"A", 1:"B", 2:"C", 3:"D"}
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("open file error", err)
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	rand.Seed(time.Now().Unix())
	for k, _ := range verL {
		write.WriteString(k+" "+labelSet[rand.Intn(4)]+"\r\n")
	}
	write.Flush()
}

func ConfigLabelForG(graphFile, labelFile string) {
	/*
	Preparing label file for the given graph
	 */
	vertices, _ := CheckVerEdge(graphFile)
	RandomGenerateLabel(vertices, labelFile)
}

func CheckGraphLabel(graphFile, labelFile string) bool {
	numLabel, _ := CheckVerEdge(graphFile)
	file, err := os.OpenFile(labelFile, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return false
	}
	defer file.Close()
	buf := bufio.NewReader(file)
	i := 0
	for {
		_, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return false
			}
		}
		i++
	}
	if i == len(numLabel) {
		return true
	}
	fmt.Println("the number of edges: ", i)
	return false
}