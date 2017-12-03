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

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/pkg/msgq"
	"github.com/verizonlabs/northstar/pkg/stats"
	"github.com/verizonlabs/northstar/docker/logger/config"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type LogLine struct {
	line string
}

type WriteLatency struct {
	count     int64
	timeTaken int64
}

//Incoming mlog compliant message format
//Marker|version|logType|mesos_taskid|processid|mon_group|appname|processname|correlationid|filename:lineno|timestamp|<msg>

type Logger struct {
	msgQ        msgq.MessageQueue
	ErrChan     chan *LogLine
	OutChan     chan *LogLine
	eSvcWg      *sync.WaitGroup
	eSvcKafkaWg *sync.WaitGroup
	logCount    uint64 //Log count in given interval
	exitReader  chan bool
}

type KafkaWriterConfig struct {
	KafkaTopic    string
	RecvChan      chan *LogLine
	TxCount       *stats.Counter
	ErrKafkaCount *stats.Counter
}

type LoggerReaderConfig struct {
	SendChan chan *LogLine
	RxCount  *stats.Counter
}

var logLevelsExcludedFromLoglimit map[string]bool
var logLatencyStatInfo WriteLatency

func populateLoglevelExcludeFromLoglimitList() {
	logLevelsExcludedFromLoglimit = make(map[string]bool)
	loglevels := strings.Split(config.LogLimitExcludeLogLevel, ",")
	for _, v := range loglevels {
		logLevelsExcludedFromLoglimit[v] = true
	}
}

func checkAndLogLatencyStats() (bMeasureLatency bool) {
	bMeasureLatency = false
	if config.LatencyPrintRate == 0 ||
		config.LatencySampleRate == 0 ||
		config.LatencyPrintRate < config.LatencySampleRate {
		// we will not measure latency
		return
	}
	count := atomic.LoadInt64(&logLatencyStatInfo.count)
	if count%int64(config.LatencySampleRate) == 0 {
		bMeasureLatency = true
	}
	if count == int64(config.LatencyPrintRate) {
		sampleCount := count / int64(config.LatencySampleRate)
		timeTaken := atomic.SwapInt64(&logLatencyStatInfo.timeTaken, 0)
		atomic.StoreInt64(&logLatencyStatInfo.count, 0)
		mlog.Info("Average Kafka write latency in us: %d, number of messages processed: %d", timeTaken/sampleCount, count)
	}
	atomic.AddInt64(&logLatencyStatInfo.count, 1)
	return
}

func updateLogLatencyStats(td int64) {
	atomic.AddInt64(&logLatencyStatInfo.timeTaken, td)
}

func (logger *Logger) stdinReader(loggerReaderConfig LoggerReaderConfig) {
	defer func() {
		logger.eSvcWg.Done()
	}()

	mlog.Debug("stdinReader: entering the method")
	ok := true
	go func() {
		<-logger.exitReader
		ok = false
	}()
	// Initial delay to allow for Kafka/Zookeper connection
	time.Sleep(time.Second * time.Duration(config.InitialLogDelay))

	reader := bufio.NewReader(os.Stdin)
	for ok {
		logLine, err := logger.readLine(reader)
		if logger.processLine(logLine, loggerReaderConfig, err) {
			mlog.Debug("stdinReader: processLine returned true, returning")
			return
		}
	}
}

func (logger *Logger) readLine(reader *bufio.Reader) (*LogLine, error) {
	var err error
	logLine := &LogLine{}
	// logging messages are based on each new line hence we should read the stream at line boundaries
	logLine.line, err = reader.ReadString('\n')
	if len(logLine.line) <= 0 || logLine.line[0] == '\n' {
		// Ignore empty line or if only character is a newline
		// This count would increase due to corresponding mlog changes which add
		// new line prior to every message to isolate spurious messages seen in error stream
		mlog.Debug("readLine: ignoring empty lines and newlines: %s", logLine.line)
		RxEmpty.Incr()
		return nil, err
	} else if config.LogLimitEnabled && len(logLine.line) > config.LogLimitMsgSize {
		mlog.Debug("readLine: ignoring large message of size %d", len(logLine.line))
		RxLargeMsgDropped.Incr()
		NotifyRxLargeMsgDropped.Incr()
		return nil, err
	}
	return logLine, err
}

