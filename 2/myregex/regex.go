package myregex

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
	Match(input string) (bool, error)
	//FindAllIter(input string) Iterator
}
type RegexDfa interface {
	Regex
	RebuildString() (string, error)
	Reverse() (Regex, error)
	Complement() (Regex, error)
}

type nfa struct{}

func Compile(pattern string) (Regex, error) {
	ast, err := buildAst(pattern)
	if err != nil {
		return nil, err
	}
	var regex Regex
	if ast.hasGroups() {
		regex, err = buildNfa(ast, pattern)
	} else {
		regex, err = buildDfa(ast)
	}
	if err != nil {
		return nil, err
	}
	//todo
	return regex, nil
}
