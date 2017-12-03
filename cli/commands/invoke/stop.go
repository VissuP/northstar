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

package invoke

import (
	"flag"

	"errors"
	"github.com/verizonlabs/northstar/cli/commands"
	"github.com/verizonlabs/northstar/cli/util"
	"github.com/verizonlabs/northstar/processing/snippets/client"
)

type StopInvocationCmd struct {
	client       *client.SnippetsClient
	cmd          *flag.FlagSet
	invocationId *string
}

func NewStopInvocation(client *client.SnippetsClient) commands.Command {
	cmd := flag.NewFlagSet("invoke-stop", flag.ExitOnError)
	invocationId := cmd.String("id", "", "The invocation id")

	return &StopInvocationCmd{client: client,
		cmd:          cmd,
		invocationId: invocationId}
}

func (d *StopInvocationCmd) Run(args []string) error {
	d.cmd.Parse(args)

	if !d.cmd.Parsed() {
		return errors.New("Failed to parse cmd")
	}

	if *d.invocationId == "" {
		return errors.New("Please set an invocation id using -id.")
	}

	err := d.client.StopSnippet(util.GetAccountID(), *d.invocationId)
	if err != nil {
		return err
	}

	return nil
}
