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

type MockSnippetsClient struct{}

func (m MockSnippetsClient) AddSnippet(accountId string,
	add *model.SnippetData) (string, *management.Error) {
	mlog.Info("Adding snippet with account ID: %s", accountId)
	return "", nil
}

func (m MockSnippetsClient) DeleteSnippet(accountId string,
	snippetId string) *management.Error {
	mlog.Info("Deleting snippet with account ID %s and snippet ID %s", accountId, snippetId)
	return nil
}

func (m MockSnippetsClient) GetSnippets(accountId string) ([]*model.SnippetData, *management.Error) {
	mlog.Info("Retreiving snippets by account ID %s", accountId)
	return []*model.SnippetData{
		&model.SnippetData{},
		&model.SnippetData{},
	}, nil
}

func (m MockSnippetsClient) GetSnippet(accountId string,
	snippetId string) (*model.SnippetData, *management.Error) {
	mlog.Info("Retrieving snippet by account ID %s and snippet ID %s", accountId, snippetId)
	return &model.SnippetData{
		Id:      snippetId,
		Name:    "snippet-name",
		Runtime: "lua",
		MainFn:  "main",
		URL:     "base64:///",
		Code:    "bG9jYWwgbnNRTCA9IHJlcXVpcmUoIm5zUUwiKQ0KDQpmdW5jdGlvbiBtYWluKCkNCiAgICBsb2NhbCBxdWVyeSA9IFtbDQogICAgICAgIElOU0VSVCBJTlRPIGFjY291bnQuaW52b2NhdGlvbnMgKGFjY291bnRpZCwgaWQsIGZpbmlzaGVkb24pDQogICAgICAgIFZBTFVFUyAgICAgICAgICAgICAgICAgICAgICAgICAgKA0KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAnOWUzYTZlNTAtMWZjOS0xMWU3LTkzYWUtOTIzNjFmMDAxOTUzJywNCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgJ2FjYTdhZTk0LTFmYzktMTFlNy05M2FlLTkyMzYxZjAwMTk1MycsDQogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICcyMDE3LTAzLTAxIDAwOjAwOjAwJw0KICAgICAgICApOw0KICAgIF1dDQogICAgbG9jYWwgc291cmNlID0gew0KICAgICAgICBQcm90b2NvbCA9ICJjYXNzYW5kcmEiLA0KICAgICAgICBIb3N0ID0gIjEwLjQ0LjkuMTEiLA0KICAgICAgICBQb3J0ID0gIjMxODMwIiwNCiAgICAgICAgQmFja2VuZCA9ICJuYXRpdmUiDQogICAgfQ0KICAgIHByb2Nlc3NRdWVyeShxdWVyeSwgc291cmNlLCB7fSkNCmVuZA0KDQpmdW5jdGlvbiBwcm9jZXNzUXVlcnkocXVlcnksIHNvdXJjZSwgb3B0aW9ucykNCiAgICBsb2NhbCByZXNwLCBlcnIgPSBuc1FMLnF1ZXJ5KHF1ZXJ5LCBzb3VyY2UsIG9wdGlvbnMpDQogICAgaWYoZXJyIH49IG5pbCkgdGhlbg0KICAgICAgICBlcnJvcihlcnIpDQogICAgZW5kDQogICAgcmV0dXJuIHJlc3ANCmVuZA==",
		Timeout: 5000,
	}, nil
}

func (m MockSnippetsClient) UpdateSnippet(accountId string,
	snippetId string,
	update *model.SnippetData) *management.Error {
	mlog.Info("Updating snippet by account ID %s and snippet ID %s", accountId, snippetId)
	return nil
}
