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

package util

import (
	"fmt"
	"os"
)

const StreamBasePath = "/stream/v1"

func GetStreamHostPort() (string, error) {
	hostPort := os.Getenv("DPE_STREAM_HOST_PORT")
	if hostPort == "" {
		return "", fmt.Errorf("Please set DPE_STREAM_HOST_PORT!")
	}

	return hostPort, nil
}

func GetStreamBaseUrl() (string, error) {
	hostPort, err := GetStreamHostPort()
	if err != nil {
		return "", err
	}

	return "http://" + hostPort, nil
}
