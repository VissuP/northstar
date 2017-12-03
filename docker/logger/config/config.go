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

package config

import (
	"flag"
	"fmt"
	"github.com/verizonlabs/northstar/pkg/config"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"strings"
)

var (
	ServiceName, _                            = config.GetString("DKT_LOGGER_SVCNAME", "Logger")
	Version, _                                = config.GetString("DKT_LOGGER_VERSION", "0.3-2.8")
	StdoutTopic, _                            = config.GetString("DKT_LOGGER_STDOUT_TOPIC", "stdoutlogging")
	StderrTopic, _                            = config.GetString("DKT_LOGGER_STDERR_TOPIC", "stderrlogging")
	AppName, _                                = config.GetString("MON_APP", "unknown")
	GroupName, _                              = config.GetString("MON_GROUP", "unknown")
	MesosTaskID, _                            = config.GetString("MESOS_TASK_ID", "unknown")
	IsKafkaEnabled, _                         = config.GetBool("DKT_LOGGER_IS_KAFKA_ENABLED", true)
	DumpMsgStdout, _                          = config.GetBool("DKT_LOGGER_DUMP_MSG_STDOUT", false)
	DumpOnWriteFailure, _                     = config.GetBool("DKT_LOGGER_DUMP_ON_WRITE_FAILURE", true)
	DisableLoggerStats, _                     = config.GetBool("DKT_LOGGER_DISABLE_STATS", true)
	EnableDebug, _                            = config.GetBool("DKT_LOGGER_ENABLE_DEBUG", false)
	InitialLogDelay, _                        = config.GetInt("DKT_LOGGER_INITIAL_DELAY", 5)
	ChannelSize, _                            = config.GetInt("DKT_LOGGER_CHANNEL_SIZE", 10000)
	MaxCpus, _                                = config.GetInt("MAX_CPUS", 1)
	MaxProcs, _                               = config.GetInt("LOGGER_MAX_PROCS", 2)
	RetryIntervalSec, _                       = config.GetInt("LOGGER_RETRY_INTERVAL_SEC", 1)
	RetryMaxIntervalSec, _                    = config.GetInt("LOGGER_RETRY_MAX_INTERVAL_SEC", 15)
	StatsInterval, _                          = config.GetInt("STATS_INTERVAL", 60)
	HostIP, _                                 = config.GetString("HOST", "")
	IsStreamStdout                            = flag.Bool("ost", true, "is stream type stdout?")
	ConnectionType                            = flag.String("st", "tcp", "connection type: tcp or uds, ignored but left for backward compatibility")
	LoggerKafkaBrokers, _                     = getLoggerKafkaBrokers()
	LoggerZookeeperBrokers, _                 = getLoggerZookeeperBrokers()
	LogLimitEnabled, _                        = config.GetBool("DKT_LOG_LIMIT_ENABLED", true)
	LogLimitIntervalSec, _                    = config.GetInt("DKT_LOG_LIMIT_INTERVAL_SEC", 1)
	LogLimitThresholdPerInterval, _           = config.GetUInt64("DKT_LOG_LIMIT_THRESHOLD_PER_INTERVAL", 1000)
	LogLimitExcludeLogLevel, _                = config.GetString("DKT_LOG_LIMIT_EXCLUDE_LOGLEVEL", "")
	LogLimitMsgSize, _                        = config.GetInt("DKT_LOG_LIMIT_MSGSIZE", 65536)
	LatencySampleRate, _                      = config.GetInt("DKT_LOGGER_LATENCY_SAMPLE_RATE", 1)
	LatencyPrintRate, _                       = config.GetInt("DKT_LOGGER_LATENCY_PRINT_RATE", 0)
	LogNotifyIntervalSec, _                   = config.GetInt("DKT_LOGGER_NOTIFY_INTERVAL_SEC", 15)
	KafkaChannelCloseIntervalSec, _           = config.GetInt("DKT_LOGGER_KAFKA_CHANNEL_CLOSE_INTERVAL_SEC", 4)
	WaitBeforeExitAfterShutdownIntervalSec, _ = config.GetInt("DKT_LOGGER_WAIT_FOR_EXIT_BEFORE_SHUTDOWN_INTERVAL_SEC", 2)
	StackDumpBackTraceBuffer, _               = config.GetInt("DKT_LOGGER_STACK_DUMP_TRACE_BUFFER", 10000)
)

func getLoggerBrokers(serviceEnv string) ([]string, error) {
	hostportInfo := []string{""}
	value, err := config.GetString(serviceEnv, "")
	if err != nil || value == "" {
		mlog.Error("Environment variable %v not set", serviceEnv)
		return nil, fmt.Errorf("Environment variable %v not set", serviceEnv)
	} else {
		hostportInfo = []string(strings.Split(value, ","))
		return hostportInfo, nil
	}
}

func getLoggerKafkaBrokers() ([]string, error) {
	return getLoggerBrokers("LOGGER_KAFKA_BROKERS_HOST_PORT")
}

func getLoggerZookeeperBrokers() ([]string, error) {
	return getLoggerBrokers("LOGGER_ZOOKEEPER_HOST_PORT")
}
