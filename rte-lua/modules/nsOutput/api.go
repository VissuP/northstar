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
	"github.com/yuin/gopher-lua/modules/gluamapper"
)

func (nsOutput *NsOutputModule) printApi(L *lua.LState) int {
	var parameters []interface{}
	for i := 1; i <= L.GetTop(); i++ {
		parameters = append(parameters, L.CheckAny(i))
	}

	err := nsOutput.print(parameters)
	if err != nil {
		nsOutput.panic(err.Error(), nil, _PRINT)
	}

	Print.Incr()
	return 0
}

func (nsOutput *NsOutputModule) printfApi(L *lua.LState) int {
	var parameters []interface{}
	for i := 1; i <= L.GetTop(); i++ {
		parameters = append(parameters, L.CheckAny(i))
	}

	err := nsOutput.printf(parameters)
	if err != nil {
		nsOutput.panic(err.Error(), nil, _PRINTF)
	}

	Printf.Incr()
	return 0
}

func (nsOutput *NsOutputModule) valueApi(L *lua.LState) int {
	output, err := nsOutput.generateFromTable(L.CheckTable(1), &Value{})
	if err != nil {
		return nsOutput.error(L, err.Error(), nil, _VALUE, 2)
	}

	L.Push(lua.LString(output))
	ValueCounter.Incr()
	return 1
}

func (nsOutput *NsOutputModule) valueDirectApi(L *lua.LState) int {
	output, err := nsOutput.generateFromTable(L.CheckTable(1), &Value{})
	if err != nil {
		return nsOutput.error(L, err.Error(), nil, _VALUE_DIRECT, 1)
	}

	nsOutput.Result = output
	ValueDirectCounter.Incr()
	return 0
}

func (nsOutput *NsOutputModule) tableApi(L *lua.LState) int {
	output, err := nsOutput.generateFromTable(L.CheckTable(1), &Table{})
	if err != nil {
		return nsOutput.error(L, err.Error(), nil, _TABLE, 2)
	}

	L.Push(lua.LString(output))
	TableCounter.Incr()
	return 1
}

func (nsOutput *NsOutputModule) tableDirectApi(L *lua.LState) int {
	output, err := nsOutput.generateFromTable(L.CheckTable(1), &Table{})
	if err != nil {
		return nsOutput.error(L, err.Error(), nil, _TABLE_DIRECT, 1)
	}

	nsOutput.Result = output
	TableDirectCounter.Incr()
	return 0
}

func (nsOutput *NsOutputModule) mapApi(L *lua.LState) int {
	output, err := nsOutput.generateFromTable(L.CheckTable(1), &Map{})
	if err != nil {
		return nsOutput.error(L, err.Error(), nil, _MAP, 2)
	}

	L.Push(lua.LString(output))
	MapCounter.Incr()
	return 1
}

func (nsOutput *NsOutputModule) mapDirectApi(L *lua.LState) int {
	output, err := nsOutput.generateFromTable(L.CheckTable(1), &Map{})
	if err != nil {
		return nsOutput.error(L, err.Error(), nil, _MAP_DIRECT, 1)
	}

	nsOutput.Result = output
	MapDirectCounter.Incr()
	return 0
}

func (nsOutput *NsOutputModule) htmlApi(L *lua.LState) int {
	output, err := nsOutput.generateOutput("text/html", L.CheckString(1))
	if err != nil {
		return nsOutput.error(L, err.Error(), nil, _HTML, 2)
	}

	L.Push(lua.LString(output))
	HTMLCounter.Incr()
	return 1
}

func (nsOutput *NsOutputModule) htmlDirectApi(L *lua.LState) int {
	output, err := nsOutput.generateOutput("text/html", L.CheckString(1))
	if err != nil {
		return nsOutput.error(L, err.Error(), nil, _HTML_DIRECT, 2)
	}

	nsOutput.Result = output
	HTMLDirectCounter.Incr()
	return 0
}

func (nsOutput *NsOutputModule) tableToCsvApi(L *lua.LState) int {
	var data string
	var table Table
	var err error

	mapper := gluamapper.NewMapper(gluamapper.Option{NameFunc: func(str string) string { return str }})
	if err = mapper.Map(L.CheckTable(1), &table); err != nil {
		return nsOutput.error(L, err.Error(), nil, _TABLE_TO_CSV, 2)
	}

	data, err = nsOutput.tableToCsv(&table)
	if err != nil {
		return nsOutput.error(L, err.Error(), nil, _TABLE_TO_CSV, 2)
	}

	TableToCsv.Incr()
	L.Push(lua.LString(data))
	return 1
}
