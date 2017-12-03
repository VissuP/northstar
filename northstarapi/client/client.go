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

package client

import (
	"fmt"

	"github.com/gin-gonic/gin"
	lb "github.com/verizonlabs/northstar/pkg/httpclientlb"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/northstarapi/model"
)

// Client defines the type used to represent a service client.
type Client struct {
	baseUrl  string
	lbClient *lb.LbClient
}

// NewClient returns a new instance of client.
func NewClient(protocol, hostAndPort string) (*Client, error) {
	baseUrl := fmt.Sprintf("%s://%s", protocol, hostAndPort)

	lbClient, err := lb.GetClient(baseUrl)
	if err != nil {
		mlog.Error("Failed to create northstarapi client lb with error: %s", err.Error())
		return nil, fmt.Errorf("Failed to create northstarapi client lb with error: %s", err.Error())
	}

	return &Client{
		baseUrl:  baseUrl,
		lbClient: lbClient,
	}, nil
}

// getResourcePath is a helper method used to get resource path.
func (client Client) getResourcePath(resource string) string {
	return fmt.Sprintf("/%s/%s/%s", model.Context, model.Version, resource)
}

// getRequestHeaders is a helper method used to get request headers.
func (client Client) getRequestHeaders(accessToken string) map[string]string {
	header := map[string]string{}
	header["Authorization"] = fmt.Sprintf("Bearer %s", accessToken)
	header["Content-Type"] = gin.MIMEJSON
	return header
}
