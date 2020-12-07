package loggo

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

func (l *Logger) itoa(i int, wid int, suffix ...byte) {
	var b [20]byte
	bp := len(b) - 1
	j := wid
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	b[bp] = byte('0' + i)

	if j < 0 {
		l.buf = append(l.buf, b[bp:]...)
	} else {
		l.buf = append(l.buf, b[bp:bp+j]...)
	}

	l.buf = append(l.buf, suffix...)
}

func (l *Logger) formatter(file, s string, line int, level LogLevel) {
	l.buf = l.buf[:0]
	l.buf = append(l.buf, '[')
	t := time.Now()
	year, month, day := t.Date()
	l.itoa(year, 4, '-')
	l.itoa(int(month), 2, '-')
	l.itoa(day, 2, ' ')
	hour, min, sec := t.Clock()
	l.itoa(hour, 2, ':')
	l.itoa(min, 2, ':')
	l.itoa(sec, 2, '.')
	se := t.Nanosecond()
	l.itoa(se, 3, ']')

	l.buf = append(l.buf, levelcode[level]...)

	l.buf = append(l.buf, '[')
	l.buf = append(l.buf, []byte(l.name)...)
	l.buf = append(l.buf, ']')

	if l.showfile {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		l.buf = append(l.buf, '[')
		l.buf = append(l.buf, short...)
		l.buf = append(l.buf, ':')
		l.itoa(line, -1, ']')
	}

	l.buf = append(l.buf, ' ')
	l.buf = append(l.buf, s...)
}

//Output Output
func (l *Logger) Output(level LogLevel, s string) {
	var file string
	var line int
	var ok bool
	if l.showfile {
		_, file, line, ok = runtime.Caller(2)
		if !ok {
			file = "???"
			line = 0
		}
	}

	mu.Lock()
	defer mu.Unlock()

	l.formatter(file, s, line, level)

	if l.fp != nil {
		l.fp.Write(l.buf)
	}

	if l.level <= level {
		l.out.Write(l.buf)
	}
}

// Println calls Output to print to the standard logger.
func (l *Logger) Println(v ...interface{}) {
	std.Output(LevelInfo, fmt.Sprintln(v...))
}

//Printf format print
func (l *Logger) Printf(format string, v ...interface{}) {
	std.Output(LevelInfo, fmt.Sprintf(format, v...)+"\n")
}

//Debug print Debug message
func (l *Logger) Debug(v ...interface{}) {
	std.Output(LevelDebug, fmt.Sprintln(v...))
}

//Info print info message
func (l *Logger) Info(v ...interface{}) {
	l.Output(LevelInfo, fmt.Sprintln(v...))
}

//Warning print Warning message
func (l *Logger) Warning(v ...interface{}) {
	l.Output(LevelWarning, fmt.Sprintln(v...))
}

//Error print Error message
func (l *Logger) Error(v ...interface{}) {
	l.Output(LevelError, trace(v...))
}

//Fatal print Fatal message and os.Exit(1)
func (l *Logger) Fatal(v ...interface{}) {
	l.Output(LevelFatal, trace(v...))
	os.Exit(1)
}

//Close 关闭该logger
func (l *Logger) Close() {
	if l.fp != nil {
		l.fp.Close()
		l.fp = nil
	}
	l.buf = l.buf[:0]
	delete(loggerMap, l.name)
}

//SetLevel 设置 logger level
func (l *Logger) SetLevel(level LogLevel) *Logger {
	l.level = level
	return l
}

//ShowFile 是否展示文件名及行号
func (l *Logger) ShowFile(show bool) {
	l.showfile = show
}

//SetServiceName 设置服务名
func (l *Logger) SetServiceName(serviceName string) {
	if l.fp != nil {
		l.fp.Close()
		l.fp = nil
	}
	var filename strings.Builder
	filename.WriteString(conf.dir)
	filename.WriteByte('/')
	filename.WriteString(conf.hostName)
	filename.WriteByte('.')
	filename.WriteString(serviceName)
	filename.WriteByte('.')
	filename.WriteString(conf.day)
	if fp, err := getFp(filename.String()); err == nil {
		l.fp = fp
	}
}
