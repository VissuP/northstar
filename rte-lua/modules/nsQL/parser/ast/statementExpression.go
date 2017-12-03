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

type StatementExpression struct{ ExpressionBase }

type SelectStatement struct {
	StatementExpression
	Select  *Select
	From    *From
	Where   Expression
	GroupBy *GroupBy
	OrderBy []Expression
	Limit   string
}

type DeleteStatement struct {
	StatementExpression
	From  *From
	Where Expression
}

type InsertStatement struct {
	StatementExpression
	Into    Expression
	Columns []Expression
	Values  []Expression
}

type UpdateStatement struct {
	StatementExpression
	Table  Expression
	Update []Expression
	Where  Expression
}

type CreateTableStatement struct {
	StatementExpression
	Table       Expression
	Description *TableDescription
	Directives  []Property
}

type DropTableStatement struct {
	StatementExpression
	Table Expression
}

type Select struct {
	Qualifier   string
	Expressions []Expression
}

type Join struct {
	Table Expression
	Type  string
	On    Expression
}

type From struct {
	Tables     []Expression
	Joins      []*Join
	References []string
}

type GroupBy struct {
	Expressions   []Expression
	Having        Expression
	Aggregators   []Expression
	OutputColumns []string
}

func (e *SelectStatement) GetAllColumns() []*IdentifierExpression {
	output := e.GetOutputColumns()
	output = append(output, e.GetFilterColumns()...)
	if e.GroupBy != nil {
		output = append(output, e.GetGroupByColumns()...)
	}
	return output
}

func (e *SelectStatement) GetOutputColumns() []*IdentifierExpression {
	var output []*IdentifierExpression
	for _, expression := range e.Select.Expressions {
		output = append(output, expression.GetColumns()...)
	}
	return output
}

func (e *SelectStatement) GetFilterColumns() []*IdentifierExpression {
	var output []*IdentifierExpression
	if e.Where != nil {
		output = append(output, e.Where.GetColumns()...)
	}
	return output
}

func (e *SelectStatement) GetGroupByColumns() []*IdentifierExpression {
	var output []*IdentifierExpression
	for _, expression := range e.OrderBy {
		output = append(output, expression.GetColumns()...)
	}
	return output
}

func (e *SelectStatement) GetOrderByColumns() []*IdentifierExpression {
	var output []*IdentifierExpression
	if e.OrderBy != nil {
		for _, expression := range e.OrderBy {
			output = append(output, expression.GetColumns()...)
		}
	}
	return output
}

func (e *SelectStatement) ToString() string {
	result := e.Select.ToString() + " "
	result += e.From.ToString()
	if e.Where != nil {
		result += " where " + e.Where.ToString()
	}
	result += ";"
	return result
}

func (s *Select) ToString() string {
	result := "select " + s.Qualifier + " "
	var expressions []string
	for _, expression := range s.Expressions {
		expressions = append(expressions, expression.ToString())
	}
	result += strings.Join(expressions, ", ")
	return result
}

func (f *From) ToString() string {
	result := "from "
	var tables []string
	for _, table := range f.Tables {
		tables = append(tables, table.ToString())
	}
	result += strings.Join(tables, ", ")
	return result
}

func (e *SelectStatement) Set(input Input) error {
	eInput, _ := input.(*SelectInput)
	e.Select = eInput.Select
	e.From = eInput.From
	e.Where = eInput.Where
	e.GroupBy = eInput.GroupBy
	e.OrderBy = eInput.OrderBy
	e.Limit = eInput.Limit
	e.Subquery = true
	return nil
}

func (e *SelectStatement) GetReference() string {
	if e.HasAlias() {
		return e.Alias
	}
	return e.ToString()
}

func (e *DeleteStatement) Set(input Input) error {
	eInput, _ := input.(*DeleteInput)
	e.From = eInput.From
	e.Where = eInput.Where
	return nil
}

func (e *DeleteStatement) GetReference() string {
	return ""
}

func (e *DeleteStatement) ToString() string {
	return ""
}

func (e *InsertStatement) Set(input Input) error {
	eInput, _ := input.(*InsertInput)
	e.Into = eInput.Into
	e.Columns = eInput.Columns
	e.Values = eInput.Values
	return nil
}

func (e *InsertStatement) GetReference() string {
	return ""
}

func (e *InsertStatement) ToString() string {
	return ""
}

func (e *UpdateStatement) Set(input Input) error {
	eInput, _ := input.(*UpdateInput)
	e.Table = eInput.Table
	e.Update = eInput.Update
	e.Where = eInput.Where
	return nil
}

func (e *UpdateStatement) GetReference() string {
	return ""
}

func (e *UpdateStatement) ToString() string {
	return ""
}

func (e *CreateTableStatement) Set(input Input) error {
	eInput, _ := input.(*CreateTableInput)
	e.Table = eInput.Table
	e.Description = eInput.Description
	e.Directives = eInput.Directives
	return nil
}

func (e *CreateTableStatement) GetReference() string {
	return ""
}

func (e *CreateTableStatement) ToString() string {
	return ""
}

func (e *DropTableStatement) Set(input Input) error {
	eInput, _ := input.(*DropTableInput)
	e.Table = eInput.Table
	return nil
}

func (e *DropTableStatement) GetReference() string {
	return ""
}

func (e *DropTableStatement) ToString() string {
	return ""
}
