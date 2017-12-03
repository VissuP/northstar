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

type ListBucketCMD struct {
	client *client.ObjectClient
	cmd    *flag.FlagSet
}

func NewListBucket(client *client.ObjectClient) commands.Command {
	cmd := flag.NewFlagSet("object-bucket-list", flag.ExitOnError)

	return &ListBucketCMD{client: client, cmd: cmd}
}

func (listBuckets *ListBucketCMD) Run(args []string) error {
	listBuckets.cmd.Parse(args)

	if !listBuckets.cmd.Parsed() {
		return errors.New("Failed to parse cmd")
	}

	buckets, err := listBuckets.client.ListBuckets(util.GetAccountID())
	if err != nil {
		return err
	}

	if len(buckets) == 0 {
		fmt.Println("No buckets found")
		return nil
	}

	printBuckets(buckets)
	return nil
}

func printBuckets(results []model.Bucket) {
	for _, data := range results {
		fmt.Printf("Name: %v, creationDate: %v\n", data.Name, data.CreationDate)
	}
}
