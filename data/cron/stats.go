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

package cron

import "github.com/verizonlabs/northstar/pkg/stats"

var (
	s          = stats.New("crondata")
	InsertJob  = s.NewCounter("InsertJob")
	GetJob     = s.NewCounter("GetJob")
	GetJobs    = s.NewCounter("GetJobs")
	GetAllJobs = s.NewCounter("GetAllJobs")
	UpdateJob  = s.NewCounter("UpdateJob")
	DelJob     = s.NewCounter("DelJob")

	ErrInsertJob  = s.NewCounter("ErrInsertJob")
	ErrGetJob     = s.NewCounter("ErrGetJob")
	ErrGetJobs    = s.NewCounter("ErrGetJobs")
	ErrGetAllJobs = s.NewCounter("ErrGetAllJobs")
	ErrUpdateJob  = s.NewCounter("ErrUpdateJob")
	ErrDelJob     = s.NewCounter("ErrDelJob")
)
