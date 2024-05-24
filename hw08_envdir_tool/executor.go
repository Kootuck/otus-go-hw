package main

import (
	"io"
	"os"
	"os/exec"
	"strings"
)

func RunCmdDefault(cmd []string, env Environment) int {
	return RunCmd(cmd, env, os.Stdout)
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment, w io.Writer) (returnCode int) {
	args := cmd[1:]
	command := cmd[0]
	comd := exec.Command(command, args...)

	// prepare env variable as slice of args
	comd.Env = append(removeFromOriginalEnv(env), addToEnv(env)...)

	comd.Stdin = os.Stdin
	comd.Stderr = os.Stderr
	comd.Stdout = w
	err := comd.Run()
	if err != nil {
		return 1
	}

	return 0
}

func addToEnv(env Environment) []string {
	arguments := make([]string, 0)
	for k, v := range env {
		if !v.NeedRemove {
			arguments = append(arguments, k+"="+v.Value)
		}
	}
	return arguments
}

func removeFromOriginalEnv(env Environment) []string {
	remove := make(map[string]bool)
	for k, v := range env {
		if v.NeedRemove {
			remove[k] = true
		}
	}

	var result []string

	for _, e := range os.Environ() {
		varName := strings.Split(e, "=")[0]
		if !remove[varName] {
			result = append(result, e)
		}
	}
	return result
}
