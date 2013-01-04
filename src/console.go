package main

import (
	"fmt"
	"github.com/sbinet/liner"
	"os"
)

import (
	"rc"
)

func main() {
	file, err := os.Open("test")
	if err == nil {
		rc.Read(file)
		file.Close()
	} else {
		fmt.Println(err)
	}

	state := liner.NewLiner()
	defer state.Close()

	done := false
	for !done {
		line, err := state.Prompt("> ")
		if err == nil {
			fmt.Println(line)
			state.AppendHistory(line)
		} else {
			done = true
		}
	}
}
