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

type GetDatasourcesCmd struct {
	client *client.DatasourcesClient
	cmd    *flag.FlagSet
	id     *string
	args   *string
}

func NewGetDatasource(client *client.DatasourcesClient) commands.Command {
	cmd := flag.NewFlagSet("datasources-get", flag.ExitOnError)
	id := cmd.String("id", "", "The datasource id")
	return &GetDatasourcesCmd{client: client, cmd: cmd, id: id}
}

func (invoke *GetDatasourcesCmd) Run(args []string) error {
	invoke.cmd.Parse(args)

	if !invoke.cmd.Parsed() {
		return errors.New("Failed to parse cmd")
	}

	if *invoke.id == "" {
		return errors.New("Please set an id using -id.")
	}

	out, err := invoke.client.GetDatasource(util.GetAccountID(), *invoke.id)
	if err != nil {
		return err
	}

	fmt.Println(out.Print())
	return nil
}
