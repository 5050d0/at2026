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
	ast, err := buildAst(pattern)
	if err != nil {
		return nil, err
	}
	if ast.hasGroups() {

	} else {

	}
	//todo
	return nil, nil
}
