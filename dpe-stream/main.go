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
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/dpe-stream/config"
	"github.com/verizonlabs/northstar/dpe-stream/master/service"
	"github.com/verizonlabs/northstar/dpe-stream/worker"
)

const (
	MASTER = "master"
	WORKER = "worker"
)

func main() {
	if len(os.Args) != 2 {
		mlog.Error("Usage: dpe-stream <master|worker>")
		os.Exit(-1)
	}

	option := os.Args[1]
	switch option {
	case MASTER:
		mlog.Info("Starting in master mode")
		service, err := service.NewSteamService()
		if err != nil {
			mlog.Error("Failed to create stream service: %v", err)
			os.Exit(-1)
		}
		service.AddRoutes()
	case WORKER:
		mlog.Info("Starting in worker mode")
		err := worker.StartWorker()
		if err != nil {
			mlog.Error("Failed to start worker: %v", err)
			os.Exit(-1)
		}
	default:
		mlog.Error("Wrong option selected: %v", option)
		os.Exit(-1)
	}

	mlog.Info("Starting web server on port: %s", config.WebPort)
	port := ":" + config.WebPort
	if err := management.Listen(port); err != nil {
		mlog.Error("Error starting web server", err)
		os.Exit(-1)
	}
}
