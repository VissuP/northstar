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

package object

import (
	"flag"

	"errors"
	"fmt"
	"github.com/verizonlabs/northstar/cli/commands"
	"github.com/verizonlabs/northstar/cli/util"
	"github.com/verizonlabs/northstar/object/client"
	"github.com/verizonlabs/northstar/object/model"
)

type ListFilesCMD struct {
	client *client.ObjectClient
	cmd    *flag.FlagSet
	name   *string
}

func NewListFiles(client *client.ObjectClient) commands.Command {
	cmd := flag.NewFlagSet("object-files-list", flag.ExitOnError)
	name := cmd.String("name", "testbucket", "The bucket name")
	return &ListFilesCMD{client: client, cmd: cmd, name: name}
}

func (listFiles *ListFilesCMD) Run(args []string) error {
	listFiles.cmd.Parse(args)

	if !listFiles.cmd.Parsed() {
		return errors.New("Failed to parse cmd")
	}

	files, err := listFiles.client.ListFiles(util.GetAccountID(), *listFiles.name)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		fmt.Println("No files found")
		return nil
	}

	printFiles(files)
	return nil
}

func printFiles(results []model.Object) {
	for _, data := range results {
		fmt.Printf("Name: %v, size: %v, etag: %v, lastModified: %v, storageClass: %v\n",
			data.Key, data.Size, data.Etag, data.LastModified, data.StorageClass)
	}
}
