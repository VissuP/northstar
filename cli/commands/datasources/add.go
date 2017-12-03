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
	"github.com/verizonlabs/northstar/data/datasources/model"
)

type AddDatasourceCmd struct {
	client       *client.DatasourcesClient
	cmd          *flag.FlagSet
	datasourceId *string
	name         *string
	description  *string
	protocol     *string
	host         *string
	port         *int
	options      *string
}

func NewAddDatasource(client *client.DatasourcesClient) commands.Command {
	cmd := flag.NewFlagSet("datasources-add", flag.ExitOnError)
	name := cmd.String("name", "", "The datasource name")
	description := cmd.String("description", "", "The datasource description")
	protocol := cmd.String("protocol", "", "The datasource protocol")
	host := cmd.String("host", "", "The datasource host")
	port := cmd.Int("port", 0, "The datasource port")
	options := cmd.String("options", "{}", "The datasource options")

	return &AddDatasourceCmd{client: client,
		cmd:         cmd,
		name:        name,
		description: description,
		protocol:    protocol,
		host:        host,
		port:        port,
		options:     options}
}

func (output *AddDatasourceCmd) Run(args []string) error {
	output.cmd.Parse(args)

	if !output.cmd.Parsed() {
		return errors.New("Failed to parse cmd")
	}

	if *output.name == "" {
		return errors.New("Please set a name using -name.")
	}

	if *output.protocol == "" {
		return errors.New("Please set a protocol using -protocol.")
	}

	if *output.host == "" {
		return errors.New("Please set a host using -host.")
	}

	if *output.port < 1 {
		return errors.New("Please set a port using -port.")
	}

	options, err := util.UnmarshalString(*output.options)
	if err != nil {
		return err
	}

	data := &model.DatasourceData{Name: *output.name,
		Description: *output.description,
		Protocol:    *output.protocol,
		Host:        *output.host,
		Port:        *output.port,
		Options:     options}
	id, mErr := output.client.AddDatasource(util.GetAccountID(), data)
	if mErr != nil {
		return mErr
	}

	fmt.Printf("Datasource %s added", id)
	return nil
}
