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

package nsSFTP

import (
	"github.com/yuin/gopher-lua"
	"github.com/verizonlabs/northstar/pkg/stats"
)

const NS_SFTP_ERROR = "nsSFTP error: "

func (m *NsSFTPModule) makeErrorMessage(msg string) string {
	return NS_SFTP_ERROR + msg
}

func (m *NsSFTPModule) error(L *lua.LState, err string, timer *stats.Timer, context string) int {
	L.Push(lua.LNil)
	L.Push(lua.LString(m.makeErrorMessage(err)))

	if timer != nil {
		timer.Stop()
	}

	switch context {
	case "connect":
		ErrConnect.Incr()
	case "disconnect":
		ErrDisconnect.Incr()
	case "mkdir":
		ErrMkdir.Incr()
	case "store":
		ErrStore.Incr()
	}

	return 2
}
