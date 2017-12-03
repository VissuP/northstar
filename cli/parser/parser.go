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

package parser

import (
	"fmt"
	"os"

	"github.com/verizonlabs/northstar/cli/commands"
	"github.com/verizonlabs/northstar/cli/commands/cron"
	"github.com/verizonlabs/northstar/cli/commands/datasets"
	"github.com/verizonlabs/northstar/cli/commands/datasources"
	"github.com/verizonlabs/northstar/cli/commands/events"
	"github.com/verizonlabs/northstar/cli/commands/invoke"
	"github.com/verizonlabs/northstar/cli/commands/kafkamngr"
	"github.com/verizonlabs/northstar/cli/commands/mappings"
	"github.com/verizonlabs/northstar/cli/commands/object"
	"github.com/verizonlabs/northstar/cli/commands/snippets"
	cronClient "github.com/verizonlabs/northstar/cron/client"
	cronDataClient "github.com/verizonlabs/northstar/data/cron/client"
	datasetsData "github.com/verizonlabs/northstar/data/datasets/client"
	datasourcesData "github.com/verizonlabs/northstar/data/datasources/client"
	eventsDataClient "github.com/verizonlabs/northstar/data/events/client"
	invocationsDataClient "github.com/verizonlabs/northstar/data/invocations/client"
	mappingsDataClient "github.com/verizonlabs/northstar/data/mappings/client"
	snippetsDataClient "github.com/verizonlabs/northstar/data/snippets/client"
	kafkaMngrClient "github.com/verizonlabs/northstar/kafkamngr/client"
	objectClient "github.com/verizonlabs/northstar/object/client"
	processEventsClient "github.com/verizonlabs/northstar/processing/events/client"
	processSnippetsClient "github.com/verizonlabs/northstar/processing/snippets/client"
)

func printError(err error) {
	fmt.Printf("Command failed: %v", err)
	os.Exit(1)
}

