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

package parser

import (
	"github.com/verizonlabs/northstar/rte-lua/modules/nsQL/parser/ast"
)

func makeSelectExpression(left, right ast.Expression, operator string, lexer *nsQLLex) *ast.SelectExpression {
	expression := &ast.SelectExpression{}
	setExpression(expression, &ast.ExpressionInput{Left: left, Right: right, Operator: operator}, lexer)
	return expression
}

func makeSelectStatement(slct *ast.Select, from *ast.From, where ast.Expression, groupBy *ast.GroupBy, orderBy []ast.Expression, limit string, lexer *nsQLLex) *ast.SelectStatement {
	expression := &ast.SelectStatement{}
	setExpression(expression, &ast.SelectInput{Select: slct, From: from, Where: where, GroupBy: groupBy, OrderBy: orderBy, Limit: limit}, lexer)

	if groupBy != nil {
		output := expression.Select.Expressions
		failOnNonAggregateNonReferenced(output, lexer)
		for _, field := range output {
			failOnDuplicateReference(field.GetReference(), expression.GroupBy.OutputColumns, lexer)
			expression.GroupBy.Aggregators = append(expression.GroupBy.Aggregators, field)
			expression.GroupBy.OutputColumns = append(expression.GroupBy.OutputColumns, field.GetReference())
			filter := expression.GroupBy.Having
			if filter != nil {
				for _, column := range filter.GetColumns() {
					if column.Owner != "" {
						lexer.Error("syntax error: unexpected PERIOD")
					}
					failOnUnknownReference(column.Name, expression.GroupBy.OutputColumns, lexer)
				}
			}
		}
		for _, field := range expression.GetOrderByColumns() {
			failOnUnknownReference(field.GetFullName(), expression.GroupBy.OutputColumns, lexer)
		}
	} else {
		failOnAggregateNonAggregateMix(expression.Select.Expressions, lexer)
	}
	return expression
}

func makeDeleteStatement(from *ast.From, where ast.Expression, lexer *nsQLLex) *ast.DeleteStatement {
	expression := &ast.DeleteStatement{}
	setExpression(expression, &ast.DeleteInput{From: from, Where: where}, lexer)
	return expression
}

func makeInsertStatement(into ast.Expression, columns, values []ast.Expression, lexer *nsQLLex) *ast.InsertStatement {
	expression := &ast.InsertStatement{}
	setExpression(expression, &ast.InsertInput{Into: into, Columns: columns, Values: values}, lexer)
	return expression
}

func makeUpdateStatement(table ast.Expression, update []ast.Expression, where ast.Expression, lexer *nsQLLex) *ast.UpdateStatement {
	expression := &ast.UpdateStatement{}
	setExpression(expression, &ast.UpdateInput{Table: table, Update: update, Where: where}, lexer)
	return expression
}

func makeCreateTableStatement(table ast.Expression, description *ast.TableDescription, directives []ast.Property, lexer *nsQLLex) *ast.CreateTableStatement {
	expression := &ast.CreateTableStatement{}
	setExpression(expression, &ast.CreateTableInput{Table: table, Description: description, Directives: directives}, lexer)
	return expression
}

func makeDropTableStatement(table ast.Expression, lexer *nsQLLex) *ast.DropTableStatement {
	expression := &ast.DropTableStatement{}
	setExpression(expression, &ast.DropTableInput{Table: table}, lexer)
	return expression
}

func makeGroupBy(expressions []ast.Expression, filter ast.Expression, lexer *nsQLLex) *ast.GroupBy {
	var outputColumns []string
	for _, expression := range expressions {
		switch expression.(type) {
		case *ast.IdentifierExpression:
			ie, _ := expression.(*ast.IdentifierExpression)
			if ie.Name == "*" {
				lexer.Error("syntax error: unexpected ASTERISK")
			}
			outputColumns = append(outputColumns, ie.GetReference())
		default:
			if expression.GetAlias() == "" {
				lexer.Error("syntax error: AS expected")
			}
			outputColumns = append(outputColumns, expression.GetReference())
		}
	}
	groupBy := &ast.GroupBy{Expressions: expressions, Having: filter, OutputColumns: outputColumns}
	return groupBy
}

