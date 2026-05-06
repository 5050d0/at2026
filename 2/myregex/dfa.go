package myregex

import "fmt"

type nodeData struct {
	nullable bool
	first    map[int]struct{}
	last     map[int]struct{}
}

func buidIndex(tree ast) ([]node, error) {
	var index []node
	tree.root.buildIndex(&index)
	return index, nil
}
func buildNfl(tree ast) (map[node]*nodeData, error) {
	nfl := make(map[node]*nodeData)
	tree.root.fillNullable(&nfl)
	tree.root.fillFirst(&nfl)
	tree.root.fillLast(&nfl)
	return nfl, nil
}

func buildDfa(tree ast) (Regex, error) {
	newRoot := nodeAnd{left: tree.root, right: &nodeEnd{}}
	tree.root = &newRoot

	index, err := buidIndex(tree)
	if err != nil {
		return nil, err
	}
	fmt.Println(index)
	nfl, err := buildNfl(tree)
	if err != nil {
		return nil, err
	}
	fmt.Println(nfl)

	//var symPos []int
	//nfl, err := buildNfl(a)

	return nil, fmt.Errorf("DFA build error")
}
