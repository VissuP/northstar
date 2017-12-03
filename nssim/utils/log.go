/*
Copyright (C) 2017 Verizon. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/verizonlabs/northstar/pkg/mlog"
)

type Step string

const StepNone = ""

type StepStat struct {
	Successes int64 `json:"successes"`
	Failures  int64 `json:"failures"`
}

//Logger implements the logging infrastructure.
type Logger struct {
	PrintLogs bool
	LastError string
	Logs      []LogEntry
	index     int
	Stats     map[string]*StepStat
}

type LogEntry struct {
	Time    time.Time `json:"time"`
	Success bool      `json:"success"`
	Step    Step      `json:"step"`
	Message string    `json:"message"`
}

//NewLogger creates a new logger instance. printLogs controls whether logs also go to mlog.
func NewLogger(printLogs bool) *Logger {
	return &Logger{
		PrintLogs: printLogs,
		index:     0,
		Stats:     make(map[string]*StepStat),
	}
}

func (logger *Logger) getEnhancedTemplate(template string) string {
	return fmt.Sprintf("%s:%s", logger.getLineAndFileNumber(), template)
}

func (logger *Logger) incrementStat(step Step, success bool) {
	statName := fmt.Sprintf("%02d %s", len(logger.Stats)+1, step)
	stat, ok := logger.Stats[statName]
	if !ok {
		stat = &StepStat{}
		logger.Stats[statName] = stat
	}

	if success {
		stat.Successes += 1
	} else {
		stat.Failures += 1

	}
}

func (logger *Logger) LogStep(step Step, success bool, template string, message ...interface{}) error {
	enhancedTemplate := logger.getEnhancedTemplate(template)
	if logger.PrintLogs {
		mlog.Error(enhancedTemplate, message...)
	}

	logEntry := LogEntry{
		Time:    time.Now(),
		Success: success,
		Step:    step,
		Message: fmt.Sprintf(enhancedTemplate, message...),
	}

	logger.incrementStat(step, success)
	logger.Logs = append(logger.Logs, logEntry)

	if !success {
		return fmt.Errorf(template, message...)
	}
	return nil
}

//LogError allows tests to log errors. Optionally it can also log to the service logger.
func (logger *Logger) LogError(template string, message ...interface{}) error {
	enhancedTemplate := logger.getEnhancedTemplate(template)
	if logger.PrintLogs {
		mlog.Error(enhancedTemplate, message...)
	}

	logger.LastError = fmt.Sprintf(template, message...)
	logEntry := LogEntry{
		Time:    time.Now(),
		Success: false,
		Step:    StepNone,
		Message: fmt.Sprintf(enhancedTemplate, message...),
	}

	logger.Logs = append(logger.Logs, logEntry)
	return fmt.Errorf(template, message...)
}
func (logger *Logger) LogDebug(template string, message ...interface{}) {
	enhancedTemplate := logger.getEnhancedTemplate(template)
	if logger.PrintLogs {
		mlog.Debug(enhancedTemplate, message...)
	}

	logEntry := LogEntry{
		Time:    time.Now(),
		Success: true,
		Step:    StepNone,
		Message: fmt.Sprintf(enhancedTemplate, message...),
	}

	logger.Logs = append(logger.Logs, logEntry)
}

//LogInfo allows tests to log info messages. Optionally it can also log to the service logger.
func (logger *Logger) LogInfo(template string, message ...interface{}) {
	enhancedTemplate := logger.getEnhancedTemplate(template)
	if logger.PrintLogs {
		mlog.Info(enhancedTemplate, message...)
	}

	logEntry := LogEntry{
		Time:    time.Now(),
		Success: true,
		Step:    StepNone,
		Message: fmt.Sprintf(enhancedTemplate, message...),
	}

	logger.Logs = append(logger.Logs, logEntry)
}

//getLineAndFileNumber allows the logging system to pick up the file and line number were a message was generated.
func (logger *Logger) getLineAndFileNumber() string {
	// get caller statistics
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???"
		line = 0
	}

	// determine the shorted-version of the filename
	// and avoid the func call of strings.SplitAfter
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short
	linx := strconv.Itoa(line)
	fileAndLine := strings.Join([]string{file, linx}, ":")
	return fileAndLine
}
