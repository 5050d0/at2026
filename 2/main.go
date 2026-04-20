package main

import (
	"5050d0/myregex"
	"fmt"
)

func main() {
	fmt.Println(myregex.Compile("(a...a...(a...a...)\\2)\\1\\2"))
}
