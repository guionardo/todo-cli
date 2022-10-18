package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/TwiN/go-color"
)

type Logger struct {
	DebugMode bool
	Output    io.Writer
}

var Log = &Logger{
	DebugMode: false,
	Output:    os.Stdout,
}

func SetLogger(debugMode bool) {
	Log.DebugMode = debugMode
}

func printf(clr string, format string, args ...interface{}) {
	fmt.Fprintf(Log.Output, color.Ize(fmt.Sprintf(format, args...), clr))
	fmt.Fprintln(Log.Output)
}

func Logf(format string, args ...interface{}) {
	fmt.Fprintf(Log.Output, format+"\n", args...)
}

func Infof(format string, args ...interface{}) {
	printf(color.Green, format, args...)
}

func Debugf(format string, args ...interface{}) {
	if Log.DebugMode {
		printf(color.Blue, format, args...)
	}
}

func Warnf(format string, args ...interface{}) {
	printf(color.Yellow, format, args...)
}

func Fatalf(format string, args ...interface{}) {
	printf(color.Red, format, args...)
	os.Exit(1)
}
