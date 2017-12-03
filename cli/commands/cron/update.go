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

type UpdateJobCmd struct {
	client      *client.CronClient
	cmd         *flag.FlagSet
	id          *string
	name        *string
	disabled    *bool
	snippetId   *string
	schedule    *string
	description *string
}

func NewUpdateJob(client *client.CronClient) commands.Command {
	cmd := flag.NewFlagSet("cron-update", flag.ExitOnError)
	id := cmd.String("id", "", "The job id")
	name := cmd.String("name", "", "The job name")
	disabled := cmd.Bool("disabled", false, "Enable/disable job")
	snippetId := cmd.String("snippetId", "", "The snippet id")
	schedule := cmd.String("schedule", "", "The schedule")
	description := cmd.String("description", "", "The description")

	return &UpdateJobCmd{client: client,
		cmd:         cmd,
		id:          id,
		name:        name,
		disabled:    disabled,
		snippetId:   snippetId,
		schedule:    schedule,
		description: description}
}

func (update *UpdateJobCmd) Run(args []string) error {
	update.cmd.Parse(args)

	if !update.cmd.Parsed() {
		return errors.New("Failed to parse cmd")
	}

	if *update.id == "" {
		return errors.New("Please set an id using -id.")
	}

	job := &model.Job{Name: *update.name,
		SnippetId:   *update.snippetId,
		Disabled:    *update.disabled,
		Schedule:    *update.schedule,
		Description: *update.description}

	mErr := update.client.UpdateJob(util.GetAccountID(), *update.id, job)
	if mErr != nil {
		return mErr
	}

	fmt.Printf("Job updated")
	return nil
}
