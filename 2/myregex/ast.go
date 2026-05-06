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
	captureGroupNumber := 0
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

		case ch == '(':
			tokens = append(tokens, fmt.Sprintf("(%d)", captureGroupNumber))
			captureGroupNumber++
			i++
		case ch == ')' || ch == '|' || ch == '$':
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
	if r.from > r.to {
		return nil, nil
	}
	mainOp, err := findMainOp(tokens, r)
	if err != nil {
		return nil, err
	}
	switch {
	case tokens[mainOp] == "|": // todo добавить ощибки при операндаъ
		left, err := nodeFromTokens(tokens, rng{r.from, mainOp - 1})
		if err != nil {
			return nil, err
		}
		right, err := nodeFromTokens(tokens, rng{mainOp + 1, r.to})
		if err != nil {
			return nil, err
		}
		return &nodeOr{left, right}, nil
	case tokens[mainOp] == "$":
		return &nodeEpsilon{}, nil
	case tokens[mainOp] == "concat":
		left, err := nodeFromTokens(tokens, rng{r.from, mainOp - 1})
		if err != nil {
			return nil, err
		}
		right, err := nodeFromTokens(tokens, rng{mainOp + 1, r.to})
		if err != nil {
			return nil, err
		}
		return &nodeAnd{left, right}, nil
	case tokens[mainOp] == "...":
		left, err := nodeFromTokens(tokens, rng{r.from, mainOp - 1})
		if err != nil {
			return nil, err
		}
		if left == nil {
			return nil, fmt.Errorf("kleene operator with no operand")
		}
		return &nodeKleene{left}, nil
	case tokens[mainOp][0] == '[':
		runes := []rune(tokens[mainOp])
		return &nodeSet{runes[1 : len(runes)-1]}, nil
	case tokens[mainOp][0] == '{':
		var val int
		_, err := fmt.Sscanf(tokens[mainOp], "{%d}", &val)
		if err != nil {
			return nil, fmt.Errorf("error parsing inside {} at pos %d", r.from)
		}
		left, err := nodeFromTokens(tokens, rng{r.from, mainOp - 1})
		if err != nil {
			return nil, err
		}
		return &nodeRepeat{child: left, number: val}, nil
	case tokens[mainOp] == "(:":
		closeIdx, err := findMatchingParen(tokens, mainOp, r.to)
		if err != nil {
			return nil, err
		}
		inner, err := nodeFromTokens(tokens, rng{mainOp + 1, closeIdx - 1})
		if err != nil {
			return nil, err
		}
		return inner, nil

	case len(tokens[mainOp]) >= 2 && tokens[mainOp][0] == '(' && tokens[mainOp][1] != ':':
		var capNum int
		_, err := fmt.Sscanf(tokens[mainOp], "(%d)", &capNum)
		if err != nil {
			return nil, fmt.Errorf("invalid capture group token %q", tokens[mainOp])
		}
		closeIdx, err := findMatchingParen(tokens, mainOp, r.to)
		if err != nil {
			return nil, err
		}
		child, err := nodeFromTokens(tokens, rng{mainOp + 1, closeIdx - 1})
		if err != nil {
			return nil, err
		}
		return &nodeGroup{child: child, index: capNum}, nil
	case tokens[mainOp][0] == '\\':
		if '0' <= tokens[mainOp][1] && tokens[mainOp][1] <= '9' {
			var val int
			_, err = fmt.Sscanf(tokens[mainOp], "\\%d", &val)
			if err != nil {
				return nil, fmt.Errorf("couldn't get capture group expr index at pos %d", r.from)
			}
			return &nodeGroupRef{val}, nil
		}
		if len(tokens[mainOp]) == 2 {
			return &nodeLiteral{[]rune(tokens[mainOp])[1]}, nil
		}
	case r.to == r.from:
		return &nodeLiteral{[]rune(tokens[mainOp])[0]}, nil
	default:
		return nil, fmt.Errorf("unknown operator '%s'", tokens[mainOp])
	}

	return nil, fmt.Errorf("AST build error")
}
func findMatchingParen(tokens []string, start int, end int) (int, error) {
	depth := 1
	i := start + 1
	for i <= end {
		tok := tokens[i]
		if len(tok) > 0 && tok[0] == '(' {
			depth++
		} else if tok == ")" {
			depth--
			if depth == 0 {
				return i, nil
			}
		}
		i++
	}
	return 0, fmt.Errorf("unmatched parenthesis at token %d", start)
}
func findMainOp(tokens []string, r rng) (int, error) {
	depth := 0
	bestIdx := -1
	bestPrio := 999

	for i := r.from; i <= r.to; i++ {
		tok := tokens[i]

		if len(tok) > 0 && tok[0] == '(' {
			depth++
		} else if tok == ")" {
			depth--
			if depth < 0 {
				return 0, fmt.Errorf("unmatched closing parenthesis at token %d", i)
			}
		} else {
			if depth == 0 {
				prio := 0
				switch {
				case tok == "|":
					prio = 1
				case tok == "concat":
					prio = 2
				case tok == "...":
					if i > r.from {
						prio = 3
					}
				case len(tok) > 0 && tok[0] == '{':
					if i > r.from {
						prio = 3
					}
				}
				if prio > 0 && prio <= bestPrio {
					bestPrio = prio
					bestIdx = i
				}
			}
		}
	}

	if depth != 0 {
		return 0, fmt.Errorf("unclosed parenthesis")
	}
	if bestIdx != -1 {
		return bestIdx, nil
	}
	return r.from, nil
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
	if len(tok) > 1 && tok[0] == '(' {
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
	pattern = "(:" + pattern + ")"
	tokens, err := tokenize(pattern)
	if err != nil {
		return ast{}, err
	}
	tokens = concatenize(tokens)
	tree := ast{}
	tree.root, err = nodeFromTokens(tokens, rng{0, len(tokens) - 1})
	//todo
	_ = tree.WriteDot("tree.dot")
	return tree, err
}
