package main

import (
	"os"
	"os/exec"
	"strings"
)

func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return -1
	}
	if len(env) == 0 {
		return -2
	}
	for key, val := range env {
		if val.NeedRemove { // нужно удалять
			os.Unsetenv(key)
		} else { // нужно добавлять
			os.Setenv(key, val.Value)
		}
	}

	//	paramsStr := ""
	//	for i := 1; i < len(cmd); i++ {
	//		paramsStr = paramsStr + cmd[i] + " "
	//	}
	//	paramsStr = strings.TrimRight(paramsStr, " ")
	paramsStr := strings.Join(cmd[1:], " ")
	myCmd := exec.Command(cmd[0], paramsStr) //nolint:gosec
	myCmd.Stdout = os.Stdout
	if err := myCmd.Run(); err != nil {
		return 1
	}
	return 0
}
