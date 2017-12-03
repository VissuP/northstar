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

type Input interface {
	InputMarker()
}

type LiteralInput struct {
	Token    int
	Value    string
	Original interface{}
}

func (i *LiteralInput) InputMarker() {}

type IdentifierInput struct {
	Owner, Name string
}

func (i *IdentifierInput) InputMarker() {}

type FunctionInput struct {
	Name       string
	Parameters []Expression
}

func (i *FunctionInput) InputMarker() {}

type ExpressionInput struct {
	Left, Right Expression
	Operator    string
}

func (i *ExpressionInput) InputMarker() {}

type ConditionalExpressionInput struct {
	ExpressionInput
	Subquery bool
}

type SelectInput struct {
	Select  *Select
	From    *From
	Where   Expression
	GroupBy *GroupBy
	OrderBy []Expression
	Limit   string
}

type DeleteInput struct {
	From  *From
	Where Expression
}

type InsertInput struct {
	Into    Expression
	Columns []Expression
	Values  []Expression
}

type UpdateInput struct {
	Table  Expression
	Update []Expression
	Where  Expression
}

type CreateTableInput struct {
	Table       Expression
	Description *TableDescription
	Directives  []Property
}

type DropTableInput struct {
	Table Expression
}

func (i *SelectInput) InputMarker()      {}
func (i *DeleteInput) InputMarker()      {}
func (i *InsertInput) InputMarker()      {}
func (i *UpdateInput) InputMarker()      {}
func (i *CreateTableInput) InputMarker() {}
func (i *DropTableInput) InputMarker()   {}
