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

package stats

import "github.com/verizonlabs/northstar/pkg/stats"

var (
	s             = stats.New("master")
	StartJob      = s.NewCounter("StartJob")
	StopJob       = s.NewCounter("StopJob")
	DataDeleteJob = s.NewCounter("DataDeleteJob")

	ErrDataDeleteJob      = s.NewCounter("ErrDataDeleteJob")
	ErrDataAddJob         = s.NewCounter("ErrDataAddJob")
	ErrDataUpdateJob      = s.NewCounter("ErrDataUpdateJob")
	ErrMarathonStartJob   = s.NewCounter("ErrMarathonStartJob")
	ErrMarathonStopJob    = s.NewCounter("ErrMarathonStopJob")
	ErrGetNumberOfWorkers = s.NewCounter("ErrGetNumberOfWorkers")
	ErrValidateJob        = s.NewCounter("ErrValidateJob")
	ErrBindJob            = s.NewCounter("ErrBindJob")
	ErrCheckAccountId     = s.NewCounter("ErrCheckAccountId")
	ErrCheckJobId         = s.NewCounter("ErrCheckJobId")
)
