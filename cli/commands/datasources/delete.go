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

package datasources

import (
	"flag"

	"errors"
	"fmt"
	"github.com/verizonlabs/northstar/cli/commands"
	"github.com/verizonlabs/northstar/cli/util"
	"github.com/verizonlabs/northstar/data/datasources/client"
)

type DeleteDatasourceCmd struct {
	client *client.DatasourcesClient
	cmd    *flag.FlagSet
	id     *string
}

func NewDeleteDatasource(client *client.DatasourcesClient) commands.Command {
	cmd := flag.NewFlagSet("datasources-delete", flag.ExitOnError)
	id := cmd.String("id", "", "The datasource id")
	return &DeleteDatasourceCmd{client: client, cmd: cmd, id: id}
}

func (d *DeleteDatasourceCmd) Run(args []string) error {
	d.cmd.Parse(args)

	if !d.cmd.Parsed() {
		return errors.New("Failed to parse cmd")
	}

	if *d.id == "" {
		return errors.New("Please set an id using -id.")
	}

	err := d.client.DeleteDatasource(util.GetAccountID(), *d.id)
	if err != nil {
		return err
	}

	fmt.Println("Datasource deleted")
	return nil
}
