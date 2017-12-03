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
	"github.com/verizonlabs/northstar/processing/env"
	"github.com/verizonlabs/northstar/processing/events"
	"github.com/verizonlabs/northstar/processing/snippets"
)

func main() {
	webPort, err := env.GetWebPort()
	if err != nil {
		mlog.Error("Failed to retrieve processing port: %v", err)
		os.Exit(-1)
	}

	snippetsService, err := snippets.NewSnippetsService()
	if err != nil {
		mlog.Error("Failed to init snippets service: %v", err)
		os.Exit(-1)
	}
	snippetsService.AddRoutes()

	eventsService := events.NewEventsService(snippetsService)
	eventsService.AddRoutes()

	port := ":" + webPort
	if err := management.Listen(port); err != nil {
		mlog.Alarm("Error starting processing service", err)
	}
}
