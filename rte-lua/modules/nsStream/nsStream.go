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

package nsStream

import (
	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack"
	"github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/modules/gluamapper"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/pkg/stats"
	"github.com/verizonlabs/northstar/dpe-stream/master/client"
	"github.com/verizonlabs/northstar/dpe-stream/master/model"
)

const (
	NS_STREAM_TYPE  = "nsStream"
	NS_STREAM_ERROR = "nsStream error: "
	CREATE          = "create"
	START           = "start"
	STOP            = "stop"
	LIMIT           = "limit"
	FOREACH         = "foreach"
	FILTER          = "filter"
	MAP             = "map"
	FOLD            = "fold"
)

type NsStreamModule struct {
	AccountId    string
	InvocationId string
	Memory       uint64
}

func NewNsStreamModule(accountId string, invocationId string, memory uint64) *NsStreamModule {
	return &NsStreamModule{AccountId: accountId, InvocationId: invocationId, Memory: memory}
}

func (nsStream *NsStreamModule) Loader(L *lua.LState) int {
	api := map[string]lua.LGFunction{
		CREATE: nsStream.createApi,
	}
	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

func (nsStream *NsStreamModule) createApi(L *lua.LState) int {
	connection := make(map[string]interface{})
	var err error

	if err = gluamapper.Map(L.CheckTable(3), &connection); err != nil {
		return nsStream.error(L, err.Error(), nil, CREATE, 2)
	}

	mt := L.NewTypeMetatable(NS_STREAM_TYPE)
	methods := map[string]lua.LGFunction{
		START:   nsStream.startApi,
		STOP:    nsStream.stopApi,
		LIMIT:   nsStream.limitApi,
		FOREACH: nsStream.foreachApi,
		FILTER:  nsStream.filterApi,
		MAP:     nsStream.mapApi,
		FOLD:    nsStream.foldApi,
	}
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), methods))

	stream := L.NewUserData()

	stream.Value = &StreamJob{
		InvocationId: nsStream.InvocationId,
		Memory:       nsStream.Memory,
		Source:       Source{Name: L.CheckString(1), Connection: connection},
		Description:  L.CheckString(2)}
	L.SetMetatable(stream, L.GetTypeMetatable(NS_STREAM_TYPE))

	Create.Incr()
	L.Push(stream)
	return 1
}

func (nsStream *NsStreamModule) startApi(L *lua.LState) int {
	stream := L.CheckUserData(1)
	streamJob, ok := stream.Value.(*StreamJob)
	if !ok {
		return nsStream.error(L, "invalid stream", nil, START, 1)
	}

	externalStreamJob := &model.StreamJob{InvocationId: streamJob.InvocationId, Memory: streamJob.Memory,
		Source: model.Source{Name: streamJob.Source.Name, Connection: streamJob.Source.Connection}}
	for _, function := range streamJob.Functions {
		externalStreamJob.Functions = append(externalStreamJob.Functions,
			model.Function{Name: function.Name,
				Parameters: function.Parameters,
				Evaluator:  function.Evaluator})
	}

	streamClient, err := client.NewStreamClient()
	if err != nil {
		return nsStream.error(L, err.Error(), nil, START, 1)
	}

	jobId, mErr := streamClient.StartJob(nsStream.AccountId, externalStreamJob)
	if mErr != nil {
		mlog.Error("Failed to start job %v: %v", externalStreamJob, mErr.Error())
		return nsStream.error(L, mErr.Error(), nil, START, 1)
	}

	streamJob.JobId = jobId
	stream.Value = streamJob
	Start.Incr()
	return 0
}

func (nsStream *NsStreamModule) stopApi(L *lua.LState) int {
	stream := L.CheckUserData(1)
	streamJob, ok := stream.Value.(*StreamJob)
	if !ok {
		return nsStream.error(L, "invalid stream", nil, STOP, 1)
	}

	streamClient, err := client.NewStreamClient()
	if err != nil {
		return nsStream.error(L, err.Error(), nil, STOP, 1)
	}

	mErr := streamClient.DeleteJob(nsStream.AccountId, streamJob.JobId)
	if mErr != nil {
		mlog.Error("Failed to delete job %v: %v", streamJob.JobId, mErr.Error())
		return nsStream.error(L, mErr.Error(), nil, STOP, 1)
	}

	Stop.Incr()
	return 0

}

