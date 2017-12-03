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
	uuid "github.com/satori/go.uuid"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/data/mappings/model"
)

type MappingClientMock struct{}

func (m MappingClientMock) AddMapping(accountId string,
	data *model.MappingsData) (string, *management.Error) {
	mlog.Info("Adding mapping with account ID: %s", accountId)
	return "", nil
}

func (m MappingClientMock) ListMappings(accountId string) ([]*model.MappingsData, *management.Error) {
	mlog.Info("Listing all mappings by accountId %s", accountId)
	return []*model.MappingsData{
		&model.MappingsData{},
		&model.MappingsData{},
	}, nil
}

func (m MappingClientMock) DeleteMapping(accountId string,
	mappingId string) *management.Error {
	mlog.Info("Deleting mapping by accountID %s and mapping ID %s", accountId, mappingId)
	return nil
}

func (m MappingClientMock) GetMapping(accountId string,
	mappingId string) (*model.MappingsData, *management.Error) {
	mlog.Info("Retrieving mapping by account ID %s and mapping ID %s", accountId, mappingId)
	return &model.MappingsData{
		Id:        mappingId,
		EventId:   "00d9eedd-1758-492b-b6c2-cbd31cbee74a",
		SnippetId: "20a955b4-f55b-42f2-b5c4-3fe775509570",
	}, nil
}

func (m MappingClientMock) GetMappingByEventId(accountId string,
	eventId string) (*model.MappingsData, *management.Error) {
	mlog.Info("Retrieving mapping by event ID %s and account ID %s", eventId, accountId)
	return &model.MappingsData{
		Id:      uuid.NewV4().String(),
		EventId: eventId,
	}, nil
}
