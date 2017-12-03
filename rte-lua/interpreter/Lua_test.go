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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/pkg/rte/repl"
	"github.com/verizonlabs/northstar/pkg/rte/rlimit"
	"testing"
)

func TestArgs(t *testing.T) {
	interpreter := NewLuaInterpreter(rlimit.MockResourceLimit{})
	params := make(map[string]interface{})
	params["param1"] = "test"

	code := `
		local output = require("nsOutput")
		function main()
			output.printf("My args are: %v", context.Args["param1"])
			return "10"
		end
	`
	input := &repl.Input{AccountId: "610140f2-2633-6e25-ef47-deda1f752cb3",
		MainFn:  "main",
		Code:    code,
		Memory:  0,
		Timeout: 1000,
		Args:    params}
	res := interpreter.DoREPL(input)
	mlog.Info("REPL test stdout: %s, result: %s", res.Stdout, res.Result)
	assert.Equal(t, "10", res.Result, "they should be equal")
}

func TestModule(t *testing.T) {
	interpreter := NewLuaInterpreter(rlimit.MockResourceLimit{})
	params := make(map[string]interface{})
	params["param1"] = "test"

	code := `
		local output = require("nsOutput")
		function main()
			output.printf("hi\n")
		end
	`
	input := &repl.Input{AccountId: "610140f2-2633-6e25-ef47-deda1f752cb3",
		MainFn:  "main",
		Code:    code,
		Memory:  0,
		Timeout: 1000,
		Args:    params}
	res := interpreter.DoREPL(input)
	mlog.Info("REPL stdout %s, result %s", res.Stdout, res.Result)
	assert.Equal(t, "hi\n", res.Stdout, "they should be equal")
}

func TestFormatter(t *testing.T) {
	interpreter := NewLuaInterpreter(rlimit.MockResourceLimit{})
	params := make(map[string]interface{})
	params["param1"] = "test"

	code := `
		local output = require("nsOutput")
		function main(args)
			output.printf("%v", package.cpath)

			local number = {
				type="int",
				value=10
			}

			return output.value(number)
		end
	`
	input := &repl.Input{AccountId: "610140f2-2633-6e25-ef47-deda1f752cb3",
		MainFn:  "main",
		Code:    code,
		Memory:  0,
		Timeout: 1000,
		Args:    params}
	res := interpreter.DoREPL(input)
	mlog.Info("REPL stdout %s, result %s", res.Stdout, res.Result)
}

func TestFastSnippet(t *testing.T) {
	interpreter := NewLuaInterpreter(rlimit.MockResourceLimit{})
	params := make(map[string]interface{})
	params["param1"] = "test"

	code := `
		function main()
			for i = 0, 10, 1 do
			end
		end
	`
	input := &repl.Input{AccountId: "610140f2-2633-6e25-ef47-deda1f752cb3",
		MainFn:  "main",
		Code:    code,
		Memory:  0,
		Timeout: 1000,
		Args:    params}
	output := interpreter.DoREPL(input)
	require.Equal(t, "FINISHED", output.Status, "should be equal")
}

func TestSlowSnippet(t *testing.T) {
	interpreter := NewLuaInterpreter(rlimit.MockResourceLimit{})
	params := make(map[string]interface{})
	params["param1"] = "test"

	code := `
		function main()
			while 1 do
			end
		end
	`
	input := &repl.Input{AccountId: "610140f2-2633-6e25-ef47-deda1f752cb3",
		MainFn:  "main",
		Code:    code,
		Memory:  0,
		Timeout: 1000,
		Args:    params}
	output := interpreter.DoREPL(input)
	require.Equal(t, "TIMED_OUT", output.Status, "should be equal")
}

func TestSnippetRuntimeError(t *testing.T) {
	interpreter := NewLuaInterpreter(rlimit.MockResourceLimit{})
	params := make(map[string]interface{})
	params["param1"] = "test"

	code := `
		local output = require("nsOutput")
		function add(a, b)
			return a + b
		end
		function main()
			output.printf("Testing Runtime Error")
			add(10)
		end
	`
	input := &repl.Input{AccountId: "610140f2-2633-6e25-ef47-deda1f752cb3",
		MainFn:  "main",
		Code:    code,
		Memory:  0,
		Timeout: 1000,
		Args:    params}
	output := interpreter.DoREPL(input)
	assert.Equal(t, "REPL_FAILED", output.Status, "should be equal")
}

func TestJsonEncode(t *testing.T) {
	interpreter := NewLuaInterpreter(rlimit.MockResourceLimit{})
	params := make(map[string]interface{})
	params["param1"] = "test"

	code := `
		local output = require("nsOutput")

		function main()
		    local raw = json.encode({some_field = 1})
		    output.print(raw)
		end
	`
	input := &repl.Input{AccountId: "610140f2-2633-6e25-ef47-deda1f752cb3",
		MainFn:  "main",
		Code:    code,
		Memory:  0,
		Timeout: 1000,
		Args:    params}
	output := interpreter.DoREPL(input)
	assert.Equal(t, "{\"some_field\":1}", output.Stdout)
}
