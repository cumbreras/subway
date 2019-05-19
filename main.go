package main

import (
	"flag"
	"github.com/cumbreras/subway/subway"
)

func main() {
	env := flag.String("env", "staging", "Environment to pull events from")
	flag.Parse()
	s := subway.New(*env)
	s.Start()
}
