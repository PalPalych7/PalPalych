package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func printVersion() { //nolint
	var (
		release   = "UNKNOWN"
		buildDate = "UNKNOWN"
		gitHash   = "UNKNOWN"
	)
	if err := json.NewEncoder(os.Stdout).Encode(struct {
		Release   string
		BuildDate string
		GitHash   string
	}{
		Release:   release,
		BuildDate: buildDate,
		GitHash:   gitHash,
	}); err != nil {
		fmt.Printf("error while decode version info: %v\n", err)
	}
}
