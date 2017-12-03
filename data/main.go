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
	"github.com/verizonlabs/northstar/data/cron"
	"github.com/verizonlabs/northstar/data/datasets"
	"github.com/verizonlabs/northstar/data/datasources"
	"github.com/verizonlabs/northstar/data/env"
	"github.com/verizonlabs/northstar/data/events"
	"github.com/verizonlabs/northstar/data/invocations"
	"github.com/verizonlabs/northstar/data/mappings"
	"github.com/verizonlabs/northstar/data/notebooks"
	"github.com/verizonlabs/northstar/data/snippets"
	"github.com/verizonlabs/northstar/data/stream"
	"github.com/verizonlabs/northstar/data/templates"
)

type DataService interface {
	AddRoutes()
}

func main() {
	dataPort, err := env.GetDataPort()
	if err != nil {
		mlog.Error("Failed to load data port", err)
		os.Exit(-1)
	}

	var dataService DataService

	dataService = new(snippets.SnippetService)
	dataService.AddRoutes()

	dataService = new(invocations.InvocationService)
	dataService.AddRoutes()

	dataService = new(events.EventsService)
	dataService.AddRoutes()

	dataService = new(mappings.MappingsService)
	dataService.AddRoutes()

	dataService = new(cron.CronService)
	dataService.AddRoutes()

	dataService = new(notebooks.NotebookService)
	dataService.AddRoutes()

	dataService = new(datasources.DatasourcesService)
	dataService.AddRoutes()

	dataService = new(datasets.DatasetsService)
	dataService.AddRoutes()

	dataService = new(templates.TemplateService)
	dataService.AddRoutes()

	dataService = new(stream.StreamService)
	dataService.AddRoutes()

	port := ":" + dataPort
	if err := management.Listen(port); err != nil {
		mlog.Error("Error starting api service", err)
		os.Exit(-1)
	}
}
