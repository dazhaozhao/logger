package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"
)

const (
	DEBUG = iota
	TRACE
	INFO
	WARN
	ERROR
	FATAL
)

var logger *Logger

type Log interface {
	// Debug logs a message at level Debug on the standard logger.
	DEBUG(message string)
	// TRACE logs a message at level Info on the standard logger.
	TRACE(message string)
	// Info logs a message at level Info on the standard logger.
	INFO(message string)
	// Warn logs a message at level Warn on the standard logger.
	WARN(message string)
	// Error logs a message at level Error on the standard logger.
	ERROR(message string)
	// Fatal logs a message at level Fatal on the standard logger.
	FATAL(message string)
}

// Logger
type Logger struct {
	level    int8
	path     string
	filename string
	stdout   bool
	file     *os.File
	time     *time.Time
}

func Init(path string, level int8, stdout bool) {
	if path == "" {
		path = "./"
	}
	log := &Logger{level: level, path: path, stdout: stdout}
	logger = log.newLogFile()
}
func checkLoggerHasInit() {
	if logger == nil {
		panic("logger not init yet")
	}
}

// implement Log interface
func Debug(message string) {
	checkLoggerHasInit()
	logger.printLog(DEBUG, message)
}

func Trace(message string) {
	checkLoggerHasInit()
	logger.printLog(TRACE, message)
}

func Info(message string) {
	checkLoggerHasInit()
	logger.printLog(INFO, message)
}

func Warn(message string) {
	checkLoggerHasInit()
	logger.printLog(WARN, message)
}
func Error(message string) {
	checkLoggerHasInit()
	logger.printLog(ERROR, message)
}

func Fatal(message string) {
	checkLoggerHasInit()
	logger.printLog(FATAL, message)
}

// log info by level
func (logger *Logger) printLog(level int8, message string) {
	if !logger.isOneDay() {
		logger = logger.newLogFile()
	}
	var _level string
	switch level {
	case DEBUG:
		_level = "DEBUG"
	case TRACE:
		_level = "TRACE"
	case INFO:
		_level = "INFO"
	case WARN:
		_level = "WARN"
	case ERROR:
		_level = "ERROR"
	case FATAL:
		_level = "FATAL"
	}

	if level >= logger.level {
		_, file, line, ok := runtime.Caller(2)
		if !ok {
			fmt.Printf("[%s] [%s] [%s:%d] %s\n", time.Now().Format("2006-01-02 15:04:05"), "ERROR", file, line, "runtime.Caller() fail")
			return
		}
		if logger.stdout {
			// 日志信息输出到控制台
			fmt.Printf("[%s] [%s] [%s:%d] %s\n", time.Now().Format("2006-01-02 15:04:05"), _level, file, line, message)
		}
		// 日志信息输出到文件
		fmt.Fprintf(logger.file, "[%s] [%s] [%s:%d] %s\n", time.Now().Format("2006-01-02 15:04:05"), _level, file, line, message)
	}
}

func (logger Logger) isOneDay() bool {
	now := time.Now()
	day := logger.time
	return now.Day() == day.Day() && now.Month() == day.Month() && now.Year() == day.Year()
}

func (logger *Logger) newLogFile() *Logger {
	now := time.Now()
	filename := now.Format("20060102") + ".log"
	road := logger.path
	newPath := path.Join(road, filename)
	exist := isExist(road)
	if !exist {
		os.MkdirAll(road, os.ModePerm)
	}
	file, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		_, f, line, _ := runtime.Caller(0)
		fmt.Fprintln(logger.file, "[", time.Now().Format("2006-01-02 15:04:05"), "]", "[ERROR]", "[", f, ":", line, "]", "create log file failed!")
		return logger
	}

	logger.filename = filename
	logger.file = file
	logger.time = &now
	return logger
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}
