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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/modules/gluamapper"
	"github.com/verizonlabs/northstar/pkg/stats"
	"github.com/verizonlabs/northstar/rte-lua/util"
	"strings"
)

const NS_OUTPUT_ERROR = "nsOutput error: "

func (nsOutput *NsOutputModule) print(lParameters []interface{}) error {
	parameters, err := nsOutput.convertParameters(lParameters)
	if err != nil {
		return err
	}

	out := fmt.Sprint(parameters...)
	if nsOutput.Rolling-len([]byte(out)) < 0 {
		return errors.New(fmt.Sprintf("%d-byte stdout limit is exceeded", nsOutput.Limit))
	}

	nsOutput.Rolling -= len([]byte(out))
	nsOutput.Stdout = append(nsOutput.Stdout, out)
	return nil
}

func (nsOutput *NsOutputModule) printf(lParameters []interface{}) error {
	switch len(lParameters) {
	case 0:
		return errors.New("printf requires at least one argument")
	default:
		parameters, err := nsOutput.convertParameters(lParameters)
		if err != nil {
			return err
		}

		parameter, ok := parameters[0].(string)
		if !ok {
			return errors.New("first argument of printf must be a string")
		}

		out := fmt.Sprintf(parameter, parameters[1:]...)

		if nsOutput.Rolling-len([]byte(out)) < 0 {
			return errors.New(fmt.Sprintf("%d-byte stdout limit is exceeded", nsOutput.Limit))
		}

		nsOutput.Rolling -= len([]byte(out))
		nsOutput.Stdout = append(nsOutput.Stdout, out)
		return nil
	}
}

func (nsOutput *NsOutputModule) convertParameters(lParameters []interface{}) ([]interface{}, error) {
	var parameters []interface{}
	for _, lparameter := range lParameters {
		if parameter, err := util.FromLua(lparameter); err != nil {
			return nil, err
		} else {
			parameters = append(parameters, parameter)
		}
	}
	return parameters, nil
}

func (nsOutput *NsOutputModule) tableToCsv(table *Table) (string, error) {
	if table.Columns == nil || table.Rows == nil {
		return "", errors.New("malformed input table")
	}

	data := []string{strings.Join(table.Columns, ",")}

	for _, row := range table.Rows {
		var fields []string
		for _, field := range row {
			fields = append(fields, fmt.Sprintf("%v", field))
		}
		data = append(data, strings.Join(fields, ","))
	}

	return strings.Join(data, "\n"), nil
}

func (nsOutput *NsOutputModule) generateFromTable(data *lua.LTable, value interface{}) (string, error) {
	mapper := gluamapper.NewMapper(gluamapper.Option{NameFunc: func(str string) string { return str }})
	if err := mapper.Map(data, value); err != nil {
		return "", err
	}

	var dataType string
	switch converted := value.(type) {
	case *Value:
		dataType = "application/vnd.vz.value"
	case *Map:
		dataType = "application/vnd.vz.map"
	case *Table:
		dataType = "application/vnd.vz.table"
		for i := 0; i < len(converted.Rows); i++ {
			for j := 0; j < len(converted.Rows[i]); j++ {
				if converted.Types == nil {
					converted.Rows[i][j] = util.DataToString(converted.Rows[i][j], "string")
				} else {
					converted.Rows[i][j] = util.DataToString(converted.Rows[i][j], converted.Types[j])
				}
			}
		}
		value = converted
	default:
		return "", fmt.Errorf("Unknown data type")
	}

	output, err := nsOutput.generateOutput(dataType, value)
	if err != nil {
		return "", err
	}

	return output, nil
}

func (nsOutput *NsOutputModule) generateOutput(dataType string, data interface{}) (string, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(Output{Type: dataType, Content: data})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (nsOutput *NsOutputModule) makeErrorMessage(msg string) string {
	return NS_OUTPUT_ERROR + msg
}

func (nsOutput *NsOutputModule) recordErrorStats(timer *stats.Timer, context string) {
	if timer != nil {
		timer.Stop()
	}

	switch context {
	case _PRINT:
		ErrPrint.Incr()
	case _PRINTF:
		ErrPrintf.Incr()
	case _VALUE:
		ErrValue.Incr()
	case _VALUE_DIRECT:
		ErrValueDirect.Incr()
	case _TABLE:
		ErrTable.Incr()
	case _TABLE_DIRECT:
		ErrTableDirect.Incr()
	case _MAP:
		ErrMap.Incr()
	case _MAP_DIRECT:
		ErrMapDirect.Incr()
	case _HTML:
		ErrHTML.Incr()
	case _HTML_DIRECT:
		ErrHTMLDirect.Incr()
	case _TABLE_TO_CSV:
		ErrTableToCsv.Incr()
	}
}

func (nsOutput *NsOutputModule) error(L *lua.LState,
	err string,
	timer *stats.Timer,
	context string,
	nRetElements int) int {
	nsOutput.recordErrorStats(timer, context)
	for i := 0; i < nRetElements-1; i++ {
		L.Push(lua.LNil)
	}
	L.Push(lua.LString(nsOutput.makeErrorMessage(err)))

	return nRetElements
}

func (nsOutput *NsOutputModule) panic(err string, timer *stats.Timer, context string) {
	nsOutput.recordErrorStats(timer, context)
	panic(nsOutput.makeErrorMessage(err))
}
