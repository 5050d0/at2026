package main

import (
	"fmt"
	"log"
	"myregex"
)

func main() {
	//fmt.Println(myregex.Compile("௸௸௸a\\|"))
	r, err := myregex.Compile("(:me...(:f|ph)[ei])")
	r.(*myregex.DFA).WriteDot("dfa.dot")

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(r.(*myregex.DFA).RebuildString())
	//fmt.Println(r.FindAll("mephimefimeeeephi"))
}
