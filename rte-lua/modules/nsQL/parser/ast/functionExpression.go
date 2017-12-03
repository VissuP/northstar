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

import "strings"

type FunctionExpression struct {
	ExpressionBase
	Name       string
	Parameters []Expression
}

func (e *FunctionExpression) SetName(name string) {
	e.Name = name
}

func (e *FunctionExpression) SetParameters(parameters []Expression) {
	e.Parameters = parameters
}

func (e *FunctionExpression) SetSubexpressions() {
	e.Subexpressions = e.Parameters
}

func (e *FunctionExpression) GetName() string {
	return e.Name
}

func (e *FunctionExpression) GetParameters() []Expression {
	return e.Parameters
}

func (e *FunctionExpression) ToString() string {
	var expressions []string
	for _, expression := range e.Parameters {
		expressions = append(expressions, expression.ToString())
	}
	result := strings.Join(expressions, ", ")
	result = e.Name + "(" + result + ")"
	return result
}

func (e *FunctionExpression) Initialize(input Input) error {
	var err error
	eInput, _ := input.(*FunctionInput)
	e.Name = eInput.Name
	e.Parameters = eInput.Parameters
	e.SetSubexpressions()
	e.SetColumns()
	e.SetSubquery()
	e.Type, err = e.InferType()
	return err
}

func (e *FunctionExpression) GetReference() string {
	if e.HasAlias() {
		return e.Alias
	}
	return e.ToString()
}
