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

type DeleteMappingCmd struct {
	client *client.MappingsClient
	cmd    *flag.FlagSet
	id     *string
}

func NewDeleteMapping(client *client.MappingsClient) commands.Command {
	cmd := flag.NewFlagSet("mappings-delete", flag.ExitOnError)
	id := cmd.String("id", "", "The mapping id")

	return &DeleteMappingCmd{client: client,
		cmd: cmd,
		id:  id}
}

func (d *DeleteMappingCmd) Run(args []string) error {
	d.cmd.Parse(args)

	if !d.cmd.Parsed() {
		return errors.New("Failed to parse cmd")
	}

	if *d.id == "" {
		return errors.New("Please set a mapping id using -id.")
	}

	err := d.client.DeleteMapping(util.GetAccountID(), *d.id)
	if err != nil {
		return err
	}

	fmt.Println("Mapping deleted")
	return err
}
