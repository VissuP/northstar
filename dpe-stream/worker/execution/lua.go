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

package execution

import (
	"errors"
	"github.com/yuin/gopher-lua"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/dpe-stream/master/cluster"
	"github.com/verizonlabs/northstar/dpe-stream/worker/events"
	"github.com/verizonlabs/northstar/pkg/rte/repl"
	"github.com/verizonlabs/northstar/rte-lua/util"
	"strings"
	"time"
)

const (
	LIMIT   = "limit"
	FOREACH = "foreach"
	FILTER  = "filter"
	MAP     = "map"
	FOLD    = "fold"
)

type LuaExecutor struct {
	eventsProducer events.EventsProducer
}

func NewLuaExecution(eventsProducer events.EventsProducer) (*LuaExecutor, error) {
	return &LuaExecutor{eventsProducer: eventsProducer}, nil
}

func (e *LuaExecutor) ExecuteJob(message []byte, job *cluster.StartJob) (bool, error) {
	var data interface{}
	var err error

	start := time.Now()

	input := &repl.Input{AccountId: job.AccountId,
		InvocationId: job.InvocationId}

	state, err := LuaPool.Get(input)
	if err != nil {
		return false, err
	}
	defer state.Clean()
	defer LuaPool.Put(state)

	data, err = util.ToLua(state.LuaState, message)
	if err != nil {
		return false, err
	}

	functions := job.Functions
L:
	for i := 0; i < len(functions); i++ {
		switch functions[i].Name {
		case LIMIT:
			data, err = e.executeLimit(data, functions[i].Parameters)
			if data == nil {
				return true, nil
			}
		case FOREACH:
			err = e.executeForeach(state.LuaState,
				data,
				functions[i].Evaluator,
				functions[i].Parameters)
		case FILTER:
			data, err = e.executeFilter(state.LuaState,
				data,
				functions[i].Evaluator,
				functions[i].Parameters)
			if data == nil {
				break L
			}
		case MAP:
			data, err = e.executeMap(state.LuaState,
				data,
				functions[i].Evaluator,
				functions[i].Parameters)
		case FOLD:
			err = e.executeFold(state.LuaState,
				data,
				functions[i].Evaluator,
				functions[i].Parameters)
		default:
			err = errors.New("unknown stream function " + functions[i].Name)
			break L
		}
	}

	elapsed := time.Since(start)
	stdout := strings.Join(state.Output.Stdout, "")

	var stderr string
	if err != nil {
		stderr = err.Error()
	}

	mlog.Debug("Processing time: %v, stdout: %v, stderr: %v, result : %v",
		elapsed, stdout, stderr, state.Output.Result)
	err = e.eventsProducer.StreamOutput(job, stdout, stderr, state.Output.Result)
	if err != nil {
		mlog.Error("Failed to send output event: %v", err)
		return false, err
	}

	return false, nil
}

func (e *LuaExecutor) executeLimit(msg interface{}, params []interface{}) (interface{}, error) {
	mlog.Debug("Executing limit")
	if len(params) != 1 {
		return nil, errors.New("limit: wrong number of arguments, exactly one argument expected")
	}

	limit, ok := params[0].(lua.LNumber)
	if !ok {
		return nil, errors.New("limit: unknown parameter type, it has to be an integer")
	}

	intLimit := int64(limit)
	if intLimit == 0 {
		return nil, nil
	}

	params[0] = lua.LNumber(intLimit - 1)
	return msg, nil
}

func (e *LuaExecutor) executeForeach(l *lua.LState,
	msg,
	eval interface{}, params []interface{}) error {
	mlog.Debug("Executing foreach")
	if err := e.execute(l, msg, eval, params, false); err != nil {
		return err
	}

	return nil
}

func (e *LuaExecutor) executeFilter(l *lua.LState,
	msg,
	eval interface{},
	params []interface{}) (interface{}, error) {
	mlog.Debug("Executing filter")
	if err := e.execute(l, msg, eval, params, true); err != nil {
		return nil, err
	}
	ret := l.Get(-1)
	l.Pop(1)

	boolean, ok := ret.(lua.LBool)
	if !ok {
		return nil, errors.New(FILTER + " evaluator has to return a boolean value")
	}

	if bool(boolean) {
		return msg, nil
	}

	return nil, nil
}

func (e *LuaExecutor) executeMap(l *lua.LState,
	msg,
	eval interface{},
	params []interface{}) (interface{}, error) {
	mlog.Debug("Executing map")
	if err := e.execute(l, msg, eval, params, true); err != nil {
		return nil, err
	}
	ret := l.Get(-1)
	l.Pop(1)

	return ret, nil
}

func (e *LuaExecutor) executeFold(l *lua.LState,
	msg,
	eval interface{},
	params []interface{}) error {
	mlog.Debug("Executing fold")
	if err := e.execute(l, msg, eval, params, true); err != nil {
		return err
	}
	ret := l.Get(-1)
	l.Pop(1)

	params[0] = ret

	return nil
}

func (e *LuaExecutor) execute(l *lua.LState,
	msg interface{},
	eval interface{},
	params []interface{},
	ret bool) error {
	function, err := makeFunction(l, eval)
	if err != nil {
		return err
	}

	parameters, err := makeParameters(msg, params)
	if err != nil {
		return err
	}

	nRet := 0
	if ret {
		nRet = 1
	}

	fn := lua.P{Fn: function, NRet: nRet, Protect: true}
	mlog.Debug("CallByParam with fn %v, parameters %v", fn, parameters)
	err = l.CallByParam(fn, parameters...)
	if err != nil {
		return err
	}

	return nil
}

func makeFunction(l *lua.LState, eval interface{}) (*lua.LFunction, error) {
	proto, ok := eval.(*lua.FunctionProto)
	if !ok {
		return nil, errors.New("invalid evaluator")
	}

	return &lua.LFunction{Proto: proto, Env: l.Env, Upvalues: nil}, nil
}

func makeParameters(msg interface{}, params []interface{}) ([]lua.LValue, error) {
	var parameters []lua.LValue

	message, ok := msg.(lua.LValue)
	if !ok {
		return nil, errors.New("invalid parameter type")
	}
	parameters = append(parameters, message)

	for _, par := range params {
		parameter, ok := par.(lua.LValue)
		if !ok {
			return nil, errors.New("invalid parameter type")
		}
		parameters = append(parameters, parameter)
	}

	return parameters, nil
}
