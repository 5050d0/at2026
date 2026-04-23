package myregex

import "fmt"

type nodeData struct {
	nullable bool
	first    map[rune]bool
	last     map[rune]bool
}

func buildNfl(tree ast) (map[node]*nodeData, error) {
	var nfl map[node]*nodeData
	tree.root.fillNullable(&nfl)
	tree.root.fillFirst(&nfl)
	tree.root.fillLast(&nfl)
	return nfl, nil
}

func buildDfa(a ast) (Regex, error) {
	newRoot := nodeAnd{left: a.root, right: nodeEnd{}}
	a.root = &newRoot
	nfl, err := buildNfl(a)

	return nil, fmt.Errorf("DFA build error")
}
