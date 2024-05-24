package main

import (
	"reflect"
	"testing"
)

func TestReadDir_Expected(t *testing.T) {
	path := "./testdata/env"
	expected := Environment{
		"BAR":   EnvValue{`bar`, false},
		"EMPTY": EnvValue{``, false},
		"FOO": EnvValue{`   foo
with new line`, false},
		"HELLO": EnvValue{`"hello"`, false},
		"UNSET": EnvValue{``, true},
	}

	env, err := ReadDir(path)
	if err != nil {
		t.Errorf("ReadDir returned error: %v", err)
	}

	if !reflect.DeepEqual(env, expected) {
		t.Errorf("\nOutput = %v; \nwant %v", env, expected)
	}
}

func TestReadDir_Nonexisting_Directory(t *testing.T) {
	_, err := ReadDir("wrong/dir")
	if err == nil {
		t.Errorf("Expected to get an error for non-existing directory")
	}
}
