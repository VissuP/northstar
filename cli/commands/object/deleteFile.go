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
)

type DeleteFileCMD struct {
	client *client.ObjectClient
	cmd    *flag.FlagSet
	bucket *string
	file   *string
}

func NewDeleteFile(client *client.ObjectClient) commands.Command {
	cmd := flag.NewFlagSet("object-file-delete", flag.ExitOnError)
	bucket := cmd.String("bucket", "testbucket", "The bucket name")
	file := cmd.String("file", "test", "The file name")

	return &DeleteFileCMD{client: client,
		cmd:    cmd,
		bucket: bucket,
		file:   file}
}

func (deleteFile *DeleteFileCMD) Run(args []string) error {
	deleteFile.cmd.Parse(args)

	if !deleteFile.cmd.Parsed() {
		return errors.New("Failed to parse cmd")
	}

	err := deleteFile.client.DeleteFile(util.GetAccountID(), *deleteFile.bucket, *deleteFile.file)
	if err != nil {
		return err
	}

	fmt.Printf("File %s deleted\n", *deleteFile.file)
	return nil
}
