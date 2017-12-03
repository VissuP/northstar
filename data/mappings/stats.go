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

package mappings

import "github.com/verizonlabs/northstar/pkg/stats"

var (
	s                   = stats.New("mappingsdata")
	InsertMapping       = s.NewCounter("InsertMapping")
	GetMapping          = s.NewCounter("GetMapping")
	GetMappings         = s.NewCounter("GetMappings")
	GetMappingByEventId = s.NewCounter("GetMappingByEventId")
	DelMapping          = s.NewCounter("DelMapping")

	ErrInsertMapping       = s.NewCounter("ErrInsertMapping")
	ErrGetMapping          = s.NewCounter("ErrGetMapping")
	ErrGetMappings         = s.NewCounter("ErrGetMappings")
	ErrGetMappingByEventId = s.NewCounter("ErrGetMappingByEventId")
	ErrDelMapping          = s.NewCounter("ErrDelMapping")
)
