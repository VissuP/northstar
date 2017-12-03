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
	"github.com/verizonlabs/northstar/kafkamngr/cluster"
	kafkaMngrEnv "github.com/verizonlabs/northstar/kafkamngr/env"
	"github.com/verizonlabs/northstar/kafkamngr/service"
)

func initBackend() (cluster.KafkaCluster, error) {
	zkUrl, err := kafkaMngrEnv.GetKafkaZkUrl()
	if err != nil {
		return nil, err
	}

	zkTimeout, err := kafkaMngrEnv.GetZKTimeout()
	if err != nil {
		return nil, err
	}

	return cluster.NewNativeKafka(zkUrl, zkTimeout)
}

func main() {
	webPort, err := kafkaMngrEnv.GetWebPort()
	if err != nil {
		mlog.Error("Failed to get API port: %v", err)
		os.Exit(-1)
	}

	backend, err := initBackend()
	if err != nil {
		mlog.Error("Failed to init backend: %v", err)
		os.Exit(-1)
	}

	kafkaMngr := new(service.KafkaMngrService)
	kafkaMngr.SetKafkaCluster(backend)
	kafkaMngr.AddRoutes()

	port := ":" + webPort
	if err := management.Listen(port); err != nil {
		mlog.Error("Error starting service", err)
	}
}
