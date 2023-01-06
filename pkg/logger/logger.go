package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/TwiN/go-color"
)

type Logger struct {
	DebugMode bool
	NoColor   bool
	Output    io.Writer
}

const (
	LOG   = "      "
	INFO  = "INFO  "
	WARN  = "WARN  "
	ERROR = "ERROR "
	DEBUG = "DEBUG "
	FATAL = "FATAL "
)

var (
	Log = &Logger{
		DebugMode: false,
		NoColor:   false,
		Output:    os.Stdout,
	}
	FileLog *log.Logger
	LogFile *os.File
)

func SetLogger(debugMode bool) {
	Log.DebugMode = debugMode
}

func SetFileOutput(fileName string) {
	var err error
	fileName, err = filepath.Abs(fileName)
	if err != nil {
		Fatalf("Failed to get absolute path for log file: %v", err)
	}
	if LogFile, err = os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		Fatalf("Failed to open log file: %v", err)
	}
	FileLog = log.New(LogFile, "", log.LstdFlags)
	Debugf("Log file: %s", fileName)
}

func CloseLogging() {
	if LogFile != nil {
		LogFile.Sync()
		LogFile.Close()
		Log.Output = os.Stdout
		LogFile = nil
	}
}

func DisableColor() {
	Log.NoColor = true
}

func printf(level, clr string, format string, args ...interface{}) {
	line := strings.TrimRight(fmt.Sprintf(format, args...), "\r\n")
	if !Log.NoColor {
		line = color.Ize(clr, line)
	}
	fmt.Fprint(Log.Output, line+"\n")

	if FileLog != nil {
		FileLog.Printf(level+" "+format, args...)
	}
}

func Logf(format string, args ...interface{}) {
	printf(LOG, color.White, format+"\n", args...)
}

func Infof(format string, args ...interface{}) {
	printf(INFO, color.Green, format, args...)
}

func Debugf(format string, args ...interface{}) {
	if Log.DebugMode {
		printf(DEBUG, color.Blue, format, args...)
	}
}

func Warnf(format string, args ...interface{}) {
	printf(WARN, color.Yellow, format, args...)
}

func Fatalf(format string, args ...interface{}) {
	printf(FATAL, color.Red, format, args...)
	CloseLogging()
	os.Exit(1)
}

func PurgeOldLogs(dataFolder string, logsToKeep int) (err error) {
	if logsToKeep < 1 {
		Warnf("Invalid logsToKeep value: %d", logsToKeep)
		return
	}
	files, err := filepath.Glob(filepath.Join(dataFolder, "*.log"))

	if err != nil {
		return err
	}
	if len(files) <= logsToKeep {
		Debugf("No logs to purge (logsToKeep: %d)", logsToKeep)
		return nil
	}
	for _, file := range files[:len(files)-logsToKeep] {
		if err := os.Remove(file); err != nil {
			return err
		}
		Debugf("Purged log file: %s", file)
	}
	return nil
}
