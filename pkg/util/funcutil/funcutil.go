package funcutil

import (
	"fmt"
	"runtime"
	"strings"
)

func IsCalledFromInit() bool {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])

	for {
		frame, more := frames.Next()
		if !more {
			break
		}

		// check if the function name is "init"
		funcName := runtime.FuncForPC(frame.PC).Name()
		fmt.Println(funcName)
		if funcName == "init" || strings.HasSuffix(funcName, ".init") || strings.Contains(funcName, ".init.") {
			return true
		}
	}

	return false
}
