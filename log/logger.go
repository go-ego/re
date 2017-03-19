// Copyright 2013 bee authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
package log

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"text/template"
	"time"
)

var errInvalidLogLevel = errors.New("logger: invalid log level")

const (
	levelCritical = iota
	levelFatal
	levelSuccess
	levelHint
	levelDebug
	levelInfo
	levelWarn
	levelError
)

var (
	sequenceNo uint64
	instance   *EgoLogger
	once       sync.Once
)

// EgoLogger logs logging records to the specified io.Writer
type EgoLogger struct {
	mu     sync.Mutex
	output io.Writer
}

// LogRecord represents a log record and contains the timestamp when the record
// was created, an increasing id, level and the actual formatted log line.
type LogRecord struct {
	ID       string
	Level    string
	Message  string
	Filename string
	LineNo   int
}

var (
	logRecordTemplate      *template.Template
	debugLogRecordTemplate *template.Template
)

// GetEgoLogger initializes the logger instance with a NewColorWriter output
// and returns a singleton
func GetEgoLogger(w io.Writer) *EgoLogger {
	once.Do(func() {
		var (
			err             error
			simpleLogFormat = `{{Now "2006/01/02 15:04:05"}} {{.Level}} ▶ {{.ID}} {{.Message}}{{EndLine}}`
			debugLogFormat  = `{{Now "2006/01/02 15:04:05"}} {{.Level}} ▶ {{.ID}} {{.Filename}}:{{.LineNo}} {{.Message}}{{EndLine}}`
		)

		// Initialize and parse logging templates
		funcs := template.FuncMap{
			"Now":     Now,
			"EndLine": EndLine,
		}
		logRecordTemplate, err = template.New("simpleLogFormat").Funcs(funcs).Parse(simpleLogFormat)
		MustCheck(err)
		debugLogRecordTemplate, err = template.New("debugLogFormat").Funcs(funcs).Parse(debugLogFormat)
		MustCheck(err)

		instance = &EgoLogger{output: NewColorWriter(w)}
	})
	return instance
}

// SetOutput sets the logger output destination
func (l *EgoLogger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = NewColorWriter(w)
}

// EndLine returns the a newline escape character
func EndLine() string {
	return "\n"
}

// MustCheck panics when the error is not nil
func MustCheck(err error) {
	if err != nil {
		panic(err)
	}
}

// IsDebugEnabled checks if DEBUG_ENABLED is set or not
func IsDebugEnabled() bool {
	debugMode := os.Getenv("DEBUG_ENABLED")
	return map[string]bool{"1": true, "0": false}[debugMode]
}

// Now returns the current local time in the specified layout
func Now(layout string) string {
	return time.Now().Format(layout)
}

func (l *EgoLogger) getLevelTag(level int) string {
	switch level {
	case levelFatal:
		return "FATAL   "
	case levelSuccess:
		return "SUCCESS "
	case levelHint:
		return "HINT    "
	case levelDebug:
		return "DEBUG   "
	case levelInfo:
		return "INFO    "
	case levelWarn:
		return "WARN    "
	case levelError:
		return "ERROR   "
	case levelCritical:
		return "CRITICAL"
	default:
		panic(errInvalidLogLevel)
	}
}

func (l *EgoLogger) getColorLevel(level int) string {
	switch level {
	case levelCritical:
		return RedBold(l.getLevelTag(level))
	case levelFatal:
		return RedBold(l.getLevelTag(level))
	case levelInfo:
		return BlueBold(l.getLevelTag(level))
	case levelHint:
		return CyanBold(l.getLevelTag(level))
	case levelDebug:
		return YellowBold(l.getLevelTag(level))
	case levelError:
		return RedBold(l.getLevelTag(level))
	case levelWarn:
		return YellowBold(l.getLevelTag(level))
	case levelSuccess:
		return GreenBold(l.getLevelTag(level))
	default:
		panic(errInvalidLogLevel)
	}
}

