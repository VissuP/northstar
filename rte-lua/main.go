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

package main

import (
	"os"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/pkg/rte"
	"github.com/verizonlabs/northstar/pkg/rte/repl"
)

// Start listening for snippets invocation requests.
// Access NS DB to fetch snippet information.
// Retrieve snippet code (source or base64 file) run it.
// Save the stdout/err to NS DB.
func main() {
	if len(os.Args) != 2 {
		mlog.Error("Usage: rte-lua <management|worker>")
		os.Exit(-1)
	}

	option := os.Args[1]
	switch option {
	case "management":
		mlog.Debug("Starting management endpoint")
		err := rte.InitManagement()
		if err != nil {
			mlog.Error("Failed to start management endpoint: %v", err)
			os.Exit(-1)
		}
	case "worker":
		err := rte.InitRTE(repl.Lua)
		if err != nil {
			mlog.Error("Failed to init worker: %v", err)
			os.Exit(-1)
		}
	default:
		mlog.Error("Wrong option selected: %v", option)
		os.Exit(-1)
	}
}
