package myregex

import "errors"

type RegexResult struct {
	Match  string
	Groups []string
}

func (r RegexResult) Group(i int) string {
	return r.Groups[i]
}
func (r RegexResult) GroupsCount() int {
	return len(r.Groups)
}

type Regex interface {
	FindAll(input string) ([]RegexResult, error)
	//FindAllIter(input string) Iterator
}
type RegexDfa interface {
	Regex
	RebuildString() (string, error)
	Invert() (Regex, error)
	Complement() (Regex, error)
}

type nfa struct{}
type dfa struct{}

func Compile(pattern string) (Regex, error) {
	if pattern == "" {
		return nil, errors.New("pattern is empty")
	}
	ast, err := buildAst(pattern)
	if err != nil {
		return nil, err
	}
	ast.hasGroups()
	//todo
	return nil, nil
}
