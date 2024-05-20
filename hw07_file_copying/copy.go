package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// 1. Reader.
	fromReader, err := os.Open(fromPath)
	if err != nil {
		return err
	}

	// 2. Writer.
	toWriter, err := os.Create(toPath)
	if err != nil {
		return err
	}

	defer func() {
		fromReader.Close()
		toWriter.Close()
	}()

	// 3. Offset via limit reader
	fileInfo, err := fromReader.Stat()
	if err != nil {
		return err
	}

	// 3.1. Check that offset is relevant for filesize
	fileSize := fileInfo.Size()
	if fileSize <= offset {
		return ErrOffsetExceedsFileSize
	}

	// 3.2. Set offset for copy
	_, err = fromReader.Seek(offset, 0)
	if err != nil {
		return err
	}

	// 4. Limit reader.
	if limit == 0 {
		limit = fileSize
	}

	lr := io.LimitReader(fromReader, limit)

	// 5. Progress bar.
	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(lr)

	// 6. Actual copy invokation.
	_, err = io.Copy(toWriter, barReader)
	if err != nil {
		panic(err)
	}

	bar.Finish()

	return nil
}
