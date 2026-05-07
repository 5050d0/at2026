package myregex

import "fmt"

type NFA struct {
	tree ast
}

func buildNfa(a ast) (Regex, error) {
	_ = NFA{a}
	return nil, fmt.Errorf("NFA is not yet implemented")
}
