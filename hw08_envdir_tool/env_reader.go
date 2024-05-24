package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
func ReadDir(dir string) (Environment, error) {
	// Extract files from /dir and parse into arguments
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	// Parse dir entries into env variables map
	env, err := mapDirEntries(dirEntries, dir)
	if err != nil {
		return nil, err
	}

	return env, nil
}

func mapDirEntries(entries []fs.DirEntry, dirPath string) (env Environment, err error) {
	env = make(Environment)

	for _, entry := range entries {
		// check if .env file
		ext := filepath.Ext(entry.Name())
		// remove extension -> variable name
		varName := strings.TrimSuffix(entry.Name(), ext)
		// empty file situation
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		if info.Size() == 0 {
			env[varName] = EnvValue{
				Value:      "",
				NeedRemove: true,
			}
			continue
		}

		// read file to extract env variable value
		filepath := fmt.Sprintf("%s/%s", dirPath, entry.Name())
		value, err := extractValueFromFile(filepath)
		if err != nil {
			return nil, err
		}

		env[varName] = EnvValue{
			Value:      value,
			NeedRemove: false,
		}
	}
	return env, nil
}

func extractValueFromFile(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}
	firstLine := scanner.Text()
	firstLine = strings.TrimRight(firstLine, " ")
	withNewline := bytes.ReplaceAll([]byte(firstLine), []byte("\x00"), []byte("\n"))

	return string(withNewline), nil
}
