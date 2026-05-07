package myregex

import "testing"

func mustBuildDfa(t *testing.T, pattern string) RegexDfa {
	t.Helper()
	dfa, err := Compile(pattern)
	if err != nil {
		t.Fatalf("buildDfa(%q) error: %v", pattern, err)
	}

	return dfa.(RegexDfa)
}

func mustMatch(t *testing.T, dfa RegexDfa, input string) bool {
	t.Helper()
	ok, err := dfa.Match(input)
	if err != nil {
		t.Fatalf("Match(%q) error: %v", input, err)
	}
	return ok
}

// ── Match ────────────────────────────────────────────────────────────────────

func TestDFA_Match(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		input   string
		want    bool
	}{
		// literal
		{"single char match", "a", "a", true},
		{"single char no match", "a", "b", false},
		{"literal string", "abc", "abc", true},
		{"literal prefix only", "abc", "ab", false},
		{"literal suffix only", "abc", "bc", false},
		{"literal extra suffix", "abc", "abcd", false},

		// epsilon / empty
		{"empty pattern empty input", "$", "", true},
		{"empty pattern nonempty input", "$", "a", false},

		// kleene
		{"kleene zero", "a...", "", true},
		{"kleene one", "a...", "a", true},
		{"kleene many", "a...", "aaaa", true},
		{"kleene wrong char", "a...", "b", false},
		{"kleene mixed", "a...", "aab", false},

		// concatenation + kleene
		{"concat kleene", "ab...", "a", true},
		{"concat kleene many b", "ab...", "abbb", true},
		{"concat kleene zero b", "ab...", "a", true},

		// or
		{"or left", "a|b", "a", true},
		{"or right", "a|b", "b", true},
		{"or neither", "a|b", "c", false},
		{"or word left", "cat|dog", "cat", true},
		{"or word right", "cat|dog", "dog", true},
		{"or word neither", "cat|dog", "cow", false},

		// character set
		{"set match first", "[abc]", "a", true},
		{"set match mid", "[abc]", "b", true},
		{"set match last", "[abc]", "c", true},
		{"set no match", "[abc]", "d", false},
		{"set kleene", "[ab]...", "aabba", true},
		{"set kleene empty", "[ab]...", "", true},

		// repeat
		{"repeat exact", "a{3}", "aaa", true},
		{"repeat too few", "a{3}", "aa", false},
		{"repeat too many", "a{3}", "aaaa", false},
		{"repeat zero", "a{0}", "", true},
		{"repeat one", "a{1}", "a", true},
		{"repeat word", "ab{2}", "abb", true},

		// combinations
		{"or then kleene", "(:a|b)...", "aabba", true},
		{"or then kleene empty", "(:a|b)...", "", true},
		{"or then kleene bad char", "(:a|b)...", "abc", false},
		{"set then literal", "[abc]d", "ad", true},
		{"set then literal wrong", "[abc]d", "dd", false},
		{"complex 1", "(:ab|cd)...ef", "ababcdef", true},
		{"complex 2", "(:ab|cd)...ef", "ef", true},
		{"complex 3", "(:ab|cd)...ef", "abef", true},
		{"complex no match", "(:ab|cd)...ef", "abcd", false},

		// escaping metachars
		{"escaped pipe", "a[|]b", "a|b", true},
		{"escaped dot literal", "a[.]b", "a.b", true},

		// unicode
		{"unicode match", "привет", "привет", true},
		{"unicode no match", "привет", "привет!", false},
		{"unicode set", "[абв]", "б", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dfa := mustBuildDfa(t, tc.pattern)
			got := mustMatch(t, dfa, tc.input)
			if got != tc.want {
				t.Errorf("pattern=%q input=%q: got %v, want %v",
					tc.pattern, tc.input, got, tc.want)
			}
		})
	}
}

// ── FindAll ──────────────────────────────────────────────────────────────────

func TestDFA_FindAll(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		input   string
		want    []string
	}{
		// no matches
		{"no match", "a", "bbb", nil},
		{"empty input", "a", "", nil},

		// single match
		{"single match start", "a", "abc", []string{"a"}},
		{"single match end", "a", "bca", []string{"a"}},
		{"single match mid", "a", "bac", []string{"a"}},
		{"word match", "cat", "mycat!", []string{"cat"}},

		// multiple matches
		{"multiple non-adjacent", "a", "abaca", []string{"a", "a", "a"}},
		{"multiple adjacent", "ab", "ababab", []string{"ab", "ab", "ab"}},
		{"overlapping impossible", "aa", "aaaa", []string{"aa", "aa"}},

		// kleene
		{"kleene greedy", "a...", "aaabaaaa", []string{"aaa", "aaaa"}},
		{"kleene single chars", "a...", "abab", []string{"a", "a"}},
		{"kleene zero width skipped", "a...", "bbb", nil},

		// or
		{"or finds both", "cat|dog", "I have a cat and a dog", []string{"cat", "dog"}},
		{"or repeated", "a|b", "abba", []string{"a", "b", "b", "a"}},

		// character set
		{"set finds all", "[aeiou]", "hello world", []string{"e", "o", "o"}},
		{"set kleene greedy", "[0123456789]...", "abc123def456", []string{"123", "456"}},

		// repeat
		{"repeat finds", "a{2}", "aabaabaaa", []string{"aa", "aa", "aa"}},
		{"repeat no partial", "a{3}", "aa", nil},

		// greedy: longest match wins
		{"greedy over short", "ab...", "xabbbbx", []string{"abbbb"}},
		{"greedy or longer", "a|ab", "ab", []string{"ab"}},

		// non-overlapping: advance past matched region
		{"non-overlapping", "aba", "abababa", []string{"aba", "aba"}},

		// unicode
		{"unicode finds", "ко[тш]", "тут кот и кош", []string{"кот", "кош"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dfa := mustBuildDfa(t, tc.pattern)
			results, err := dfa.FindAll(tc.input)
			if err != nil {
				t.Fatalf("FindAll(%q) error: %v", tc.input, err)
			}

			got := make([]string, len(results))
			for i, r := range results {
				got[i] = r.Match
			}

			if len(got) != len(tc.want) {
				t.Errorf("pattern=%q input=%q\n  got  %v\n  want %v",
					tc.pattern, tc.input, got, tc.want)
				return
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("pattern=%q input=%q match[%d]: got %q want %q",
						tc.pattern, tc.input, i, got[i], tc.want[i])
				}
			}
		})
	}
}
