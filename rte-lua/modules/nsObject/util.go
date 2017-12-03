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

package nsObject

import (
	"github.com/yuin/gopher-lua"
	"github.com/verizonlabs/northstar/pkg/stats"
)

const NS_OBJECT_ERROR = "nsObject error: "

func (nsObject *NsObjectModule) makeErrorMessage(msg string) string {
	return NS_OBJECT_ERROR + msg
}

func (nsObject *NsObjectModule) error(L *lua.LState, err string, timer *stats.Timer, context string) int {
	mode := 1

	if timer != nil {
		timer.Stop()
	}

	switch context {
	case "createBucket":
		ErrCreateBucket.Incr()
	case "deleteBucket":
		ErrDeleteBucket.Incr()
	case "listBuckets":
		mode = 2
		ErrListBuckets.Incr()
	case "uploadFile":
		ErrUploadFile.Incr()
	case "downloadFile":
		mode = 2
		ErrDownloadFile.Incr()
	case "deleteFile":
		ErrDeleteFile.Incr()
	case "listFiles":
		mode = 2
		ErrListFiles.Incr()
	}

	if mode == 1 {
		L.Push(lua.LString(nsObject.makeErrorMessage(err)))
		return 1
	} else {
		L.Push(lua.LNil)
		L.Push(lua.LString(nsObject.makeErrorMessage(err)))
		return 2
	}
}
