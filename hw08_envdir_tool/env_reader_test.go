package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	_, myEr := ReadDir("bad dir")
	require.Error(t, myEr)
	require.Truef(t, errors.Is(myEr, ErrBadDir), "actual error %q", myEr)
}
