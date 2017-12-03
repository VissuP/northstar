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
	"github.com/verizonlabs/northstar/data/mappings/model"
)

type AddMappingCmd struct {
	client    *client.MappingsClient
	cmd       *flag.FlagSet
	data      *model.MappingsData
	id        *string
	eventId   *string
	snippetId *string
}

func NewAddMapping(client *client.MappingsClient) commands.Command {
	cmd := flag.NewFlagSet("mappings-add", flag.ExitOnError)
	id := cmd.String("id", "", "The mapping id")
	eventId := cmd.String("eventId", "", "The event id")
	snippetId := cmd.String("snippetId", "", "The snippet id")

	return &AddMappingCmd{client: client,
		cmd:       cmd,
		id:        id,
		eventId:   eventId,
		snippetId: snippetId}
}

func (output *AddMappingCmd) Run(args []string) error {
	output.cmd.Parse(args)

	if !output.cmd.Parsed() {
		return errors.New("Failed to parse cmd")
	}

	if *output.eventId == "" {
		return errors.New("Please set a event id using -eventId.")
	}

	if *output.snippetId == "" {
		return errors.New("Please set a snippet id using -snippetId.")
	}

	data := &model.MappingsData{Id: *output.id,
		EventId:   *output.eventId,
		SnippetId: *output.snippetId}

	id, err := output.client.AddMapping(util.GetAccountID(), data)
	if err != nil {
		return err
	}

	fmt.Printf("Mapping %s created", id)
	return nil
}
