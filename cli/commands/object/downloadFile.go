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
	"io/ioutil"
	"github.com/verizonlabs/northstar/cli/commands"
	"github.com/verizonlabs/northstar/cli/util"
	"github.com/verizonlabs/northstar/object/client"
)

type DownloadFileCMD struct {
	client      *client.ObjectClient
	cmd         *flag.FlagSet
	bucket      *string
	file        *string
	destination *string
}

func NewDownloadFile(client *client.ObjectClient) commands.Command {
	cmd := flag.NewFlagSet("object-file-download", flag.ExitOnError)
	bucket := cmd.String("bucket", "testbucket", "The bucket name")
	file := cmd.String("file", "test.csv", "The file name")
	destination := cmd.String("destination", "/tmp/test.csv", "The file destination")

	return &DownloadFileCMD{client: client,
		cmd:         cmd,
		bucket:      bucket,
		file:        file,
		destination: destination,
	}
}

func (download *DownloadFileCMD) Run(args []string) error {
	download.cmd.Parse(args)

	if !download.cmd.Parsed() {
		return errors.New("Failed to parse cmd")
	}

	data, mErr := download.client.DownloadFile(util.GetAccountID(), *download.bucket, *download.file)
	if mErr != nil {
		return mErr
	}

	err := ioutil.WriteFile(*download.destination, data.Payload, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("File %s with content type %s downloaded to %s\n",
		*download.file, data.ContentType, *download.destination)
	return nil
}
