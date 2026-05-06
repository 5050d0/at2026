package main

import (
	"fmt"
	"log"
	"myregex"
)

func main() {
	//fmt.Println(myregex.Compile("௸௸௸a\\|"))
	r, err := myregex.Compile("a|$bc")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(r.Match("bc"))
}