func InitParser(kafkaMngr *kafkaMngrClient.KafkaMngrClient,
	snippetsProcessing *processSnippetsClient.SnippetsClient,
	eventsProcessing *processEventsClient.EventsClient,
	snippetsData *snippetsDataClient.SnippetsClient,
	invocationData *invocationsDataClient.InvocationClient,
	eventsData *eventsDataClient.EventsClient,
	mappingsData *mappingsDataClient.MappingsClient,
	datasetsData *datasetsData.DatasetsClient,
	datasourcesData *datasourcesData.DatasourcesClient,
	cronClient *cronClient.CronClient,
	cronDataClient *cronDataClient.CronClient,
	objectClient *objectClient.ObjectClient) {
	if len(os.Args) == 1 {
		commands.PrintHelp()
		return
	}

	var err error

	// Kafkamngr cmd
	addTopic := kafkamngr.NewAddTopics(kafkaMngr)
	listTopics := kafkamngr.NewListTopics(kafkaMngr)
	updateTopic := kafkamngr.NewUpdateTopics(kafkaMngr)

	// Snippets cmd
	addSnippet := snippets.NewAddSnippet(snippetsData)
	invokeSnippetDirect := snippets.NewInvokeSnippetDirect(snippetsProcessing)
	invokeSnippetByName := snippets.NewInvokeSnippetById(snippetsProcessing)
	listSnippet := snippets.NewListSnippets(snippetsData)
	deleteSnippet := snippets.NewDeleteSnippet(snippetsData)
	updateSnippet := snippets.NewUpdateSnippet(snippetsData)

	// List cmd
	getInvocation := invoke.NewGetInvocation(invocationData)
	listInvocations := invoke.NewListInvocation(invocationData)
	deleteInvocation := invoke.NewDeleteInvocation(invocationData)
	stopInvocation := invoke.NewStopInvocation(snippetsProcessing)

	// Events cmd
	addEvent := events.NewAddEvent(eventsData)
	invokeEvent := events.NewInvokeEvent(eventsProcessing)
	listEvents := events.NewListEvents(eventsData)
	deleteEvent := events.NewDeleteEvent(eventsData)

	//Mapping cmd
	addMapping := mappings.NewAddMapping(mappingsData)
	listMapping := mappings.NewListMappings(mappingsData)
	deleteMapping := mappings.NewDeleteMapping(mappingsData)

	//Datasources cmd
	addDatasource := datasources.NewAddDatasource(datasourcesData)
	getDatasource := datasources.NewGetDatasource(datasourcesData)
	listDatasources := datasources.NewListDatasources(datasourcesData)
	deleteDatasource := datasources.NewDeleteDatasource(datasourcesData)

	//Datasets cmd
	addDataset := datasets.NewAddDataset(datasetsData)
	getDatasetById := datasets.NewGetDatasetById(datasetsData)
	getDatasetByName := datasets.NewGetDatasetByName(datasetsData)
	listDatasets := datasets.NewListDatasets(datasetsData)
	deleteDataset := datasets.NewDeleteDataset(datasetsData)

	// Cron cmd
	addCron := cron.NewAddCronJob(cronClient)
	deleteCron := cron.NewDeleteJob(cronClient)
	updateCron := cron.NewUpdateJob(cronClient)
	listCron := cron.NewListJobs(cronDataClient)

	// Object buckets cmd
	createBucket := object.NewCreateBucket(objectClient)
	listBuckets := object.NewListBucket(objectClient)
	deleteBucket := object.NewDeleteBucket(objectClient)

	// Object file cmd
	uploadFile := object.NewUploadFile(objectClient)
	downloadFile := object.NewDownloadFile(objectClient)
	listFiles := object.NewListFiles(objectClient)
	deleteFile := object.NewDeleteFile(objectClient)

	switch os.Args[1] {
	case "object-bucket-create":
		err = createBucket.Run(os.Args[2:])
	case "object-bucket-list":
		err = listBuckets.Run(os.Args[2:])
	case "object-bucket-delete":
		err = deleteBucket.Run(os.Args[2:])
	case "object-file-upload":
		err = uploadFile.Run(os.Args[2:])
	case "object-file-download":
		err = downloadFile.Run(os.Args[2:])
	case "object-file-list":
		err = listFiles.Run(os.Args[2:])
	case "object-file-delete":
		err = deleteFile.Run(os.Args[2:])
	case "topics-add":
		err = addTopic.Run(os.Args[2:])
	case "topics-list":
		err = listTopics.Run(os.Args[2:])
	case "topics-update":
		err = updateTopic.Run(os.Args[2:])
	case "snippets-add":
		err = addSnippet.Run(os.Args[2:])
	case "snippets-update":
		err = updateSnippet.Run(os.Args[2:])
	case "snippets-invoke-direct":
		err = invokeSnippetDirect.Run(os.Args[2:])
	case "snippets-invoke-by-id":
		err = invokeSnippetByName.Run(os.Args[2:])
	case "snippets-list":
		err = listSnippet.Run(os.Args[2:])
	case "snippets-delete":
		err = deleteSnippet.Run(os.Args[2:])
	case "cron-add":
		err = addCron.Run(os.Args[2:])
	case "cron-update":
		err = updateCron.Run(os.Args[2:])
	case "cron-list":
		err = listCron.Run(os.Args[2:])
	case "cron-delete":
		err = deleteCron.Run(os.Args[2:])
	case "invoke-get":
		err = getInvocation.Run(os.Args[2:])
	case "invoke-list":
		err = listInvocations.Run(os.Args[2:])
	case "invoke-stop":
		err = stopInvocation.Run(os.Args[2:])
	case "invoke-delete":
		err = deleteInvocation.Run(os.Args[2:])
	case "events-add":
		err = addEvent.Run(os.Args[2:])
	case "events-invoke":
		err = invokeEvent.Run(os.Args[2:])
	case "events-list":
		err = listEvents.Run(os.Args[2:])
	case "events-delete":
		err = deleteEvent.Run(os.Args[2:])
	case "mappings-add":
		err = addMapping.Run(os.Args[2:])
	case "mappings-list":
		err = listMapping.Run(os.Args[2:])
	case "mappings-delete":
		err = deleteMapping.Run(os.Args[2:])
	case "datasources-add":
		err = addDatasource.Run(os.Args[2:])
	case "datasources-get":
		err = getDatasource.Run(os.Args[2:])
	case "datasources-list":
		err = listDatasources.Run(os.Args[2:])
	case "datasources-delete":
		err = deleteDatasource.Run(os.Args[2:])
	case "datasets-add":
		err = addDataset.Run(os.Args[2:])
	case "datasets-get-by-id":
		err = getDatasetById.Run(os.Args[2:])
	case "datasets-get-by-name":
		err = getDatasetByName.Run(os.Args[2:])
	case "datasets-list":
		err = listDatasets.Run(os.Args[2:])
	case "datasets-delete":
		err = deleteDataset.Run(os.Args[2:])
	default:
		commands.PrintHelp()
		os.Exit(2)
	}

	if err != nil {
		printError(err)
	}
}
