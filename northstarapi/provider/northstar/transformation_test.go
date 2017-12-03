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

package northstar

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	snippetsModel "github.com/verizonlabs/northstar/data/snippets/model"
	"github.com/verizonlabs/northstar/northstarapi/model"
)

func TestFromExternal(t *testing.T) {

	Convey("Test fromExternal()", t, func() {
		uut := &NorthStarTransformationProvider{}

		// Code
		externalSnippet := &snippetsModel.SnippetData{
			Id:        "s123",
			Name:      "name",
			CreatedOn: time.Now(),
			UpdatedOn: time.Now(),
			Runtime:   "runtime",
			MainFn:    "main",
			Timeout:   1000,
			Code:      "code",
			URL:       "base64://code",
			EventType: snippetsModel.DeviceEventType,
			EventId:   "d123",
		}

		transformation, err := uut.fromExternal(externalSnippet)

		So(transformation, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(transformation.Code.Type, ShouldEqual, model.SourceCodeType)
		So(transformation.Code.Url, ShouldEqual, "base64://code")
		So(transformation.Code.Value, ShouldEqual, "code")
		So(transformation.Scheduled, ShouldEqual, true)
		So(transformation.SchedulerType, ShouldEqual, model.DeviceEvent)
		So(transformation.SchedulerId, ShouldEqual, "d123")

		// S3
		externalSnippet = &snippetsModel.SnippetData{
			Id:        "s123",
			Name:      "name",
			CreatedOn: time.Now(),
			UpdatedOn: time.Now(),
			Runtime:   "runtime",
			MainFn:    "main",
			Timeout:   1000,
			Code:      "code",
			URL:       "s3://bucket/container",
			EventType: snippetsModel.TimerEventType,
			EventId:   "t123",
		}

		transformation, err = uut.fromExternal(externalSnippet)

		So(transformation, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(transformation.Code.Type, ShouldEqual, model.ObjectCodeType)
		So(transformation.Code.Url, ShouldEqual, "code")
		So(transformation.Code.Value, ShouldEqual, "code")
		So(transformation.Scheduled, ShouldEqual, true)
		So(transformation.SchedulerType, ShouldEqual, model.TimerEvent)
		So(transformation.SchedulerId, ShouldEqual, "t123")

		// Empty Event Type
		externalSnippet = &snippetsModel.SnippetData{
			Id:        "s123",
			Name:      "name",
			CreatedOn: time.Now(),
			UpdatedOn: time.Now(),
			Runtime:   "runtime",
			MainFn:    "main",
			Timeout:   1000,
			Code:      "code",
			URL:       "s3://bucket/container",
		}

		transformation, err = uut.fromExternal(externalSnippet)

		So(transformation, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(transformation.Code.Type, ShouldEqual, model.ObjectCodeType)
		So(transformation.Code.Url, ShouldEqual, "code")
		So(transformation.Code.Value, ShouldEqual, "code")
		So(transformation.Scheduled, ShouldEqual, false)
		So(transformation.SchedulerType, ShouldEqual, "")
		So(transformation.SchedulerId, ShouldEqual, "")

		// None Event Type
		externalSnippet = &snippetsModel.SnippetData{
			Id:        "s123",
			Name:      "name",
			CreatedOn: time.Now(),
			UpdatedOn: time.Now(),
			Runtime:   "runtime",
			MainFn:    "main",
			Timeout:   1000,
			Code:      "code",
			URL:       "s3://bucket/container",
			EventType: snippetsModel.NoneEventType,
		}

		transformation, err = uut.fromExternal(externalSnippet)

		So(transformation, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(transformation.Code.Type, ShouldEqual, model.ObjectCodeType)
		So(transformation.Code.Url, ShouldEqual, "code")
		So(transformation.Code.Value, ShouldEqual, "code")
		So(transformation.Scheduled, ShouldEqual, false)
		So(transformation.SchedulerType, ShouldEqual, "")
		So(transformation.SchedulerId, ShouldEqual, "")
	})
}

func TestToExternal(t *testing.T) {

	Convey("Test toExternal()", t, func() {
		uut := &NorthStarTransformationProvider{}

		// No schedule
		transformation := &model.Transformation{
			Id:          "s123",
			Name:        "name",
			Description: "description",
			Language:    "lua",
			EntryPoint:  "main",
			Timeout:     1000,
			Scheduled:   false,
			Code: model.Code{
				Type:  model.SourceCodeType,
				Value: "code",
			},
		}

		externalSnippet := uut.toExternal(transformation)

		So(externalSnippet.EventType, ShouldEqual, snippetsModel.NoneEventType)
		So(externalSnippet.EventId, ShouldEqual, "")
		So(externalSnippet.URL, ShouldEqual, "base64://code")

		// Timer Schedule
		transformation = &model.Transformation{
			Id:            "s123",
			Name:          "name",
			Description:   "description",
			Language:      "lua",
			EntryPoint:    "main",
			Timeout:       1000,
			Scheduled:     true,
			SchedulerType: model.TimerEvent,
			SchedulerId:   "t123",
			Code: model.Code{
				Type:  model.SourceCodeType,
				Value: "code",
			},
		}

		externalSnippet = uut.toExternal(transformation)

		So(externalSnippet.EventType, ShouldEqual, snippetsModel.TimerEventType)
		So(externalSnippet.EventId, ShouldEqual, "t123")
		So(externalSnippet.URL, ShouldEqual, "base64://code")

		// Event Schedule
		transformation = &model.Transformation{
			Id:            "s123",
			Name:          "name",
			Description:   "description",
			Language:      "lua",
			EntryPoint:    "main",
			Timeout:       1000,
			Scheduled:     true,
			SchedulerType: model.DeviceEvent,
			SchedulerId:   "d123",
			Code: model.Code{
				Type:  model.SourceCodeType,
				Value: "code",
			},
		}

		externalSnippet = uut.toExternal(transformation)

		So(externalSnippet.EventType, ShouldEqual, snippetsModel.DeviceEventType)
		So(externalSnippet.EventId, ShouldEqual, "d123")
		So(externalSnippet.URL, ShouldEqual, "base64://code")
	})
}
