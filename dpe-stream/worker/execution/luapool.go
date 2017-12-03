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
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/dpe-stream/config"
	"github.com/verizonlabs/northstar/pkg/rte/repl"
	"github.com/verizonlabs/northstar/rte-lua/interpreter"
	"sync"
)

type lStatePool struct {
	m     sync.Mutex
	saved []*interpreter.State
}

func (pl *lStatePool) Get(input *repl.Input) (*interpreter.State, error) {
	pl.m.Lock()
	defer pl.m.Unlock()
	n := len(pl.saved)
	if n == 0 {
		state, err := pl.New(input)
		if err != nil {
			return nil, err
		}
		return state, nil
	}

	mlog.Debug("Returning existing state")
	x := pl.saved[n-1]
	pl.saved = pl.saved[0 : n-1]
	return x, nil
}

func (pl *lStatePool) New(input *repl.Input) (*interpreter.State, error) {
	mlog.Debug("Creating new state")

	state, err := interpreter.CreateState(input)
	if err != nil {
		return nil, err
	}

	return state, nil
}

func (pl *lStatePool) Put(L *interpreter.State) {
	pl.m.Lock()
	defer pl.m.Unlock()
	pl.saved = append(pl.saved, L)
}

func (pl *lStatePool) Shutdown() {
	for _, state := range pl.saved {
		state.Close()
	}
}

var LuaPool *lStatePool = nil

func init() {
	LuaPool = &lStatePool{
		saved: make([]*interpreter.State, 0, config.NumThreads*2),
	}
}
