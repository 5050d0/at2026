package main

import (
	"fmt"
	"myregex"
)

func main() {
	//fmt.Println(myregex.Compile("௸௸௸a\\|"))
	fmt.Println(myregex.Compile("[abc]me(:(:ph)|f)i"))
}
