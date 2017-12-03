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

package spark

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/rte-lua/modules/nsQL/compiler"
	"github.com/verizonlabs/northstar/rte-lua/modules/nsQL/compiler/spark/udfs"
	"github.com/verizonlabs/northstar/rte-lua/modules/nsQL/constants"
	"github.com/verizonlabs/northstar/rte-lua/modules/nsQL/parser"
	"github.com/verizonlabs/northstar/rte-lua/modules/nsQL/parser/ast"
	"github.com/verizonlabs/northstar/rte-lua/util"
	"strconv"
	"strings"
)

const (
	VALUE     = 0
	DATAFRAME = 1
	SUFFIX    = ";"
)

var EXECUTE_PATH = "/execute"

type SparkCompiler struct {
	SparkHostPort       string
	DataSource          *compiler.DataSource
	CassandraFetchLimit int
	State               *State
}

type State struct {
	Context    int
	Chunk      *Chunk
	References map[string]string
	Udfs       map[string]bool
}

type Chunk struct {
	Converter            int
	Variable, Code       string
	LookBack1, LookBack2 string
}

func NewSparkCompiler(sparkHostPort string, dataSource *compiler.DataSource) *SparkCompiler {
	return &SparkCompiler{SparkHostPort: sparkHostPort, DataSource: dataSource}
}

func (c *SparkCompiler) Run(code string, options *compiler.Options) (interface{}, error) {
	compiled, err := c.compile(code, options)
	if err != nil {
		return nil, err
	}

	var request interface{}
	err = json.Unmarshal([]byte(compiled), &request)
	if err != nil {
		return nil, err
	}

	resp, mErr := management.PostJSON("http://"+c.SparkHostPort, EXECUTE_PATH, request)
	if mErr != nil {
		mlog.Error("PostJSON failed, path: %v, request: %v, err: %v",
			EXECUTE_PATH, request, mErr.Error())
		return nil, errors.New(fmt.Sprintf("nsQL spark transcompiler communication error: %v",
			mErr.Error()))
	}

	body := struct{ MimeType, Result, Status, ErrorDescr string }{}
	err = json.Unmarshal(resp, &body)
	if err != nil {
		return nil, err
	}

	if body.ErrorDescr != "" {
		errM := fmt.Sprintf("nsQL spark transcompiler data access error: %v", body.ErrorDescr)
		mlog.Error(errM)
		return nil, errors.New(errM)
	}

	result := make(map[string]interface{})
	err = json.Unmarshal([]byte(body.Result), &result)
	if err != nil {
		return nil, err
	}

	var types []string
	rawTypes, _ := result["types"].([]interface{})
	for _, rawType := range rawTypes {
		t, _ := rawType.(string)
		types = append(types, t)
	}

	var rows [][]interface{}
	rawRows, _ := result["rows"].([]interface{})

	if options.ReturnTyped {
		for _, rawRow := range rawRows {
			var row []interface{}
			rawCells, _ := rawRow.([]interface{})
			for i, rawCell := range rawCells {
				stringData, _ := rawCell.(string)
				data, err := util.StringToData(stringData, types[i])
				if err != nil {
					return nil, errors.New("nsQL spark transcompiler error: " + err.Error())
				}
				row = append(row, data)
			}
			rows = append(rows, row)
		}

	} else {
		for _, rawRow := range rawRows {
			var row []interface{}
			rawCells, _ := rawRow.([]interface{})
			for i, rawCell := range rawCells {
				stringData, _ := rawCell.(string)
				data, err := util.StringToData(stringData, types[i])
				if err != nil {
					return nil, errors.New("nsQL spark transcompiler error: " + err.Error())
				}
				stringData = util.DataToString(data, types[i])
				row = append(row, stringData)
			}
			rows = append(rows, row)
		}
	}

	result["rows"] = rows
	return result, nil
}

