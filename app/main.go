package main

import (
	"os"
)

func main() {
	args := os.Args
	var dir string
	for i, arg := range args {
		if arg == "--directory" {
			dir = args[i+1]
		}
	}
	s := Server{}
	s.start(dir)
}
