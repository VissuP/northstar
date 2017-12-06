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
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/verizonlabs/northstar/northstarapi/config"
)

/*
	Type       CellType  `json:"type"`
	Language   string    `json:"language"`
	Arguments  Arguments `json:"arguments,omitempty"`
	EntryPoint string    `json:"entryPoint"`
	Body       string    `json:"body"`
	Timeout    int       `json:"timeout"`
*/
func TestValidateNotebook(t *testing.T) {
	Convey("Test -- Empty notebook is valid.", t, func() {
		notebook := Notebook{
			Id:          "5f691b96-e0c8-478d-bc58-a5514b7c8c96",
			Name:        "Test notebook",
			Permissions: OwnerPermissions,
			CreatedOn:   time.Now().String(),
		}
		err := notebook.Validate()

		So(err, ShouldBeNil)
	})

	Convey("Test -- Valid notebook validates successfully", t, func() {
		notebook := Notebook{
			Id:          "5f691b96-e0c8-478d-bc58-a5514b7c8c96",
			Name:        "Test notebook",
			Permissions: OwnerPermissions,
			CreatedOn:   time.Now().String(),
			Cells: []Cell{
				{
					Id:   "5f691b96-e0c8-478d-bc58-a5514b7c8c95",
					Name: "Test cell",
					Input: Input{
						Type:       CodeCellType,
						Language:   "Lua",
						EntryPoint: "main",
						Body:       "YourCodeHere",
						Timeout:    180,
					},
					Settings: Settings{},
				},
			},
		}
		err := notebook.Validate()

		So(err, ShouldBeNil)
	})

	Convey("Test -- Body cannot be empty", t, func() {
		notebook := Notebook{
			Id:          "5f691b96-e0c8-478d-bc58-a5514b7c8c96",
			Name:        "Test notebook",
			Permissions: OwnerPermissions,
			CreatedOn:   time.Now().String(),
			Cells: []Cell{
				{
					Id:   "5f691b96-e0c8-478d-bc58-a5514b7c8c95",
					Name: "Test cell",
					Input: Input{
						Type:       CodeCellType,
						Language:   "Lua",
						EntryPoint: "main",
						Timeout:    180,
					},
					Settings: Settings{},
				},
			},
		}
		err := notebook.Validate()

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "The input body is missing.")
	})

	Convey("Test -- Timeout cannot be greater than 180 seconds", t, func() {
		notebook := Notebook{
			Id:          "5f691b96-e0c8-478d-bc58-a5514b7c8c96",
			Name:        "Test notebook",
			Permissions: OwnerPermissions,
			CreatedOn:   time.Now().String(),
			Cells: []Cell{
				{
					Id:   "5f691b96-e0c8-478d-bc58-a5514b7c8c95",
					Name: "Test cell",
					Input: Input{
						Type:       CodeCellType,
						Language:   "Lua",
						EntryPoint: "main",
						Body:       "YourCodeHere",
						Timeout:    181,
					},
					Settings: Settings{},
				},
			},
		}
		err := notebook.Validate()

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "The timeout value is out of range (181 of 180).")
	})

	Convey("Test -- Cell name cannot be empty", t, func() {
		notebook := Notebook{
			Id:          "5f691b96-e0c8-478d-bc58-a5514b7c8c96",
			Name:        "Test notebook",
			Permissions: OwnerPermissions,
			CreatedOn:   time.Now().String(),
			Cells: []Cell{
				{
					Id: "5f691b96-e0c8-478d-bc58-a5514b7c8c95",
					Input: Input{
						Type:       CodeCellType,
						Language:   "Lua",
						EntryPoint: "main",
						Body:       "YourCodeHere",
						Timeout:    180,
					},
					Settings: Settings{},
				},
			},
		}
		err := notebook.Validate()

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "The cell name is empty.")
	})

	Convey("Test -- Cannot set greater than max memory", t, func() {
		notebook := Notebook{
			Id:          "5f691b96-e0c8-478d-bc58-a5514b7c8c96",
			Name:        "Test notebook",
			Permissions: OwnerPermissions,
			CreatedOn:   time.Now().String(),
			Cells: []Cell{
				{
					Id:   "5f691b96-e0c8-478d-bc58-a5514b7c8c95",
					Name: "Test cell",
					Input: Input{
						Type:       CodeCellType,
						Language:   "Lua",
						EntryPoint: "main",
						Body:       "YourCodeHere",
						Timeout:    180,
					},
					Settings: Settings{
						Memory: config.Configuration.MaxMemory + 1,
					},
				},
			},
		}
		err := notebook.Validate()

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "Requested memory (1) is greater than the max of: 0.")
	})
}
