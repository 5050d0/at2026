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
			tokens = append(tokens, "(:")
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

type rng struct {
	from int
	to   int
}

func nodeFromTokens(tokens []string, r rng) (node, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty token")
	}
	if r.from >= r.to {
		return nil, nil
	}
	main_op, err := findMainOp(tokens, r)
	if err != nil {
		return nil, err
	}
	switch {
	case tokens[main_op] == "|":
		left, err := nodeFromTokens(tokens, rng{r.from, main_op})
		if err != nil {
			return nil, err
		}
		right, err := nodeFromTokens(tokens, rng{main_op + 1, r.to})
		if err != nil {
			return nil, err
		}
		return nodeOr{left, right}, nil
	case tokens[main_op] == "concat":
		left, err := nodeFromTokens(tokens, rng{r.from, main_op})
		if err != nil {
			return nil, err
		}
		right, err := nodeFromTokens(tokens, rng{main_op + 1, r.to})
		if err != nil {
			return nil, err
		}
		return nodeAnd{left, right}, nil
	case tokens[main_op] == "...":
		left, err := nodeFromTokens(tokens, rng{r.from, main_op})
		if err != nil {
			return nil, err
		}
		return nodeKleene{left}, nil
	case tokens[main_op][0] == '[':
		runes := []rune(tokens[main_op])
		return nodeSet{runes[1 : len(runes)-1]}, nil
	case tokens[main_op][0] == '{':
		var val int
		_, err := fmt.Sscanf(tokens[main_op], "{%d}", &val)
		if err != nil {
			return nil, fmt.Errorf("error parsing inside {} at pos %d", r.from)
		}
		left, err := nodeFromTokens(tokens, rng{r.from, main_op})
		if err != nil {
			return nil, err
		}
		return nodeRepeat{child: left, number: val}, nil
	case tokens[main_op][0] == '(' && tokens[main_op][1] == ':':
	//todo
	case tokens[main_op][0] == '(':
	//todo
	case tokens[main_op][0] == '\\':
		//todo
	default:
		return nil, fmt.Errorf("unknown operator '%s'", tokens[main_op])
	}

	return nil, nil
}

func findMainOp(tokens []string, r rng) (int, error) {
	//todo
	return 0, nil
}
func canBeLeft(tok string) bool {
	switch tok {
	case "(", "(:", "|":
		return false
	case ")", "$", "...":
		return true
	}
	if len(tok) > 1 && tok[0] == '{' {
		return true
	}
	if len(tok) > 1 && tok[0] == '[' {
		return true
	}
	runes := []rune(tok)
	if len(runes) >= 2 && runes[0] == '\\' && runes[1] >= '1' && runes[1] <= '9' {
		return true
	}
	if len(runes) == 2 && runes[0] == '\\' {
		return true
	}
	if len(runes) == 1 {
		return true
	}
	return false
}

func canBeRight(tok string) bool {
	switch tok {
	case ")", "|", "...", "{":
		return false
	case "(", "(:", "$":
		return true
	}
	if len(tok) > 1 && tok[0] == '{' {
		return false
	}
	if len(tok) > 1 && tok[0] == '[' {
		return true
	}
	runes := []rune(tok)
	if len(runes) >= 2 && runes[0] == '\\' && runes[1] >= '1' && runes[1] <= '9' {
		return true
	}
	if len(runes) == 2 && runes[0] == '\\' {
		return true
	}
	if len(runes) == 1 {
		return true
	}
	return false
}
func concatenize(tokens []string) []string {
	if len(tokens) == 0 {
		return tokens
	}

	result := make([]string, 0, len(tokens)*2)

	for i := 0; i < len(tokens)-1; i++ {
		result = append(result, tokens[i])
		if canBeLeft(tokens[i]) && canBeRight(tokens[i+1]) {
			result = append(result, "concat")
		}
	}
	result = append(result, tokens[len(tokens)-1])

	return result
}
func buildAst(pattern string) (ast, error) {
	pattern = "(" + pattern + ")"
	tokens, err := tokenize(pattern)
	if err != nil {
		return ast{}, err
	}
	tokens = concatenize(tokens)
	tree := ast{}
	tree.root, err = nodeFromTokens(tokens, rng{0, len(tokens)})
	return tree, err
}