func (nsStream *NsStreamModule) limitApi(L *lua.LState) int {
	stream, streamJob, err := nsStream.getStream(L)
	if err != nil {
		nsStream.panic(err.Error(), nil, START)
	}

	if err := nsStream.validateChain(streamJob.Functions, LIMIT); err != nil {
		nsStream.panic(err.Error(), nil, START)
	}

	limit := L.CheckInt(2)
	if limit < 0 {
		nsStream.panic("negative limit", nil, START)
	}

	encoded, err := msgpack.Marshal(limit)
	if err != nil {
		nsStream.panic(err.Error(), nil, START)
	}

	streamJob.Functions = append(streamJob.Functions, Function{Name: LIMIT, Parameters: []interface{}{encoded}})
	stream.Value = streamJob
	L.Push(stream)
	return 1

}

func (nsStream *NsStreamModule) foreachApi(L *lua.LState) int {
	if err := nsStream.makeFunctionWithEvaluator(L, FOREACH); err != nil {
		nsStream.panic(err.Error(), nil, START)
	}
	return 1
}

func (nsStream *NsStreamModule) filterApi(L *lua.LState) int {
	if err := nsStream.makeFunctionWithEvaluator(L, FILTER); err != nil {
		nsStream.panic(err.Error(), nil, START)
	}
	return 1
}

func (nsStream *NsStreamModule) mapApi(L *lua.LState) int {
	if err := nsStream.makeFunctionWithEvaluator(L, MAP); err != nil {
		nsStream.panic(err.Error(), nil, START)
	}
	return 1
}

func (nsStream *NsStreamModule) foldApi(L *lua.LState) int {
	if err := nsStream.makeFunctionWithEvaluator(L, FOLD); err != nil {
		nsStream.panic(err.Error(), nil, START)
	}
	return 1
}

func (nsStream *NsStreamModule) getStream(L *lua.LState) (*lua.LUserData, *StreamJob, error) {
	stream := L.CheckUserData(1)
	sj, ok := stream.Value.(*StreamJob)
	if !ok {
		return nil, nil, errors.New("invalid stream")

	}

	return stream, sj, nil
}

func (nsStream *NsStreamModule) validateChain(functions []Function, functionType string) error {
	if len(functions) == 0 {
		return nil
	}

	if functionName := functions[len(functions)-1].Name; functionName == FOREACH || functionName == FOLD {
		return errors.New("invalid stream processing chain: " + functionType + " cannot be preceeded by " +
			functionName)
	}

	return nil
}

func (nsStream *NsStreamModule) makeFunctionWithEvaluator(L *lua.LState, functionType string) error {
	stream, sj, err := nsStream.getStream(L)
	if err != nil {
		return err
	}

	if err := nsStream.validateChain(sj.Functions, functionType); err != nil {
		return err
	}

	proto, err := msgpack.Marshal(lua.Proto2Container(L.CheckFunction(2).Proto))
	if err != nil {
		return err
	}

	if functionType == FOLD && L.GetTop() < 3 {
		return errors.New("fold requires an accumulator to be the second parameter to its evaluator")
	}

	function := Function{Name: functionType, Evaluator: proto}
	for i := 3; i <= L.GetTop(); i++ {
		encoded, err := msgpack.Marshal(L.CheckAny(i))
		if err != nil {
			return err
		}
		function.Parameters = append(function.Parameters, encoded)
	}

	sj.Functions = append(sj.Functions, function)
	stream.Value = sj

	L.Push(stream)
	return nil
}

func (nsStream *NsStreamModule) makeErrorMessage(msg string) string {
	return NS_STREAM_ERROR + msg
}

func (nsStream *NsStreamModule) recordErrorStats(timer *stats.Timer, context string) {
	if timer != nil {
		timer.Stop()
	}

	switch context {
	case CREATE:
		ErrCreate.Incr()
	case START:
		ErrStart.Incr()
	case STOP:
		ErrStop.Incr()
	}
}

func (nsStream *NsStreamModule) error(L *lua.LState,
	err string,
	timer *stats.Timer,
	context string,
	nRetElements int) int {
	nsStream.recordErrorStats(timer, context)
	for i := 0; i < nRetElements-1; i++ {
		L.Push(lua.LNil)
	}
	L.Push(lua.LString(nsStream.makeErrorMessage(err)))

	return nRetElements
}

func (nsStream *NsStreamModule) panic(err string, timer *stats.Timer, context string) {
	nsStream.recordErrorStats(timer, context)
	panic(nsStream.makeErrorMessage(err))
}
