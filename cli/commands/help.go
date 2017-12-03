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

package commands

import (
	"fmt"
)

func PrintHelp() {
	fmt.Println("usage: nor-cli <command> [<args>]")
	fmt.Println("The available commands are: ")
	fmt.Println("	object-bucket-create            Creates a bucket")
	fmt.Println("	object-bucket-list              List buckets")
	fmt.Println("	object-bucket-delete            Delete bucket")
	fmt.Println("	object-file-upload              Upload file")
	fmt.Println("	object-file-download            File download")
	fmt.Println("	object-file-list                List files")
	fmt.Println("	object-file-delete              Delete file")
	fmt.Println("	topics-add                      Add topic")
	fmt.Println("	topics-list                     List topics")
	fmt.Println("	topics-update                   Update topic")
	fmt.Println("	snippets-add                    Add snippet")
	fmt.Println("	snippets-update                 Update snippet")
	fmt.Println("	snippets-invoke-direct          Invoke snippet directly")
	fmt.Println("	snippets-invoke-by-id           Invoke snippet by id")
	fmt.Println("	snippets-list                   List snippets")
	fmt.Println("	snippets-delete                 Delete snippet")
	fmt.Println("	cron-add                        Add cron job")
	fmt.Println("	cron-update                     Update cron job")
	fmt.Println("	cron-list                       List cron jobs")
	fmt.Println("	cron-delete                     Delete cron job")
	fmt.Println("	invoke-get                      Get invocation")
	fmt.Println("	invoke-list                     List invocations")
	fmt.Println("	invoke-delete                   Delete invocation")
	fmt.Println("	invoke-stop                     Stop invocation")
	fmt.Println("	events-add                      Add event")
	fmt.Println("	events-invoke                   Invoke event")
	fmt.Println("	events-list                     Lists events")
	fmt.Println("	events-delete                   Delete event")
	fmt.Println("	mappings-add                    Add mapping")
	fmt.Println("	mappings-list                   Lists mappings")
	fmt.Println("	mappings-delete                 Delete mapping")
	fmt.Println("	datasources-add                 Add datasources")
	fmt.Println("	datasources-get                 Get datasources")
	fmt.Println("	datasources-list                Lists datasources")
	fmt.Println("	datasources-delete              Delete datasources")
	fmt.Println("	datasets-add                    Add dataset")
	fmt.Println("	datasets-get-by-id              Get dataset by id")
	fmt.Println("	datasets-get-by-name            Get dataset by name")
	fmt.Println("	datasets-list                   Lists datasets")
	fmt.Println("	datasets-delete                 Delete dataset")
}
