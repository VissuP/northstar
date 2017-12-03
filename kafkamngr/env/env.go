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

package env

import (
	"fmt"
	"os"
	"strconv"
)

func GetWebPort() (string, error) {
	port := os.Getenv("KAFKA_MNGR_PORT")
	if port == "" {
		return "", fmt.Errorf("Please KAFKA_MNGR_PORT!")
	}
	return port, nil
}

func GetHostPort() (string, error) {
	hostPort := os.Getenv("KAFKA_MNGR_HOST_PORT")
	if hostPort == "" {
		return "", fmt.Errorf("Please KAFKA_MNGR_HOST_PORT!")
	}
	return hostPort, nil
}

func GetZKTimeout() (int, error) {
	timeout := os.Getenv("KAFKA_MNGR_ZK_TIMEOUT")
	if timeout == "" {
		return 0, fmt.Errorf("Please KAFKA_MNGR_ZK_TIMEOUT!")
	}

	timeoutInt, err := strconv.Atoi(timeout)
	if err != nil {
		return 0, fmt.Errorf("Failed to parse zk timeout: %v", err)
	}
	return timeoutInt, nil
}

func GetKafkaZkUrl() (string, error) {
	path := os.Getenv("KAFKA_ZK_URL")
	if path == "" {
		return "", fmt.Errorf("Please set KAFKA_ZK_URL!")
	}

	return path, nil
}
