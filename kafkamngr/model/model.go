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

import "errors"

type Topic struct {
	Name        string `json:"name,omitempty"`
	Partitions  int    `json:"partitions,omitempty"`
	Replication int    `json:"replication,omitempty"`
}

func (topic Topic) Validate() error {
	if topic.Name == "" {
		return errors.New("Name is empty")
	}

	if topic.Partitions < 1 {
		return errors.New("Partitions less than one")
	}

	if topic.Replication < 1 {
		return errors.New("Replication less than one")
	}

	return nil
}
