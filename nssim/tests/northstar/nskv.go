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

package northstar

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	northstarApiModel "github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/nssim/utils"
)

//KVFunctionality provides support for Redis capabilities
type KVFunctionality struct {
	*NorthstarApiBaseTest
}

//NewKVFunctionality creates a new instance of the NSQLFunctionality library
func NewKVFunctionality(nsAPIBase *NorthstarApiBaseTest) *KVFunctionality {
	return &KVFunctionality{
		NorthstarApiBaseTest: nsAPIBase,
	}
}

//KVSet sets a given value for a key
func (kv *KVFunctionality) KVSet(logs *utils.Logger, callbacks *utils.ThreadSafeMap, key string, value string, ttl int) error {
	logs.LogDebug("KVSet")

	code := `local nsKV = require("nsKV")
			local output = require("nsOutput")

			function main()
    			err = nsKV.set("` + key + `", "` + value + `", ` + fmt.Sprintf("%v", ttl) + `)
    			if err ~= nil then
       			 error(err)
				end
			end`

	if _, err := kv.executeAndValidate(logs, callbacks, code); err != nil {
		return err
	}
	return nil
}

//KVGet gets a given value for a key
func (kv *KVFunctionality) KVGet(logs *utils.Logger, callbacks *utils.ThreadSafeMap, key string) (string, error) {
	logs.LogDebug("KVGet")

	code := `local nsKV = require("nsKV")
			local output = require("nsOutput")
			function main()
    			val, err = nsKV.get("` + key + `")
    			if err ~= nil then
       			 error(err)
				end
				output.printf("%v",val)
			end`

	output, err := kv.executeAndValidate(logs, callbacks, code)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", output.ExecutionOutput), nil
}

//KVSetNx sets a key if it does not exist
func (kv *KVFunctionality) KVSetNx(logs *utils.Logger, callbacks *utils.ThreadSafeMap, key string, value string, ttl int) (bool, error) {
	logs.LogDebug("KVSetNx")
	code := `local nsKV = require("nsKV")
			local output = require("nsOutput")

			function main()
    			success, err = nsKV.setnx("` + key + `", "` + value + `", ` + fmt.Sprintf("%v", ttl) + `)
    			if err ~= nil then
       			 error(err)
				end
				output.printf("%v", success)
			end`

	output, err := kv.executeAndValidate(logs, callbacks, code)
	if err != nil {
		return false, err
	}

	success, err := strconv.ParseBool(output.ExecutionOutput)
	if err != nil {
		return false, fmt.Errorf("Could not parse setnx result.")
	}

	return success, nil
}

//KVIncr increments a number stored at a key
func (kv *KVFunctionality) KVIncr(logs *utils.Logger, callbacks *utils.ThreadSafeMap, key string) (int, error) {
	logs.LogDebug("KVIncr")
	code := `local nsKV = require("nsKV")
			local output = require("nsOutput")

			function main()
    			value, err = nsKV.incr("` + key + `")
    			if err ~= nil then
       			 error(err)
				end
				output.printf("%v", value)
			end`

	output, err := kv.executeAndValidate(logs, callbacks, code)
	if err != nil {
		return 0, err
	}

	value, err := strconv.ParseInt(output.ExecutionOutput, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Could not parse increment result.")
	}

	return int(value), nil
}

//KVDecr decrements a number stored at a key
func (kv *KVFunctionality) KVDecr(logs *utils.Logger, callbacks *utils.ThreadSafeMap, key string) (int, error) {
	logs.LogDebug("KVDecr")
	code := `local nsKV = require("nsKV")
			local output = require("nsOutput")

			function main()
    			value, err = nsKV.decr("` + key + `")
    			if err ~= nil then
       			 error(err)
				end
				output.printf("%v", value)
			end`

	output, err := kv.executeAndValidate(logs, callbacks, code)
	if err != nil {
		return 0, err
	}

	value, err := strconv.ParseInt(output.ExecutionOutput, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Could not parse decrement result.")
	}

	return int(value), nil
}

