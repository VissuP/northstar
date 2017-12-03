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
	"fmt"
	"net"
	"os"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"strconv"
)

const (
	//ServiceName is the name of our service
	ServiceName = "Portal"

	//DefaultServicePort is the port for the REST endpoint to listen om
	DefaultServicePort = "8080"

	//DefaultAuthScopes are the privileges requested from the Thingspace Auth service
	DefaultAuthScopes = "ts.user.ro ts.transformation ts.runtime ts.object"
)

const (
	//AdvertisedPortEnv is the environment variable used to change the auth port
	AdvertisedPortEnv = "ADVERTISED_PORT"

	//NorthstarAPIProtocolEnv is the environment variable used to set the protocol (HTTP/HTTPS) used to talk to the portal API.
	NorthstarAPIProtocolEnv = "NORTHSTARAPI_PROTOCOL"

	//NorthstarAPIHostPortEnv is environment variable used to set the the URI used to access the NorthstarAPI service.
	NorthstarAPIHostPortEnv = "NORTHSTARAPI_HOST_PORT"

	//ThingSpaceProtocolEnv is the environment variable used to set the protocol (HTTP/HTTPS) used to talk to Thingspace.
	ThingSpaceProtocolEnv = "THINGSPACE_PROTOCOL"

	//ThingSpaceUserHostPortEnv is the environment variable used to set the URI used to access the Thingspace User service.
	ThingSpaceUserHostPortEnv = "THINGSPACE_HOST_PORT"

	//ThingSpaceAuthHostPortEnv is the environment variable used to set the URI used to access the Thingspace Auth service.
	ThingSpaceAuthHostPortEnv = "THINGSPACE_AUTH_HOST_PORT"

	//ThingSpaceClientIDEnv is the environment variable used to set the username used for the service to access Thingspace.
	ThingSpaceClientIDEnv = "THINGSPACE_CLIENT_ID"

	//ThingSpaceSecretEnv is the environment variable used to set the password used for the service to access Thingspace.
	ThingSpaceSecretEnv = "THINGSPACE_SECRET"

	//ThingSpaceUserScopes is the environment variable used to set the privileges requested for a user from the Thingspace service.
	ThingSpaceUserScopes = "THINGSPACE_USER_SCOPES"

	//ThingSpaceClientScopes is the environment variable used to set the privileges requested for the application from the Thingspace service.
	ThingSpaceClientScopes = "THINGSPACE_CLIENT_SCOPES"

	//ConnectionBufferSizeEnv is the environment variable used to set the buffer size used for the Portal<->Angular websockets.
	ConnectionBufferSizeEnv = "CONNECTION_BUFFER_SIZE"
)

var (
	//Configuration defines the service configuration.
	Configuration = new(configuration)
)

type configuration struct {
	ServiceName            string
	ServiceHostPort        string
	Port                   string
	NorthstarAPIProtocol   string
	NorthstarAPIHostPort   string
	ThingspaceProtocol     string
	ThingSpaceUserHostPort string
	ThingSpaceAuthHostPort string
	ThingSpaceClientID     string
	ThingSpaceClientSecret string
	ThingSpaceUserScopes   string
	ThingSpaceClientScopes string
	ConnectionBufferSize   int
}

// Load loads the configuration from the environment variables.
func Load() (err error) {
	mlog.Debug("Load configuration variables")

	Configuration.ServiceName = ServiceName

	// Get host and port assignment from the environment variable.
	if Configuration.Port = os.Getenv(AdvertisedPortEnv); Configuration.Port == "" {
		mlog.Info("Warning, %s environment variable not set. Using default service port %s.", AdvertisedPortEnv, DefaultServicePort)
		Configuration.Port = DefaultServicePort
	}

	// Get the interface IP address to generate service host and port.
	if interfaceIP := getInterfaceIP(); interfaceIP != "" {
		Configuration.ServiceHostPort = fmt.Sprintf("%s:%s", interfaceIP, Configuration.Port)
	}

	// Get the Portal API (Service Tier) Protocol.
	if Configuration.NorthstarAPIProtocol = os.Getenv(NorthstarAPIProtocolEnv); Configuration.NorthstarAPIProtocol == "" {
		mlog.Error("Error, %s environment variable not set.", NorthstarAPIProtocolEnv)
		return fmt.Errorf("%s environment variable not set", NorthstarAPIProtocolEnv)
	}

	// Get the Portal API (Service Tier) Host and Port.
	if Configuration.NorthstarAPIHostPort = os.Getenv(NorthstarAPIHostPortEnv); Configuration.NorthstarAPIHostPort == "" {
		mlog.Error("Error, %s environment variable not set.", NorthstarAPIHostPortEnv)
		return fmt.Errorf("%s environment variable not set", NorthstarAPIHostPortEnv)
	}

	// Get environment variables needed for ThingSpace communication.
	if Configuration.ThingspaceProtocol = os.Getenv(ThingSpaceProtocolEnv); Configuration.ThingspaceProtocol == "" {
		return fmt.Errorf("Error, %s environment variable not set.", ThingSpaceProtocolEnv)
	}

	if Configuration.ThingSpaceUserHostPort = os.Getenv(ThingSpaceUserHostPortEnv); Configuration.ThingSpaceUserHostPort == "" {
		return fmt.Errorf("Error, %s environment variable not set.", ThingSpaceUserHostPortEnv)
	}

	if Configuration.ThingSpaceAuthHostPort = os.Getenv(ThingSpaceAuthHostPortEnv); Configuration.ThingSpaceAuthHostPort == "" {
		return fmt.Errorf("Error, %s environment variable not set.", ThingSpaceAuthHostPortEnv)
	}

	if Configuration.ThingSpaceClientID = os.Getenv(ThingSpaceClientIDEnv); Configuration.ThingSpaceClientID == "" {
		return fmt.Errorf("Error, %s environment variable not set.", ThingSpaceClientIDEnv)
	}

	if Configuration.ThingSpaceClientSecret = os.Getenv(ThingSpaceSecretEnv); Configuration.ThingSpaceClientSecret == "" {
		return fmt.Errorf("Error, %s environment variable not set.", ThingSpaceSecretEnv)
	}

	if Configuration.ThingSpaceUserScopes = os.Getenv(ThingSpaceUserScopes); Configuration.ThingSpaceUserScopes == "" {
		return fmt.Errorf("Error, %s environment variable not set.", ThingSpaceUserScopes)
	}

	if Configuration.ThingSpaceClientScopes = os.Getenv(ThingSpaceClientScopes); Configuration.ThingSpaceClientScopes == "" {
		return fmt.Errorf("Error, %s environment variable not set.", ThingSpaceClientScopes)
	}

	ConnectionBufferSize := os.Getenv(ConnectionBufferSizeEnv)
	if ConnectionBufferSize == "" {
		return fmt.Errorf("Error, %s environment variable not set.", ConnectionBufferSizeEnv)
	}

	bufferSize, err := strconv.ParseInt(ConnectionBufferSize, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid integer value for environment variable %s. Error was: %s", ConnectionBufferSizeEnv, err.Error())
	}

	Configuration.ConnectionBufferSize = int(bufferSize)

	mlog.Debug("Loaded Service Configuration: %v", Configuration)

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