func (c *SparkCompiler) compile(code string, options *compiler.Options) (string, error) {
	if options == nil {
		return "", errors.New("nsQL spark transcompiler error: mandatory options missing")
	}

	c.buildState()
	c.buildInitializer()

	query, err := parser.Parse(code)
	if err != nil {
		return "", err
	}

	c.CassandraFetchLimit = options.CassandraFetchLimit

	err = c.buildQuery(query)
	if err != nil {
		return "", err
	}

	return c.buildResponse(), nil
}

func (c *SparkCompiler) buildQuery(query ast.Expression) error {
	var err error = nil
	switch query.(type) {
	case *ast.SelectStatement:
		statement, _ := query.(*ast.SelectStatement)
		err = c.buildSelectStatement(statement)
	case *ast.SelectExpression:
		expression, _ := query.(*ast.SelectExpression)
		err = c.buildSelectExpression(expression)
	default:
		err = errors.New("nsQL spark transcompiler error: unsupported statement type")
	}
	return err
}

func (c *SparkCompiler) buildSelectStatement(statement *ast.SelectStatement) error {
	err := c.buildFrom(statement.From)
	if err != nil {
		return err
	}
	err = c.buildWhere(statement.Where)
	if err != nil {
		return err
	}
	err = c.buildGroupBy(statement.GroupBy)
	if err != nil {
		return err
	}
	if statement.GroupBy == nil {
		err = c.buildSelect(statement.Select)
		if err != nil {
			return err
		}
	}
	err = c.buildOrderBy(statement.OrderBy)
	if err != nil {
		return err
	}
	c.buildLimit(statement.Limit)

	return nil
}

func (c *SparkCompiler) buildSelectExpression(expression *ast.SelectExpression) error {
	err := c.buildQuery(expression.Left)
	if err != nil {
		return err
	}
	left := c.State.Chunk.Variable

	err = c.buildQuery(expression.Right)
	if err != nil {
		return err
	}
	right := c.State.Chunk.Variable

	c.State.Chunk.Converter = DATAFRAME
	c.State.Chunk.Variable = c.getNewDataframe()

	prefix := "var " + c.State.Chunk.Variable + " = " + left
	switch expression.Operator {
	case "union":
		c.State.Chunk.Code += prefix + ".unionAll(" + right + ").dropDuplicates()"
	case "union all":
		c.State.Chunk.Code += prefix + ".unionAll(" + right + ")"
	default:
		c.State.Chunk.Code += prefix + ".intersect(" + right + ")"
	}
	c.State.Chunk.Code += SUFFIX

	return nil
}

func (c *SparkCompiler) buildSelect(slct *ast.Select) error {
	var selection []string
	var qualifier string

	for _, expression := range slct.Expressions {
		build, err := c.buildExpression(expression)
		if err != nil {
			return err
		}
		if build == "" {
			return nil
		}

		alias := expression.GetAlias()
		if alias != "" {
			alias = ".as(\\\"" + alias + "\\\")"
		}

		build += alias

		selection = append(selection, build)
	}

	if slct.Qualifier == "distinct" {
		qualifier = ".dropDuplicates()"
	}

	dataframe := c.State.Chunk.Variable

	c.State.Chunk.Converter = DATAFRAME
	c.State.Chunk.Variable = c.getNewDataframe()

	c.State.Chunk.Code += "var " + c.State.Chunk.Variable + " = " + dataframe
	c.State.Chunk.Code += ".select(" + strings.Join(selection, ", ") + ")" + qualifier + SUFFIX

	return nil
}