// The method processLine returns true only when application it is tied to has exited
func (logger *Logger) processLine(logLine *LogLine, loggerReaderConfig LoggerReaderConfig, err error) bool {
	sendChan := loggerReaderConfig.SendChan
	if err != nil {
		if err == io.EOF {
			mlog.Error("processLine: main process exited, terminating")
			logger.closeKafkaChannels()
			return true
		} else {
			mlog.Error("processLine: failed to read logs from stdin, terminating: %v", err)
		}
		os.Exit(0)
	}
	if logLine == nil {
		mlog.Debug("processLine: logLine is nil, returning false")
		return false
	}
	//We need an efficient way to determine if this is mlog compliant message or not
	//A mlog compliant message is atleast mlog.NumHeaderFields characters long - 1 ('*' marker character and 10 '|' separator characters)
	compliantLog := false
	if len(logLine.line) >= mlog.NumHeaderFields && logLine.line[0] == mlog.MarkerChar && logLine.line[1] == mlog.SeparatorChar {
		compliantLog = true
	}
	if config.LogLimitEnabled && config.IsKafkaEnabled == true {
		limitLog := logger.checkToLimitLog(logLine.line, compliantLog)
		if limitLog {
			RxRateLimitDropped.Incr()
			NotifyRxRateLimitDropped.Incr()
			return false
		}
	}
	mlog.Debug("processLine: message received: %s", logLine.line)

	if compliantLog {
		logLine.line = strings.Join([]string{config.HostIP, logLine.line}, mlog.Separator)
		mlog.Debug("processLine: added host meta data: %s", logLine.line)
	} else {
		// Downgrade non-compliant messages to stdout
		sendChan = logger.OutChan
		timestamp := time.Now().UTC().Format("2006/01/02 15:04:05.999999999")
		logLine.line = strings.Join([]string{config.HostIP, mlog.MarkerTP, config.MesosTaskID, config.GroupName, config.AppName, timestamp, logLine.line}, mlog.Separator)
		mlog.Debug("processLine: added host meta data to non-compliant message: %s", logLine.line)
	}
	loggerReaderConfig.RxCount.Incr()
	if config.IsKafkaEnabled && logger.msgQ != nil {
		mlog.Debug("processLine: message after processing: %s", logLine.line)
		select {
		case sendChan <- logLine:
		default:
			// We reach this situation when the channel gets full because the other
			// end is not reading fast enough. We will drop this to stdout
			if config.DumpOnWriteFailure {
				mlog.Info("processLine: congestion, dropping message to stdout")
				fmt.Println(logLine.line)
			} else {
				mlog.Debug("processLine: congestion, dropping message to stdout")
				fmt.Println(logLine.line)
			}
			TxCongestionDropped.Incr()
			NotifyTxCongestionDropped.Incr()
		}
	} else {
		if config.DumpMsgStdout {
			fmt.Println(logLine.line)
		}
		// mlog.Debug("processLine: message after processing: %s", logLine.line)
	}
	return false
}

func (logger *Logger) startNotifyTimer() {
	timer := time.NewTicker(time.Duration(config.LogNotifyIntervalSec) * time.Second)
	for {
		select {
		case <-timer.C:
			count := NotifyTxCongestionDropped.Reset()
			if count > 0 {
				mlog.Alarm("Messages dropped due to congestion: %d", count)
			}
			count = NotifyRxLargeMsgDropped.Reset()
			if count > 0 {
				mlog.Alarm("Messages dropped due to large size: %d", count)
			}
			count = NotifyRxRateLimitDropped.Reset()
			if count > 0 {
				mlog.Alarm("Messages dropped due to rate limit: %d", count)
			}
			count = NotifyErrKafkaCount.Reset()
			if count > 0 {
				mlog.Alarm("Error sending messages to Kafka: %d", count)
			}
		}
	}
}

func (logger *Logger) startLogLimitTimer() {
	timer := time.NewTicker(time.Duration(config.LogLimitIntervalSec) * time.Second)
	for {
		select {
		case <-timer.C:
			count := atomic.SwapUint64(&logger.logCount, 0)
			if count > config.LogLimitThresholdPerInterval {
				dropCount := count - config.LogLimitThresholdPerInterval
				mlog.Debug("startLogLimitTimer: number of log messages dropped due to limit on log rate: %d", dropCount)
			}
		}
	}
}

func (logger *Logger) checkToLimitLog(logLine string, compliantLog bool) (limit bool) {
	// Excluded messages should not be dropped
	if compliantLog {
		bExcludedMsg := logger.isMessageExcludedFromLoglimit(logLine)
		if bExcludedMsg == true {
			return false
		}
	}
	return logger.shouldDropMessage()
}

