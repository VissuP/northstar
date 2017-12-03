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

package datasets

import (
	"flag"

	"errors"
	"fmt"
	"github.com/verizonlabs/northstar/cli/commands"
	"github.com/verizonlabs/northstar/cli/util"
	"github.com/verizonlabs/northstar/data/datasets/client"
	"github.com/verizonlabs/northstar/data/datasets/model"
)

type AddDatasetCmd struct {
	client       *client.DatasetsClient
	cmd          *flag.FlagSet
	datasourceId *string
	name         *string
	description  *string
	tables       *string
}

func NewAddDataset(client *client.DatasetsClient) commands.Command {
	cmd := flag.NewFlagSet("datasets-add", flag.ExitOnError)
	datasourceId := cmd.String("datasourceId", "", "The datasource id")
	name := cmd.String("name", "", "The dataset name")
	description := cmd.String("description", "test description", "The dataset description")
	tables := cmd.String("tables", "{}", "The dataset tables")

	return &AddDatasetCmd{client: client,
		cmd:          cmd,
		datasourceId: datasourceId,
		name:         name,
		description:  description,
		tables:       tables}
}

func (output *AddDatasetCmd) Run(args []string) error {
	output.cmd.Parse(args)

	if !output.cmd.Parsed() {
		return errors.New("Failed to parse cmd")
	}

	if *output.name == "" {
		return errors.New("Please set a name using -name.")
	}

	tables, err := util.UnmarshalTables(*output.tables)
	if err != nil {
		return err
	}

	if len(tables) == 0 {
		return errors.New("Please set tables using -tables.")
	}

	data := &model.DatasetData{DatasourceId: *output.datasourceId,
		Name:        *output.name,
		Description: *output.description,
		Tables:      tables}
	id, mErr := output.client.AddDataset(util.GetAccountID(), data)
	if mErr != nil {
		return mErr
	}

	fmt.Printf("Dataset %s added", id)
	return nil
}
