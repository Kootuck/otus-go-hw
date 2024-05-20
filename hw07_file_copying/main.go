package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	if from == "" {
		fmt.Println("error: from is required")
		os.Exit(1)
	}
	if to == "" {
		fmt.Println("error: to is required")
		os.Exit(1)
	}
	if offset < 0 {
		fmt.Println("error: offset must be non-negative")
		os.Exit(1)
	}
	if limit < 0 {
		fmt.Println("error: limit must be non-negative")
		os.Exit(1)
	}
	fmt.Println("DoCopy args=", from, to, offset, limit)
	err := Copy(from, to, offset, limit)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
