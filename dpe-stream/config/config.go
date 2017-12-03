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

package config

import (
	"io/ioutil"
	"github.com/verizonlabs/northstar/pkg/config"
	"github.com/verizonlabs/northstar/pkg/mlog"
)

var (
	WebPort, _             = config.GetString("DPE_STREAM_PORT", "8080")
	MsgInterval, _         = config.GetInt("DPE_STREAM_WORKER_MSG_INTERVAL", 120)
	NumThreads, _          = config.GetInt("DPE_STREAM_WORKER_NUM_THREADS", 200)
	WorkerQueueCapacity, _ = config.GetInt("DPE_STREAM_WORKER_BUFFER_CAPACITY", 5000)
	WorkerMarathonJson, _  = config.GetString("DPE_STREAM_WORKER_MARATHON_JSON", readLocalMarathonFile())
)

func readLocalMarathonFile() string {
	file, err := ioutil.ReadFile("./marathon.json")
	if err != nil {
		mlog.Debug("Could not read local marathon file: %v", err)
		return ""
	}

	return string(file)
}
