package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"testing"
)

type OffsetLimit struct {
	offset        int64
	limit         int64
	checkFilePath string
}

var offsetLimit = []OffsetLimit{
	{
		offset:        0,
		limit:         0,
		checkFilePath: "./testdata/out_offset0_limit0.txt",
	},
	{
		offset:        0,
		limit:         10,
		checkFilePath: "./testdata/out_offset0_limit10.txt",
	},
	{
		offset:        0,
		limit:         1000,
		checkFilePath: "./testdata/out_offset0_limit1000.txt",
	},
	{
		offset:        0,
		limit:         10000,
		checkFilePath: "./testdata/out_offset0_limit10000.txt",
	},
	{
		offset:        100,
		limit:         1000,
		checkFilePath: "./testdata/out_offset100_limit1000.txt",
	},
	{
		offset:        6000,
		limit:         1000,
		checkFilePath: "./testdata/out_offset6000_limit1000.txt",
	},
}

func TestCopy(t *testing.T) {
	t.Run("Just test that files are successfully copied", TestOffsetLimitCombinations)
	t.Run("Offset larger than filesize is not allowed", TestOffsetErrorHandled)
}

func TestOffsetLimitCombinations(t *testing.T) {
	sourceFilepath := "./testdata/input.txt"
	targetFilepath := "./testdata/out.txt"

	for _, tCase := range offsetLimit {
		err := Copy(sourceFilepath, targetFilepath, tCase.offset, tCase.limit)
		if err != nil {
			t.Errorf("DoCopy(%q, %d) returned error: %v", tCase.offset, tCase.limit, err)
		}

		expected, err := os.Open(tCase.checkFilePath)
		if err != nil {
			t.Fatal(err)
		}
		defer expected.Close()

		copyResult, err := os.Open(targetFilepath)
		if err != nil {
			t.Fatal(err)
		}
		defer copyResult.Close()

		buf1 := bufio.NewReader(expected)
		buf2 := bufio.NewReader(copyResult)

		// read byte-by-byte and compare files
		for {
			b1, err1 := buf1.ReadByte()
			b2, err2 := buf2.ReadByte()

			if err1 == io.EOF && err2 == io.EOF {
				break
			}

			if err1 != nil {
				t.Logf("case: offset=%vlimit=%v", tCase.offset, tCase.limit)
				t.Error("expected file read error:", err1)
				break
			}

			if err2 != nil {
				t.Logf("case: offset=%vlimit=%v", tCase.offset, tCase.limit)
				t.Error("copy file read error:", err2)
				break
			}

			if b1 != b2 {
				t.Logf("case: offset=%vlimit=%v", tCase.offset, tCase.limit)
				t.Errorf("Files are different")
				break
			}
		}
	}
	if err := os.Remove(targetFilepath); err != nil {
		t.Errorf("couldn't delete file: %v", err)
	}
}

func TestOffsetErrorHandled(t *testing.T) {
	sourceFilepath := "./testdata/input.txt"
	targetFilepath := "./testdata/out.txt"

	sourceFile, err := os.Open(sourceFilepath)
	if err != nil {
		t.Fatal(err)
	}

	fileInfo, err := sourceFile.Stat()
	if err != nil {
		t.Fatal(err)
	}

	exceedingOffset := fileInfo.Size() + 1

	err = Copy(sourceFilepath, targetFilepath, exceedingOffset, 0)

	if err == nil {
		t.Error("Must return error")
		return
	}
	if !errors.Is(err, ErrOffsetExceedsFileSize) {
		t.Errorf("Expected ErrOffsetExceedsFileSize, but got %v", err)
	}

	if err := os.Remove(targetFilepath); err != nil {
		t.Errorf("couldn't delete file: %v", err)
	}
}
