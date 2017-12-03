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

package thingspace

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/verizonlabs/northstar/pkg/thingspace/api"
	"github.com/verizonlabs/northstar/northstarapi/model"
)

func TestIsEmail(t *testing.T) {

	Convey("Test isEmail()", t, func() {
		results := isEmail("first.last@verizon.com")
		So(results, ShouldEqual, true)

		results = isEmail("First.Last@verizon.com")
		So(results, ShouldEqual, true)

		results = isEmail("first.last@verizon")
		So(results, ShouldEqual, false)

		results = isEmail("first.last@")
		So(results, ShouldEqual, false)

		results = isEmail("First Last")
		So(results, ShouldEqual, false)
	})
}

func TestContains(t *testing.T) {

	Convey("Test contains()", t, func() {
		results := contains("first.last@verizon.com", "first")
		So(results, ShouldEqual, true)

		results = contains("first.last@verizon.com", "first.last")
		So(results, ShouldEqual, true)

		results = contains("first.last@verizon.com", "last")
		So(results, ShouldEqual, true)

		results = contains("first.last@verizon.com", "unknown")
		So(results, ShouldEqual, false)

		results = contains("", "unknown")
		So(results, ShouldEqual, false)

		results = contains("first.last@verizon.com", "")
		So(results, ShouldEqual, false)

		results = contains("", "")
		So(results, ShouldEqual, false)
	})
}

func TestSearch(t *testing.T) {
	users := []api.User{
		api.User{
			Id:          "1",
			ForeignId:   "1234",
			DisplayName: "John Smith",
			Email:       "john.smith@verizon.com",
		},
		api.User{
			Id:          "2",
			ForeignId:   "1234",
			DisplayName: "Paul Smith",
			Email:       "paul.jones@verizon.com",
		},
		api.User{
			Id:          "3",
			ForeignId:   "1234",
			DisplayName: "Alex Williams",
			Email:       "alex.1.williams@verizon.com",
		},
		api.User{
			Id:          "4",
			ForeignId:   "1234",
			DisplayName: "Alex Brown",
			Email:       "alex.2.brown@verizon.com",
		},
		api.User{
			Id:          "5",
			ForeignId:   "1234",
			DisplayName: "Alex Taylor",
			Email:       "alex.3.taylor@verizon.com",
		},
		api.User{
			Id:          "6",
			ForeignId:   "1234",
			DisplayName: "Alex Miller",
			Email:       "alex.4.miller@verizon.com",
		},
		api.User{
			Id:          "7",
			ForeignId:   "1234",
			DisplayName: "Rich Martin",
			Email:       "rich.martin@verizon.com",
		},
		api.User{
			Id:          "8",
			ForeignId:   "1234",
			DisplayName: "Marta Ramos",
			Email:       "marta.ramos@verizon.com",
		},
		api.User{
			Id:          "9",
			ForeignId:   "1234",
			DisplayName: "Cristy Clark",
			Email:       "cristy.clark@verizon.com",
		},
		api.User{
			Id:          "10",
			ForeignId:   "1234",
			DisplayName: "Maria Ramos",
			Email:       "maria.santos@verizon.com",
		},
	}

	Convey("Test search() - One Match", t, func() {
		numMatches := 2
		maxRoutines := 2
		criteria := &model.User{
			DisplayName: "John",
		}

		matches := search(numMatches, maxRoutines, criteria, users)
		So(len(matches), ShouldEqual, 1)

		numMatches = 10
		maxRoutines = 3
		criteria = &model.User{
			DisplayName: "Maria",
		}

		matches = search(numMatches, maxRoutines, criteria, users)
		So(len(matches), ShouldEqual, 1)
	})

	Convey("Test search() - Several Match", t, func() {
		numMatches := 2
		maxRoutines := 2
		criteria := &model.User{
			DisplayName: "Alex",
		}

		matches := search(numMatches, maxRoutines, criteria, users)
		So(len(matches), ShouldEqual, 2)

		numMatches = 10
		maxRoutines = 2
		criteria = &model.User{
			DisplayName: "Alex",
		}

		matches = search(numMatches, maxRoutines, criteria, users)
		So(len(matches), ShouldEqual, 4)
	})
}
