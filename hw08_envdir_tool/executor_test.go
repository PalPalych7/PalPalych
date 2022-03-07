package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	myMap := make(Environment)
	parmList := []string{"/test.sh", "parm1"}
	myRes := RunCmd(parmList, myMap)
	require.Equal(t, -2, myRes)

	parmList = []string{}
	myMap["key1"] = EnvValue{"val1", true}
	myRes = RunCmd(parmList, myMap)
	require.Equal(t, -1, myRes)
}