//KVExists checks if a key exists
func (kv *KVFunctionality) KVExists(logs *utils.Logger, callbacks *utils.ThreadSafeMap, key string) (bool, error) {
	logs.LogDebug("KVExists")
	code := `local nsKV = require("nsKV")
			local output = require("nsOutput")

			function main()
    			value, err = nsKV.exists("` + key + `")
    			if err ~= nil then
       			 error(err)
				end
				output.printf("%v", value)
			end`

	output, err := kv.executeAndValidate(logs, callbacks, code)
	if err != nil {
		return false, err
	}

	exists, err := strconv.ParseBool(output.ExecutionOutput)
	if err != nil {
		return false, fmt.Errorf("Could not parse exists result.")
	}

	return exists, nil
}

//KVExpire sets an expiration on a key
func (kv *KVFunctionality) KVExpire(logs *utils.Logger, callbacks *utils.ThreadSafeMap, key string, ttl int) (bool, error) {
	logs.LogDebug("KVExpire")
	code := `local nsKV = require("nsKV")
			local output = require("nsOutput")

			function main()
    			success, err = nsKV.expire("` + key + `",` + fmt.Sprintf("%d", ttl) + `)
    			if err ~= nil then
       			 error(err)
				end
				output.printf("%v", success)
			end`

	output, err := kv.executeAndValidate(logs, callbacks, code)
	if err != nil {
		return false, err
	}

	success, err := strconv.ParseBool(output.ExecutionOutput)
	if err != nil {
		return false, fmt.Errorf("Could not parse exists result.")
	}

	return success, nil
}

//KVAppend appends to a value
func (kv *KVFunctionality) KVAppend(logs *utils.Logger, callbacks *utils.ThreadSafeMap, key string, value string) (int, error) {
	logs.LogDebug("KVAppend")
	code := `local nsKV = require("nsKV")
			local output = require("nsOutput")

			function main()
    			length, err = nsKV.append("` + key + `","` + value + `")
    			if err ~= nil then
       			 error(err)
				end
				output.printf("%v", length)
			end`

	output, err := kv.executeAndValidate(logs, callbacks, code)
	if err != nil {
		return 0, err
	}

	length, err := strconv.ParseInt(output.ExecutionOutput, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("Could not parse exists result.")
	}

	return int(length), nil
}

//KVGetRange returns the substring of the string value stored at key
func (kv *KVFunctionality) KVGetRange(logs *utils.Logger, callbacks *utils.ThreadSafeMap, key string, start int, end int) (string, error) {
	logs.LogDebug("KVGetRange")
	code := `local nsKV = require("nsKV")
			local output = require("nsOutput")

			function main()
    			length, err = nsKV.getrange("` + key + `",` + fmt.Sprintf("%d", start) + `,` + fmt.Sprintf("%d", end) + `)
    			if err ~= nil then
       			 error(err)
				end
				output.printf("%v", length)
			end`

	output, err := kv.executeAndValidate(logs, callbacks, code)
	if err != nil {
		return "", err
	}

	return output.ExecutionOutput, nil
}

//KVSetRange overwrites part of the string stored at key, starting at the specified offset, for the entire length of value
func (kv *KVFunctionality) KVSetRange(logs *utils.Logger, callbacks *utils.ThreadSafeMap, key string, offset int, value string) (int, error) {
	logs.LogDebug("KVSetRange")
	code := `local nsKV = require("nsKV")
			local output = require("nsOutput")

			function main()
    			length, err = nsKV.setrange("` + key + `",` + fmt.Sprintf("%d", offset) + `,"` + value + `")
    			if err ~= nil then
       			 error(err)
				end
				output.printf("%v", length)
			end`

	output, err := kv.executeAndValidate(logs, callbacks, code)
	if err != nil {
		return 0, err
	}
	length, err := strconv.ParseInt(output.ExecutionOutput, 10, 32)
	if err != nil {
		return 0, logs.LogError("Error, could not parse length. Error: %s", err.Error())
	}

	return int(length), nil
}

//KVHGet returns all fields and values of the hash stored at key
func (kv *KVFunctionality) KVHGet(logs *utils.Logger, callbacks *utils.ThreadSafeMap, key string, field string) (string, error) {
	logs.LogDebug("KVHGet")
	code := `local nsKV = require("nsKV")
			local output = require("nsOutput")

			function main()
    			value, err = nsKV.hget("` + key + `","` + field + `")
    			if err ~= nil then
       			 error(err)
				end
				output.printf("%v", value)
			end`

	output, err := kv.executeAndValidate(logs, callbacks, code)
	if err != nil {
		return "", err
	}

	return output.ExecutionOutput, nil
}

