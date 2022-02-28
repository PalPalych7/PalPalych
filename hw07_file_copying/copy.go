package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrLimitLes0             = errors.New("limiters <0")
	ErrOffsetLes0            = errors.New("offset <0")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if limit < 0 {
		return ErrLimitLes0
	}
	if offset < 0 {
		return ErrOffsetLes0
	}
	fmt.Println(fromPath, toPath, limit, offset)
	inputFile, err := os.Open(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer inputFile.Close()
	stat, _ := inputFile.Stat()

	if offset > stat.Size() {
		return ErrOffsetExceedsFileSize
	}
	if limit == 0 || limit > stat.Size()-offset {
		limit = stat.Size() - offset
	}

	fmt.Println(inputFile, err)
	outputFile, err := os.Create(toPath)
	if err != nil {
		err = ErrUnsupportedFile
	}
	defer outputFile.Close()
	fmt.Println(&outputFile, outputFile, err)
	inputFile.Seek(offset, 0)
	/*
		_, err := io.CopyN(outputFile, inputFile, limit)
		if err == io.EOF {
			err = nil
		}
	*/
	byteN := 0
	var bufLen int64 = 1
	buf := make([]byte, bufLen)
	bar := pb.StartNew(int(limit))
	for byteN < int(limit) {
		_, err = inputFile.Read(buf)
		byteN += int(bufLen)
		if errors.Is(err, io.EOF) {
			err = nil
			break
		}
		outputFile.Write(buf)
		bar.Increment()
		time.Sleep(time.Millisecond)
	}
	bar.FinishPrint("The End!")

	return err
}
