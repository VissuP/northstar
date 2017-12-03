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
	"github.com/verizonlabs/northstar/data/invocations/model"
)

type MockInvocationClientFail struct{}

func (i MockInvocationClientFail) AddInvocation(accountId string,
	add *model.InvocationData) (string, *management.Error) {
	mlog.Info("Adding invocation")
	return "", management.GetInternalError("")
}

func (i MockInvocationClientFail) UpdateInvocation(accountId string,
	invocationId string,
	update *model.InvocationData) *management.Error {
	mlog.Info("Updating invocation")
	return management.GetInternalError("")
}

func (i MockInvocationClientFail) GetInvocation(accountId string,
	invocationId string) (*model.InvocationData, *management.Error) {
	mlog.Info("Retrieving invocation")
	return nil, management.GetInternalError("")
}

func (i MockInvocationClientFail) GetInvocationsByAccountId(accountId string,
	limit int) ([]*model.InvocationData,
	*management.Error) {
	mlog.Info("Retrieving all invocations with account ID %s", accountId)
	return nil, management.GetInternalError("")
}

func (i MockInvocationClientFail) GetInvocationResults(accountId string,
	snippetId string,
	limit int) ([]*model.InvocationData, *management.Error) {
	mlog.Info("Retreiving all invocation results with account ID %s, snippet ID %s, and limit %d",
		accountId, snippetId, limit)
	return nil, management.GetInternalError("")
}

func (i MockInvocationClientFail) DeleteInvocation(accountId string,
	invocationId string) *management.Error {
	mlog.Info("Deleting invocation with account ID %s and invocation ID %s",
		accountId, invocationId)
	return management.GetInternalError("")
}
