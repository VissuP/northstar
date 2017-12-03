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

	lb "github.com/verizonlabs/northstar/pkg/httpclientlb"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/data/templates/model"
	"github.com/verizonlabs/northstar/data/util"
)

type Client interface {
	CreateTemplate(template *model.Template) (*model.Template, *management.Error)
	UpdateTemplate(templateId string, template *model.Template) (*model.Template, *management.Error)
	GetTemplate(templateId string) (*model.Template, *management.Error)
	QueryTemplate(query *model.Query) ([]model.Template, *management.Error)
	DeleteTemplate(templateId string) *management.Error
}

type TemplatesClient struct {
	lbClient *lb.LbClient
}

// Returns a new template client.
func NewTemplatesClient() (*TemplatesClient, error) {
	url, err := util.GetDataBaseUrl()

	if err != nil {
		return nil, fmt.Errorf("Get data base url returned error: %s", err.Error())
	}

	lbClient, err := lb.GetClient(url)

	if err != nil {
		return nil, fmt.Errorf("Get client returned error: %s", err.Error())
	}

	client := &TemplatesClient{
		lbClient: lbClient,
	}

	return client, nil
}
