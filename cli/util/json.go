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

package util

import (
	"encoding/json"
	"github.com/verizonlabs/northstar/data/datasets/model"
)

func UnmarshalString(data string) (map[string]string, error) {
	kv := make(map[string]string)
	err := json.Unmarshal([]byte(data), &kv)
	return kv, err
}

func UnmarshalInterface(data string) (map[string]interface{}, error) {
	kv := make(map[string]interface{})
	err := json.Unmarshal([]byte(data), &kv)
	return kv, err
}

func UnmarshalTables(data string) (map[string]model.Table, error) {
	kv := make(map[string]model.Table)
	err := json.Unmarshal([]byte(data), &kv)
	return kv, err
}