func finalizeFrom(from *ast.From, lexer *nsQLLex) {
	for _, table := range from.Tables {
		failOnInvalidTable(table, from.References, lexer)
		from.References = append(from.References, table.GetReference())
	}
}

func makeLogicalExpression(left, right ast.Expression, operator string, lexer *nsQLLex) *ast.LogicalExpression {
	expression := &ast.LogicalExpression{}
	setExpression(expression, &ast.ExpressionInput{Left: left, Right: right, Operator: operator}, lexer)
	return expression
}

func makeConditionalExpression(left, right ast.Expression, operator string, subquery bool, lexer *nsQLLex) *ast.ConditionalExpression {
	expression := &ast.ConditionalExpression{}
	setExpression(expression, &ast.ConditionalExpressionInput{ExpressionInput: ast.ExpressionInput{Left: left, Right: right, Operator: operator}, Subquery: subquery}, lexer)
	return expression
}

func makeTemporalExpression(left, right ast.Expression, operator string, lexer *nsQLLex) *ast.TemporalExpression {
	expression := &ast.TemporalExpression{}
	setExpression(expression, &ast.ExpressionInput{Left: left, Right: right, Operator: operator}, lexer)
	return expression
}

func makeNumericExpression(left, right ast.Expression, operator string, lexer *nsQLLex) *ast.NumericExpression {
	expression := &ast.NumericExpression{}
	setExpression(expression, &ast.ExpressionInput{Left: left, Right: right, Operator: operator}, lexer)
	return expression
}

func makeColumnExpression(left, right ast.Expression, operator string, lexer *nsQLLex) *ast.ColumnExpression {
	expression := &ast.ColumnExpression{}
	setExpression(expression, &ast.ExpressionInput{Left: left, Right: right, Operator: operator}, lexer)
	return expression
}

func makeSignedLiteralExpression(left, right ast.Expression, operator string, lexer *nsQLLex) *ast.SignedLiteralExpression {
	expression := &ast.SignedLiteralExpression{}
	setExpression(expression, &ast.ExpressionInput{Left: left, Right: right, Operator: operator}, lexer)
	return expression
}

func makeTableAggregator(name string, parameters []string, lexer *nsQLLex) *ast.TableAggregator {
	fInput := &ast.FunctionInput{Name: name}
	for _, parameter := range parameters {
		fInput.Parameters = append(fInput.Parameters, &ast.IdentifierExpression{Name: parameter})
	}
	function := &ast.TableAggregator{}
	setExpression(function, fInput, lexer)
	return function
}

func makeToColumnAggregator(name string, parameters []ast.Expression, lexer *nsQLLex) *ast.ToColumnAggregator {
	function := &ast.ToColumnAggregator{}
	setExpression(function, &ast.FunctionInput{Name: name, Parameters: parameters}, lexer)
	return function
}

func makeToNumericAggregator(name string, parameters []ast.Expression, lexer *nsQLLex) *ast.ToNumericAggregator {
	function := &ast.ToNumericAggregator{}
	setExpression(function, &ast.FunctionInput{Name: name, Parameters: parameters}, lexer)
	return function
}

func makeToNumericTransformer(name string, parameters []ast.Expression, lexer *nsQLLex) *ast.ToNumericTransformer {
	function := &ast.ToNumericTransformer{}
	setExpression(function, &ast.FunctionInput{Name: name, Parameters: parameters}, lexer)
	return function
}

func makeToTemporalTransformer(name string, parameters []ast.Expression, lexer *nsQLLex) *ast.ToTemporalTransformer {
	function := &ast.ToTemporalTransformer{}
	setExpression(function, &ast.FunctionInput{Name: name, Parameters: parameters}, lexer)
	return function
}

func makeToStringTransformer(name string, parameters []ast.Expression, lexer *nsQLLex) *ast.ToStringTransformer {
	function := &ast.ToStringTransformer{}
	setExpression(function, &ast.FunctionInput{Name: name, Parameters: parameters}, lexer)
	return function
}

func makeTableName(owner, name string, lexer *nsQLLex) *ast.IdentifierExpression {
	return makeName(owner, name, lexer)
}

