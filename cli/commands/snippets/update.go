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
	"errors"
	"flag"

	"code.cloudfoundry.org/bytefmt"
	"fmt"
	"github.com/verizonlabs/northstar/cli/commands"
	"github.com/verizonlabs/northstar/cli/util"
	"github.com/verizonlabs/northstar/data/snippets/client"
	"github.com/verizonlabs/northstar/data/snippets/model"
	"github.com/verizonlabs/northstar/pkg/b64"
)

type UpdateSnippetCmd struct {
	client      *client.SnippetsClient
	cmd         *flag.FlagSet
	id          *string
	name        *string
	runtime     *string
	mainfn      *string
	timeout     *int
	callback    *string
	mem         *string
	url         *string
	file        *string
	description *string
}

func NewUpdateSnippet(client *client.SnippetsClient) commands.Command {
	cmd := flag.NewFlagSet("snippets-update", flag.ExitOnError)
	id := cmd.String("id", "", "The snippet id")
	name := cmd.String("name", "", "The snippet name")
	runtime := cmd.String("runtime", "", "The snippet runtime")
	mainfn := cmd.String("mainfn", "", "The snippet main function")
	timeout := cmd.Int("timeout", 0, "The snippet timeout in milliseconds")
	callback := cmd.String("callback", "", "The snippet callback")
	mem := cmd.String("mem", "", "The amount of memory (e.g., 1M, 1G)")
	url := cmd.String("url", "", "The snippet URL (e.g., s3:///bucket/<your key>, base64:///")
	file := cmd.String("file", "", "The snippet file")
	description := cmd.String("description", "", "The snippet description")

	return &UpdateSnippetCmd{client: client,
		cmd:         cmd,
		id:          id,
		name:        name,
		runtime:     runtime,
		mainfn:      mainfn,
		timeout:     timeout,
		callback:    callback,
		mem:         mem,
		url:         url,
		file:        file,
		description: description,
	}
}

func (update *UpdateSnippetCmd) Run(args []string) error {
	update.cmd.Parse(args)

	if !update.cmd.Parsed() {
		return errors.New("Failed to parse cmd")
	}

	var err error
	var base string
	if *update.file != "" {
		base, err = b64.EncodeFileToString(*update.file)
		if err != nil {
			return err
		}
	}

	if *update.id == "" {
		return errors.New("Please set an id using -id.")
	}

	var memory uint64
	if *update.mem != "" {
		memory, err = bytefmt.ToBytes(*update.mem)
		if err != nil {
			return err
		}
	}

	req := &model.SnippetData{Name: *update.name,
		Runtime:     *update.runtime,
		MainFn:      *update.mainfn,
		Timeout:     *update.timeout,
		Memory:      memory,
		Callback:    *update.callback,
		URL:         *update.url,
		Code:        base,
		Description: *update.description}

	mErr := update.client.UpdateSnippet(util.GetAccountID(), *update.id, req)
	if mErr != nil {
		return mErr
	}

	fmt.Println("Snippet updated")
	return nil
}
