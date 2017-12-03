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

package northstar

import (
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/northstarapi/client"
	"github.com/verizonlabs/northstar/portal/config"
	"github.com/verizonlabs/northstar/portal/utils/thingspace"
)

// NorthStarPortalProvider defines the type that implements the PortalProvider.
type NorthStarPortalProvider struct {
	northstarApiClient *client.Client
	thingSpaceClient   *thingspace.UserClient
}

// NewNorthStarPortalProvider returns a new NorthStar Portal Provider.
func NewNorthStarPortalProvider(protocol string, hostAndPort string) (*NorthStarPortalProvider, error) {
	mlog.Debug("NewNorthStarPortalProvider")

	northstarApiClient, err := client.NewClient(protocol, hostAndPort)
	if err != nil {
		return nil, err
	}

	provider := &NorthStarPortalProvider{
		northstarApiClient: northstarApiClient,
		thingSpaceClient:   thingspace.NewUserClient(config.Configuration.ThingspaceProtocol, config.Configuration.ThingSpaceUserHostPort),
	}

	return provider, nil
}
