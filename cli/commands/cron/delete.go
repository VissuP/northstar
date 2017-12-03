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

package cron

import (
	"flag"

	"errors"
	"fmt"
	"github.com/verizonlabs/northstar/cli/commands"
	"github.com/verizonlabs/northstar/cli/util"
	"github.com/verizonlabs/northstar/cron/client"
)

type DeleteJobCmd struct {
	client *client.CronClient
	cmd    *flag.FlagSet
	id     *string
}

func NewDeleteJob(client *client.CronClient) commands.Command {
	cmd := flag.NewFlagSet("job-delete", flag.ExitOnError)
	id := cmd.String("id", "", "The job id")

	return &DeleteJobCmd{client: client,
		cmd: cmd,
		id:  id}
}

func (delete *DeleteJobCmd) Run(args []string) error {
	delete.cmd.Parse(args)

	if !delete.cmd.Parsed() {
		return errors.New("Failed to parse cmd")
	}

	if *delete.id == "" {
		return errors.New("Please set a job id using -id.")
	}

	err := delete.client.DeleteJob(util.GetAccountID(), *delete.id)
	if err != nil {
		return err
	}

	fmt.Println("Job deleted")
	return nil
}
