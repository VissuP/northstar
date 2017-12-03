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

type Column struct {
	Name     string `cql:"name"`
	DataType string `cql:"datatype"`
}

type Table struct {
	Name    string            `cql:"name"`
	Columns map[string]Column `cql:"columns"`
}

type DatasetData struct {
	Id           string           `cql:"id"`
	AccountId    string           `cql:"accountid"`
	DatasourceId string           `cql:"datasourceid"`
	Name         string           `cql:"name"`
	Description  string           `cql:"description"`
	Tables       map[string]Table `cql:"tables"`
	CreatedOn    time.Time        `cql:"createdon"`
	UpdatedOn    time.Time        `cql:"updatedon"`
}

func (d *DatasetData) Validate() error {
	if d.Name == "" {
		return fmt.Errorf("Name is empty")
	}

	if len(d.Tables) == 0 {
		return fmt.Errorf("Tables is empty")
	}

	return nil
}

func (d *DatasetData) Print() string {
	return fmt.Sprintf("ID: %s, "+
		"DatasourceID: %s, "+
		"Name: %s, "+
		"Description: %s, "+
		"Tables: %v, "+
		"CreatedOn: %s, "+
		"UpdatedOn: %s",
		d.Id,
		d.DatasourceId,
		d.Name,
		d.Description,
		d.Tables,
		d.CreatedOn,
		d.UpdatedOn)
}
