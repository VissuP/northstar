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

package model

// Defines the type that represents device hardware model information.
type Model struct {
	Id                            string  `json:"id,omitempty"`
	Name                          string  `json:"name,omitempty"`
	Description                   string  `json:"description,omitempty"`
	Dimensions                    string  `json:"dimensions,omitempty"`
	ProductModel                  string  `json:"productModel,omitempty"`
	HardwareVersion               string  `json:"hardwareVersion,omitempty"`
	Voltage                       float32 `json:"voltage,omitempty"`
	Lifetime                      int     `json:"lifetime,omitempty"`
	InitialBatteryLevel           int     `json:"initialBatteryLevel,omitempty"`
	EnableInitialBatteryLevelTest bool    `json:"enableInitialBatteryLevelTest,omitempty"`
	LowBatteryLevel               int     `json:"lowBatteryLevel,omitempty"`
	CriticalBatteryLevel          int     `json:"criticalBatteryLevel,omitempty"`
	DefaultReportingInterval      int     `json:"defaultReportingInterval,omitempty" `
	FactoryReportingInterval      int     `json:"factoryReportingInterval,omitempty" `
	Configurable                  bool    `json:"configurable,omitempty"`
	MaxNumberOfPeripherals        int     `json:"maxNumberOfPeripherals,omitempty"`
	Ttl                           int     `json:"ttl,omitempty"`
	PollingInterval               int     `json:"pollingInterval,omitempty"`
	RateLimitInterval             int     `json:"rateLimitInterval,omitempty"`
	RateLimit                     int     `json:"rateLimit,omitempty"`
	Smscapable                    bool    `json:"smscapable,omitempty"`
}


