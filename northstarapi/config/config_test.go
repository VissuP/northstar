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

package config

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {

	Convey("Test Load()", t, func() {
		os.Setenv(ScopesFileNameEnv, "./scopes.test.json")

		os.Setenv(ThingSpaceProtocolEnv, "http")
		os.Setenv(ThingSpaceSouthHostPortEnv, "127.0.0.1:8080")
		os.Setenv(ThingSpaceAuthHostPortEnv, "127.0.0.1:8080")
		os.Setenv(ThingSpaceDataHostPortEnv, "127.0.0.1:8080")
		os.Setenv(ThingSpaceClientIdEnv, "id")
		os.Setenv(ThingSpaceSecretEnv, "secret")
		os.Setenv(KafkaBrokersHostPortEnv, "127.0.0.1:8080")
		os.Setenv(ZookeperHostPortEnd, "127.0.0.1:8080")

		err := Load()

		So(err, ShouldBeNil)
		So(Configuration.ServiceName, ShouldEqual, ServiceName)
		So(len(Configuration.Scopes["scope"]), ShouldEqual, 1)
		So(len(Configuration.Scopes["scope"][0].Methods), ShouldEqual, 1)
		So(Configuration.Scopes["scope"][0].Methods[0], ShouldEqual, "GET")
		So(len(Configuration.Scopes["scope"][0].Paths), ShouldEqual, 1)
		So(Configuration.Scopes["scope"][0].Paths[0], ShouldEqual, "/path")
	})
}
