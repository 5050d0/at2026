package myregex

import (
	"fmt"
	"sort"
	"strings"
)

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
			result = &nodeAnd{left: result, right: child.copy()}
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

func buildFollowpos(n node, nfl map[node]*nodeData) map[int]map[int]struct{} {
	followpos := make(map[int]map[int]struct{})
	fillFollowpos(n, nfl, followpos)
	return followpos
}

func fillFollowpos(n node, nfl map[node]*nodeData, followpos map[int]map[int]struct{}) {
	if n == nil {
		return
	}
	switch v := n.(type) {
	case *nodeAnd:
		fillFollowpos(v.left, nfl, followpos)
		fillFollowpos(v.right, nfl, followpos)
		for i := range nfl[v.left].last {
			if followpos[i] == nil {
				followpos[i] = map[int]struct{}{}
			}
			for j := range nfl[v.right].first {
				followpos[i][j] = struct{}{}
			}
		}
	case *nodeKleene:
		fillFollowpos(v.child, nfl, followpos)
		for i := range nfl[n].last {
			if followpos[i] == nil {
				followpos[i] = map[int]struct{}{}
			}
			for j := range nfl[n].first {
				followpos[i][j] = struct{}{}
			}
		}
	case *nodeOr:
		fillFollowpos(v.left, nfl, followpos)
		fillFollowpos(v.right, nfl, followpos)
	}
}

type dfaState struct {
	positions   map[int]struct{}
	transitions map[rune]int
	isAccept    bool
}

type DFA struct {
	startState int
	states     []*dfaState
}

func hashSet(pos map[int]struct{}) string {
	keys := make([]int, 0, len(pos))
	for k := range pos {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	var sb strings.Builder
	for i, k := range keys {
		if i > 0 {
			sb.WriteRune(';')
		}
		sb.WriteString(fmt.Sprint(k))
	}
	t := sb.String()
	return t
}

func buildDfa(tree ast) (RegexDfa, error) {
	dfa := &DFA{}
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
	followpos := buildFollowpos(tree.root, nfl)
	endPos := len(index) - 1

	stateMap := make(map[string]int)

	startPositions := nfl[tree.root].first
	startHash := hashSet(startPositions)

	startState := &dfaState{
		positions:   startPositions,
		transitions: make(map[rune]int),
		isAccept:    false,
	}

	dfa.startState = 0
	dfa.states = append(dfa.states, startState)
	stateMap[startHash] = 0

	queue := make([]int, 1)
	for len(queue) > 0 {
		currID := queue[0]
		queue = queue[1:]
		currState := dfa.states[currID]

		if _, found := currState.positions[endPos]; found {
			currState.isAccept = true
		}

		symbolToPositions := make(map[rune][]int)

		for p := range currState.positions {
			n := index[p]
			switch v := n.(type) {
			case *nodeLiteral:
				symbolToPositions[v.value] = append(symbolToPositions[v.value], p)
			case *nodeSet:
				for _, char := range v.values {
					symbolToPositions[char] = append(symbolToPositions[char], p)
				}
			}
		}

		for char, positions := range symbolToPositions {
			nextPos := make(map[int]struct{})

			for _, p := range positions {
				for fp := range followpos[p] {
					nextPos[fp] = struct{}{}
				}
			}

			nextHash := hashSet(nextPos)
			nextID, exists := stateMap[nextHash]

			if !exists {
				nextID = len(dfa.states)
				newState := &dfaState{
					positions:   nextPos,
					transitions: make(map[rune]int),
				}
				dfa.states = append(dfa.states, newState)
				stateMap[nextHash] = nextID
				queue = append(queue, nextID)
			}

			currState.transitions[char] = nextID
		}
	}
	//dfa.WriteDot("dfa.dot")

	dfa.Minimize()
	return dfa, nil
}

func (d *DFA) FindAll(input string) ([]RegexResult, error) {
	runes := []rune(input)
	var results []RegexResult
	i := 0

	for i < len(runes) {
		currentStateID := d.startState
		lastAcceptEnd := -1
		j := i

		for j < len(runes) {
			nextStateID, ok := d.states[currentStateID].transitions[runes[j]]
			if !ok {
				break
			}
			currentStateID = nextStateID
			j++
			if d.states[currentStateID].isAccept {
				lastAcceptEnd = j
			}
		}

		if lastAcceptEnd == -1 {
			i++
		} else {
			results = append(results, RegexResult{
				Match:  string(runes[i:lastAcceptEnd]),
				Groups: []string{},
			})
			i = lastAcceptEnd
		}
	}

	return results, nil
}
func (d *DFA) Match(input string) (bool, error) {
	currentStateID := d.startState

	for _, char := range input {
		currentState := d.states[currentStateID]

		nextStateID, found := currentState.transitions[char]
		if !found {
			return false, nil
		}

		currentStateID = nextStateID
	}

	return d.states[currentStateID].isAccept, nil

}
func (d *DFA) RebuildString() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DFA) Reverse() (Regex, error) {
	pattern, err := d.RebuildString()
	if err != nil {
		return nil, err
	}
	ast, err := buildAst(pattern)
	if err != nil {
		return nil, err
	}
	ast.root = ast.root.reverse()
	return buildDfa(ast)
}

func (d *DFA) Complement() (Regex, error) {
	//TODO implement me
	panic("implement me")
}

func (dfa *DFA) Minimize() {
	//todo
}
