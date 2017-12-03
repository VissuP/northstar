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
	"context"
	"github.com/yuin/gluare"
	"github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/modules/extended"
	"github.com/yuin/gopher-lua/modules/gluahttp"
	"github.com/yuin/gopher-lua/modules/gopher-luar"
	"github.com/verizonlabs/northstar/pkg/config"
	"github.com/verizonlabs/northstar/pkg/lualib"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	pkgCfg "github.com/verizonlabs/northstar/pkg/rte/config"
	"github.com/verizonlabs/northstar/pkg/rte/repl"
	"github.com/verizonlabs/northstar/rte-lua/modules/nsFTP"
	"github.com/verizonlabs/northstar/rte-lua/modules/nsObject"
	"github.com/verizonlabs/northstar/rte-lua/modules/nsOutput"
	"github.com/verizonlabs/northstar/rte-lua/modules/nsQL"
	"github.com/verizonlabs/northstar/rte-lua/modules/nsSFTP"
	"github.com/verizonlabs/northstar/rte-lua/modules/nsStream"
	"github.com/verizonlabs/northstar/rte-lua/modules/nsUtil"
	"time"
)

var (
	EnableHttp, _           = config.GetBool("ENABLE_HTTP", false)
	EnableNSQL, _           = config.GetBool("ENABLE_NSQL", false)
	EnableNSOutput, _       = config.GetBool("ENABLE_NSOUTPUT", true)
	EnableNSFTP, _          = config.GetBool("ENABLE_NSFTP", false)
	EnableNSSFTP, _         = config.GetBool("ENABLE_NSSFTP", false)
	EnableNSObject, _       = config.GetBool("ENABLE_NSOBJECT", false)
	EnableNSKV, _           = config.GetBool("ENABLE_NSKV", false)
	EnableNSStream, _       = config.GetBool("ENABLE_NSSTREAM", false)
	EnableNSUtil, _         = config.GetBool("ENABLE_NSUTIL", true)
)

type ExecutionContext struct {
	Args map[string]interface{}
}

type State struct {
	cancel   context.CancelFunc
	LuaState *lua.LState
	Output   *nsOutput.NsOutputModule
	NSQL     *nsQL.NsQLModule
}

func CreateState(input *repl.Input) (*State, error) {
	luaState := lua.NewState(lua.Options{SkipOpenLibs: true, IncludeGoStackTrace: false})

	ctx, cancel := createContext(input.Timeout)
	luaState.SetContext(ctx)
	luaState.SetGlobal("context", luar.New(luaState, ExecutionContext{Args: input.Args}))

	// Base modules
	luaextended.OpenLibs(luaState)
	lua.OpenString(luaState)
	lua.OpenMath(luaState)
	lua.OpenTable(luaState)
	lua.OpenSecureBase(luaState)
	lua.OpenSecureOs(luaState)
	lua.OpenPackage(luaState)
	luaState.PreloadModule("re", gluare.Loader)

	output := &State{LuaState: luaState,
		cancel: cancel}

	if EnableHttp {
		mlog.Debug("Loading HTTP module")
		luaState.PreloadModule("http", gluahttp.NewHttpModule(management.NewHttpClient()).Loader)
	}

	if EnableNSQL {
		mlog.Debug("Loading nsQL module")
		output.NSQL = nsQL.NewNSQLModule()
		luaState.PreloadModule("nsQL", output.NSQL.Loader)
	}

	if EnableNSOutput {
		mlog.Debug("Loading nsOutput module")
		output.Output = nsOutput.NewNsOutputModule()
		luaState.PreloadModule("nsOutput", output.Output.Loader)
	}

	if EnableNSFTP {
		mlog.Debug("Loading FTP module")
		luaState.PreloadModule("nsFTP", nsFTP.NewNsFTPModule().Loader)
	}

	if EnableNSSFTP {
		mlog.Debug("Loading SFTP module")
		luaState.PreloadModule("nsSFTP", nsSFTP.NewNsSFTPModule().Loader)
	}

	if EnableNSObject {
		mlog.Debug("Loading nsObject module")
		nsObjectModule, err := nsObject.NewNsObjectModule(input.AccountId)
		if err != nil {
			return nil, err
		}
		luaState.PreloadModule("nsObject", nsObjectModule.Loader)
	}

	if EnableNSKV {
		mlog.Debug("Loading nsKV module")
		redisCluster, err := pkgCfg.CreateRedisCluster()
		if err != nil {
			return nil, err
		}

		options := &lualib.Options{KeyPrefix: input.AccountId}
		luaState.PreloadModule("nsKV", lualib.NewRedisModuleWithOptions(redisCluster, options).Loader)
	}

	if EnableNSStream {
		mlog.Debug("Loading nsStream module")
		luaState.PreloadModule("nsStream", nsStream.NewNsStreamModule(input.AccountId,
			input.InvocationId,
			input.Memory).Loader)
	}

	if EnableNSUtil {
		mlog.Debug("Loading nsUtil module")
		luaState.PreloadModule("nsUtil", nsUtil.NewNsUtilModule().Loader)
	}

	return output, nil
}

func createContext(timeout int) (context.Context, context.CancelFunc) {
	if timeout > 0 {
		return context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	}

	return context.WithCancel(context.Background())
}

func (s *State) Close() {
	s.Clean()
	if s.cancel != nil {
		mlog.Debug("Canceling context")
		s.cancel()
	}

	if s.LuaState != nil {
		s.LuaState.Close()
	}
}

func (s *State) Clean() {
	if s.Output != nil {
		s.Output.Reset()
	}

	if s.NSQL != nil {
		s.NSQL.Reset()
	}
}
