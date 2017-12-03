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
	"net"
	"os"
	"strings"
	"time"

	"github.com/verizonlabs/northstar/pkg/config"
	"github.com/verizonlabs/northstar/pkg/mlog"
)

const (
	//ServiceName is the name of our service
	ServiceName = "NSSim"

	// Defines the default advertised host name.
	DEFAULT_ADVERTISED_HOST            = "0.0.0.0"
	DEFAULT_ENVIRONMENT                = "Unknown"
	DEFAULT_CASSANDRA_SLEEP            = 20 * time.Second
	DEFAULT_TIMEOUT                    = 5
	DEFAULT_EXECUTION_RESPONSE_TIMEOUT = 60
	DEFAULT_REPEAT_EXECUTION_SLEEP     = 2 * time.Second
)

const (
	//AdvertisedPortEnv is the environment variable used to change the auth port
	AdvertisedPortEnv     = "ADVERTISED_PORT"
	AdvertisedHostPortEnv = "ADVERTISED_HOSTNAME"
	IterationDelayEnv     = "ITERATION_DELAY_IN_SEC"

	// Environment is the environment that the service is running under
	EnvironmentEnv              = "CONFIGURED_ENVIRONMENT"
	ConfigFileEnv               = "CONFIG_FILE"
	HttpTimeoutEnv              = "HTTP_TIMEOUT"
	TsUserEnv                   = "THINGSPACE_USER_HOST_PORT"
	TsSouthEnv                  = "THINGSPACE_SOUTH_HOST_PORT"
	ThingspaceAuthEnv           = "THINGSPACE_AUTH_HOST_PORT"
	ThingspaceProtocolEnv       = "THINGSPACE_PROTOCOL"
	SimModeEnv                  = "SIMULATOR_MODE"
	EnabledGroupsEnv            = "SIMULATOR_GROUPS"
	CassandraHostEnv            = "CASSANDRA_HOST"
	CassandraPortEnv            = "CASSANDRA_NATIVE_TRANSPORT_PORT"
	ExecutionResponseTimeoutEnv = "EXECUTION_RESPONSE_TIMEOUT"

	NorthstarProtocolEnv = "NORTHSTAR_PROTOCOL"
	NorthstarApiHostEnv  = "NORTHSTARAPI_HOST_PORT"

	DakotaProtocolEnv = "DAKOTA_PROTOCOL"
	UtProvisionEnv    = "UTPROVISION_HOST_PORT"
)

const (
	// Define the different test types.
	LOAD_TEST       string = "Load"
	ERROR_TEST      string = "Error"
	FUNCTIONAL_TEST string = "Functional"
)
const (
	//Define the different supported test names
	Notebook                      string = "Notebook CRUD Operations"
	Execution                     string = "Generic execution operations"
	TableExecution                string = "Notebook Table Execution"
	ValueExecution                string = "Notebook Value Execution"
	MapExecution                  string = "Notebook Map Execution"
	TemplateCell                  string = "Cell Template CRUD Operations"
	TemplateNotebook              string = "Notebook Template CRUD Operations"
	SingleDeviceMultiBarExecution string = "Template Single Device Multi Bar Execution"
	Transformation                string = "Transformation CRUD Operations"
	TimerExecution                string = "Transformation Timer Execution"
	ExternalEventExecution        string = "Transformation External Event Execution"
	ObjectRead                    string = "NS Object Operations"
	NSQLNativeTable               string = "NSQL Native Table Operations"
	NSQLNativeCrud                string = "NSQL Native CRUD Operations"
	NSQLNativeTypedCrud           string = "NSQL Native Typed CRUD Operations"
	NSQLNativeJSONFetch           string = "NSQL Native JSON Fetch"
	NSQLNativeMapBlob             string = "NSQL Native Map Blob"
	NSQLSparkJSONFetch            string = "NSQL Spark JSON Fetch"
	NSQLSparkMabBlob              string = "NSQL Spark Map Blob"
	NSQLSparkRead                 string = "NSQL Spark Read Operations"
	NSQLSparkTypedRead            string = "NSQL Spark Typed Read Operations"
	NSKVCore                      string = "NS Key-Value Operations"
)

const (
	// Define supported modes.
	AUTORUN_MODE string = "auto"
	MANUAL_MODE  string = "manual"
)

var (
	// Defines the service configuration.
	Configuration = new(configuration)
)

