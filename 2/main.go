package main

import (
	"fmt"
	"myregex"
)

func main() {
	//fmt.Println(myregex.Compile("௸௸௸a\\|"))
	fmt.Println(myregex.Compile("(asda({333})(:s)dva...|a...(a...a...)\\2)\\1\\2"))
}
