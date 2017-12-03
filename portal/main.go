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
	"os"
	"os/signal"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/portal/config"
	"github.com/verizonlabs/northstar/portal/service"
	"syscall"
)

const (
	Success       int = 0
	CreationError int = 1
	StartupError  int = 2
)

var (
	signalChannel chan os.Signal
)

// Defines service main entry point.
func main() {
	signalChannel = make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGABRT)
	go signalHander()

	// Setup function used to recover from panics.
	defer func() {
		if r := recover(); r != nil {
			mlog.Error("%s failed with panic %v", config.Configuration.ServiceName, r)
		}
	}()

	// Create the service
	mlog.Event("%s starting", config.Configuration.ServiceName)
	service, err := service.NewService()
	if err != nil {
		mlog.Error("Failed to start service with error %s.\n", err.Error())
		os.Exit(CreationError)
	}

	// Start the service.
	if err = service.Start(); err != nil {
		mlog.Error("Failed to start service with error %s.\n", err.Error())
		os.Exit(StartupError)
	}

	mlog.Event("%s shutdown", config.Configuration.ServiceName)
}

// Helper method used to handle abort and termination signals.
func signalHander() {
	sig := <-signalChannel

	switch sig {
	case os.Interrupt:
		fallthrough
	case syscall.SIGABRT:
		fallthrough
	case syscall.SIGTERM:
		os.Exit(0)
	}
}