func makeColumnName(owner, name string, lexer *nsQLLex) *ast.IdentifierExpression {
	return makeName(owner, name, lexer)
}

func makeName(owner, name string, lexer *nsQLLex) *ast.IdentifierExpression {
	expression := &ast.IdentifierExpression{}
	setExpression(expression, &ast.IdentifierInput{Owner: owner, Name: name}, lexer)
	return expression
}

func makeLiteralExpression(token int, value string, original interface{}, lexer *nsQLLex) *ast.LiteralExpression {
	expression := &ast.LiteralExpression{}
	setExpression(expression, &ast.LiteralInput{Token: token, Value: value, Original: original}, lexer)
	return expression
}

func setExpression(expression ast.Expression, input ast.Input, lexer *nsQLLex) {
	err := expression.Set(input)
	if err != nil {
		lexer.Error(err.Error())
	}
}

func failOnNoColumnName(expressions []ast.Expression, lexer *nsQLLex) {
	for _, expression := range expressions {
		if len(expression.GetColumns()) == 0 {
			lexer.Error("syntax error: IDENTIFIER expected")
		}
	}
}

func failOnNonColumnName(expression ast.Expression, lexer *nsQLLex) {
	if _, ok := expression.(*ast.IdentifierExpression); !ok {
		lexer.Error("syntax error: IDENTIFIER expected")
	}
}

func failOnColumnExpression(expression ast.Expression, lexer *nsQLLex) {
	if _, ok := expression.(*ast.ColumnExpression); ok {
		lexer.Error("syntax error: unexpected " + expression.ToString())
	}
}

func failOnNonSelectStatement(expression ast.Expression, lexer *nsQLLex) {
	if _, ok := expression.(*ast.SelectStatement); !ok {
		lexer.Error("syntax error: IDENTIFIER expected")
	}
}

func failOnSubquery(expressions []ast.Expression, lexer *nsQLLex) {
	for _, expression := range expressions {
		if expression.GetSubquery() {
			lexer.Error("syntax error: unexpected subquery")
		}
	}
}

func failOnInvalidTable(table ast.Expression, references []string, lexer *nsQLLex) {
	alias := table.GetAlias()
	switch table.(type) {
	case *ast.IdentifierExpression:
		table, _ := table.(*ast.IdentifierExpression)
		if alias == "" {
			alias = table.GetName()
		}
	case *ast.SelectStatement:
		if alias == "" {
			lexer.Error("syntax error: AS expected")
		}
	case *ast.SelectExpression:
		lexer.Error("syntax error: unexpected LEFT_PARANTHESIS")
	}
	failOnDuplicateReference(alias, references, lexer)
}

func failOnDuplicateReference(reference string, references []string, lexer *nsQLLex) {
	for _, ref := range references {
		if reference == ref {
			lexer.Error("syntax error: duplicate reference")
		}
	}
}

func failOnUnknownReference(reference string, references []string, lexer *nsQLLex) {
	for _, ref := range references {
		if reference == ref {
			return
		}
	}
	lexer.Error("syntax error: unknown column name")
}

func failOnNonAggregateNonReferenced(expressions []ast.Expression, lexer *nsQLLex) {
	for _, expression := range expressions {
		if !expression.IsAggregate() {
			lexer.Error("syntax error: aggregate function or expression expected")
		}
		if expression.GetAlias() == "" {
			lexer.Error("syntax error: AS expected")
		}
	}
}

func failOnAggregateNonAggregateMix(expressions []ast.Expression, lexer *nsQLLex) {
	var previous bool
	for i, expression := range expressions {
		if i == 0 {
			previous = expression.IsAggregate()
		} else {
			if previous != expression.IsAggregate() {
				lexer.Error("syntax error: unexpected aggregation expression")
			}
		}
	}
}

func failOnNonTableName(expressions []ast.Expression, lexer *nsQLLex) {
	if len(expressions) != 1 {
		lexer.Error("syntax error: table expected")
	}
	if _, ok := expressions[0].(*ast.IdentifierExpression); !ok {
		lexer.Error("syntax error: table expected")
	}
}
