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

type AddSnippetCmd struct {
	client      *client.SnippetsClient
	cmd         *flag.FlagSet
	id          *string
	name        *string
	runtime     *string
	mainfn      *string
	timeout     *int
	mem         *string
	callback    *string
	url         *string
	file        *string
	description *string
}

func NewAddSnippet(client *client.SnippetsClient) commands.Command {
	cmd := flag.NewFlagSet("snippets-add", flag.ExitOnError)
	id := cmd.String("id", "", "The snippet id")
	name := cmd.String("name", "", "The snippet name")
	runtime := cmd.String("runtime", "", "The snippet runtime")
	mainfn := cmd.String("mainfn", "", "The snippet main function")
	timeout := cmd.Int("timeout", 5000, "The snippet timeout in milliseconds")
	url := cmd.String("url", "", "The snippet URL (e.g., s3:///bucket/<your key>, base64:///")
	mem := cmd.String("mem", "", "The amount of memory (e.g., 1M, 1G)")
	callback := cmd.String("callback", "", "The snippet callback")
	file := cmd.String("file", "", "The snippet file")
	description := cmd.String("description", "Test description", "The snippet description")

	return &AddSnippetCmd{client: client,
		cmd:         cmd,
		id:          id,
		name:        name,
		runtime:     runtime,
		mainfn:      mainfn,
		timeout:     timeout,
		mem:         mem,
		callback:    callback,
		url:         url,
		file:        file,
		description: description,
	}
}

func (create *AddSnippetCmd) Run(args []string) error {
	create.cmd.Parse(args)

	if !create.cmd.Parsed() {
		return errors.New("Failed to parse cmd")
	}

	var err error
	var base string
	if *create.file != "" {
		base, err = b64.EncodeFileToString(*create.file)
		if err != nil {
			return err
		}
	}

	var memory uint64
	if *create.mem != "" {
		memory, err = bytefmt.ToBytes(*create.mem)
		if err != nil {
			return err
		}
	}

	req := &model.SnippetData{Id: *create.id,
		Name:        *create.name,
		Runtime:     *create.runtime,
		MainFn:      *create.mainfn,
		Timeout:     *create.timeout,
		Memory:      memory,
		Callback:    *create.callback,
		URL:         *create.url,
		Code:        base,
		Description: *create.description}

	id, mErr := create.client.AddSnippet(util.GetAccountID(), req)
	if mErr != nil {
		return mErr
	}

	fmt.Println(id)
	return nil
}