// mustLog logs the message according to the specified level and arguments.
// It panics in case of an error.
func (l *EgoLogger) mustLog(level int, message string, args ...interface{}) {
	// Acquire the lock
	l.mu.Lock()
	defer l.mu.Unlock()

	// Create the logging record and pass into the output
	record := LogRecord{
		ID:      fmt.Sprintf("%04d", atomic.AddUint64(&sequenceNo, 1)),
		Level:   l.getColorLevel(level),
		Message: fmt.Sprintf(message, args...),
	}

	err := logRecordTemplate.Execute(l.output, record)
	MustCheck(err)
}

// mustLogDebug logs a debug message only if debug mode
// is enabled. i.e. DEBUG_ENABLED="1"
func (l *EgoLogger) mustLogDebug(message string, file string, line int, args ...interface{}) {
	if !IsDebugEnabled() {
		return
	}

	// Change the output to Stderr
	l.SetOutput(os.Stderr)

	// Create the log record
	record := LogRecord{
		ID:       fmt.Sprintf("%04d", atomic.AddUint64(&sequenceNo, 1)),
		Level:    l.getColorLevel(levelDebug),
		Message:  fmt.Sprintf(message, args...),
		LineNo:   line,
		Filename: filepath.Base(file),
	}
	err := debugLogRecordTemplate.Execute(l.output, record)
	MustCheck(err)
}

// Debug outputs a debug log message
func (l *EgoLogger) Debug(message string, file string, line int) {
	l.mustLogDebug(message, file, line)
}

// Debugf outputs a formatted debug log message
func (l *EgoLogger) Debugf(message string, file string, line int, vars ...interface{}) {
	l.mustLogDebug(message, file, line, vars...)
}

// Info outputs an information log message
func (l *EgoLogger) Info(message string) {
	l.mustLog(levelInfo, message)
}

// Infof outputs a formatted information log message
func (l *EgoLogger) Infof(message string, vars ...interface{}) {
	l.mustLog(levelInfo, message, vars...)
}

// Warn outputs a warning log message
func (l *EgoLogger) Warn(message string) {
	l.mustLog(levelWarn, message)
}

// Warnf outputs a formatted warning log message
func (l *EgoLogger) Warnf(message string, vars ...interface{}) {
	l.mustLog(levelWarn, message, vars...)
}

// Error outputs an error log message
func (l *EgoLogger) Error(message string) {
	l.mustLog(levelError, message)
}

// Errorf outputs a formatted error log message
func (l *EgoLogger) Errorf(message string, vars ...interface{}) {
	l.mustLog(levelError, message, vars...)
}

// Fatal outputs a fatal log message and exists
func (l *EgoLogger) Fatal(message string) {
	l.mustLog(levelFatal, message)
	os.Exit(255)
}

// Fatalf outputs a formatted log message and exists
func (l *EgoLogger) Fatalf(message string, vars ...interface{}) {
	l.mustLog(levelFatal, message, vars...)
	os.Exit(255)
}

// Success outputs a success log message
func (l *EgoLogger) Success(message string) {
	l.mustLog(levelSuccess, message)
}

// Successf outputs a formatted success log message
func (l *EgoLogger) Successf(message string, vars ...interface{}) {
	l.mustLog(levelSuccess, message, vars...)
}

// Hint outputs a hint log message
func (l *EgoLogger) Hint(message string) {
	l.mustLog(levelHint, message)
}

// Hintf outputs a formatted hint log message
func (l *EgoLogger) Hintf(message string, vars ...interface{}) {
	l.mustLog(levelHint, message, vars...)
}

// Critical outputs a critical log message
func (l *EgoLogger) Critical(message string) {
	l.mustLog(levelCritical, message)
}

// Criticalf outputs a formatted critical log message
func (l *EgoLogger) Criticalf(message string, vars ...interface{}) {
	l.mustLog(levelCritical, message, vars...)
}
