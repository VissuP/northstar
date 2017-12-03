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

type BasicExpression struct {
	ExpressionBase
	Left, Right Expression
	Operator    string
}

func (e *BasicExpression) SetLeft(left Expression) {
	e.Left = left
}

func (e *BasicExpression) SetRight(right Expression) {
	e.Right = right
}

func (e *BasicExpression) SetOperator(operator string) {
	e.Operator = operator
}

func (e *BasicExpression) SetSubexpressions() {
	if e.Left != nil {
		e.Subexpressions = append(e.Subexpressions, e.Left)
	}
	e.Subexpressions = append(e.Subexpressions, e.Right)
}

func (e *BasicExpression) GetLeft() Expression {
	return e.Left
}

func (e *BasicExpression) GetRight() Expression {
	return e.Right
}

func (e *BasicExpression) GetOperator() string {
	return e.Operator
}

func (e *BasicExpression) Initialize(input Input) error {
	var err error
	eInput, _ := input.(*ExpressionInput)
	e.Left = eInput.Left
	e.Right = eInput.Right
	e.Operator = eInput.Operator
	e.SetSubexpressions()
	e.SetColumns()
	e.SetSubquery()
	e.Type, err = e.InferType()
	return err
}

func (e *BasicExpression) ToString() string {
	var result string
	if e.Left != nil {
		result += e.Left.ToString()
	}
	result += " " + e.Operator + " " + e.Right.ToString()
	result = "(" + result + ")"
	return result
}

func (e *BasicExpression) GetReference() string {
	if e.HasAlias() {
		return e.Alias
	}
	return e.ToString()
}
