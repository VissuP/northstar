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
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/verizonlabs/northstar/pkg/config"
	"github.com/verizonlabs/northstar/pkg/middleware"
	"github.com/verizonlabs/northstar/pkg/mlog"
)

const (
	// Defines the service name.
	ServiceName = "NorthstarAPI"
)

const (
	// Define expected environment variable names.
	ScopesFileNameEnv = "SCOPES_FILENAME"

	MaxMemoryEnv     = "EXECUTION_MEMORY_MAX"
	DefaultMemoryEnv = "EXECUTION_MEMORY_DEFAULT"

	ThingSpaceProtocolEnv      = "THINGSPACE_PROTOCOL"
	ThingSpaceSouthHostPortEnv = "THINGSPACE_SOUTH_HOST_PORT"
	ThingSpaceAuthHostPortEnv  = "THINGSPACE_AUTH_HOST_PORT"
	ThingSpaceDataHostPortEnv  = "THINGSPACE_DATA_HOST_PORT"
	ThingSpaceClientIdEnv      = "THINGSPACE_CLIENT_ID"
	ThingSpaceSecretEnv        = "THINGSPACE_SECRET"

	KafkaBrokersHostPortEnv = "KAFKA_BROKERS_HOST_PORT"
	ZookeperHostPortEnd     = "KAFKA_ZK_URL"
)

//Get our max values. Note that errors are returned by the library when default values are used. Ignore them.
var (
	MaxArgCount, _     = config.GetInt("NORTHSTARAPI_MAX_ARG_COUNT", 10)
	MaxTimeout, _      = config.GetInt("NORTHSTARAPI_MAX_TIMEOUT", 180)
	MaxCodeSize, _     = config.GetInt("NORTHSTARAPI_MAX_CODE_SIZE", 68916)
	EnforceChecksum, _ = config.GetBool("NORTHSTARAPI_ENFORCE_CHECKSUM", false)
)

var (
	// Defines the service configuration.
	Configuration = new(configuration)
)

// Defines the type used to represent service configuration.
type configuration struct {
	ServiceName          string
	Scopes               middleware.AllowedServices
	DataHostPort         string
	ProcessingHostPort   string
	CronHostPort         string
	KafkaManagerHostPort string
	KafkaBrokersHostPort []string
	ZookeeperHostPort    string
	MaxMemory            uint64
	DefaultMemory        uint64

	// ThingSpace Information
	ThingspaceProtocol      string
	ThingSpaceSouthHostPort string
	ThingSpaceAuthHostPort  string
	ThingSpaceDataHostPort  string
	ThingSpaceClientId      string
	ThingSpaceClientSecret  string
}

// Loads service configuration from the environment variables.
func Load() (err error) {
	mlog.Info("Load")

	// Set service name.
	Configuration.ServiceName = ServiceName

	// Get scopes.
	scopesFileName := os.Getenv(ScopesFileNameEnv)
	if scopesFileName == "" {
		return fmt.Errorf("The %s environment variable is missing.", ScopesFileNameEnv)
	}

	// Open the access configuration file
	scopesFile, err := os.Open(scopesFileName)
	if err != nil {
		return fmt.Errorf("Failed to load the scopes file %s with error: %v.", scopesFileName, err)
	}

	// Decode expected json.
	decoder := json.NewDecoder(scopesFile)
	if err := decoder.Decode(&Configuration.Scopes); err != nil {
		return fmt.Errorf("Failed to decode scopes file %s with error: %v", scopesFileName, err)
	}

	if Configuration.ThingspaceProtocol, err = config.GetString(ThingSpaceProtocolEnv, ""); err != nil {
		return fmt.Errorf("Error, %s environment variable not set.", ThingSpaceProtocolEnv)
	}

	if Configuration.ThingSpaceSouthHostPort, err = config.GetString(ThingSpaceSouthHostPortEnv, ""); err != nil {
		return fmt.Errorf("Error, %s environment variable not set.", ThingSpaceSouthHostPortEnv)
	}

	if Configuration.ThingSpaceAuthHostPort, err = config.GetString(ThingSpaceAuthHostPortEnv, ""); err != nil {
		return fmt.Errorf("Error, %s environment variable not set.", ThingSpaceAuthHostPortEnv)
	}

	if Configuration.ThingSpaceDataHostPort, err = config.GetString(ThingSpaceDataHostPortEnv, ""); err != nil {
		return fmt.Errorf("Error, %s environment variable not set.", ThingSpaceDataHostPortEnv)
	}

	if Configuration.ThingSpaceClientId, err = config.GetString(ThingSpaceClientIdEnv, ""); err != nil {
		return fmt.Errorf("Error, %s environment variable not set.", ThingSpaceClientIdEnv)
	}

	if Configuration.ThingSpaceClientSecret, err = config.GetString(ThingSpaceSecretEnv, ""); err != nil {
		return fmt.Errorf("Error, %s environment variable not set.", ThingSpaceSecretEnv)
	}

	kafkaBrokersHostPort, err := config.GetString(KafkaBrokersHostPortEnv, "")
	if kafkaBrokersHostPort == "" {
		return fmt.Errorf("Error, %s environment variable not set.", KafkaBrokersHostPortEnv)
	}

	Configuration.KafkaBrokersHostPort = strings.Split(kafkaBrokersHostPort, ",")

	if Configuration.ZookeeperHostPort, err = config.GetString(ZookeperHostPortEnd, ""); err != nil {
		return fmt.Errorf("Error, %s environment variable not set.", ZookeperHostPortEnd)
	}

	Configuration.MaxMemory, _ = config.GetUInt64(MaxMemoryEnv, 6000)
	Configuration.DefaultMemory, _ = config.GetUInt64(DefaultMemoryEnv, 1000)

	mlog.Info("Using Configuration: %+v", Configuration)
	return nil
}
