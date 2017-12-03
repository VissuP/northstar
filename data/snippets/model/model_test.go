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

package model

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestValidateOnAdd(t *testing.T) {

	Convey("Test ValidateOnAdd()", t, func() {
		snippet := SnippetData{
			Name:      "name",
			CreatedOn: time.Now(),
			Runtime:   "runtime",
			MainFn:    "main",
			URL:       "url",
			Timeout:   100,
			EventType: TimerEventType,
			EventId:   "123",
		}

		err := snippet.Validate()
		So(err, ShouldBeNil)

		// Missing name.
		errSnipper := snippet
		errSnipper.Name = ""
		err = errSnipper.Validate()
		So(err, ShouldNotBeNil)

		// Missing runtime.
		errSnipper = snippet
		errSnipper.Runtime = ""
		err = errSnipper.Validate()
		So(err, ShouldNotBeNil)

		// Missing mainfn.
		errSnipper = snippet
		errSnipper.MainFn = ""
		err = errSnipper.Validate()
		So(err, ShouldNotBeNil)

		// Missing url.
		errSnipper = snippet
		errSnipper.URL = ""
		err = errSnipper.Validate()
		So(err, ShouldNotBeNil)

		// Missing timeout.
		errSnipper = snippet
		errSnipper.Timeout = 0
		err = errSnipper.Validate()
		So(err, ShouldNotBeNil)

		// Missing event id.
		errSnipper = snippet
		errSnipper.EventId = ""
		err = errSnipper.Validate()
		So(err, ShouldNotBeNil)
	})
}
