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

package nsUtil

import (
	"bytes"
	"encoding/binary"
	"github.com/yuin/gopher-lua"
	"reflect"
	"github.com/verizonlabs/northstar/pkg/stats"
	"github.com/verizonlabs/northstar/rte-lua/util"
)

const (
	NS_UTIL_ERROR        = "nsUtil error: "
	READ_FROM_BYTE_ARRAY = "readFromByteArray"
)

type NsUtilModule struct{}

func NewNsUtilModule() *NsUtilModule {
	return &NsUtilModule{}
}

func (nsUtil *NsUtilModule) Loader(L *lua.LState) int {
	api := map[string]lua.LGFunction{
		READ_FROM_BYTE_ARRAY: nsUtil.readFromByteArray,
	}
	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

func (nsUtil *NsUtilModule) readFromByteArray(L *lua.LState) int {
	table, err := util.FromLua(L.CheckTable(1))
	if err != nil {
		return nsUtil.error(L, err.Error(), nil, READ_FROM_BYTE_ARRAY, 2)
	}

	rawData, ok := table.([]interface{})
	if !ok {
		return nsUtil.error(L, "invalid data", nil, READ_FROM_BYTE_ARRAY, 2)
	}

	data := make([]byte, len(rawData))
	for i, element := range rawData {
		e, ok := element.(float64)
		if !ok {
			return nsUtil.error(L, "invalid data", nil, READ_FROM_BYTE_ARRAY, 2)
		}
		data[i] = byte(int(e))
	}

	start := L.CheckInt64(2)

	var order binary.ByteOrder
	switch L.CheckString(3) {
	case "BigEndian":
		order = binary.BigEndian
	case "LittleEndian":
		order = binary.LittleEndian
	default:
		return nsUtil.error(L, "unknown byte order", nil, READ_FROM_BYTE_ARRAY, 2)
	}

	var reading interface{}
	switch L.CheckString(4) {
	case "uint8":
		reading = new(uint8)
	case "uint16":
		reading = new(uint16)
	case "uint", "uint32":
		reading = new(uint32)
	case "uint64":
		reading = new(uint64)
	case "int8":
		reading = new(int8)
	case "int16":
		reading = new(int16)
	case "int", "int32":
		reading = new(int32)
	case "int64":
		reading = new(int64)
	case "float32":
		reading = new(float32)
	case "float64":
		reading = new(float64)
	default:
		return nsUtil.error(L, "unknown data type", nil, READ_FROM_BYTE_ARRAY, 2)
	}

	if err := binary.Read(bytes.NewReader(data[start:]), order, reading); err != nil {
		return nsUtil.error(L, err.Error(), nil, READ_FROM_BYTE_ARRAY, 2)
	}

	result, err := util.ToLua(L, reflect.ValueOf(reading).Elem().Interface())
	if err != nil {
		return nsUtil.error(L, err.Error(), nil, READ_FROM_BYTE_ARRAY, 2)
	}

	ReadFromByteArray.Incr()
	L.Push(result)
	return 1
}

func (nsUtil *NsUtilModule) makeErrorMessage(msg string) string {
	return NS_UTIL_ERROR + msg
}

func (nsUtil *NsUtilModule) recordErrorStats(timer *stats.Timer, context string) {
	if timer != nil {
		timer.Stop()
	}

	switch context {
	case READ_FROM_BYTE_ARRAY:
		ErrReadFromByteArray.Incr()
	}
}

func (nsUtil *NsUtilModule) error(L *lua.LState,
	err string,
	timer *stats.Timer,
	context string,
	nRetElements int) int {
	nsUtil.recordErrorStats(timer, context)
	for i := 0; i < nRetElements-1; i++ {
		L.Push(lua.LNil)
	}
	L.Push(lua.LString(nsUtil.makeErrorMessage(err)))

	return nRetElements
}

func (nsUtil *NsUtilModule) panic(err string, timer *stats.Timer, context string) {
	nsUtil.recordErrorStats(timer, context)
	panic(nsUtil.makeErrorMessage(err))
}
