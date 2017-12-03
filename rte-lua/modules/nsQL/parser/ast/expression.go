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

import (
	"errors"
)

const (
	LITERAL   = 0
	AGGREGATE = 1
	COLUMN    = 1
)

type Expression interface {
	SetAlias(Alias string)
	SetType(Type int)
	SetColumns()
	SetSubquery()
	GetAlias() string
	GetSubexpressions() []Expression
	GetType() int
	GetColumns() []*IdentifierExpression
	GetSubquery() bool
	InferType() (int, error)
	IsLiteral() bool
	IsColumn() bool
	IsAggregate() bool
	HasAlias() bool
	ToString() string
	Negate()
	Set(Input Input) error
	GetReference() string
}

type ExpressionBase struct {
	Subexpressions []Expression
	Type           int
	Columns        []*IdentifierExpression
	Subquery       bool
	Alias          string
}

func (e *ExpressionBase) HasAlias() bool {
	return e.Alias != ""
}

func (e *ExpressionBase) SetAlias(alias string) {
	e.Alias = alias
}

func (e *ExpressionBase) SetType(typ int) {
	e.Type = typ
}

func (e *ExpressionBase) SetColumns() {
	for _, expression := range e.Subexpressions {
		e.Columns = append(e.Columns, expression.GetColumns()...)
	}
}

func (e *ExpressionBase) SetSubquery() {
	for _, expression := range e.Subexpressions {
		if expression.GetSubquery() {
			e.Subquery = true
			return
		}
	}
	e.Subquery = false
}

func (e *ExpressionBase) GetAlias() string {
	return e.Alias
}

func (e *ExpressionBase) GetSubexpressions() []Expression {
	return e.Subexpressions
}

func (e *ExpressionBase) GetType() int {
	return e.Type
}

func (e *ExpressionBase) GetColumns() []*IdentifierExpression {
	return e.Columns
}

func (e *ExpressionBase) GetSubquery() bool {
	return e.Subquery
}

func (e *ExpressionBase) InferType() (int, error) {
	var current int
	for i, expression := range e.Subexpressions {
		if i == 0 {
			current = expression.GetType()
		} else {
			if current == LITERAL {
				current = expression.GetType()
			} else {
				if expression.GetType() != LITERAL && current != expression.GetType() {
					return -1, errors.New("syntax error: unexpected aggregation function")
				}
			}
		}
	}
	return current, nil
}

func (e *ExpressionBase) IsLiteral() bool {
	return e.Type == LITERAL
}

func (e *ExpressionBase) IsColumn() bool {
	return e.Type == COLUMN
}

func (e *ExpressionBase) IsAggregate() bool {
	return e.Type == AGGREGATE
}

func (e *ExpressionBase) Negate() {}
