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
		{"escaped pipe", "a\\|b", "a|b", true},
		{"escaped dot literal", "a\\...b", "a...b", true},

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

func mustReverse(t *testing.T, pattern string) Regex {
	t.Helper()
	dfa := mustBuildDfa(t, pattern)
	rev, err := dfa.Reverse()
	if err != nil {
		t.Fatalf("Reverse() on pattern %q error: %v", pattern, err)
	}
	return rev
}

func matchRev(t *testing.T, rev Regex, input string) bool {
	t.Helper()
	ok, err := rev.Match(input)
	if err != nil {
		t.Fatalf("Match(%q) on reversed regex error: %v", input, err)
	}
	return ok
}

// ── Match on reversed DFA ────────────────────────────────────────────────────

func TestDFA_Reverse_Match(t *testing.T) {
	tests := []struct {
		name    string
		pattern string // original pattern
		input   string // input to feed to the REVERSED regex
		want    bool
	}{
		// single character — reverse of one char is itself
		{"single char match", "a", "a", true},
		{"single char no match", "a", "b", false},

		// literal reversal
		{"literal reversed", "abc", "cba", true},
		{"literal not reversed", "abc", "abc", false},
		{"literal partial", "abc", "cb", false},
		{"literal extra", "abc", "dcba", false},

		// two char swap
		{"two chars reversed", "ab", "ba", true},
		{"two chars not swapped", "ab", "ab", false},

		// palindromes — reverse equals original
		{"palindrome 1", "aba", "aba", true},
		{"palindrome 2", "abba", "abba", true},
		{"palindrome single", "a", "a", true},

		// epsilon — reverse of empty is empty
		{"epsilon match", "$", "", true},
		{"epsilon no match", "$", "a", false},

		// kleene — reverse of a* is a*
		{"kleene empty", "a...", "", true},
		{"kleene one", "a...", "a", true},
		{"kleene many", "a...", "aaaa", true},
		{"kleene wrong char", "a...", "b", false},

		// kleene with concat — (ab)* reversed is (ba)*
		{"concat kleene empty", "(:ab)...", "", true},
		{"concat kleene one", "(:ab)...", "ba", true},
		{"concat kleene two", "(:ab)...", "baba", true},
		{"concat kleene not reversed", "(:ab)...", "ab", false},
		{"concat kleene partial", "(:ab)...", "bab", false},

		// or — reverse of (a|b) is (a|b), reverse of (ab|cd) is (ba|dc)
		{"or single chars", "a|b", "a", true},
		{"or single chars right", "a|b", "b", true},
		{"or single chars neither", "a|b", "c", false},
		{"or words reversed left", "cat|dog", "tac", true},
		{"or words reversed right", "cat|dog", "god", true},
		{"or words not reversed", "cat|dog", "cat", false},
		{"or words not reversed right", "cat|dog", "dog", false},

		// character set — [abc] reversed is still [abc] (single char)
		{"set match", "[abc]", "a", true},
		{"set match b", "[abc]", "b", true},
		{"set no match", "[abc]", "d", false},

		// repeat — reverse of a{3} is a{3} (same chars)
		{"repeat same", "a{3}", "aaa", true},
		{"repeat too few", "a{3}", "aa", false},
		{"repeat too many", "a{3}", "aaaa", false},

		// repeat of word — reverse of (ab){3} is (ba){3}
		{"repeat word reversed", "(:ab){3}", "bababa", true},
		{"repeat word not reversed", "(:ab){3}", "ababab", false},

		// longer concat reversal
		{"long literal", "abcde", "edcba", true},
		{"long literal wrong order", "abcde", "abcde", false},

		// mixed: concat + or
		{"concat or reversed", "(:a|b)c", "ca", true},
		{"concat or reversed b", "(:a|b)c", "cb", true},
		{"concat or not reversed", "(:a|b)c", "ac", false},

		// double reverse cancels out: rev(rev(L)) = L
		{"double reverse abc match", "abc", "abc", false}, // rev matches "cba", not "abc"
		{"double reverse cba match", "abc", "cba", true},

		// unicode
		{"unicode reversed", "привет", "тевирп", true},
		{"unicode not reversed", "привет", "привет", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rev := mustReverse(t, tc.pattern)
			got := matchRev(t, rev, tc.input)
			if got != tc.want {
				t.Errorf("Reverse(%q).Match(%q) = %v, want %v",
					tc.pattern, tc.input, got, tc.want)
			}
		})
	}
}

// ── Double reverse ────────────────────────────────────────────────────────────
// Reversing twice must produce a regex equivalent to the original.

