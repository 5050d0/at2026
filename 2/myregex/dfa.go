package myregex

import "fmt"

type nodeData struct {
	nullable bool
	pos      int
	first    map[int]struct{}
	last     map[int]struct{}
}

func expandNodeRepeat(n node) node {
	switch v := n.(type) {
	case *nodeRepeat:
		child := expandNodeRepeat(v.child)
		if v.number == 0 {
			return &nodeEpsilon{}
		}
		result := child
		for i := 1; i < v.number; i++ {
			result = &nodeAnd{left: result, right: child}
		}
		return result
	case *nodeAnd:
		return &nodeAnd{expandNodeRepeat(v.left), expandNodeRepeat(v.right)}
	case *nodeOr:
		return &nodeOr{expandNodeRepeat(v.left), expandNodeRepeat(v.right)}
	case *nodeKleene:
		return &nodeKleene{expandNodeRepeat(v.child)}
	case *nodeGroup, *nodeGroupRef:
		panic("unsopported")
	default:
		return n
	}
}
func unionMaps(a, b map[int]struct{}) map[int]struct{} {
	result := make(map[int]struct{}, len(a)+len(b))
	for k := range a {
		result[k] = struct{}{}
	}
	for k := range b {
		result[k] = struct{}{}
	}
	return result
}

func copyMap(a map[int]struct{}) map[int]struct{} {
	result := make(map[int]struct{}, len(a))
	for k := range a {
		result[k] = struct{}{}
	}
	return result
}

func buidIndex(tree ast) ([]node, error) {
	var index []node
	tree.root.buildIndex(&index)
	return index, nil
}

func buildNfl(tree ast, index []node) (map[node]*nodeData, error) {
	nfl := make(map[node]*nodeData)
	for i, n := range index {
		if _, ok := nfl[n]; !ok {
			nfl[n] = &nodeData{}
		}
		nfl[n].pos = i
	}
	tree.root.fillNullable(&nfl)
	tree.root.fillFirst(&nfl)
	tree.root.fillLast(&nfl)
	return nfl, nil
}

func buildDfa(tree ast) (Regex, error) {
	newRoot := nodeAnd{left: tree.root, right: &nodeEnd{}}
	tree.root = &newRoot
	tree.root = expandNodeRepeat(tree.root)
	index, err := buidIndex(tree)
	if err != nil {
		return nil, err
	}
	nfl, err := buildNfl(tree, index)
	if err != nil {
		return nil, err
	}
	fmt.Println(nfl)
	return nil, fmt.Errorf("DFA build error")
}
