package loggo

import (
	"fmt"
	"os"
)

// Println calls Output to print to the standard logger.
func Println(v ...interface{}) {
	std.Output(LevelInfo, fmt.Sprintln(v...))
}

//Printf format print
func Printf(format string, v ...interface{}) {
	std.Output(LevelInfo, fmt.Sprintf(format, v...))
}

//Debug print Debug message
func Debug(v ...interface{}) {
	std.Output(LevelDebug, fmt.Sprintln(v...))
}

//Info print info message
func Info(v ...interface{}) {
	std.Output(LevelInfo, fmt.Sprintln(v...))
}

//Warning print Warning message
func Warning(v ...interface{}) {
	std.Output(LevelWarning, fmt.Sprintln(v...))
}

//Error print Error message
func Error(v ...interface{}) {
	std.Output(LevelError, trace(v...))
}

//Fatal print Fatal message and os.Exit(1)
func Fatal(v ...interface{}) {
	std.Output(LevelFatal, trace(v...))
	os.Exit(1)
}
