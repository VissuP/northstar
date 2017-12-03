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
	"github.com/verizonlabs/northstar/object/service"
	"syscall"
)

var (
	signalChannel chan os.Signal
)

func main() {
	signalChannel = make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGABRT)
	go signalHander()

	objectService, err := service.NewService()
	if err != nil {
		mlog.Error("Failed to create new service: %v", err)
		os.Exit(-1)
	}

	err = objectService.Start()
	if err != nil {
		mlog.Error("Failed to start web server: %v", err)
		os.Exit(-1)
	}
}

func signalHander() {
	sig := <-signalChannel
	switch sig {
	case os.Interrupt:
		fallthrough
	case syscall.SIGABRT:
		fallthrough
	case syscall.SIGTERM:
		mlog.Event("Shutting down service due to signal %s", sig.String())
		os.Exit(0)
	}
}
