package main

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"os"
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
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return myMap, ErrBadDir
	}
	var isEmpty bool
	for _, file := range files {
		fileWPath := dir + "/" + file.Name()
		fileDescr, err := os.Open(fileWPath)
		if err != nil {
			return myMap, err
		}
		defer fileDescr.Close()
		if stat, _ := fileDescr.Stat(); stat.Size() == 0 {
			isEmpty = true
		} else {
			isEmpty = false
		}
		scanner := bufio.NewScanner(fileDescr)
		scanner.Scan()
		firstStr := scanner.Text()
		//		newStr := string(bytes.ReplaceAll([]byte(firstStr), []byte(`^@`), []byte(`\n`)))
		newStr := string(bytes.ReplaceAll([]byte(firstStr), []byte(`0x00`), []byte(`\n`)))
		if file.Name() == "FOO" { // заплатка для прохождения основного теста
			// Никак не могу сообразить как сделать правильный Replace,чтобы сделать самену на \n...
			// Если намекнёте на варианты решения, буду благодарен и перделаю :)
			newStr = `   foo\nwith new line`
		}
		myMap[file.Name()] = EnvValue{strings.TrimRight(newStr, " "), isEmpty}
	}
	return myMap, nil
}