func (c *SparkCompiler) buildFrom(from *ast.From) error {
	var err error
	for i, table := range from.Tables {
		if i == 0 {
			err = c.buildLoad(table)
		} else {
			err = c.buildJoin(from.Joins[i-1])

		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *SparkCompiler) buildWhere(where ast.Expression) error {
	if where == nil {
		return nil
	}
	filter, err := c.buildExpression(where)
	if err != nil {
		return err
	}
	if filter != "" {
		dataframe := c.State.Chunk.Variable
		c.State.Chunk.Converter = DATAFRAME
		c.State.Chunk.Variable = c.getNewDataframe()

		c.State.Chunk.Code += "var " + c.State.Chunk.Variable + " = " + dataframe
		c.State.Chunk.Code += ".filter(" + filter + ")" + SUFFIX
	}
	return nil
}

func (c *SparkCompiler) buildGroupBy(groupBy *ast.GroupBy) error {
	var expressions, aggregators []string

	if groupBy == nil {
		return nil
	}

	for _, expression := range groupBy.Expressions {
		build, err := c.buildExpression(expression)
		if err != nil {
			return err
		}

		alias := expression.GetAlias()
		if alias != "" {
			alias = ".as(\\\"" + alias + "\\\")"
		}

		build += alias

		expressions = append(expressions, build)
	}

	for _, expression := range groupBy.Aggregators {
		build, err := c.buildExpression(expression)
		if err != nil {
			return err
		}

		alias := expression.GetAlias()
		if alias != "" {
			alias = ".as(\\\"" + alias + "\\\")"
		}

		build += alias

		aggregators = append(aggregators, build)
	}

	dataframe := c.State.Chunk.Variable
	c.State.Chunk.Converter = DATAFRAME
	c.State.Chunk.Variable = c.getNewDataframe()

	c.State.Chunk.Code += "var " + c.State.Chunk.Variable + " = " + dataframe
	c.State.Chunk.Code += ".groupBy(" + strings.Join(expressions, ", ") + ")"
	c.State.Chunk.Code += ".agg(" + strings.Join(aggregators, ", ") + ")" + SUFFIX

	if groupBy.Having != nil {
		hBuild, err := c.buildExpression(groupBy.Having)
		if err != nil {
			return err
		}

		dataframe = c.State.Chunk.Variable
		c.State.Chunk.Converter = DATAFRAME
		c.State.Chunk.Variable = c.getNewDataframe()

		c.State.Chunk.Code += "var " + c.State.Chunk.Variable + " = " + dataframe
		c.State.Chunk.Code += ".filter(" + hBuild + ")" + SUFFIX

	}

	return nil
}

func (c *SparkCompiler) buildOrderBy(orderBy []ast.Expression) error {
	var expressions []string

	if orderBy == nil {
		return nil
	}

	for _, expression := range orderBy {
		built, err := c.buildExpression(expression)
		if err != nil {
			return err
		}
		expressions = append(expressions, built)
	}

	dataframe := c.State.Chunk.Variable
	c.State.Chunk.Converter = DATAFRAME
	c.State.Chunk.Variable = c.getNewDataframe()
	c.State.Chunk.Code += "var " + c.State.Chunk.Variable + " = " + dataframe
	c.State.Chunk.Code += ".orderBy(" + strings.Join(expressions, ", ") + ")" + SUFFIX

	return nil
}

func (c *SparkCompiler) buildLimit(limit string) {
	if limit == "" {
		return
	}

	dataframe := c.State.Chunk.Variable
	c.State.Chunk.Converter = DATAFRAME
	c.State.Chunk.Variable = c.getNewDataframe()

	c.State.Chunk.Code += "var " + c.State.Chunk.Variable + " = " + dataframe + ".limit(" + limit + ")" + SUFFIX
}

func (c *SparkCompiler) buildLoad(tbl ast.Expression) error {
	switch tbl.(type) {
	case *ast.IdentifierExpression:
		table, _ := tbl.(*ast.IdentifierExpression)
		if internal, ok := c.State.References[table.Name]; ok {
			if table.Alias != "" {
				c.State.References[table.Alias] = internal
			}
			c.State.Chunk.Variable = internal
		} else {
			c.State.Chunk.Converter = DATAFRAME
			c.State.Chunk.Variable = c.getNewDataframe()
			c.State.References[table.Name] = c.State.Chunk.Variable
			if table.Alias != "" {
				c.State.References[table.Alias] = c.State.Chunk.Variable
			}

			if c.DataSource.Protocol == constants.CASSANDRA {
				c.State.Chunk.Code += "var " + c.State.Chunk.Variable + " = "
				c.State.Chunk.Code += "sqlContext.read.format(\\\"org.apache.spark.sql.cassandra\\\")"
				c.State.Chunk.Code += ".options(Map(\\\"cluster\\\" -> \\\"ClusterOne\\\","
				c.State.Chunk.Code += "\\\"keyspace\\\" -> \\\"" + table.Owner + "\\\","
				c.State.Chunk.Code += "\\\"table\\\" -> \\\"" + table.Name + "\\\")).load()"
				if c.CassandraFetchLimit > 0 {
					c.State.Chunk.Code += ".limit(" + strconv.Itoa(c.CassandraFetchLimit) + ")"
				}
				c.State.Chunk.Code += ".cache()" + SUFFIX
				c.State.Chunk.Code += c.State.Chunk.Variable + ".count()"
				c.State.Chunk.Code += SUFFIX
			}
		}
	default:
		statement, _ := tbl.(*ast.SelectStatement)
		err := c.buildSelectStatement(statement)
		if err != nil {
			return err
		}
		if statement.Alias != "" {
			c.State.References[statement.Alias] = c.State.Chunk.Variable
		}
	}
	return nil
}

func (c *SparkCompiler) buildJoin(join *ast.Join) error {
	var on string
	left := c.State.Chunk.Variable
	err := c.buildLoad(join.Table)
	if err != nil {
		return err
	}

	if join.On != nil {
		on, err = c.buildExpression(join.On)
		if err != nil {
			return err
		}
		on = ", " + on
	}

	joinType := c.buildJoinType(join.Type)
	if joinType != "" {
		joinType = ", " + joinType
	}

	right := c.State.Chunk.Variable

	c.State.Chunk.Converter = DATAFRAME
	c.State.Chunk.Variable = c.getNewDataframe()

	c.State.Chunk.Code += "var " + c.State.Chunk.Variable + " = " + left
	c.State.Chunk.Code += ".join(" + right + on + joinType + ")" + SUFFIX

	return nil
}

func (c *SparkCompiler) buildJoinType(typ string) string {
	switch typ {
	case "inner":
		return ""
	case "left_outer", "right_outer":
		return "\\\"" + typ + "\\\""
	case "full_outer":
		return "\\\"outer\\\""
	default:
		return "\\\"leftsemi\\\""
	}
}

func (c *SparkCompiler) buildExpression(expression ast.Expression) (string, error) {
	var build string
	var err error

	switch expression.(type) {
	case *ast.IdentifierExpression:
		identifier, _ := expression.(*ast.IdentifierExpression)
		build = c.buildIdentifier(identifier)
	case *ast.LiteralExpression:
		literal, _ := expression.(*ast.LiteralExpression)
		build, err = c.buildLiteral(literal)
	case *ast.TableAggregator:
		tableAggregator, _ := expression.(*ast.TableAggregator)
		c.buildTableAggregator(tableAggregator)
	case *ast.ToColumnAggregator, *ast.ToNumericAggregator,
		*ast.ToNumericTransformer, *ast.ToTemporalTransformer, *ast.ToStringTransformer:
		build, err = c.buildFunction(expression)
	default:
		build, err = c.buildBasicExpression(expression)
	}

	return build, err
}

func (c *SparkCompiler) buildLiteral(literal *ast.LiteralExpression) (string, error) {
	var build string
	var err error

	switch literal.Token {
	case parser.STRING, parser.UUID:
		build = "lit(\\\"" + literal.Value + "\\\")"
	case parser.TIMESTAMP:
		build = "lit(java.sql.Timestamp.valueOf(\\\"" + literal.Value + "\\\"))"
	case parser.DATE:
		build = "lit(java.sql.Date.valueOf(\\\"" + literal.Value + "\\\"))"
	case parser.TIME:
		build = "lit(java.sql.Time.valueOf(\\\"" + literal.Value + "\\\").getTime()*1000000)"
	case parser.TIME_INTERVAL:
		build = "expr(\\\"" + literal.Value + "\\\")"
	default:
		build = "lit(" + literal.Value + ")"
	}

	return build, err
}

func (c *SparkCompiler) buildIdentifier(identifier *ast.IdentifierExpression) string {
	if identifier.Owner == "" {
		build := c.State.Chunk.Variable + "(\\\"" + identifier.Name + "\\\")"
		return build
	}
	build := c.State.References[identifier.Owner] + "(\\\"" + identifier.Name + "\\\")"
	return build
}

func (c *SparkCompiler) buildTableAggregator(function *ast.TableAggregator) {
	var parameters []string

	c.State.Chunk.Converter = VALUE

	for _, parameter := range function.Parameters {
		identifier, _ := parameter.(*ast.IdentifierExpression)
		parameters = append(parameters, identifier.Name)
	}
	pBuild := strings.Join(parameters, ", ")

	switch function.Name {
	case "tcount":
		c.State.Chunk.Code += "var nm = " + c.State.Chunk.Variable + ".count()"
	case "tcorr":
		c.State.Chunk.Code = "var nm = " + c.State.Chunk.Variable + "stat.corr(" + pBuild + ")"
	default:
		c.State.Chunk.Code = "var nm = " + c.State.Chunk.Variable + "stat.cov(" + pBuild + ")"
	}

	c.State.Chunk.Code += SUFFIX
}

func (c *SparkCompiler) buildFunction(fnctn ast.Expression) (string, error) {
	var name, build string
	var parameters []ast.Expression
	var pBuilds []string
	var err error

	switch fnctn.(type) {
	case *ast.ToColumnAggregator:
		function, _ := fnctn.(*ast.ToColumnAggregator)
		name = function.Name
		parameters = function.Parameters
	case *ast.ToNumericAggregator:
		function, _ := fnctn.(*ast.ToNumericAggregator)
		name = function.Name
		parameters = function.Parameters
	case *ast.ToNumericTransformer:
		function, _ := fnctn.(*ast.ToNumericTransformer)
		name = function.Name
		switch function.Name {
		case "day":
			name = "dayofmonth"
		case "subtract_timestamps":
			if _, ok := c.State.Udfs[function.Name]; !ok {
				c.State.Chunk.Code += udfs.SUBTRACT_TIMESTAMPS + SUFFIX
				c.State.Udfs[function.Name] = true
			}
		}
		parameters = function.Parameters
	case *ast.ToStringTransformer:
		function, _ := fnctn.(*ast.ToStringTransformer)
		name = function.Name
		switch function.Name {
		case "map_blob_json_fetch":
			if _, ok := c.State.Udfs[function.Name]; !ok {
				c.State.Chunk.Code += udfs.MAP_BLOB_JSON_FETCH + SUFFIX
				c.State.Udfs[function.Name] = true
			}
		case "json_fetch":
			if _, ok := c.State.Udfs[function.Name]; !ok {
				c.State.Chunk.Code += udfs.JSON_FETCH + SUFFIX
				c.State.Udfs[function.Name] = true
			}
		}
		parameters = function.Parameters
	default:
		name = "current_timestamp"
	}

	for _, parameter := range parameters {
		pBuild, err := c.buildExpression(parameter)
		if err != nil {
			return "", err
		}
		pBuilds = append(pBuilds, pBuild)
	}

	build = name + "(" + strings.Join(pBuilds, ", ") + ")"

	return build, err
}

func (c *SparkCompiler) buildBasicExpression(exprssn ast.Expression) (string, error) {
	var build string
	var err error
	switch exprssn.(type) {
	case *ast.SignedLiteralExpression:
		expression, _ := exprssn.(*ast.SignedLiteralExpression)
		right, err := c.buildExpression(expression.Right)
		if err != nil {
			return "", err
		}
		build = "(" + expression.Operator + " " + right + ")"
	case *ast.ColumnExpression:
		var left string
		expression, _ := exprssn.(*ast.ColumnExpression)
		if expression.Left != nil {
			left, err = c.buildExpression(expression.Left)
			if err != nil {
				return "", err
			}
		}
		right, err := c.buildExpression(expression.Right)
		if err != nil {
			return "", err
		}
		build = "(" + left + " " + expression.Operator + " " + right + ")"
	case *ast.NumericExpression:
		var left string
		expression, _ := exprssn.(*ast.NumericExpression)
		if expression.Left != nil {
			left, err = c.buildExpression(expression.Left)
			if err != nil {
				return "", err
			}
		}
		right, err := c.buildExpression(expression.Right)
		if err != nil {
			return "", err
		}
		build = "(" + left + " " + expression.Operator + " " + right + ")"
	case *ast.TemporalExpression:
		var left string
		expression, _ := exprssn.(*ast.TemporalExpression)
		if expression.Left != nil {
			left, err = c.buildExpression(expression.Left)
			if err != nil {
				return "", err
			}
		}
		right, err := c.buildExpression(expression.Right)
		if err != nil {
			return "", err
		}
		build = "(" + left + " " + expression.Operator + " " + right + ")"
	case *ast.ConditionalExpression:
		expression, _ := exprssn.(*ast.ConditionalExpression)
		build, err = c.buildConditionalExpression(expression)
		if err != nil {
			return "", err
		}
	default:
		expression, _ := exprssn.(*ast.LogicalExpression)
		build, err = c.buildLogicalExpression(expression)
		if err != nil {
			return "", err
		}
	}
	return build, nil
}

func (c *SparkCompiler) buildConditionalExpression(exprssn *ast.ConditionalExpression) (string, error) {
	var left, right string
	var dataframe1, dataframe2 string
	var build, iBuild string
	var err error

	if exprssn.Left != nil {
		left, err = c.buildExpression(exprssn.Left)
		if err != nil {
			return "", err
		}
	}

	dataframe1 = c.State.Chunk.Variable

	statement, ok := exprssn.Right.(*ast.SelectStatement)
	if ok {
		err = c.buildSelectStatement(statement)
		if err != nil {
			return "", err
		}
		dataframe2 = c.State.Chunk.Variable
		iBuild = c.buildIdentifier(statement.GetOutputColumns()[0])
	} else {
		right, err = c.buildExpression(exprssn.Right)
		if err != nil {
			return "", err
		}
	}

	switch exprssn.Operator {
	case "in", "not in":
		c.State.Chunk.Converter = DATAFRAME
		c.State.Chunk.Variable = c.getNewDataframe()

		c.State.Chunk.Code += "var " + c.State.Chunk.Variable + " = "
		c.State.Chunk.Code += dataframe1
		c.State.Chunk.Code += ".join(" + dataframe2 + ", " + left + "===" + iBuild + ", left_outer)"

		if exprssn.Operator == "in" {
			c.State.Chunk.Code += ".filter(" + iBuild + ".isNotNull)"
		} else {
			c.State.Chunk.Code += ".filter(" + iBuild + ".isNull)"
		}
		c.State.Chunk.Code += SUFFIX

	case "is":
		build = left + ".isNull"
	case "is not":
		build = left + ".isNotNull"
	case "==":
		build = left + " === " + right
	case "!=":
		build = left + " !== " + right
	default:
		build = left + exprssn.Operator + right
	}

	return build, nil
}

func (c *SparkCompiler) buildLogicalExpression(exprssn *ast.LogicalExpression) (string, error) {
	var build string
	var err error

	def0 := c.State.Chunk.Variable

	left, err := c.buildExpression(exprssn.Left)
	if err != nil {
		return "", err
	}

	def1 := c.State.Chunk.Variable
	c.State.Chunk.Variable = def0

	right, err := c.buildExpression(exprssn.Right)
	if err != nil {
		return "", err
	}
	def2 := c.State.Chunk.Variable

	if left == "" && right == "" {
		c.State.Chunk.Converter = DATAFRAME
		c.State.Chunk.Variable = c.getNewDataframe()

		if exprssn.Operator == "or" {
			c.State.Chunk.Code += "var " + c.State.Chunk.Variable + " = "
			c.State.Chunk.Code += def1 + ".union(" + def2 + ").dropDuplicates()"
		} else {
			c.State.Chunk.Code += "var " + c.State.Chunk.Variable + " = "
			c.State.Chunk.Code += def1 + ".intersect(" + def2 + ")"
		}
	} else if left == "" && right != "" {
		c.State.Chunk.Converter = DATAFRAME
		c.State.Chunk.Variable = c.getNewDataframe()

		if exprssn.Operator == "or" {
			c.State.Chunk.Code += "var " + c.State.Chunk.Variable + " = "
			c.State.Chunk.Code += def1 + ".union(" + def0 + ".filter(" + right + ")).dropDuplicates()"
		} else {
			c.State.Chunk.Code += "var " + c.State.Chunk.Variable + " = "
			c.State.Chunk.Code += def1 + ".intersect(" + def0 + ".filter(" + right + "))"
		}
	} else if left != "" && right == "" {
		c.State.Chunk.Converter = DATAFRAME
		c.State.Chunk.Variable = c.getNewDataframe()

		if exprssn.Operator == "or" {
			c.State.Chunk.Code += "var " + c.State.Chunk.Variable + " = "
			c.State.Chunk.Code += def2 + ".union(" + def0 + ".filter(" + left + ")).dropDuplicates()"
		} else {
			c.State.Chunk.Code += "var " + c.State.Chunk.Variable + " = "
			c.State.Chunk.Code += def2 + ".intersect(" + def0 + ".filter(" + left + "))"
		}
	} else {
		operator := exprssn.Operator
		if operator == "or" {
			operator = "||"
		} else {
			operator = "&&"
		}
		build = left + " " + operator + " " + right
	}
	c.State.Chunk.Code += SUFFIX

	return build, err
}

func (c *SparkCompiler) buildState() {
	chunk := &Chunk{Converter: -1}
	c.State = &State{Context: 0, Chunk: chunk, References: make(map[string]string), Udfs: make(map[string]bool)}
}

func (c *SparkCompiler) buildInitializer() {
	if c.DataSource.Protocol == constants.CASSANDRA {
		c.State.Chunk.Code = "import org.apache.spark.sql.functions._" + SUFFIX +
			"sqlContext.setConf(\\\"ClusterOne/spark.cassandra.connection.host\\\", \\\"" +
			c.DataSource.Connection.Host + "\\\")" + SUFFIX +
			"sqlContext.setConf(\\\"ClusterOne/spark.cassandra.connection.port\\\", \\\"" +
			c.DataSource.Connection.Port + "\\\")" + SUFFIX +
			"sqlContext.setConf(\\\"ClusterOne/spark.cassandra.input.consistency.level\\\", " +
			"\\\"LOCAL_QUORUM\\\")" + SUFFIX +
			"sqlContext.setConf(\\\"ClusterOne/spark.cassandra.auth.username\\\", \\\"" +
			c.DataSource.Connection.Username + "\\\")" + SUFFIX +
			"sqlContext.setConf(\\\"ClusterOne/spark.cassandra.auth.password\\\", \\\"" +
			c.DataSource.Connection.Password + "\\\")" + SUFFIX
	}
}

func (c *SparkCompiler) buildResponse() string {
	response := "{" + c.buildConverter() + ", " + c.buildConstraint() + ", " + c.buildBody() + "}"
	return response
}

func (c *SparkCompiler) buildConverter() string {
	switch c.State.Chunk.Converter {
	case VALUE:
		return "\"converter\":\"number\""
	default:
		return "\"converter\":\"dataframe\""
	}
}

func (c *SparkCompiler) buildConstraint() string {
	limit := "\"limit\":0"
	return limit
}

func (c *SparkCompiler) buildBody() string {
	body := "\"statements\":\"" + c.State.Chunk.Code + "\""
	return body
}

func (c *SparkCompiler) getNewDataframe() string {
	c.State.Context += 1
	dataframe := "df" + strconv.Itoa(c.State.Context)
	return dataframe
}
