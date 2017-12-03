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

type InvocationGetCmd struct {
	client       *client.InvocationClient
	cmd          *flag.FlagSet
	invocationId *string
}

func NewGetInvocation(client *client.InvocationClient) commands.Command {
	cmd := flag.NewFlagSet("invoke-get", flag.ExitOnError)
	invocationId := cmd.String("id", "", "The invocation id")

	return &InvocationGetCmd{client: client,
		cmd:          cmd,
		invocationId: invocationId}
}

func (output *InvocationGetCmd) Run(args []string) error {
	output.cmd.Parse(args)

	if !output.cmd.Parsed() {
		return errors.New("Failed to parse cmd")
	}

	if *output.invocationId == "" {
		return errors.New("Please set a invocation id using -id.")
	}

	result, err := output.client.GetInvocation(util.GetAccountID(), *output.invocationId)
	if err != nil {
		return err
	}

	fmt.Println(result.Print())
	return nil
}
