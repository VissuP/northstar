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
	"errors"
	"flag"
	"fmt"
	"github.com/verizonlabs/northstar/cli/commands"
	"github.com/verizonlabs/northstar/cli/util"
	"github.com/verizonlabs/northstar/data/datasets/client"
)

type GetDatasetByNameCmd struct {
	client *client.DatasetsClient
	cmd    *flag.FlagSet
	name   *string
}

func NewGetDatasetByName(client *client.DatasetsClient) commands.Command {
	cmd := flag.NewFlagSet("datasets-get-by-name", flag.ExitOnError)
	name := cmd.String("name", "", "The dataset name")
	return &GetDatasetByNameCmd{client: client, cmd: cmd, name: name}
}

func (g *GetDatasetByNameCmd) Run(args []string) error {
	g.cmd.Parse(args)

	if !g.cmd.Parsed() {
		return errors.New("Failed to parse cmd")
	}

	if *g.name == "" {
		return errors.New("Please set an name using -name.")
	}

	out, err := g.client.GetDatasetByName(util.GetAccountID(), *g.name)
	if err != nil {
		return err
	}

	fmt.Println(out.Print())
	return nil
}