//KVHSet sets fields in the hash stored at key to value
func (kv *KVFunctionality) KVHSet(logs *utils.Logger, callbacks *utils.ThreadSafeMap, key string, field string, value string) (bool, error) {
	logs.LogDebug("KVHSet")
	code := `local nsKV = require("nsKV")
			local output = require("nsOutput")

			function main()
    			value, err = nsKV.hset("` + key + `","` + field + `","` + value + `")
    			if err ~= nil then
       			 error(err)
				end
				output.printf("%v", value)
			end`

	output, err := kv.executeAndValidate(logs, callbacks, code)
	if err != nil {
		return false, err
	}

	updated, err := strconv.ParseBool(output.ExecutionOutput)
	if err != nil {
		return false, fmt.Errorf("could not parse kvhset result.")
	}

	return updated, nil
}

//KVHGetAll returns all fields and values of the hash stored at key
func (kv *KVFunctionality) KVHGetAll(logs *utils.Logger, callbacks *utils.ThreadSafeMap, key string) (map[string]tableRow, error) {
	logs.LogDebug("KVHGetAll")

	code := `
	local nsKV = require("nsKV")
	local output = require("nsOutput")

	function main()
  		hashes, err = nsKV.hgetall("` + key + `")
  		if err ~= nil then
    		error(err)
  		end

 		 local results = {
    		columns = {},
    		rows = {}
  		}

 		local index = 0
  		for key, value in pairs(hashes) do
   		 	table.insert(results.columns,key)
    		row = {}

    		for i = 0, index - 1 do table.insert(row, "") end
   			table.insert(row, value)
   			table.insert(results.rows,row)
    		index = index + 1
   		end

  		local out, err = output.table(results)
    	if err ~= nil then
        	error(err)
    	end

    	return out
	end`

	output, err := kv.executeAndValidate(logs, callbacks, code)
	if err != nil {
		return nil, err
	}

	if output.ExecutionResults == nil {
		return nil, logs.LogError("Error, empty execution results returned")
	}

	results, err := kv.TableToMap(logs, output.ExecutionResults, nil)
	if err != nil {
		return nil, logs.LogError("Error, can't convert results to map: %s", err.Error())
	}

	logs.LogInfo("GetAll results: %+v", results)

	return results, nil
}

//KVDel deletes the specified keys
func (kv *KVFunctionality) KVDel(logs *utils.Logger, callbacks *utils.ThreadSafeMap, keys ...string) (int, error) {
	logs.LogDebug("KVDel")

	deleteKeys := ""
	for _, key := range keys {
		deleteKeys += key + ","
	}
	deleteKeys = strings.TrimSuffix(deleteKeys, ",")

	code := `local nsKV = require("nsKV")
			local output = require("nsOutput")

			function main()
    			count, err = nsKV.del("` + deleteKeys + `")
    			if err ~= nil then
       			 error(err)
				end
				output.printf("%v", count)
			end`

	output, err := kv.executeAndValidate(logs, callbacks, code)
	if err != nil {
		return 0, err
	}

	count, err := strconv.ParseInt(output.ExecutionOutput, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Could not parse setnx result.")
	}

	return int(count), nil
}

func (kv *KVFunctionality) executeAndValidate(logs *utils.Logger, callbacks *utils.ThreadSafeMap, code string) (*northstarApiModel.Output, error) {
	response, err := kv.ExecuteCell(logs, callbacks, code)
	if err != nil {
		return nil, err
	}

	output := northstarApiModel.Output{}
	err = json.Unmarshal([]byte(*response), &output)
	if err != nil {
		return nil, err
	}

	logs.LogInfo("Execution output: %+v results: %+v", output, output.ExecutionResults)

	if output.Status != northstarApiModel.OutputSuccessStatus {
		return nil, logs.LogError("Failed to execute with error: %s.", output.StatusDescription)
	}

	return &output, nil
}
