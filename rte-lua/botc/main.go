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
	"github.com/verizonlabs/northstar/pkg/mlog"
	kafkaClient "github.com/verizonlabs/northstar/kafkamngr/client"
	"github.com/verizonlabs/northstar/kafkamngr/model"
	marathonClient "github.com/verizonlabs/northstar/pkg/marathon"
	"github.com/verizonlabs/northstar/rte-lua/botc/config"
	"strconv"
)

func main() {
	mlog.Info("Starting RTE-Lua Botc")
	mClient, err := marathonClient.NewMarathonClient()
	if err != nil {
		mlog.Error("Failed to create client: %v", err)
		os.Exit(-1)
	}

	kClient, err := kafkaClient.NewKafkaMngrClient()
	if err != nil {
		mlog.Error("Failed to create kafka mngr client: %v", err)
		os.Exit(-1)
	}

	app, err := marathonClient.GetApplicationFromJson(config.RTELuaMarathonJson)
	if err != nil {
		mlog.Error("Failed to get application from json: %v", err)
		os.Exit(-1)
	}

	nWorkers, err := getNumberOfWorkers(*app.Env)
	if err != nil {
		mlog.Error("Failed to get number of workers: %v", err)
		os.Exit(-1)
	}

	err = mClient.CreateApplication(app)
	if err != nil {
		mlog.Error("Failed to create application: %v", err)
	}

	nPartitions := *app.Instances * nWorkers
	mlog.Debug("Number of instances %v, workers %v and partitions %v",
		*app.Instances, nWorkers, nPartitions)

	topic := model.Topic{Partitions: nPartitions}
	mErr := kClient.UpdateTopic(config.RTEServiceName, config.RTELuaCtrlTopic, &topic)
	if mErr != nil {
		mlog.Error("Failed to update topic: %s", mErr.Error())
	}

	err = mClient.DeleteApplication(config.RTEBotcAppName)
	if err != nil {
		mlog.Error("Failed to cleanup myself: %v", err)
	}
}

func getNumberOfWorkers(env map[string]string) (int, error) {
	numWorkers := env["NUM_WORKERS"]
	if numWorkers == "" {
		return 0, fmt.Errorf("Please set NUM_WORKERS")
	}

	nWorkers, err := strconv.Atoi(numWorkers)
	if err != nil {
		return 0, fmt.Errorf("Failed to convert number of workers to int: %v", err)
	}

	if nWorkers < 1 {
		return 0, fmt.Errorf("Number of workers must be greater than zero!")
	}

	return nWorkers, nil
}
