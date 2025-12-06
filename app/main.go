package main

import (
	"os"
)

func main() {
	args := os.Args
	s := Server{Address: "0.0.0.0:4221"}
	s.start(getDirName(args, "--directory"))
}
