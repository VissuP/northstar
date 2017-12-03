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

	"fmt"

	"code.cloudfoundry.org/bytefmt"
	"github.com/verizonlabs/northstar/cli/commands"
	"github.com/verizonlabs/northstar/cli/util"
	"github.com/verizonlabs/northstar/pkg/b64"
	"github.com/verizonlabs/northstar/processing/snippets/client"
	"github.com/verizonlabs/northstar/processing/snippets/model"
)

type InvokeSnippetDirectCmd struct {
	client   *client.SnippetsClient
	cmd      *flag.FlagSet
	runtime  *string
	mainfn   *string
	url      *string
	file     *string
	timeout  *int
	callback *string
	mem      *string
	args     *string
}

func NewInvokeSnippetDirect(client *client.SnippetsClient) commands.Command {
	cmd := flag.NewFlagSet("snippets-invoke-direct", flag.ExitOnError)
	runtime := cmd.String("runtime", "lua", "The snippet runtime")
	mainfn := cmd.String("mainfn", "main", "The snippet main function")
	timeout := cmd.Int("timeout", 5000, "The snippet timeout in seconds")
	url := cmd.String("url", "", "The snippet URL (e.g., s3:///bucket/<your key>, base64:///")
	file := cmd.String("file", "", "The snippet file")
	callback := cmd.String("callback", "", "The snippet callback")
	mem := cmd.String("mem", "", "The amount of memory (e.g., 1M, 1G)")
	args := cmd.String("args", "{}", "The snippet arguments")

	return &InvokeSnippetDirectCmd{client: client,
		cmd:      cmd,
		runtime:  runtime,
		mainfn:   mainfn,
		timeout:  timeout,
		url:      url,
		file:     file,
		callback: callback,
		mem:      mem,
		args:     args,
	}
}

func (invoke *InvokeSnippetDirectCmd) Run(args []string) error {
	invoke.cmd.Parse(args)

	if !invoke.cmd.Parsed() {
		return errors.New("Failed to parse cmd")
	}

	var err error
	var base string
	if *invoke.file != "" {
		base, err = b64.EncodeFileToString(*invoke.file)
		if err != nil {
			return err
		}
	}

	arguments, err := util.UnmarshalInterface(*invoke.args)
	if err != nil {
		return err
	}

	var memory uint64
	if *invoke.mem != "" {
		memory, err = bytefmt.ToBytes(*invoke.mem)
		if err != nil {
			return err
		}
	}

	req := &model.Snippet{
		Runtime: *invoke.runtime,
		MainFn:  *invoke.mainfn,
		Timeout: *invoke.timeout,
		URL:     *invoke.url,
		Code:    base,
		Options: model.Options{Callback: *invoke.callback,
			Memory: memory,
			Args:   arguments}}

	out, mErr := invoke.client.StartSnippet(util.GetAccountID(), req)
	if mErr != nil {
		return mErr
	}

	fmt.Println(out)
	return nil
}
