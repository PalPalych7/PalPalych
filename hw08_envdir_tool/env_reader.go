package main

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

var ErrBadDir = errors.New("bad dir")

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	myMap := make(Environment)
	files, err := os.ReadDir(dir)
	if err != nil {
		return myMap, ErrBadDir
	}
	var isEmpty bool
	for _, file := range files {
		fileWPath := filepath.Join(dir, file.Name())
		fileDescr, err := os.Open(fileWPath)
		var newStr string
		if err != nil {
			return myMap, err
		}
		defer fileDescr.Close()
		if stat, _ := fileDescr.Stat(); stat.Size() == 0 {
			isEmpty = true
		} else {
			isEmpty = false
			scanner := bufio.NewScanner(fileDescr)
			scanner.Scan()
			firstStr := scanner.Bytes()
			var old []byte
			old = append(old, 0x00)
			newStr = string(bytes.ReplaceAll(firstStr, old, []byte(`\n`)))
		}
		myMap[file.Name()] = EnvValue{strings.TrimRight(newStr, " "), isEmpty}
	}
	return myMap, nil
}
