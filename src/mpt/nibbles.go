package mpt

type Nibble byte

func PrefixMatchedLen(node1 []Nibble, node2 []Nibble) int {
	matched := 0
	for i := 0; i < len(node1) && i < len(node2); i++ {
		n1, n2 := node1[i], node2[i]
		if n1 == n2 {
			matched++
		} else {
			break
		}
	}

	return matched
}