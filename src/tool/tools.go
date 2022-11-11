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
var Labels = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

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

func RandomGenerateLabel(verL map[string]bool, labelSet []string, filePath string) {
	/*
	Randomly choosing a label from the set {'A', 'B', 'C', 'D'} to a vertex
	 */
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("open file error", err)
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	rand.Seed(time.Now().Unix())
	for k, _ := range verL {
		write.WriteString(k+" "+labelSet[rand.Intn(len(labelSet))]+"\r\n")
	}
	write.Flush()
}

func ConfigLabelForG(graphFile, labelFile string, labelSet []string) {
	/*
	Preparing label file for the given graph
	 */
	vertices, _ := CheckVerEdge(graphFile)
	RandomGenerateLabel(vertices,  labelSet, labelFile)
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

//fmt.Println("----------------Generating Label----------------")
//labels := tool.Labels[:20]
//tool.ConfigLabelForG("./data/amazon.txt", "./data/amazon_label.txt", labels)
//tool.ConfigLabelForG("./data/query/query"+strconv.Itoa(workload)+".txt", "./data/query/query"+strconv.Itoa(workload)+"_label.txt", labels)