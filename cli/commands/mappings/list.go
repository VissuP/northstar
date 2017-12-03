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

package mappings

import (
	"flag"

	"errors"
	"fmt"
	"github.com/verizonlabs/northstar/cli/commands"
	"github.com/verizonlabs/northstar/cli/util"
	"github.com/verizonlabs/northstar/data/mappings/client"
)

type ListMappingsCmd struct {
	client *client.MappingsClient
	cmd    *flag.FlagSet
}

func NewListMappings(client *client.MappingsClient) commands.Command {
	cmd := flag.NewFlagSet("mappings-list", flag.ExitOnError)
	return &ListMappingsCmd{client: client,
		cmd: cmd}
}

func (list *ListMappingsCmd) Run(args []string) error {
	list.cmd.Parse(args)

	if !list.cmd.Parsed() {
		return errors.New("Failed to parse cmd")
	}

	result, err := list.client.ListMappings(util.GetAccountID())
	if err != nil {
		return err
	}

	if len(result) == 0 {
		fmt.Println("No mappings defined")
		return nil
	}

	for _, mapping := range result {
		fmt.Println(mapping.Print())
	}
	return nil
}
