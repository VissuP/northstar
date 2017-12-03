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

package ast

type IdentifierExpression struct {
	ExpressionBase
	Owner, Name string
}

func (e *IdentifierExpression) SetOwner(owner string) {
	e.Owner = owner
}

func (e *IdentifierExpression) SetName(name string) {
	e.Name = name
}

func (e *IdentifierExpression) Set(input Input) error {
	eInput, _ := input.(*IdentifierInput)
	e.Owner = eInput.Owner
	e.Name = eInput.Name
	e.Columns = append(e.Columns, e)
	return nil
}

func (e *IdentifierExpression) GetOwner() string {
	return e.Owner
}

func (e *IdentifierExpression) GetFullName() string {
	var result string
	if e.Owner != "" {
		result += e.Owner + "."
	}
	result += e.Name
	return result
}

func (e *IdentifierExpression) GetName() string {
	return e.Name
}

func (e *IdentifierExpression) GetReference() string {
	if e.HasAlias() {
		return e.Alias
	}
	return e.Name
}

func (e *IdentifierExpression) ToString() string {
	return e.GetFullName()
}
