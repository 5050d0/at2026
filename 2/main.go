package main

import (
	"fmt"
	"log"
	"myregex"
)

func main() {
	//fmt.Println(myregex.Compile("௸௸௸a\\|"))
	r, err := myregex.Compile("a...")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(r.Match(""))
}
