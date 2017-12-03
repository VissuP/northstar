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

package nsQL

import (
	"github.com/pkg/errors"
	"github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/modules/gluamapper"
	"github.com/verizonlabs/northstar/pkg/stats"
	"github.com/verizonlabs/northstar/rte-lua/modules/nsQL/compiler"
	"github.com/verizonlabs/northstar/rte-lua/modules/nsQL/compiler/cassandra"
	"github.com/verizonlabs/northstar/rte-lua/util"
)

const (
	NSQL_TYPE    = "nsQL"
	NSQL_ERROR   = "nsQL error: "
	CONNECT      = "connect"
	DISCONNECT   = "disconnect"
	QUERY        = "query"
	QUERY_DIRECT = "queryDirect"
)

type NsQLModule struct {
	Compiler compiler.Compiler
}

func NewNSQLModule() *NsQLModule {
	return &NsQLModule{}
}

func (nsQL *NsQLModule) Loader(L *lua.LState) int {
	api := map[string]lua.LGFunction{
		CONNECT: nsQL.connect,
		QUERY:   nsQL.queryDirect,
	}
	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

func (nsQL *NsQLModule) Reset() {
	if nsQL.Compiler == nil {
		return
	}

	if compiler, ok := nsQL.Compiler.(*cassandra.CassandraCompiler); ok {
		compiler.Disconnect()
	}
}

func (nsQL *NsQLModule) connect(L *lua.LState) int {
	timer := NsQL.NewTimer("ConnectTimer")
	var source compiler.Source
	if err := gluamapper.Map(L.CheckTable(1), &source); err != nil {
		return nsQL.error(L, err.Error(), timer, CONNECT, 2)
	}

	processing := &compiler.Processing{
		Backend: source.Backend,
		DataSource: &compiler.DataSource{
			Protocol: source.Protocol,
			Connection: &compiler.Connection{
				Host:     source.Host,
				Port:     source.Port,
				Username: source.Username,
				Password: source.Password,
				Version:  source.Version,
			},
		},
	}

	compiler, err := getCompiler(processing)
	if err != nil {
		return nsQL.error(L, err.Error(), timer, CONNECT, 2)
	}

	cassandraCompiler, ok := compiler.(*cassandra.CassandraCompiler)
	if !ok {
		return nsQL.error(L, "invalid backend or protocol", timer, CONNECT, 2)
	}
	nsQL.Compiler = cassandraCompiler

	cassandraCompiler.Connect()

	mt := L.NewTypeMetatable(NSQL_TYPE)
	methods := map[string]lua.LGFunction{
		DISCONNECT: nsQL.disconnect,
		QUERY:      nsQL.query,
	}
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), methods))

	ql := L.NewUserData()
	ql.Value = cassandraCompiler
	L.SetMetatable(ql, L.GetTypeMetatable(NSQL_TYPE))

	L.Push(ql)
	timer.Stop()
	Connect.Incr()
	return 1
}

func (nsQL *NsQLModule) disconnect(L *lua.LState) int {
	ql := L.CheckUserData(1)
	cassandraCompiler, ok := ql.Value.(*cassandra.CassandraCompiler)
	if !ok {
		return nsQL.error(L, "invalid transcompiler", nil, DISCONNECT, 1)
	}
	cassandraCompiler.Disconnect()
	Disconnect.Incr()
	return 0
}

func (nsQL *NsQLModule) query(L *lua.LState) int {
	timer := NsQL.NewTimer("QueryTimer")
	ql := L.CheckUserData(1)
	cassandraCompiler, ok := ql.Value.(*cassandra.CassandraCompiler)
	if !ok {
		return nsQL.error(L, "invalid transcompiler", timer, QUERY, 2)
	}

	query := L.CheckString(2)
	options, err := nsQL.getOptions(L)
	if err != nil {
		return nsQL.error(L, err.Error(), timer, QUERY, 2)
	}

	response, err := cassandraCompiler.Run(query, options)
	if err != nil {
		return nsQL.error(L, err.Error(), timer, QUERY, 2)
	}

	converted, err := util.ToLua(L, response)
	if err != nil {
		return nsQL.error(L, err.Error(), timer, QUERY, 2)
	}

	L.Push(converted)
	timer.Stop()
	Query.Incr()
	return 1
}

func (nsQL *NsQLModule) queryDirect(L *lua.LState) int {
	timer := NsQL.NewTimer("QueryDirectTimer")
	query := L.CheckString(1)

	var source compiler.Source
	if err := gluamapper.Map(L.CheckTable(2), &source); err != nil {
		return nsQL.error(L, err.Error(), timer, QUERY_DIRECT, 2)
	}

	processing := &compiler.Processing{
		Backend: source.Backend,
		DataSource: &compiler.DataSource{
			Protocol: source.Protocol,
			Connection: &compiler.Connection{
				Host:     source.Host,
				Port:     source.Port,
				Username: source.Username,
				Password: source.Password,
				Version:  source.Version,
			},
		},
	}

	comp, err := getCompiler(processing)
	if err != nil {
		return nsQL.error(L, err.Error(), timer, QUERY_DIRECT, 2)
	}

	options, err := nsQL.getOptions(L)
	if err != nil {
		return nsQL.error(L, err.Error(), timer, QUERY_DIRECT, 2)
	}

	response, err := comp.Run(query, options)
	if err != nil {
		return nsQL.error(L, err.Error(), timer, QUERY_DIRECT, 2)
	}

	converted, err := util.ToLua(L, response)
	if err != nil {
		return nsQL.error(L, err.Error(), timer, QUERY_DIRECT, 2)
	}

	L.Push(converted)
	timer.Stop()
	QueryDirect.Incr()
	return 1
}

func (nsQL *NsQLModule) getOptions(L *lua.LState) (*compiler.Options, error) {
	opts := L.CheckAny(3)
	var options compiler.Options
	switch opts.(type) {
	case *lua.LNilType:
	case *lua.LTable:
		if err := gluamapper.Map(L.CheckTable(3), &options); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unknown input type")
	}

	return &options, nil
}

func (nsQL *NsQLModule) makeErrorMessage(msg string) string {
	return NSQL_ERROR + msg
}

func (nsQL *NsQLModule) recordErrorStats(timer *stats.Timer, context string) {
	if timer != nil {
		timer.Stop()
	}

	switch context {
	case CONNECT:
		ErrConnect.Incr()
	case DISCONNECT:
		ErrDisconnect.Incr()
	case QUERY:
		ErrQuery.Incr()
	case QUERY_DIRECT:
		ErrQueryDirect.Incr()
	}
}

func (nsQL *NsQLModule) error(L *lua.LState,
	err string,
	timer *stats.Timer,
	context string,
	nRetElements int) int {
	nsQL.recordErrorStats(timer, context)
	for i := 0; i < nRetElements-1; i++ {
		L.Push(lua.LNil)
	}
	L.Push(lua.LString(nsQL.makeErrorMessage(err)))

	return nRetElements
}

func (nsQL *NsQLModule) panic(err string, timer *stats.Timer, context string) {
	nsQL.recordErrorStats(timer, context)
	panic(nsQL.makeErrorMessage(err))
}
