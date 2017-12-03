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
	"github.com/verizonlabs/northstar/processing/snippets/client"
	"github.com/verizonlabs/northstar/processing/snippets/model"
)

type InvokeSnippetByIdCmd struct {
	client *client.SnippetsClient
	cmd    *flag.FlagSet
	id     *string
	args   *string
}

func NewInvokeSnippetById(client *client.SnippetsClient) commands.Command {
	cmd := flag.NewFlagSet("snippets-invoke-by-id", flag.ExitOnError)
	id := cmd.String("id", "", "The snippet id")
	args := cmd.String("args", "{}", "The snippet arguments")

	return &InvokeSnippetByIdCmd{client: client,
		cmd:  cmd,
		id:   id,
		args: args}
}

func (invoke *InvokeSnippetByIdCmd) Run(args []string) error {
	invoke.cmd.Parse(args)

	if !invoke.cmd.Parsed() {
		return errors.New("Failed to parse cmd")
	}

	arguments, err := util.UnmarshalInterface(*invoke.args)
	if err != nil {
		return err
	}

	if *invoke.id == "" {
		return errors.New("Please set an id using -id.")
	}

	req := &model.Snippet{SnippetId: *invoke.id,
		Options: model.Options{Args: arguments}}

	out, mErr := invoke.client.StartSnippet(util.GetAccountID(), req)
	if mErr != nil {
		return mErr
	}

	fmt.Println(out)
	return nil
}
