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
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/nssim/service"
)

var (
	signalChannel chan os.Signal
)

func main() {
	signalChannel = make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGABRT)
	go signalHander()

	// set-up recover
	defer func() {
		if r := recover(); r != nil {
			if r := recover(); r != nil {
				log.Alarm("ALARM: panic occurred %v", r)
			}
		}
	}()

	log.Info("Creating Services")

	// Create Service
	service, err := service.NewService()

	// In case of error, stop execution
	if err != nil {
		errorMessage := fmt.Sprintf("Error, failed to start service with error %s.\n", err.Error())
		log.Error("%s", errorMessage)
		os.Exit(1)
	}

	// Run Service
	log.Debug("Starting Service...")
	err = service.Start()

	// If any error occurs, stop execution
	if err != nil {
		errorMessage := fmt.Sprintf("Error, failed to start service with error %s.\n", err.Error())
		log.Error("Error, ", errorMessage)
		os.Exit(1)
	}

	log.Debug("Shutting down service...")
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
