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
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	dataStreamClient "github.com/verizonlabs/northstar/data/stream/client"
	dataStreamModel "github.com/verizonlabs/northstar/data/stream/model"
	dpeStreamClient "github.com/verizonlabs/northstar/dpe-stream/master/client"
	"github.com/verizonlabs/northstar/northstarapi/model"
)

// Defines the type used to support operations on NorthStar streams.
type NorthstarStreamProvider struct {
	streamDataClient *dataStreamClient.StreamClient
	streamDpeClient  *dpeStreamClient.StreamClient
}

// Returns a new NorthStar stream provider.
func NewNorthstarStreamsProvider() (*NorthstarStreamProvider, error) {
	mlog.Info("NewNorthStarJobsProvider")

	streamDataClient, err := dataStreamClient.NewStreamClient()
	if err != nil {
		return nil, err
	}

	streamDpeClient, err := dpeStreamClient.NewStreamClient()
	if err != nil {
		return nil, err
	}

	// Create the provider.
	provider := &NorthstarStreamProvider{
		streamDataClient: streamDataClient,
		streamDpeClient:  streamDpeClient,
	}

	return provider, nil
}

func (provider *NorthstarStreamProvider) ListStreams(accountId string) ([]model.Stream, *management.Error) {
	mlog.Info("ListStreams")

	externalJobs, mErr := provider.streamDataClient.GetJobs(accountId)
	if mErr != nil {
		return nil, mErr
	}

	streams := []model.Stream{}
	for _, externalJob := range externalJobs {
		streams = append(streams, *provider.fromExternalStream(externalJob))
	}

	return streams, nil
}

func (provider *NorthstarStreamProvider) GetStream(accountId string, jobId string) (*model.Stream, *management.Error) {
	mlog.Info("GetStream")

	externalJob, mErr := provider.streamDataClient.GetJob(accountId, jobId)
	if mErr != nil {
		return nil, mErr
	}

	return provider.fromExternalStream(externalJob), nil
}

func (provider *NorthstarStreamProvider) RemoveStream(accountId string, jobId string) *management.Error {
	mlog.Info("RemoveStream")

	mErr := provider.streamDpeClient.DeleteJob(accountId, jobId)
	if mErr != nil {
		return mErr
	}
	return nil
}

func (provider *NorthstarStreamProvider) fromExternalStream(externalJob *dataStreamModel.JobData) *model.Stream {
	mlog.Info("fromExternalStream")

	stream := &model.Stream{
		Id:          externalJob.Id,
		ExecutionId: externalJob.InvocationId,
		Memory:      externalJob.Memory,
		Source: model.Source{
			Name:       externalJob.Source.Name,
			Connection: externalJob.Source.Connection,
		},
		CreatedOn:   externalJob.CreatedOn,
		UpdatedOn:   externalJob.UpdatedOn,
		Status:      externalJob.Status,
		ErrorDescr:  externalJob.ErrorDescr,
		Description: externalJob.Description,
	}

	var functions []model.Function
	for _, externalFunction := range externalJob.Functions {
		functions = append(functions, model.Function{
			Name:       externalFunction.Name,
			Parameters: externalFunction.Parameters,
		})
	}
	stream.Functions = functions

	return stream
}
