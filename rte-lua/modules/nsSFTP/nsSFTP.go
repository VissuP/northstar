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

package nsSFTP

import (
	"github.com/pkg/sftp"
	"github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/modules/gluamapper"
	"golang.org/x/crypto/ssh"
	"github.com/verizonlabs/northstar/pkg/config"
	"github.com/verizonlabs/northstar/pkg/mlog"
)

const (
	SFTP_CONN_TYPE = "SFTPConnection"
)

var (
	NsSFTPConnectionLimit, _ = config.GetInt("NS_SFTP_CONNECTION_LIMIT", 3)
)

type NsSFTPModule struct {
	Limit int
}

func NewNsSFTPModule() *NsSFTPModule {
	return &NsSFTPModule{Limit: NsSFTPConnectionLimit}
}

func (nsSFTP *NsSFTPModule) Loader(L *lua.LState) int {
	api := map[string]lua.LGFunction{
		"connect": nsSFTP.connect,
	}

	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

func (nsSFTP *NsSFTPModule) connect(L *lua.LState) int {
	if nsSFTP.Limit == 0 {
		return nsSFTP.error(L, "connection limit exceeded", nil, "")
	}

	var destination Destination
	if err := gluamapper.Map(L.CheckTable(1), &destination); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return nsSFTP.error(L, err.Error(), nil, "connect")
	}

	config := &ssh.ClientConfig{
		User:            destination.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(destination.Password),
		},
	}

	config.SetDefaults()
	sshConn, err := ssh.Dial("tcp", destination.HostPort, config)
	if err != nil {
		return nsSFTP.error(L, err.Error(), nil, "connect")
	}

	client, err := sftp.NewClient(sshConn)
	if err != nil {
		return nsSFTP.error(L, err.Error(), nil, "connect")
	}

	nsSFTP.Limit--

	mt := L.NewTypeMetatable(SFTP_CONN_TYPE)
	methods := map[string]lua.LGFunction{
		"store":      nsSFTP.store,
		"mkdir":      nsSFTP.mkdir,
		"disconnect": nsSFTP.disconnect,
	}

	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), methods))

	connection := L.NewUserData()
	connection.Value = &Clients{SSH: sshConn, SFTP: client}
	L.SetMetatable(connection, L.GetTypeMetatable(SFTP_CONN_TYPE))

	Connect.Incr()
	L.Push(connection)

	mlog.Debug("SSH connection established to: %s", destination.HostPort)
	return 1
}

func (nsSFTP *NsSFTPModule) mkdir(L *lua.LState) int {
	conn := L.CheckUserData(1)
	if clients, ok := conn.Value.(*Clients); ok {
		path := L.CheckString(2)
		if err := clients.SFTP.Mkdir(path); err != nil {
			return nsSFTP.error(L, err.Error(), nil, "mkdir")
		}
		Mkdir.Incr()
		mlog.Debug("Directory %s created", path)
		return 0
	}
	return nsSFTP.error(L, "unknown connection handle", nil, "mkdir")
}

func (nsSFTP *NsSFTPModule) store(L *lua.LState) int {
	timer := NsSFTP.NewTimer("SFTPStoreTimer")
	conn := L.CheckUserData(1)

	if clients, ok := conn.Value.(*Clients); ok {
		filename := L.CheckString(2)
		data := L.CheckString(3)

		mlog.Debug("Creating file: %s", filename)
		file, err := clients.SFTP.Create(filename)
		if err != nil {
			return nsSFTP.error(L, err.Error(), timer, "store")
		}
		defer file.Close()

		mlog.Debug("Writing data: %s", data)
		if _, err := file.Write([]byte(data)); err != nil {
			return nsSFTP.error(L, err.Error(), timer, "store")
		}

		mlog.Debug("Data written")
		Store.Incr()
		return 0
	}

	return nsSFTP.error(L, "unknown connection handle", nil, "store")
}

func (nsSFTP *NsSFTPModule) disconnect(L *lua.LState) int {
	conn := L.CheckUserData(1)
	if connection, ok := conn.Value.(*Clients); ok {
		if err := connection.SSH.Close(); err != nil {
			return nsSFTP.error(L, err.Error(), nil, "disconnect")
		}
		if err := connection.SFTP.Close(); err != nil {
			return nsSFTP.error(L, err.Error(), nil, "disconnect")
		}
		nsSFTP.Limit++
		Disconnect.Incr()
		mlog.Debug("Disconnected from server: %v", connection.SSH.RemoteAddr())
		return 0
	}

	return nsSFTP.error(L, "unknown connection handle", nil, "disconnect")
}
