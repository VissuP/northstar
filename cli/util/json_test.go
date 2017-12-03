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
	"fmt"
	"github.com/verizonlabs/northstar/data/datasets/model"
	"strings"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	columns := make(map[string]model.Column)
	columns["col1"] = model.Column{Name: "col1", DataType: "int"}
	columns["col2"] = model.Column{Name: "col2", DataType: "string"}

	tables := make(map[string]model.Table)
	tables["table1"] = model.Table{Name: "table1", Columns: columns}
	b, err := json.Marshal(tables)
	if err != nil {
		fmt.Println(err)
		return
	}

	test := string(b)
	fmt.Println(test)
	fmt.Println(strings.Replace(test, "\"", "'", -1))
}
