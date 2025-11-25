package main

import (
	"os"
)

func main() {
	args := os.Args
	s := Server{}
	s.start(getDirName(args, "--directory"))
}
