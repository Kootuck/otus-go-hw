package main

import (
	"log"
	"os"
)

func main() {
	// 1. Pass args an proceed next
	cmdArgs := os.Args[2:]
	envPath := os.Args[1]
	rc := run(envPath, cmdArgs)
	if rc != 0 {
		log.Fatalf("error while executing command: %v", cmdArgs[1])
	}
}

func run(envPath string, args []string) (rc int) {
	// 2. Read directory to create env values map
	env, err := ReadDir(envPath)
	if err != nil {
		log.Fatal("error reading env dir:", err)
	}
	// 3. Run programm with 2.; call e.g.:
	// $ go-envdir /path/to/env/dir command arg1 arg2
	return RunCmdDefault(args, env)
}
