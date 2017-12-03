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
	"fmt"
	"github.com/verizonlabs/northstar/cli/commands"
	"github.com/verizonlabs/northstar/cli/util"
	"github.com/verizonlabs/northstar/data/invocations/client"
)

type DeleteInvocationCmd struct {
	client       *client.InvocationClient
	cmd          *flag.FlagSet
	invocationId *string
}

func NewDeleteInvocation(client *client.InvocationClient) commands.Command {
	cmd := flag.NewFlagSet("invoke-delete", flag.ExitOnError)
	invocationId := cmd.String("id", "", "The invocation id")

	return &DeleteInvocationCmd{client: client,
		cmd:          cmd,
		invocationId: invocationId}
}

func (d *DeleteInvocationCmd) Run(args []string) error {
	d.cmd.Parse(args)

	if !d.cmd.Parsed() {
		return errors.New("Failed to parse cmd")
	}

	if *d.invocationId == "" {
		return errors.New("Please set an id using -id.")
	}

	err := d.client.DeleteInvocation(util.GetAccountID(), *d.invocationId)
	if err != nil {
		return err
	}

	fmt.Println("Invocation deleted")
	return nil
}