func (logger *Logger) isMessageExcludedFromLoglimit(logLine string) bool {
	bExcludedMsg := false
	if len(logLevelsExcludedFromLoglimit) == 0 {
		return false
	}
	//Compliant Log Format: <marker>|<version>|<severity>|<mesos_task_id>|<process_id>|<group>|<app>|<exe_name>|<corelation_id>|<filename:line>|<timestamp>
	sepCount, start, end := 0, 0, 0
	//Find log severity/level
	for i, c := range logLine {
		if c == mlog.SeparatorChar {
			sepCount++
			if sepCount == 2 {
				start = i + 1
			} else if sepCount == 3 {
				end = i
				break
			}
		}
	}
	//Check if it is an excluded message
	if end != 0 {
		severity := logLine[start:end]
		if logLevelsExcludedFromLoglimit[severity] {
			bExcludedMsg = true
		}
	}
	return bExcludedMsg
}

func (logger *Logger) shouldDropMessage() bool {
	dropMsg := false
	//Drop the log if total log threshold reached
	count := atomic.AddUint64(&logger.logCount, 1)
	if count > config.LogLimitThresholdPerInterval {
		dropMsg = true
	}
	return dropMsg
}

func (logger *Logger) shouldResendMessageToKafka(err error) bool {
	serr := fmt.Sprintf("%s", err)
	if strings.Contains(serr, "kafka server: Message was too large") ||
		strings.Contains(serr, "kafka server: The request included message batch larger") {
		return false
	}
	return true
}

func (logger *Logger) startKafkaWriter(kwConfig KafkaWriterConfig) {
	var prod msgq.MsgQProducer
	var err error

	defer func() {
		if prod != nil {
			prod.Close()
		}
		logger.eSvcWg.Done()
	}()

	done := make(chan bool, 1)
	no_error := new(AtomBool)

	// Wait for the msgq library to establish initial connection with Kafka/Zookeeper
	logger.eSvcKafkaWg.Wait()

	// if replication factor is more than the number of up Kafka brokers,
	// the call to create producer may fail when the topic does not exist
	// we should retry with exponential backoff till this call succeeds
	retryInterval := config.RetryIntervalSec
	for {
		prod, err = logger.msgQ.NewProducer(
			&msgq.ProducerConfig{
				TopicName:   kwConfig.KafkaTopic,
				Partitioner: msgq.RoundRobinPartitioner,
				NotifyError: true})
		if err != nil || prod == nil {
			// in some cases, msgq library returns nil producer and no error
			// we need to check for both cases
			mlog.Error("startKafkaWriter: failed to create producer for %s, error = %v", kwConfig.KafkaTopic, err)
			time.Sleep(time.Duration(retryInterval) * time.Second)
			retryInterval *= 2
			if retryInterval > config.RetryMaxIntervalSec {
				retryInterval = config.RetryMaxIntervalSec
			}
		} else {
			break
		}
	}
	// Look for errors in sending messages; and set the error state when it happens
	// Once the error channel is empty, the error state can be cleared
	go func() {
		ok := true
		for ok {
			select {
			case evt := <-prod.ReceiveErrors():
				no_error.Set(false)
				kwConfig.ErrKafkaCount.Incr()
				NotifyErrKafkaCount.Incr()
				// We are attempting to resend the message here
				if logger.shouldResendMessageToKafka(evt.Err) {
					var t0 time.Time
					bMeasureLatency := checkAndLogLatencyStats()
					if bMeasureLatency {
						t0 = time.Now()
					}
					prod.SendMsg(evt.Msg)
					if bMeasureLatency {
						t1 := time.Since(t0).Nanoseconds() / 1000
						updateLogLatencyStats(t1)
					}
					mlog.Debug("startKafkaWriter: resending message to kafka: %s", evt.Msg)
				} else {
					mlog.Debug("startKafkaWriter: not resending message to kafka: %v", evt.Err)
				}
			case <-time.After(time.Second):
				// clear the error state
				no_error.Set(true)
				// The following, if uncommented, prints message every second
				// mlog.Debug("startKafkaWriter: set error to true, resume reading data from logger")
			case <-done:
				ok = false
			}
		}
	}()
	for {
		// When there are no errors, read the messages from the receive channel
		// Write the log message to Kafka
		// msgq library never returns error, so there is no need to check for error code
		// All messages that could not be sent, should be received through ReceiveErrors call
		// We should attempt to sent those messages again
		if no_error.Get() {
			if logLine, ok := <-kwConfig.RecvChan; ok {
				var t0 time.Time
				bMeasureLatency := checkAndLogLatencyStats()
				if bMeasureLatency {
					t0 = time.Now()
				}
				prod.Send([]byte(logLine.line))
				if bMeasureLatency {
					t1 := time.Since(t0).Nanoseconds() / 1000
					updateLogLatencyStats(t1)
				}
				mlog.Debug("startKafkaWriter: sending message to kafka: %s", logLine.line)
				kwConfig.TxCount.Incr()
			} else {
				mlog.Info("Logger stdout[%v] topic[%s] channel closed", *config.IsStreamStdout, kwConfig.KafkaTopic)
				done <- true
				break
			}
		}
	}
}

