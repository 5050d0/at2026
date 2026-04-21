package myregex

import (
	"fmt"
	"slices"
	"testing"
)

func TestConcatenize(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		// Basic literals
		{
			name:  "two chars",
			input: []string{"a", "b"},
			want:  []string{"a", "concat", "b"},
		},
		{
			name:  "three chars",
			input: []string{"a", "b", "c"},
			want:  []string{"a", "concat", "b", "concat", "c"},
		},
		{
			name:  "single char",
			input: []string{"a"},
			want:  []string{"a"},
		},
		{
			name:  "empty",
			input: []string{},
			want:  []string{},
		},

		// Alternation — no concat around |
		{
			name:  "alternation",
			input: []string{"a", "|", "b"},
			want:  []string{"a", "|", "b"},
		},
		{
			name:  "alternation with concat on sides",
			input: []string{"a", "b", "|", "c", "d"},
			want:  []string{"a", "concat", "b", "|", "c", "concat", "d"},
		},

		// Kleene closure
		{
			name:  "kleene then char",
			input: []string{"a", "...", "b"},
			want:  []string{"a", "...", "concat", "b"},
		},
		{
			name:  "char then kleene",
			input: []string{"a", "b", "..."},
			want:  []string{"a", "concat", "b", "..."},
		},
		{
			name:  "kleene then kleene",
			input: []string{"a", "...", "b", "..."},
			want:  []string{"a", "...", "concat", "b", "..."},
		},

		// Groups
		{
			name:  "concat after closing paren",
			input: []string{"(", "a", ")", "b"},
			want:  []string{"(", "a", ")", "concat", "b"},
		},
		{
			name:  "concat before opening paren",
			input: []string{"a", "(", "b", ")"},
			want:  []string{"a", "concat", "(", "b", ")"},
		},
		{
			name:  "no concat inside parens around pipe",
			input: []string{"(", "a", "|", "b", ")"},
			want:  []string{"(", "a", "|", "b", ")"},
		},
		{
			name:  "two groups concatenated",
			input: []string{"(", "a", ")", "(", "b", ")"},
			want:  []string{"(", "a", ")", "concat", "(", "b", ")"},
		},

		// Non-capturing group
		{
			name:  "non-capturing group with concat after",
			input: []string{"(:", "a", ")", "b"},
			want:  []string{"(:", "a", ")", "concat", "b"},
		},
		{
			name:  "char before non-capturing group",
			input: []string{"a", "(:", "b", ")"},
			want:  []string{"a", "concat", "(:", "b", ")"},
		},

		// Character class
		{
			name:  "class then char",
			input: []string{"[ab]", "c"},
			want:  []string{"[ab]", "concat", "c"},
		},
		{
			name:  "char then class",
			input: []string{"a", "[bc]"},
			want:  []string{"a", "concat", "[bc]"},
		},
		{
			name:  "class then class",
			input: []string{"[ab]", "[cd]"},
			want:  []string{"[ab]", "concat", "[cd]"},
		},
		{
			name:  "class then kleene",
			input: []string{"[ab]", "..."},
			want:  []string{"[ab]", "..."},
		},
		{
			name:  "class kleene then char",
			input: []string{"[ab]", "...", "c"},
			want:  []string{"[ab]", "...", "concat", "c"},
		},

		// Repeat {n}
		{
			name:  "repeat then char",
			input: []string{"a", "{3}", "b"},
			want:  []string{"a", "{3}", "concat", "b"},
		},
		{
			name:  "char then repeat — no concat (repeat binds left)",
			input: []string{"a", "{3}"},
			want:  []string{"a", "{3}"},
		},
		{
			name:  "repeat then group",
			input: []string{"a", "{3}", "(", "b", ")"},
			want:  []string{"a", "{3}", "concat", "(", "b", ")"},
		},

		// Backreferences
		{
			name:  "backreference after group",
			input: []string{"(", "a", ")", `\1`},
			want:  []string{"(", "a", ")", "concat", `\1`},
		},
		{
			name:  "backreference before char",
			input: []string{`\1`, "a"},
			want:  []string{`\1`, "concat", "a"},
		},

		// Escaped symbols
		{
			name:  "escaped char then char",
			input: []string{`\|`, "a"},
			want:  []string{`\|`, "concat", "a"},
		},
		{
			name:  "char then escaped char",
			input: []string{"a", `\|`},
			want:  []string{"a", "concat", `\|`},
		},

		// Dollar (empty string)
		{
			name:  "dollar then char",
			input: []string{"$", "a"},
			want:  []string{"$", "concat", "a"},
		},
		{
			name:  "char then dollar",
			input: []string{"a", "$"},
			want:  []string{"a", "concat", "$"},
		},

		// Complex
		{
			name:  "full expression",
			input: []string{"(:", "a", "|", "b", ")", "{2}", "[cd]", "..."},
			want:  []string{"(:", "a", "|", "b", ")", "{2}", "concat", "[cd]", "..."},
		},
		{
			name:  "nested groups",
			input: []string{"(", "a", "(", "b", "c", ")", ")"},
			want:  []string{"(", "a", "concat", "(", "b", "concat", "c", ")", ")"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := concatenize(tt.input)
			if !slices.Equal(got, tt.want) {
				t.Errorf("insertConcat(%v)\n got  %v\n want %v", tt.input, got, tt.want)
			}
		})
	}
}
func TestTokenize(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		// Basic literals
		{
			name:  "single char",
			input: "a",
			want:  []string{"a"},
		},
		{
			name:  "multiple chars",
			input: "abc",
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "empty string",
			input: "",
			want:  nil,
		},

		// Metacharacters
		{
			name:  "pipe",
			input: "a|b",
			want:  []string{"a", "|", "b"},
		},
		{
			name:  "dollar",
			input: "$",
			want:  []string{"$"},
		},
		{
			name:  "parens",
			input: "(ab)",
			want:  []string{"(0)", "a", "b", ")"},
		},
		{
			name:  "empty parens",
			input: "()",
			want:  []string{"(0)", ")"},
		},

		// Kleene closure '...'
		{
			name:  "kleene on char",
			input: "a...",
			want:  []string{"a", "..."},
		},
		{
			name:  "kleene on group",
			input: "(ab)...",
			want:  []string{"(0)", "a", "b", ")", "..."},
		},
		{
			name:  "single dot is literal",
			input: "a.b",
			want:  []string{"a", ".", "b"},
		},
		{
			name:  "two dots are literals",
			input: "a..b",
			want:  []string{"a", ".", ".", "b"},
		},

		// Repeat {n}
		{
			name:  "repeat single digit",
			input: "a{3}",
			want:  []string{"a", "{3}"},
		},
		{
			name:  "repeat multi digit",
			input: "a{12}",
			want:  []string{"a", "{12}"},
		},
		{
			name:  "repeat range",
			input: "a{2,5}",
			want:  []string{"a", "{2,5}"},
		},
		{
			name:    "unclosed brace",
			input:   "a{3",
			wantErr: true,
		},

		// Character class [...]
		{
			name:  "simple class",
			input: "[abc]",
			want:  []string{"[abc]"},
		},
		{
			name:  "class with escaped char",
			input: `[a\]b]`,
			want:  []string{`[a\]b]`},
		},
		{
			name:  "class followed by repeat",
			input: "[abc]{3}",
			want:  []string{"[abc]", "{3}"},
		},
		{
			name:    "unclosed bracket",
			input:   "[abc",
			wantErr: true,
		},

		// Non-capturing group (:
		{
			name:  "non-capturing group",
			input: "(:ab)",
			want:  []string{"(:", "a", "b", ")"},
		},
		{
			name:  "non-capturing vs capturing",
			input: "(:a)|(b)",
			want:  []string{"(:", "a", ")", "|", "(0)", "b", ")"},
		},

		// Escape sequences
		{
			name:  "escaped pipe",
			input: `a\|b`,
			want:  []string{"a", `\|`, "b"},
		},
		{
			name:  "escaped backslash",
			input: `a\\b`,
			want:  []string{"a", `\\`, "b"},
		},
		{
			name:  "escaped dot",
			input: `a\.b`,
			want:  []string{"a", `\.`, "b"},
		},
		{
			name:    "trailing backslash",
			input:   `a\`,
			wantErr: true,
		},

		// Backreferences
		{
			name:  "backreference \\1",
			input: `(a)\1`,
			want:  []string{"(0)", "a", ")", `\1`},
		},
		{
			name:  "backreference \\12",
			input: `(a)\12`,
			want:  []string{"(0)", "a", ")", `\12`},
		},

		// Complex expressions
		{
			name:  "alternation with groups",
			input: "(ab)|(cd)",
			want:  []string{"(0)", "a", "b", ")", "|", "(1)", "c", "d", ")"},
		},
		{
			name:  "nested groups",
			input: "(a(bc))",
			want:  []string{"(0)", "a", "(1)", "b", "c", ")", ")"},
		},
		{
			name:  "class with kleene",
			input: "[ab]...",
			want:  []string{"[ab]", "..."},
		},
		{
			name:  "full expression",
			input: "(:a|b){2}[cd]...",
			want:  []string{"(:", "a", "|", "b", ")", "{2}", "[cd]", "..."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tokenize(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("tokenize(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && !slices.Equal(got, tt.want) {
				t.Errorf("tokenize(%q)\n got  %v\n want %v", tt.input, got, tt.want)
			}
		})
	}
}

func nodeString(n node) string {
	if n == nil {
		return "nil"
	}
	switch v := n.(type) {
	case nodeLiteral:
		return "lit(" + string(v.value) + ")"
	case nodeOr:
		return "or(" + nodeString(v.left) + ", " + nodeString(v.right) + ")"
	case nodeAnd:
		return "and(" + nodeString(v.left) + ", " + nodeString(v.right) + ")"
	case nodeKleene:
		return "kleene(" + nodeString(v.child) + ")"
	case nodeSet:
		return "set(" + string(v.values) + ")"
	case nodeRepeat:
		return "repeat(" + nodeString(v.child) + ", " + itoa(v.number) + ")"
	case nodeGroup:
		return "group(" + nodeString(v.child) + ", " + itoa(v.index) + ")"
	case nodeGroupRef:
		return "ref(" + itoa(v.index) + ")"
	default:
		return "unknown"
	}
}

func itoa(n int) string { return fmt.Sprintf("%d", n) }

// assertNodeEqual сравнивает два узла с подробным сообщением об ошибке.
func assertNodeEqual(t *testing.T, got, want node, context string) {
	t.Helper()
	if !nodeEqual(got, want) {
		t.Errorf("%s:\n got  %s\n want %s", context, nodeString(got), nodeString(want))
	}
}

func nodeEqual(a, b node) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	switch va := a.(type) {
	case nodeLiteral:
		vb, ok := b.(nodeLiteral)
		return ok && va.value == vb.value
	case nodeOr:
		vb, ok := b.(nodeOr)
		return ok && nodeEqual(va.left, vb.left) && nodeEqual(va.right, vb.right)
	case nodeAnd:
		vb, ok := b.(nodeAnd)
		return ok && nodeEqual(va.left, vb.left) && nodeEqual(va.right, vb.right)
	case nodeKleene:
		vb, ok := b.(nodeKleene)
		return ok && nodeEqual(va.child, vb.child)
	case nodeSet:
		vb, ok := b.(nodeSet)
		return ok && string(va.values) == string(vb.values)
	case nodeRepeat:
		vb, ok := b.(nodeRepeat)
		return ok && va.number == vb.number && nodeEqual(va.child, vb.child)
	case nodeGroup:
		vb, ok := b.(nodeGroup)
		return ok && va.index == vb.index && nodeEqual(va.child, vb.child)
	case nodeGroupRef:
		vb, ok := b.(nodeGroupRef)
		return ok && va.index == vb.index
	default:
		return false
	}
}

// ========== Тесты построения AST ==========

func TestBuildAst(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		want    node
		wantErr bool
	}{
		{
			name:    "single literal",
			pattern: "a",
			want:    nodeLiteral{'a'},
		},
		{
			name:    "concatenation",
			pattern: "ab",
			want:    nodeAnd{nodeLiteral{'a'}, nodeLiteral{'b'}},
		},
		{
			name:    "alternation",
			pattern: "a|b",
			want:    nodeOr{nodeLiteral{'a'}, nodeLiteral{'b'}},
		},
		{
			name:    "kleene star",
			pattern: "a...",
			want:    nodeKleene{nodeLiteral{'a'}},
		},
		{
			name:    "character set",
			pattern: "[abc]",
			want:    nodeSet{[]rune("abc")},
		},
		{
			name:    "repeat",
			pattern: "a{3}",
			want:    nodeRepeat{child: nodeLiteral{'a'}, number: 3},
		},
		{
			name:    "capturing group",
			pattern: "(a)",
			want:    nodeGroup{child: nodeLiteral{'a'}, index: 0},
		},
		{
			name:    "non-capturing group",
			pattern: "(:a)",
			want:    nodeLiteral{'a'},
		},
		{
			name:    "backreference",
			pattern: `(a)\0`,
			want: nodeAnd{
				nodeGroup{child: nodeLiteral{'a'}, index: 0},
				nodeGroupRef{0},
			},
		},
		{
			name:    "complex expression with precedence",
			pattern: "ab|cd",
			want: nodeOr{
				nodeAnd{nodeLiteral{'a'}, nodeLiteral{'b'}},
				nodeAnd{nodeLiteral{'c'}, nodeLiteral{'d'}},
			},
		},
		{
			name:    "nested groups",
			pattern: "(a(b))",
			want: nodeGroup{
				child: nodeAnd{
					nodeLiteral{'a'},
					nodeGroup{child: nodeLiteral{'b'}, index: 1},
				},
				index: 0,
			},
		},
		{
			name:    "empty string anchor $",
			pattern: "a$",
			want:    nodeAnd{nodeLiteral{'a'}, nodeLiteral{'$'}},
		},
		{
			name:    "complex nesting with kleene",
			pattern: "(a|b)...",
			want: nodeKleene{
				nodeGroup{
					child: nodeOr{nodeLiteral{'a'}, nodeLiteral{'b'}},
					index: 0,
				},
			},
		},
		{
			name:    "repeat on group",
			pattern: "(ab){2}",
			want: nodeRepeat{
				child: nodeGroup{
					child: nodeAnd{nodeLiteral{'a'}, nodeLiteral{'b'}},
					index: 0,
				},
				number: 2,
			},
		},
		{
			name:    "multiple backreferences",
			pattern: `(a)(b)\2\1`,
			want: nodeAnd{
				nodeAnd{
					nodeAnd{
						nodeGroup{child: nodeLiteral{'a'}, index: 0},
						nodeGroup{child: nodeLiteral{'b'}, index: 1},
					},
					nodeGroupRef{2}},
				nodeGroupRef{1},
			},
		},

		{
			name:    "nested non-capturing groups",
			pattern: "(:a(:b))",
			want: nodeAnd{
				nodeLiteral{'a'},
				nodeLiteral{'b'},
			},
		},
		{
			name:    "character set with hyphen",
			pattern: "[a-z]",
			want:    nodeSet{[]rune("a-z")},
		},
		{
			name:    "escaped bracket inside set",
			pattern: `[a\]b]`,
			want:    nodeSet{[]rune(`a\]b`)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := buildAst(tt.pattern)
			if (err != nil) != tt.wantErr {
				t.Fatalf("buildAst(%q) error = %v, wantErr %v", tt.pattern, err, tt.wantErr)
			}
			if !tt.wantErr {
				assertNodeEqual(t, ast.root, tt.want, "AST mismatch")
			}
		})
	}
}
