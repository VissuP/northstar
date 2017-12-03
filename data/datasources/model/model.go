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

package model

import (
	"fmt"
	"time"
)

type DatasourceData struct {
	Id          string            `cql:"id"`
	AccountId   string            `cql:"accountid"`
	Name        string            `cql:"name"`
	Description string            `cql:"description"`
	Protocol    string            `cql:"protocol"`
	Host        string            `cql:"host"`
	Port        int               `cql:"port"`
	Options     map[string]string `cql:"options"`
	CreatedOn   time.Time         `cql:"createdon"`
	UpdatedOn   time.Time         `cql:"updatedon"`
}

func (d *DatasourceData) Validate() error {
	if d.Protocol == "" {
		return fmt.Errorf("Protocol is empty")
	}

	if d.Name == "" {
		return fmt.Errorf("Name is empty")
	}

	if d.Host == "" {
		return fmt.Errorf("Host is empty")
	}

	if d.Port < 1 {
		return fmt.Errorf("Port is less than one")
	}

	return nil
}

func (d *DatasourceData) Print() string {
	return fmt.Sprintf("ID: %s, "+
		"Name: %s, "+
		"Description: %s, "+
		"Protocol: %s, "+
		"Host: %s, "+
		"Port: %d, "+
		"Options: %v, "+
		"CreatedOn: %s, "+
		"UpdatedOn: %s",
		d.Id,
		d.Name,
		d.Description,
		d.Protocol,
		d.Host,
		d.Port,
		d.Options,
		d.CreatedOn,
		d.UpdatedOn)
}
