package loggo

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

//LogLevel LogLevel
type LogLevel int8

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
)

var levelcode = map[LogLevel]string{
	LevelDebug:   "[DEBUG]",
	LevelInfo:    "[INFO]",
	LevelWarning: "[WARNING]",
	LevelError:   "[ERROR]",
	LevelFatal:   "[FATAL]",
}

type config struct {
	serviceName string
	hostName    string
	dir         string
	level       LogLevel
	showfile    bool
	day         string
}

//Logger Logger
type Logger struct {
	out      io.Writer // destination for output
	buf      []byte    // for accumulating text to write
	level    LogLevel
	fp       *os.File
	name     string //logger name
	showfile bool   //是否显示文件
}

var (
	conf      *config
	mu        sync.Mutex
	loggerMap map[string]*Logger
)

func init() {
	hostName, _ := os.Hostname()
	conf = &config{
		level:       LevelInfo,
		hostName:    hostName,
		serviceName: "",
		dir:         "./log",
		showfile:    true,
		day:         time.Now().Format("2006-01-02"),
	}
}

//New return new logger
func New(name string) *Logger {

	// 如果logger 已经存在， 则直接返回
	for n, l := range loggerMap {
		if n == name {
			return l
		}
	}

	logger := &Logger{
		level:    conf.level,
		out:      os.Stderr,
		name:     name,
		showfile: conf.showfile,
	}

	if fp, err := conf.open(); err == nil {
		logger.fp = fp
	}

	return logger
}

//SetDir 设置日志文件
func SetDir(dir string) {
	conf.dir = dir
}

//SetLevel 设置默认level, default="INFO"
func SetLevel(level LogLevel) {
	conf.level = level
}

//mkdir 创建日志目录
func (cf *config) mkdir() error {
	if exists(cf.dir) && !isDir(cf.dir) {
		err := os.Remove(cf.dir)
		if err != nil {
			return err
		}
	}
	return os.MkdirAll(cf.dir, os.ModePerm)
}

//open 获取文件指针
func (cf *config) open() (*os.File, error) {
	if cf.dir == "" {
		return nil, errors.New("No filePath")
	}

	if err := cf.mkdir(); err != nil {
		return nil, err
	}

	var filename strings.Builder

	filename.WriteString(conf.dir)
	filename.WriteByte('/')
	filename.WriteString(conf.hostName)
	if conf.serviceName == "" {
		filename.WriteByte('.')
		filename.WriteString(conf.serviceName)
	}
	filename.WriteByte('.')
	filename.WriteString(conf.day)
	filename.WriteString(".log")

	return getFp(filename.String())
}

func getFp(filename string) (*os.File, error) {
	if !exists(filename) {
		return os.Create(filename)
	}
	if isDir(filename) {
		err := os.Remove(filename)
		if err != nil {
			return nil, err
		}
		return os.Create(filename)
	}
	return os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
}

func isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

//Clear 关闭所有logger
func Clear() {
	for _, l := range loggerMap {
		l.Close()
	}
	loggerMap = make(map[string]*Logger)
}

//ShowFile 是否展示文件名及行号
func ShowFile(show bool) {
	conf.showfile = show
}

func trace(v ...interface{}) string {
	message := fmt.Sprintln(v...)
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])

	var str strings.Builder
	str.WriteString(message + "Traceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	str.WriteString("\n")
	return str.String()
}

//SetServiceName 设置服务名
func SetServiceName(serviceName string) {
	conf.serviceName = serviceName
}
