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
	"github.com/verizonlabs/northstar/cron/model"
)

type AddJobCMD struct {
	client      *client.CronClient
	cmd         *flag.FlagSet
	name        *string
	disabled    *bool
	snippetId   *string
	schedule    *string
	description *string
}

func NewAddCronJob(client *client.CronClient) commands.Command {
	cmd := flag.NewFlagSet("cron-add", flag.ExitOnError)
	name := cmd.String("name", "cron_test", "The cron name")
	disabled := cmd.Bool("disabled", false, "Enable/disable job")
	snippetId := cmd.String("snippetId", "", "The cron name")
	schedule := cmd.String("schedule", "0 * * * * *", "The cron schedule")
	description := cmd.String("description", "", "The cron description")

	return &AddJobCMD{client: client,
		cmd:         cmd,
		name:        name,
		disabled:    disabled,
		snippetId:   snippetId,
		schedule:    schedule,
		description: description}
}

func (add *AddJobCMD) Run(args []string) error {
	add.cmd.Parse(args)

	if !add.cmd.Parsed() {
		return errors.New("Failed to parse cmd")
	}

	if *add.snippetId == "" {
		return errors.New("Please set a snippet id using -snippetId.")
	}

	job := &model.Job{Name: *add.name,
		SnippetId:   *add.snippetId,
		Disabled:    *add.disabled,
		Schedule:    *add.schedule,
		Description: *add.description}

	id, mErr := add.client.AddJob(util.GetAccountID(), job)
	if mErr != nil {
		return mErr
	}

	fmt.Printf("Job %s created", id)
	return nil
}
