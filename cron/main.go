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
	"github.com/verizonlabs/northstar/cron/env"
	"github.com/verizonlabs/northstar/cron/service"
	cronDataClient "github.com/verizonlabs/northstar/data/cron/client"
	snippetsDataClient "github.com/verizonlabs/northstar/data/snippets/client"
	processingClient "github.com/verizonlabs/northstar/processing/snippets/client"
)

func main() {
	webPort, err := env.GetWebPort()
	if err != nil {
		mlog.Error("Failed to get cron port")
		os.Exit(-1)
	}

	cronDataClient, err := cronDataClient.NewCronClient()
	if err != nil {
		mlog.Error("Failed to init cron data client: %v", err)
		os.Exit(-1)
	}

	snippetsDataClient, err := snippetsDataClient.NewSnippetClient()
	if err != nil {
		mlog.Error("Failed to init snippets data client: %v", err)
		os.Exit(-1)
	}

	processingClient, err := processingClient.NewSnippetsClient()
	if err != nil {
		mlog.Error("Failed to init processing client: %v", err)
		os.Exit(-1)
	}

	cronService := service.NewCronService(processingClient, snippetsDataClient, cronDataClient)
	cronService.AddRoutes()
	err = cronService.StartScheduler()
	if err != nil {
		mlog.Error("Failed to start scheduler: %v", err)
		os.Exit(-1)
	}

	port := ":" + webPort
	if err := management.Listen(port); err != nil {
		mlog.Error("Error starting cron service", err)
	}
}
