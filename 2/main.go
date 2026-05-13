package main

import (
	"fmt"
	"log"
	"myregex"
)

func main() {
	//fmt.Println(myregex.Compile("௸௸௸a\\|"))
	r, err := myregex.Compile("me...(:f|ph)i")
	r.(*myregex.DFA).WriteDot("dfa.dot")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(r.FindAll("mephimefimeeeephi"))
}
