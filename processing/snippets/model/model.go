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

package model

import (
	"fmt"
)

type Snippet struct {
	SnippetId string  `json:"snippetId,omitempty"`
	Runtime   string  `json:"runtime,omitempty"`
	MainFn    string  `json:"mainfn,omitempty"`
	URL       string  `json:"url,omitempty"`
	Code      string  `json:"code,omitempty"`
	Timeout   int     `json:"timeout,omitempty"`
	Options   Options `json:"options,omitempty"`
}

type Options struct {
	Callback string                 `json:"callback,omitempty"`
	Memory   uint64                 `json:"memory,omitempty"`
	Args     map[string]interface{} `json:"args,omitempty"`
}

func (snippet *Snippet) Validate() error {
	if snippet.Runtime == "" {
		return fmt.Errorf("Runtime type is empty")
	}

	if snippet.MainFn == "" {
		return fmt.Errorf("MainFN is empty")
	}

	if snippet.URL == "" {
		return fmt.Errorf("URL is empty")
	}

	if snippet.Timeout <= 0 {
		return fmt.Errorf("Timeout needs to be greater than zero")
	}

	return nil
}
