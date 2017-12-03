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

package api

import (
	"bytes"
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"testing"
)

func TestDeviceSerialization(t *testing.T) {
	Convey("Test Device Serialization", t, func() {
		mlog.EnableDebug(true)
		mlog.Debug(" DEBUG log level enabled ")

		var device Device
		updates := "{\"kind\": \"ts.device.ultratag\", \"provisionedOn\":\"2015-10-16T19:20:30.45+01:00\", \"fields\": {\"reportinginterval\": {\"value\":900}, \"firmwareversion\":\"1.0\"}, \"iccid\":\"01234567890123456789\", \"state\":\"xyz\"}"

		decoder := json.NewDecoder(bytes.NewReader([]byte(updates)))

		err := decoder.Decode(&device)

		So(err, ShouldBeNil)
		So(device.Kind, ShouldEqual, "ts.device.ultratag")
		So(device.ProvisionedOn, ShouldEqual, "2015-10-16T19:20:30.45+01:00")
		So(device.Fields["reportinginterval"], ShouldNotBeNil)
		So(device.Fields["firmwareversion"], ShouldNotBeNil)
		So(device.State, ShouldEqual, "xyz")
	})
}
