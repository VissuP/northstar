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

package nsFTP

import "github.com/verizonlabs/northstar/pkg/stats"

var (
	NsFTP         = stats.New("nsFTP")
	Connect       = NsFTP.NewCounter("Connect")
	Disconnect    = NsFTP.NewCounter("Disconnect")
	Login         = NsFTP.NewCounter("Login")
	Logout        = NsFTP.NewCounter("Logout")
	Mkdir         = NsFTP.NewCounter("Mkdir")
	Store         = NsFTP.NewCounter("Store")
	ErrConnect    = NsFTP.NewCounter("ErrConnect")
	ErrDisconnect = NsFTP.NewCounter("ErrDisconnect")
	ErrLogin      = NsFTP.NewCounter("ErrLogin")
	ErrLogout     = NsFTP.NewCounter("ErrLogout")
	ErrMkdir      = NsFTP.NewCounter("ErrMkdir")
	ErrStore      = NsFTP.NewCounter("ErrStore")
)
