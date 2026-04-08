package myregex

import (
	"slices"
	"testing"
)

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
			want:  []string{"(", "a", "b", ")"},
		},
		{
			name:  "empty parens",
			input: "()",
			want:  []string{"(", ")"},
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
			want:  []string{"(", "a", "b", ")", "..."},
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
			want:  []string{"(?:", "a", "b", ")"},
		},
		{
			name:  "non-capturing vs capturing",
			input: "(:a)|(b)",
			want:  []string{"(?:", "a", ")", "|", "(", "b", ")"},
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
			want:  []string{"(", "a", ")", `\1`},
		},
		{
			name:  "backreference \\12",
			input: `(a)\12`,
			want:  []string{"(", "a", ")", `\12`},
		},

		// Complex expressions
		{
			name:  "alternation with groups",
			input: "(ab)|(cd)",
			want:  []string{"(", "a", "b", ")", "|", "(", "c", "d", ")"},
		},
		{
			name:  "nested groups",
			input: "(a(bc))",
			want:  []string{"(", "a", "(", "b", "c", ")", ")"},
		},
		{
			name:  "class with kleene",
			input: "[ab]...",
			want:  []string{"[ab]", "..."},
		},
		{
			name:  "full expression",
			input: "(:a|b){2}[cd]...",
			want:  []string{"(?:", "a", "|", "b", ")", "{2}", "[cd]", "..."},
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
