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
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/vmihailenco/msgpack"
	"github.com/yuin/gopher-lua"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/rte-lua/util"
)

type Source struct {
	Name       string      `json:"name,omitempty"`
	Connection interface{} `json:"connection,omitempty"`
}

type Function struct {
	Name       string        `json:"name,omitempty"`
	Parameters []interface{} `json:"parameters,omitempty"`
	Evaluator  interface{}   `json:"evaluator,omitempty"`
}

func (function *Function) Decode() error {
	state := lua.NewState()
	defer state.Close()

	if function.Evaluator != nil {
		evaluator, ok := function.Evaluator.(string)
		if !ok {
			return errors.New("nsStream: unknown encoding for the evaluator")
		}

		decoded, err := base64.StdEncoding.DecodeString(evaluator)
		if err != nil {
			return errors.New("nsStream: unable to decode the evaluator")
		}

		var protoContainer lua.FunctionProtoContainer
		err = msgpack.Unmarshal(decoded, &protoContainer)
		if err != nil {
			return errors.New("nsStream: unable to decode the evaluator")
		}

		function.Evaluator = lua.Container2Proto(&protoContainer)
	}

	for i := 0; i < len(function.Parameters); i++ {
		parameter, ok := function.Parameters[i].(string)
		if !ok {
			return errors.New("nsStream: unknown encoding for evaluator parameter")
		}

		decoded, err := base64.StdEncoding.DecodeString(parameter)
		if err != nil {
			return errors.New("nsStream: unable to decode evaluator parameter")
		}

		var value interface{}
		err = msgpack.Unmarshal(decoded, &value)
		if err != nil {
			return errors.New("nsStream: unable to decode evaluator parameter")
		}

		lValue, err := util.ToLua(state, value)
		if err != nil {
			return errors.New("nsStream: " + err.Error())
		}

		function.Parameters[i] = lValue
	}

	return nil
}

type StreamJob struct {
	InvocationId string     `json:"invocationId,omitempty"`
	Memory       uint64     `json:"memory,omitempty"`
	Source       Source     `json:"source,omitempty"`
	Functions    []Function `json:"functions,omitempty"`
	Description  string     `json:"description,omitempty"`
}

func (j *StreamJob) Validate() error {
	if j.InvocationId == "" {
		return fmt.Errorf("Invocation id is empty")
	}

	if !isSourceSupported(j.Source.Name) {
		return fmt.Errorf("Source not supported")
	}

	if len(j.Functions) < 1 {
		return fmt.Errorf("Number of functions less than one")
	}

	return nil
}

func isSourceSupported(source string) bool {
	switch source {
	case SOURCE_KAFKA:
		return true
	default:
		mlog.Error("Unknown source selected: %v", source)
		return false
	}
}
