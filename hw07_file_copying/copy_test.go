package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	// исходный файл не существует
	myEr := Copy("1.txt", "2.txt", 0, 0)
	require.Error(t, myEr)

	// offset<0
	myEr = Copy("testdata/input.txt", "2.txt", -2, 0)
	require.Error(t, myEr)
	require.Truef(t, errors.Is(myEr, ErrOffsetLes0), "actual error %q", myEr)

	// limit<0
	myEr = Copy("testdata/input.txt", "2.txt", 0, -10)
	require.Error(t, myEr)
	require.Truef(t, errors.Is(myEr, ErrLimitLes0), "actual error %q", myEr)

	// ofset > len
	myEr = Copy("testdata/input.txt", "2.txt", 10000, 0)
	require.Error(t, myEr)
	require.Truef(t, errors.Is(myEr, ErrOffsetExceedsFileSize), "actual error %q", myEr)
	return
}
