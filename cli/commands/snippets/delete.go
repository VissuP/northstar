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

package snippets

import (
	"flag"

	"errors"
	"fmt"
	"github.com/verizonlabs/northstar/cli/commands"
	"github.com/verizonlabs/northstar/cli/util"
	"github.com/verizonlabs/northstar/data/snippets/client"
)

type DeleteSnippetCmd struct {
	client    *client.SnippetsClient
	cmd       *flag.FlagSet
	snippetId *string
}

func NewDeleteSnippet(client *client.SnippetsClient) commands.Command {
	cmd := flag.NewFlagSet("snippets-delete", flag.ExitOnError)
	snippetId := cmd.String("id", "", "The snippet id")

	return &DeleteSnippetCmd{client: client,
		cmd:       cmd,
		snippetId: snippetId}
}

func (d *DeleteSnippetCmd) Run(args []string) error {
	d.cmd.Parse(args)

	if !d.cmd.Parsed() {
		return errors.New("Failed to parse cmd")
	}

	if *d.snippetId == "" {
		return errors.New("Please set an name using -name.")
	}

	err := d.client.DeleteSnippet(util.GetAccountID(), *d.snippetId)
	if err != nil {
		return err
	}

	fmt.Println("Snippet deleted")
	return nil
}
