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

package client

import (
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/data/snippets/model"
)

type MockSnippetsClientFail struct{}

func (m MockSnippetsClientFail) AddSnippet(accountId string,
	add *model.SnippetData) (string, *management.Error) {
	mlog.Info("Adding snippet with account ID: %s", accountId)
	return "", management.GetInternalError("")
}

func (m MockSnippetsClientFail) DeleteSnippet(accountId string,
	snippetId string) *management.Error {
	mlog.Info("Deleting snippet with account ID %s and snippet ID %s", accountId, snippetId)
	return management.GetInternalError("")
}

func (m MockSnippetsClientFail) GetSnippets(accountId string) ([]*model.SnippetData,
	*management.Error) {
	mlog.Info("Retreiving snippets by account ID %s", accountId)
	return nil, management.GetInternalError("")
}

func (m MockSnippetsClientFail) GetSnippet(accountId string,
	snippetId string) (*model.SnippetData, *management.Error) {
	mlog.Info("Retrieving snippet by account ID %s and snippet ID %s", accountId, snippetId)
	return nil, management.GetInternalError("")
}

func (m MockSnippetsClientFail) UpdateSnippet(accountId string,
	snippetId string,
	update *model.SnippetData) *management.Error {
	mlog.Info("Updating snippet by account ID %s and snippet ID %s", accountId, snippetId)
	return management.GetInternalError("")
}
