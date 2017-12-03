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

package interpreter

import (
	"github.com/yuin/gopher-lua"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/pkg/rte/repl"
	"github.com/verizonlabs/northstar/pkg/rte/rlimit"
	"github.com/verizonlabs/northstar/pkg/rte/config"
	"strings"
	"time"
)

type LuaInterpreter struct {
	State  *State
	rLimit rlimit.ResourceLimit
}

func NewLuaInterpreter(rLimit rlimit.ResourceLimit) repl.Interpreter {
	return &LuaInterpreter{rLimit: rLimit}
}

func (i *LuaInterpreter) DoREPL(input *repl.Input) *repl.Output {
	timer := Lua.NewTimer("DoREPLTimer")
	startedOn := time.Now()
	mlog.Debug("Running main: %s, code: %s, args: %s, timeout: %d, memory: %v, accountId: %s, "+
		"invocationId: %s", input.MainFn, input.Code, input.Args, input.Timeout, input.Memory,
		input.AccountId, input.InvocationId)

	state, err := CreateState(input)
	if err != nil {
		mlog.Error("Failed to create state: %v", err.Error())
		timer.Stop()
		ErrDoREPL.Incr()
		return &repl.Output{StartedOn: startedOn,
			Status:     repl.STATE_CREATE_FAILED,
			ErrorDescr: err.Error()}
	}

	i.State = state
	defer state.Close()

	var rLimitErr error
	if config.EnableRLimit {
		msg, err := i.rLimit.StartMonitoring(&rlimit.Resources{Memory: input.Memory})
		if err != nil {
			mlog.Error("Failed to start resource monitoring: %v", err.Error())
			timer.Stop()
			ErrDoREPL.Incr()
			return &repl.Output{StartedOn: startedOn,
				Status:     repl.START_MONITORING_FAILED,
				ErrorDescr: err.Error()}
		}

		go func() {
		L:
			for {
				select {
				case rLimitErr = <-msg:
					mlog.Debug("Received resource limit msg: %v", rLimitErr)
					i.State.Close()
					break L
				default:
				}
			}

			mlog.Debug("Lua rlimit goroutine stopped")
		}()
	}

	err = state.LuaState.DoString(input.Code)
	if err != nil {
		mlog.Error("DoString error: %v", err)
		i.rLimit.StopMonitoring()
		timer.Stop()
		ErrDoREPL.Incr()
		execError := repl.GetExecutionError(err, rLimitErr)
		return &repl.Output{StartedOn: startedOn,
			Status:     execError.Status,
			ErrorDescr: execError.Description}
	}
	finishedOn := time.Now()
	elapsedTime := finishedOn.Sub(startedOn)

	err = state.LuaState.CallByParam(lua.P{
		Fn:      state.LuaState.GetGlobal(input.MainFn),
		NRet:    1,
		Protect: true,
	})

	stdout := strings.Join(state.Output.Stdout, "")
	if err != nil {
		mlog.Error("CallByParam error: %v", err)
		i.rLimit.StopMonitoring()
		timer.Stop()
		ErrDoREPL.Incr()
		execError := repl.GetExecutionError(err, rLimitErr)
		return &repl.Output{StartedOn: startedOn,
			FinishedOn: finishedOn,
			Stdout:     stdout,
			Status:     execError.Status,
			ErrorDescr: execError.Description}
	}

	result := lua.LVAsString(state.LuaState.Get(-1))
	state.LuaState.Pop(1)
	if state.Output.Result != "" {
		result = state.Output.Result
	}

	i.rLimit.StopMonitoring()
	timer.Stop()
	DoREPL.Incr()

	output := &repl.Output{
		StartedOn:   startedOn,
		FinishedOn:  finishedOn,
		ElapsedTime: elapsedTime,
		Stdout:      stdout,
		Result:      result,
		Status:      repl.SNIPPET_RUN_FINISHED,
		ErrorDescr:  ""}
	mlog.Debug("REPL output: %v", output)
	return output
}

func (i *LuaInterpreter) Terminate() {
	mlog.Debug("Terminating")

	if i.State != nil {
		i.State.Close()
	}
}
