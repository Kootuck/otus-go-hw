package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunningCommandWithArgumentsAndEnv(t *testing.T) {
	t.Run("Test main", func(t *testing.T) {
		dir, err := os.Getwd()
		if err != nil {
			t.Fatal("was not able to get workdir")
		}
		envPath := filepath.Join(dir, "/testdata/env")
		args := []string{"/bin/bash", filepath.Join(dir, "/testdata/echo.sh"), "arg1=1", "arg2=2"}

		rc := run(envPath, args)

		wantRc := 0
		if rc != wantRc {
			t.Errorf("run() returned %v, want %v", rc, wantRc)
		}
	})
}
