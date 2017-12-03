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

package nsOutput

type Output struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

type Table struct {
	Columns []string        `json:"columns"`
	Types   []string        `json:"types"`
	Rows    [][]interface{} `json:"rows"`
}

type Value struct {
	Type       string `json:"type"`
	Value      string `json:"value"`
	Semantic   string `json:"semantic"`
	Background string `json:"background"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Item struct {
	Label     string      `json:"label"`
	Locations []*Location `json:"locations"`
}

type Map struct {
	Type   string    `json:"type"`
	Center *Location `json:"center"`
	Zoom   int       `json:"zoom"`
	Items  []*Item   `json:"items"`
}