func TestDFA_Reverse_DoubleReverse(t *testing.T) {
	patterns := []struct {
		pattern string
		accept  []string
		reject  []string
	}{
		{"abc", []string{"abc"}, []string{"cba", "ab", ""}},
		{"a...", []string{"", "a", "aaa"}, []string{"b", "ab"}},
		{"a|b", []string{"a", "b"}, []string{"c", "ab"}},
		{"(:ab)...", []string{"", "ab", "abab"}, []string{"ba", "a"}},
		{"[abc]", []string{"a", "b", "c"}, []string{"d", "ab"}},
		{"a{3}", []string{"aaa"}, []string{"aa", "aaaa"}},
	}

	for _, tc := range patterns {
		t.Run(tc.pattern, func(t *testing.T) {
			dfa := mustBuildDfa(t, tc.pattern)

			rev1, err := dfa.Reverse()
			if err != nil {
				t.Fatalf("first Reverse() error: %v", err)
			}
			rev1Dfa, ok := rev1.(*DFA)
			if !ok {
				t.Fatal("Reverse() did not return a *DFA")
			}
			rev2, err := rev1Dfa.Reverse()
			if err != nil {
				t.Fatalf("second Reverse() error: %v", err)
			}

			for _, s := range tc.accept {
				got, err := rev2.Match(s)
				if err != nil {
					t.Fatalf("double-reverse Match(%q) error: %v", s, err)
				}
				if !got {
					t.Errorf("double-reverse of %q should accept %q", tc.pattern, s)
				}
			}
			for _, s := range tc.reject {
				got, err := rev2.Match(s)
				if err != nil {
					t.Fatalf("double-reverse Match(%q) error: %v", s, err)
				}
				if got {
					t.Errorf("double-reverse of %q should reject %q", tc.pattern, s)
				}
			}
		})
	}
}

// ── FindAll on reversed DFA ───────────────────────────────────────────────────

func TestDFA_Reverse_FindAll(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		input   string
		want    []string
	}{
		// reverse of "ab" matches "ba"
		{"find ba in string", "ab", "xbayba", []string{"ba", "ba"}},
		{"find reversed word", "cat", "I saw tactac", []string{"tac", "tac"}},

		// reverse of a* is a* — finds runs of a's greedily
		{"kleene finds runs", "a...", "aaabaa", []string{"aaa", "aa"}},

		// reverse of (ab|cd) matches ba or dc
		{"or reversed findall", "ab|cd", "xbadcx", []string{"ba", "dc"}},

		// reverse of single char is itself
		{"single char findall", "x", "axbxcx", []string{"x", "x", "x"}},

		// no matches
		{"no match", "abc", "abc", nil}, // rev matches "cba", not "abc"

		// unicode
		{"unicode findall", "кот", "токкоттоккот", []string{"ток", "ток"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rev := mustReverse(t, tc.pattern)
			results, err := rev.FindAll(tc.input)
			if err != nil {
				t.Fatalf("FindAll(%q) error: %v", tc.input, err)
			}

			got := make([]string, len(results))
			for i, r := range results {
				got[i] = r.Match
			}

			if len(got) != len(tc.want) {
				t.Errorf("Reverse(%q).FindAll(%q)\n  got  %v\n  want %v",
					tc.pattern, tc.input, got, tc.want)
				return
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("match[%d]: got %q want %q", i, got[i], tc.want[i])
				}
			}
		})
	}
}

// ── Edge cases ────────────────────────────────────────────────────────────────

func TestDFA_Reverse_EdgeCases(t *testing.T) {
	t.Run("reverse returns non-nil", func(t *testing.T) {
		dfa := mustBuildDfa(t, "abc")
		rev, err := dfa.Reverse()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rev == nil {
			t.Fatal("Reverse() returned nil")
		}
	})

	t.Run("reverse of epsilon is epsilon", func(t *testing.T) {
		rev := mustReverse(t, "$")
		if !matchRev(t, rev, "") {
			t.Error("reverse of epsilon should match empty string")
		}
		if matchRev(t, rev, "a") {
			t.Error("reverse of epsilon should not match non-empty string")
		}
	})

	t.Run("reverse of single char is idempotent", func(t *testing.T) {
		rev := mustReverse(t, "a")
		if !matchRev(t, rev, "a") {
			t.Error("reverse of 'a' should still match 'a'")
		}
	})

	t.Run("reverse does not mutate original dfa", func(t *testing.T) {
		dfa := mustBuildDfa(t, "abc")
		_, err := dfa.Reverse()
		if err != nil {
			t.Fatalf("Reverse() error: %v", err)
		}
		// original must still work correctly after Reverse() is called
		ok, err := dfa.Match("abc")
		if err != nil {
			t.Fatalf("Match after Reverse() error: %v", err)
		}
		if !ok {
			t.Error("original DFA was mutated by Reverse()")
		}
		ok, err = dfa.Match("cba")
		if err != nil {
			t.Fatalf("Match after Reverse() error: %v", err)
		}
		if ok {
			t.Error("original DFA was mutated to accept reversed strings")
		}
	})
}
