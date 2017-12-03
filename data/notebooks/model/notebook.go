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

import (
	"encoding/json"
	"fmt"
	"time"
)

// Defines the type that represents a notebook.
type Notebook struct {
	Id        string    `json:"id,omitempty"`
	CreatedOn time.Time `json:"createdOn,omitempty"`
	Version   string    `json:"version,omitempty"`
	Data      string    `json:"data,omitempty"`
}

// Defines internal type for marshaling/unmarshaling.
type typeNotebook Notebook

// Helper method used to unmarshal notebook while validating required fields.
func (notebook *Notebook) UnmarshalJSON(data []byte) error {
	var value typeNotebook

	// Unmarshal to the internal type
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	*notebook = Notebook(value)

	// Validate the notebook data.
	if notebook.Data == "" {
		return fmt.Errorf("The data is missing.")
	}

	return nil
}
