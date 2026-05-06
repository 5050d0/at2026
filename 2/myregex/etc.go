package myregex

import (
	"fmt"
	"os"
	"strconv"
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
