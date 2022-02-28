package main

import (
	"os"
)

func main() {
	if len(os.Args) < 4 {
		return
	}
	myDir := os.Args[1]
	myMap, err := ReadDir(myDir)
	if err != nil {
		return
	}
	RunCmd(os.Args[3:], myMap)
}