func (logger *Logger) closeKafkaChannels() {
	var wg sync.WaitGroup
	closeCh := func(ch chan *LogLine) {
		defer wg.Done()
		timer := time.NewTimer(time.Duration(config.KafkaChannelCloseIntervalSec) * time.Second)
		select {
		case _, ok := <-ch:
			if ok == true {
				mlog.Debug("Logger stdout[%v] channel open,closing it", *config.IsStreamStdout)
				close(ch)
			} else {
				mlog.Debug("Logger stdout[%v] channel already closed", *config.IsStreamStdout)
			}
		case <-timer.C:
			mlog.Debug("Logger stdout[%v] timeout::channel open,closing it [%v]", *config.IsStreamStdout, ch)
			close(ch)
		}
		timer.Stop()
	}
	wg.Add(2)
	closeCh(logger.OutChan)
	closeCh(logger.ErrChan)
	wg.Wait()
}

func (logger *Logger) sigHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c)

	go func() {
		for {
			// Block until a signal is received.
			s := <-c
			switch s {
			case syscall.SIGINT, syscall.SIGTERM:
				// Lets make logger dump message on stdout so that they can be collected from platform logging pipeline
				config.IsKafkaEnabled = false
				config.DumpMsgStdout = true
				logger.closeKafkaChannels()
				logger.exitReader <- true
				mlog.Info("Logger stdout[%v] received signal [%v].Logs dumped locally", *config.IsStreamStdout, s)
			default:
				mlog.Info("Logger stdout[%v] received unknown signal [%v]", *config.IsStreamStdout, s)
			}
		}
	}()
}

func bye(svc string) {
	mlog.Event("%s SHUTDOWN", svc)
	time.Sleep(time.Duration(config.WaitBeforeExitAfterShutdownIntervalSec) * time.Second)
}

func main() {
	flag.Parse()
	svc := fmt.Sprintf("%s::Stdout[%v]", config.ServiceName, *config.IsStreamStdout)

	defer func() {
		//Check for exiting main because of run time panic
		if r := recover(); r != nil {
			buf := make([]byte, config.StackDumpBackTraceBuffer)
			runtime.Stack(buf, false)
			mlog.Info("Run-time panic resulting in %s exiting. Error[%s] StackDump follows: %s", svc, r, string(buf))
		}
		//Emit bye message from service
		bye(svc)
	}()

	runtime.GOMAXPROCS(config.MaxProcs)
	mlog.Event("%s STARTING, version %s", svc, config.Version)
	if config.DisableLoggerStats {
		s.Disable()
	}
	mlog.EnableDebug(config.EnableDebug)
	populateLoglevelExcludeFromLoglimitList()
	var eSvcWg sync.WaitGroup
	var eSvcKafkaWg sync.WaitGroup

	logger := &Logger{
		eSvcWg:      &eSvcWg,
		eSvcKafkaWg: &eSvcKafkaWg,
		ErrChan:     make(chan *LogLine, config.ChannelSize),
		OutChan:     make(chan *LogLine, config.ChannelSize),
		exitReader:  make(chan bool, 1),
	}
	logger.sigHandler()
	eSvcKafkaWg.Add(1)
	if config.IsKafkaEnabled {
		logger.eSvcWg.Add(2)
		go logger.startKafkaWriter(KafkaWriterConfig{config.StdoutTopic, logger.OutChan, TxStdout, ErrKafkaStdout})
		go logger.startKafkaWriter(KafkaWriterConfig{config.StderrTopic, logger.ErrChan, TxStderr, ErrKafkaStderr})
	}

	//Start timer for log limiting
	if config.LogLimitEnabled {
		go logger.startLogLimitTimer()
	}
	go logger.startNotifyTimer()

	//Read from stdin
	eSvcWg.Add(1)
	if *config.IsStreamStdout {
		go logger.stdinReader(LoggerReaderConfig{logger.OutChan, RxStdout})
	} else {
		go logger.stdinReader(LoggerReaderConfig{logger.ErrChan, RxStderr})
	}

	if config.IsKafkaEnabled {
		msgQ, err := msgq.NewMsgQWithStatsOptional(config.ServiceName, config.LoggerKafkaBrokers, config.LoggerZookeeperBrokers, config.DisableLoggerStats, []string{}, true)
		if err != nil {
			mlog.Error("%s: failed to instantiate msgq client: %v", config.ServiceName, err)
			return
		}
		// Connection established with Kafka/Zookeeper
		logger.msgQ = msgQ
	}
	eSvcKafkaWg.Done()

	mlog.Event("%s READY", svc)
	logger.eSvcWg.Wait()

}
