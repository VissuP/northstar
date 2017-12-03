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

package nsOutput

import (
	"github.com/yuin/gopher-lua"
	"github.com/verizonlabs/northstar/pkg/config"
)

const (
	_PRINT        = "print"
	_PRINTF       = "printf"
	_VALUE        = "value"
	_VALUE_DIRECT = "valueDirect"
	_TABLE        = "table"
	_TABLE_DIRECT = "tableDirect"
	_MAP          = "map"
	_MAP_DIRECT   = "mapDirect"
	_HTML         = "html"
	_HTML_DIRECT  = "htmlDirect"
	_TABLE_TO_CSV = "tableToCsv"
)

var (
	NsOutputPrintLimit, _ = config.GetInt("NS_OUTPUT_PRINT_LIMIT", 10000)
)

type NsOutputModule struct {
	Limit   int
	Rolling int
	Stdout  []string
	Result  string
}

func NewNsOutputModule() *NsOutputModule {
	return &NsOutputModule{Limit: NsOutputPrintLimit,
		Rolling: NsOutputPrintLimit,
		Stdout:  []string{}}
}

func (nsOutput *NsOutputModule) Loader(L *lua.LState) int {
	api := map[string]lua.LGFunction{
		_PRINT:        nsOutput.printApi,
		_PRINTF:       nsOutput.printfApi,
		_VALUE:        nsOutput.valueApi,
		_VALUE_DIRECT: nsOutput.valueDirectApi,
		_TABLE:        nsOutput.tableApi,
		_TABLE_DIRECT: nsOutput.tableDirectApi,
		_MAP:          nsOutput.mapApi,
		_MAP_DIRECT:   nsOutput.mapDirectApi,
		_HTML:         nsOutput.htmlApi,
		_HTML_DIRECT:  nsOutput.htmlDirectApi,
		_TABLE_TO_CSV: nsOutput.tableToCsvApi,
	}
	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

func (nsOutput *NsOutputModule) Reset() {
	nsOutput.Rolling = NsOutputPrintLimit
	nsOutput.Stdout = nil
	nsOutput.Result = ""
}
