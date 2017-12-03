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
	"github.com/verizonlabs/northstar/processing/snippets/model"
)

type SnippetsClientMock struct{}

func (s SnippetsClientMock) StartSnippet(accountId string,
	snippet *model.Snippet) (string, *management.Error) {
	mlog.Info("Starting snippet with accountId %s", accountId)
	return "", nil
}

func (s SnippetsClientMock) StopSnippet(accountId string, invocationId string) *management.Error {
	mlog.Info("Stopping snippet with accountId %s and invocationId %s", accountId, invocationId)
	return nil
}
