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
	"github.com/verizonlabs/northstar/object/model"
)

type UploadFileCMD struct {
	client      *client.ObjectClient
	cmd         *flag.FlagSet
	bucket      *string
	contentType *string
	localFile   *string
	remoteFile  *string
}

func NewUploadFile(client *client.ObjectClient) commands.Command {
	cmd := flag.NewFlagSet("object-file-upload", flag.ExitOnError)
	bucket := cmd.String("bucket", "testbucket", "The bucket name")
	contentType := cmd.String("contentType", "text/plain", "The content type")
	localFile := cmd.String("localFile", "test.csv", "The local file name")
	remoteFile := cmd.String("remoteFile", "test.csv", "The local file name")

	return &UploadFileCMD{client: client,
		cmd:         cmd,
		bucket:      bucket,
		contentType: contentType,
		localFile:   localFile,
		remoteFile:  remoteFile}
}

func (upload *UploadFileCMD) Run(args []string) error {
	upload.cmd.Parse(args)

	if !upload.cmd.Parsed() {
		return errors.New("Failed to parse cmd")
	}

	byteArr, err := ioutil.ReadFile(*upload.localFile)
	if err != nil {
		return err
	}

	data := model.UploadData{FileName: *upload.remoteFile,
		ContentType: *upload.contentType,
		Payload:     byteArr}
	_, mErr := upload.client.UploadFile(util.GetAccountID(),
		*upload.bucket,
		&data)
	if mErr != nil {
		return mErr
	}

	fmt.Printf("File %s uploaded as %s\n", *upload.localFile, *upload.remoteFile)
	return nil
}
