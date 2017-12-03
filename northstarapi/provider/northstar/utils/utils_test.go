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

package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/pkg/rte/repl"
)

func TestGetOutputStatus(t *testing.T) {

	Convey("Test GetOutputStatus() - Valid status/message", t, func() {
		for key, value := range outputStatusCodeMap {
			code, description := GetOutputStatus(key, "Message")
			So(code, ShouldEqual, value)
			So(description, ShouldEqual, "Message")
		}
	})

	Convey("Test GetOutputStatus() - Empty message", t, func() {
		code, description := GetOutputStatus(repl.SNIPPET_RUN_FINISHED, "")
		So(code, ShouldEqual, model.OutputSuccessStatus)
		So(description, ShouldEqual, model.DefaultOutputStatusDescriptions[model.OutputSuccessStatus])
	})

	Convey("Test GetOutputStatus() - Empty status", t, func() {
		code, description := GetOutputStatus("", "Should be ignored!")
		So(code, ShouldEqual, model.OutputUnknownStatus)
		So(description, ShouldEqual, model.DefaultOutputStatusDescriptions[model.OutputUnknownStatus])
	})
}
