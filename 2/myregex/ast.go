package myregex

import (
	"fmt"
)

type ast struct {
	root node
}

func (a ast) hasGroups() bool {
	return hasGroups(a.root)
}

func tokenize(pattern string) ([]string, error) {
	var tokens []string
	i := 0
	runes := []rune(pattern)
	n := len(runes)

	for i < n {
		ch := runes[i]

		switch {
		case ch == '\\':
			if i+1 >= n {
				return nil, fmt.Errorf("empty escape char at pos %d", i)
			}
			next := runes[i+1]
			if next >= '1' && next <= '9' {
				j := i + 1
				for j < n && runes[j] >= '0' && runes[j] <= '9' {
					j++
				}
				tokens = append(tokens, string(runes[i:j]))
				i = j
			} else {
				tokens = append(tokens, string(runes[i:i+2]))
				i += 2
			}

		case ch == '.' && i+2 < n && runes[i+1] == '.' && runes[i+2] == '.':
			tokens = append(tokens, "...")
			i += 3

		case ch == '[':
			j := i + 1
			for j < n && runes[j] != ']' {
				if runes[j] == '\\' {
					j++
				}
				j++
			}
			if j >= n {
				return nil, fmt.Errorf("unclosed '[' at pos %d", i)
			}
			tokens = append(tokens, string(runes[i:j+1]))
			i = j + 1

		case ch == '{':
			j := i + 1
			for j < n && runes[j] != '}' {
				j++
			}
			if j >= n {
				return nil, fmt.Errorf("unclosed '{' at pos %d", i)
			}
			tokens = append(tokens, string(runes[i:j+1]))
			i = j + 1

		case ch == '(' && i+1 < n && runes[i+1] == ':':
			tokens = append(tokens, "(?:")
			i += 2

		case ch == '(' || ch == ')' || ch == '|' || ch == '$':
			tokens = append(tokens, string(ch))
			i++

		case ch == ' ' || ch == '\t' || ch == '\n':
			i++

		default:
			tokens = append(tokens, string(ch))
			i++
		}
	}

	return tokens, nil
}
func node_from_tokens(tokens []string) node {

}
func buildAst(pattern string) (ast, error) {
	pattern = "(" + pattern + ")"
	tokens, err := tokenize(pattern)
	if err != nil {
		return ast{}, err
	}
	tree := ast{}
	tree.root = node_from_tokens(tokens)
	return ast{}, nil
}
