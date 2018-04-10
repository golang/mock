package gomock

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const skipFrames = 2

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func stackTraceStringFromError(err error, skipFrames int) string {
	if err, ok := err.(stackTracer); ok {
		return stackTraceString(err.StackTrace(), skipFrames)
	}
	return ""
}

func currentStackTrace(skipFrames int) string {
	err := errors.New("fake error just to get stack")
	if err, ok := err.(stackTracer); ok {
		return stackTraceString(err.StackTrace(), skipFrames)
	}
	return ""
}

func stackTraceString(stackTrace errors.StackTrace, skipFrames int) string {
	buffer := bytes.NewBufferString("")
	for i := skipFrames + 1; i < len(stackTrace); i++ {
		frame := stackTrace[i]
		buffer.WriteString(fmt.Sprintf("%+v\n", frame))
		filename := fmt.Sprintf("%s", frame)
		if strings.Contains(filename, "_test.go") {
			break
		}
	}
	return buffer.String()
}
