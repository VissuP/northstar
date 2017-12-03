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

package nsOutput

import "github.com/verizonlabs/northstar/pkg/stats"

var (
	nsOutput           = stats.New("nsOutput")
	Print              = nsOutput.NewCounter("Print")
	Printf             = nsOutput.NewCounter("Printf")
	ValueCounter       = nsOutput.NewCounter("Value")
	ValueDirectCounter = nsOutput.NewCounter("ValueDirect")
	TableCounter       = nsOutput.NewCounter("Table")
	TableDirectCounter = nsOutput.NewCounter("TableDirect")
	MapCounter         = nsOutput.NewCounter("Map")
	MapDirectCounter   = nsOutput.NewCounter("MapDirect")
	HTMLCounter        = nsOutput.NewCounter("HTML")
	HTMLDirectCounter  = nsOutput.NewCounter("HTMLDirect")
	TableToCsv         = nsOutput.NewCounter("TableToCsv")
	ErrTableToCsv      = nsOutput.NewCounter("ErrTableToCsv")
	ErrPrint           = nsOutput.NewCounter("ErrPrint")
	ErrPrintf          = nsOutput.NewCounter("ErrPrintf")
	ErrValue           = nsOutput.NewCounter("ErrValue")
	ErrValueDirect     = nsOutput.NewCounter("ErrValueDirect")
	ErrTable           = nsOutput.NewCounter("ErrTable")
	ErrTableDirect     = nsOutput.NewCounter("ErrTableDirect")
	ErrMap             = nsOutput.NewCounter("ErrMap")
	ErrMapDirect       = nsOutput.NewCounter("ErrMapDirect")
	ErrHTML            = nsOutput.NewCounter("ErrHTML")
	ErrHTMLDirect      = nsOutput.NewCounter("ErrHTMLDirect")
)
