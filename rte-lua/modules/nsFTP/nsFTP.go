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

import (
	"bytes"
	"github.com/jlaffaye/ftp"
	"github.com/yuin/gopher-lua"
	"github.com/verizonlabs/northstar/pkg/config"
)

const (
	FTP_CONN_TYPE = "FtpConnection"
)

var (
	NsFTPConnectionLimit, _ = config.GetInt("NS_FTP_CONNECTION_LIMIT", 3)
)

type NsFTPModule struct {
	Limit int
}

func NewNsFTPModule() *NsFTPModule {
	return &NsFTPModule{Limit: NsFTPConnectionLimit}
}

func (nsFTP *NsFTPModule) Loader(L *lua.LState) int {
	api := map[string]lua.LGFunction{
		"connect": nsFTP.connect,
	}
	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

func (nsFTP *NsFTPModule) connect(L *lua.LState) int {
	if nsFTP.Limit == 0 {
		return nsFTP.error(L, "connection limit exceeded", nil, "connect")
	}

	hostport := L.CheckString(1)
	conn, err := ftp.Connect(hostport)
	if err != nil {
		return nsFTP.error(L, err.Error(), nil, "connect")
	}

	nsFTP.Limit--

	mt := L.NewTypeMetatable(FTP_CONN_TYPE)
	methods := map[string]lua.LGFunction{
		"disconnect": nsFTP.disconnect,
		"login":      nsFTP.login,
		"logout":     nsFTP.logout,
		"mkdir":      nsFTP.mkdir,
		"store":      nsFTP.store,
	}

	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), methods))

	connection := L.NewUserData()
	connection.Value = conn
	L.SetMetatable(connection, L.GetTypeMetatable(FTP_CONN_TYPE))

	Connect.Incr()
	L.Push(connection)
	return 1
}

func (nsFTP *NsFTPModule) disconnect(L *lua.LState) int {
	conn := L.CheckUserData(1)
	if connection, ok := conn.Value.(*ftp.ServerConn); ok {
		if err := connection.Quit(); err != nil {
			return nsFTP.error(L, err.Error(), nil, "disconnect")
		}
		nsFTP.Limit++
		Disconnect.Incr()
		return 0
	}
	return nsFTP.error(L, "unknown connection handle", nil, "disconnect")
}

func (nsFTP *NsFTPModule) login(L *lua.LState) int {
	conn := L.CheckUserData(1)
	if connection, ok := conn.Value.(*ftp.ServerConn); ok {
		user := L.CheckString(2)
		password := L.CheckString(3)
		if err := connection.Login(user, password); err != nil {
			return nsFTP.error(L, "unable to login", nil, "login")
		}
		Login.Incr()
		return 0
	}
	return nsFTP.error(L, "unknown connection handle", nil, "login")
}

func (nsFTP *NsFTPModule) logout(L *lua.LState) int {
	conn := L.CheckUserData(1)
	if connection, ok := conn.Value.(*ftp.ServerConn); ok {
		if err := connection.Logout(); err != nil {
			return nsFTP.error(L, err.Error(), nil, "logout")
		}
		Logout.Incr()
		return 0
	}
	return nsFTP.error(L, "unknown connection handle", nil, "logout")
}

func (nsFTP *NsFTPModule) mkdir(L *lua.LState) int {
	conn := L.CheckUserData(1)
	if connection, ok := conn.Value.(*ftp.ServerConn); ok {
		path := L.CheckString(2)
		if err := connection.MakeDir(path); err != nil {
			return nsFTP.error(L, err.Error(), nil, "mkdir")
		}
		Mkdir.Incr()
		return 0
	}
	return nsFTP.error(L, "unknown connection handle", nil, "mkdir")
}

func (nsFTP *NsFTPModule) store(L *lua.LState) int {
	timer := NsFTP.NewTimer("StoreTimer")
	conn := L.CheckUserData(1)
	if connection, ok := conn.Value.(*ftp.ServerConn); ok {
		filename := L.CheckString(2)
		data := L.CheckString(3)
		if err := connection.Stor(filename, bytes.NewBufferString(data)); err != nil {
			return nsFTP.error(L, err.Error(), timer, "store")
		}
		Store.Incr()
		return 0
	}
	return nsFTP.error(L, "unknown connection handle", nil, "store")
}
