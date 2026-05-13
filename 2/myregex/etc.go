package myregex

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

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

func (a ast) WriteDot(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString("digraph AST {\n")
	f.WriteString("  node [shape=box, fontname=\"monospace\"];\n")
	f.WriteString("  rankdir=TB;\n\n")

	counter := 0
	var traverse func(n node) string
	traverse = func(n node) string {
		if n == nil {
			return ""
		}
		counter++
		id := "n" + strconv.Itoa(counter)

		var label string
		var children []node

		switch v := n.(type) {
		case *nodeLiteral:
			ch := v.value
			switch ch {
			case '\\':
				label = "\\\\"
			case '"':
				label = "\\\""
			case '\n':
				label = "\\n"
			case '\t':
				label = "\\t"
			default:
				label = string(ch)
			}
			label = fmt.Sprintf("'%s'", label)
		case *nodeOr:
			label = "|"
			children = []node{v.left, v.right}
		case *nodeAnd:
			label = "concat"
			children = []node{v.left, v.right}
		case *nodeKleene:
			label = "..."
			children = []node{v.child}
		case *nodeSet:
			label = fmt.Sprintf("[%s]", string(v.values))
		case *nodeRepeat:
			label = fmt.Sprintf("{%d}", v.number)
			children = []node{v.child}
		case *nodeGroup:
			label = fmt.Sprintf("group #%d", v.index)
			children = []node{v.child}
		case *nodeGroupRef:
			label = fmt.Sprintf("\\%d", v.index)
		default:
			label = "unknown"
		}

		fmt.Fprintf(f, "  %s [label=%q];\n", id, label)
		for _, child := range children {
			childID := traverse(child)
			if childID != "" {
				fmt.Fprintf(f, "  %s -> %s;\n", id, childID)
			}
		}
		return id
	}

	traverse(a.root)
	f.WriteString("}\n")
	return nil
}
func (d *DFA) WriteDot(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString("digraph DFA {\n")
	f.WriteString("  node [fontname=\"monospace\"];\n")
	f.WriteString("  rankdir=LR;\n\n")

	// Invisible start arrow pointing to start state
	f.WriteString("  __start [shape=point, width=0.2];\n")
	fmt.Fprintf(f, "  __start -> %d;\n\n", d.startState)

	// Emit each state node
	for t, state := range d.states {
		if state.isAccept {
			fmt.Fprintf(f, "  %d [shape=doublecircle, label=%q];\n", t, strconv.Itoa(t))
		} else {
			fmt.Fprintf(f, "  %d [shape=circle, label=%q];\n", t, strconv.Itoa(t))
		}
	}

	f.WriteString("\n")

	// Collect and merge transitions: (from, to) -> []rune, so parallel edges
	// get a single label like "a,b,c" instead of three separate arrows.
	type edgeKey struct{ from, to int }
	edgeLabels := make(map[edgeKey][]rune)

	for _, state := range d.states {
		// Sort runes for deterministic output
		chars := make([]rune, 0, len(state.transitions))
		for ch := range state.transitions {
			chars = append(chars, ch)
		}
		sort.Slice(chars, func(i, j int) bool { return chars[i] < chars[j] })

		for t, ch := range chars {
			to := state.transitions[ch]
			key := edgeKey{t, to}
			edgeLabels[key] = append(edgeLabels[key], ch)
		}
	}

	// Emit one edge per (from, to) pair with a combined label
	type edge struct {
		from, to int
		chars    []rune
	}
	edges := make([]edge, 0, len(edgeLabels))
	for k, chars := range edgeLabels {
		edges = append(edges, edge{k.from, k.to, chars})
	}
	sort.Slice(edges, func(i, j int) bool {
		if edges[i].from != edges[j].from {
			return edges[i].from < edges[j].from
		}
		return edges[i].to < edges[j].to
	})

	for _, e := range edges {
		var sb strings.Builder
		for i, ch := range e.chars {
			if i > 0 {
				sb.WriteRune(',')
			}
			switch ch {
			case '\\':
				sb.WriteString(`\\`)
			case '"':
				sb.WriteString(`\"`)
			case '\n':
				sb.WriteString(`\n`)
			case '\t':
				sb.WriteString(`\t`)
			default:
				sb.WriteRune(ch)
			}
		}
		fmt.Fprintf(f, "  %d -> %d [label=%q];\n", e.from, e.to, sb.String())
	}

	f.WriteString("}\n")
	return nil
}
