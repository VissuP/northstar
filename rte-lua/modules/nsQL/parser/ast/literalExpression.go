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

type LiteralExpression struct {
	ExpressionBase
	Token    int
	Value    string
	Original interface{}
}

func (e *LiteralExpression) SetToken(token int) {
	e.Token = token
}

func (e *LiteralExpression) SetValue(value string) {
	e.Value = value
}

func (e *LiteralExpression) Set(input Input) error {
	eInput, _ := input.(*LiteralInput)
	e.Token = eInput.Token
	e.Value = eInput.Value
	e.Original = eInput.Original
	return nil
}

func (e *LiteralExpression) GetToken() int {
	return e.Token
}

func (e *LiteralExpression) GetValue() string {
	return e.Value
}

func (e *LiteralExpression) ToString() string {
	return e.Value
}

func (e *LiteralExpression) GetReference() string {
	if e.HasAlias() {
		return e.Alias
	}
	return e.ToString()
}