// Defines the type used to represent a test to be executed.
type Test struct {
	Type        string `json:"type" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Group       string `json:"group" binding:"required"`
	Concurrency int    `json:"concurrency" binding:"required"`
	Data        string `json:"data" binding:"optional"`
	Verbose     bool   `json:"verbose" binding:"required"`
}

type Credentials struct {
	Client Client `json:"client" binding:"required"`
	User   *User  `json:"user" binding:"optional"`
}
type Client struct {
	Id     string `json:"id" binding:"required"`
	Secret string `json:"secret" binding:"required"`
	Scope  string `json:"scope" binding:"required"`
}

type User struct {
	Id              string `json:"id" binding:"required"`
	Secret          string `json:"secret" binding:"required"`
	Scope           string `json:"scope" binding:"required"`
	AccountPassword string `json:"password" binding:"required"`
}

type configuration struct {
	ServiceName     string
	Host            string
	Port            string
	Environment     string
	ServiceHostPort string

	// From environment variables
	HttpTimeout         int    `json:"-"`
	DakotaProtocol      string `json:"-"`
	UtProvisionHostPort string `json:"-"`

	NorthstarProtocol    string `json:"-"`
	NorthstarApiHostPort string `json:"-"`

	ThingspaceProtocol      string   `json:"-"`
	ThingspaceUserHostPort  string   `json:"-"`
	ThingspaceSouthHostPort string   `json:"-"`
	ThingspaceAuthHostPort  string   `json:"-"`
	Mode                    string   `json:"-"`
	EnabledGroups           []string `json:"-"`
	CassandraPort           string   `json:"-"`
	CassandraHost           string   `json:"-"`

	// From configuration file
	Credentials Credentials `json:"credentials" binding:"required"`
	Tests       []Test      `json:"tests" binding:"required"`
	Devices     []Device    `json:"devices"`

	IterationDelayInSec      int           `json:"iterationDelayInSec" binding:"required"`
	MaxExecutionTime         int           `json:"maxExecutionTime" binding:"required"`
	ExecutionResponseTimeout time.Duration `json:"-"`
	RepeatExecutionSleep     time.Duration `json:"-"`
}

// Defines the type used to represent a test device.
type Device struct {
	Kind       string   `json:"kind" binding:"required"`
	ModelId    string   `json:"modelId" binding:"required"`
	ProviderId string   `json:"providerId" binding:"required"`
	Fields     []string `json:"fields" binding:"required"`
}

//getEnableGroups gets the enabled test groups
func getEnabledGroups() []string {
	groups := make([]string, 0)
	strGroups := strings.ToLower(os.Getenv(EnabledGroupsEnv))
	for _, s := range strings.Split(strGroups, ",") {
		groups = append(groups, strings.TrimSpace(s))
	}

	return groups
}

// Is given test enabled
func (config *configuration) IsTestEnabled(test *Test) bool {
	if test != nil {
		for _, group := range config.EnabledGroups {
			if group == strings.ToLower(test.Group) {
				return true
			}
		}
	}
	return false

}

// Returns device for the specified kind from the configuraiton.
func (config *configuration) FindDeviceByKind(kind string) (*Device, error) {
	for _, device := range config.Devices {
		if device.Kind == kind {
			return &device, nil
		}
	}
	return nil, fmt.Errorf("No device found with kind %s", kind)
}

// Returns device fields for the specified kind from the configuraiton.
func (config *configuration) FindDeviceFieldsByKind(kind string) ([]string, error) {
	device, err := config.FindDeviceByKind(kind)
	if err != nil {
		return nil, err
	}
	return device.Fields, nil
}

// Loads configuration from the environment variables.
func Load() (err error) {
	mlog.Debug("Load configuration variables")

	configFile := os.Getenv(ConfigFileEnv)
	if configFile == "" {
		return fmt.Errorf("%s environment variable not valid.", ConfigFileEnv)
	}

	mlog.Info("Loading configuration file %s.", configFile)
	var file *os.File

	// Open the configuration file
	if file, err = os.Open(configFile); err != nil {
		return err
	}

	// Parse the configuration file
	mlog.Info("Parsing configuration file %s.", configFile)
	parser := json.NewDecoder(file)

	if err = parser.Decode(Configuration); err != nil {
		mlog.Info("Failed to parse config with errors: %+v", err)
		return err
	}

	mlog.Info("Configuration Tests = %+v", Configuration.Tests)

	Configuration.ServiceName = ServiceName

	// Get host and port assignment from the environment variable.
	if Configuration.Port, err = config.GetString(AdvertisedPortEnv, ""); err != nil {
		return fmt.Errorf("Error, expected %s environment variable not found or invalid.", AdvertisedPortEnv)
	}

	if Configuration.Host, err = config.GetString(AdvertisedHostPortEnv, DEFAULT_ADVERTISED_HOST); err != nil {
		mlog.Info("Warning, %s environment variable not found. Using default %s", AdvertisedHostPortEnv, DEFAULT_ADVERTISED_HOST)
	}

	if Configuration.Environment, err = config.GetString(EnvironmentEnv, DEFAULT_ENVIRONMENT); err != nil {
		mlog.Info("Warning, %s environment variable not found. Using default: %s", EnvironmentEnv, DEFAULT_ENVIRONMENT)
	}

	//Use environment as our username for email so that we can easily differentiate between datacenters.
	if Configuration.Environment != "" {
		Configuration.Credentials.User.Id = Configuration.Environment + "nssim"
	}
	if Configuration.ThingspaceProtocol, err = config.GetString(ThingspaceProtocolEnv, ""); err != nil {
		return fmt.Errorf("Error, expected %s environment variable not found or invalid.", ThingspaceProtocolEnv)
	}
	Configuration.ThingspaceProtocol = strings.ToLower(Configuration.ThingspaceProtocol)

	if Configuration.ThingspaceUserHostPort, err = config.GetString(TsUserEnv, ""); err != nil {
		return fmt.Errorf("Error, expected %s environment variable not found or invalid.", TsUserEnv)
	}

	if Configuration.ThingspaceSouthHostPort, err = config.GetString(TsSouthEnv, ""); err != nil {
		return fmt.Errorf("Error, expected %s environment variable not found or invalid.", TsUserEnv)
	}

	if Configuration.ThingspaceAuthHostPort, err = config.GetString(ThingspaceAuthEnv, ""); err != nil {
		return fmt.Errorf("Error, expected %s environment variable not found or invalid.", ThingspaceAuthEnv)
	}

	if Configuration.NorthstarProtocol, err = config.GetString(NorthstarProtocolEnv, ""); err != nil {
		return fmt.Errorf("Error, expected %s environment variable not found or invalid.", NorthstarProtocolEnv)
	}

	Configuration.NorthstarProtocol = strings.ToLower(Configuration.NorthstarProtocol)

	if Configuration.NorthstarApiHostPort, err = config.GetString(NorthstarApiHostEnv, ""); err != nil {
		return fmt.Errorf("Error, expected %s environment variable not found or invalid.", NorthstarApiHostEnv)
	}

	if Configuration.DakotaProtocol, err = config.GetString(DakotaProtocolEnv, ""); err != nil {
		mlog.Info("WARNING: %s environment variable not set.", DakotaProtocolEnv)
	}

	if Configuration.UtProvisionHostPort, err = config.GetString(UtProvisionEnv, ""); err != nil {
		mlog.Info("WARNING: %s environment variable not set.", UtProvisionEnv)
	}

	if Configuration.CassandraHost, err = config.GetString(CassandraHostEnv, ""); err != nil {
		mlog.Info("WARNING: %s environment variable not set.", CassandraHostEnv)
	}

	if Configuration.CassandraPort, err = config.GetString(CassandraPortEnv, ""); err != nil {
		mlog.Info("WARNING: %s environment variable not set.", CassandraPortEnv)
	}

	if Configuration.Mode = strings.ToLower(os.Getenv(SimModeEnv)); Configuration.Mode == "" {
		return fmt.Errorf("Error, expected %s environment variable not found.", SimModeEnv)
	}

	if Configuration.Mode != AUTORUN_MODE && Configuration.Mode != MANUAL_MODE {
		return fmt.Errorf("Error, expected %s environment variable invalid.", SimModeEnv)
	}

	if Configuration.HttpTimeout, err = config.GetInt(HttpTimeoutEnv, DEFAULT_TIMEOUT); err != nil {
		mlog.Info("Warning: %s environment variable not set.", HttpTimeoutEnv)
	}

	//Override the iteration delay if it's specified in marathon.json
	Configuration.IterationDelayInSec, _ = config.GetInt(IterationDelayEnv, Configuration.IterationDelayInSec)

	Configuration.RepeatExecutionSleep = DEFAULT_REPEAT_EXECUTION_SLEEP

	ExecutionResponseTimeout, _ := config.GetInt(ExecutionResponseTimeoutEnv, DEFAULT_EXECUTION_RESPONSE_TIMEOUT)
	Configuration.ExecutionResponseTimeout = time.Duration(ExecutionResponseTimeout) * time.Second

	Configuration.EnabledGroups = getEnabledGroups()
	mlog.Info("Enabled groups: %s", Configuration.EnabledGroups)
	if len(Configuration.EnabledGroups) == 0 {
		return fmt.Errorf("No test groups enabled")
	}

	// Get the interface IP address to generate service host and port.
	if interfaceIP := getInterfaceIP(); interfaceIP != "" {
		Configuration.ServiceHostPort = fmt.Sprintf("%s:%s", interfaceIP, Configuration.Port)
	}

	mlog.Debug("Loaded Service Configuration: %+v", Configuration)

	return nil
}

// getInterfaceIp is a helper method used to get interface ip address.
func getInterfaceIP() string {
	mlog.Debug("getInterfaceIp")

	// Determine the service IP address from the eth interface.
	interfaces, err := net.Interfaces()

	if err != nil {
		mlog.Error("Net interfaces returned error: %+v", err)
		return ""
	}

	for _, ethInterface := range interfaces {
		mlog.Debug("Interface: %+v", ethInterface)

		// Check if interface is the loop back.
		if ethInterface.Flags&net.FlagLoopback == 0 {

			// Get the address.
			if addresses, err := ethInterface.Addrs(); err == nil {
				for _, address := range addresses {
					// Our IPv4 address will have a '.' in it.
					if ip, _, err := net.ParseCIDR(address.String()); err == nil {
						if ip4 := ip.To4(); ip4 != nil {
							return ip4.String()
						}
					}
				}
			}
		}
	}

	return ""
}
