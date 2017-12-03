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

package cassandra

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"reflect"
	"github.com/verizonlabs/northstar/pkg/config"
	"github.com/verizonlabs/northstar/pkg/mlog"
	dktUtils "github.com/verizonlabs/northstar/pkg/utils"
	"github.com/verizonlabs/northstar/rte-lua/modules/nsQL/compiler"
	"github.com/verizonlabs/northstar/rte-lua/modules/nsQL/parser"
	"github.com/verizonlabs/northstar/rte-lua/modules/nsQL/parser/ast"
	"github.com/verizonlabs/northstar/rte-lua/util"
	"strconv"
	"strings"
	"time"
)

const (
	UNSET = iota
	GROUP
	SINGULAR
	IDENTITY            = "identity"
	JSON_FETCH          = "json_fetch"
	MAP_BLOB_JSON_FETCH = "map_blob_json_fetch"
	SUBTRACT_TIMESTAMPS = "subtract_timestamps"
	TCOUNT              = "tcount"
)

var CassandraProtoVersion, _ = config.GetInt("CASSANDRA_PROTO_VERSION", 3)

type Transformer struct {
	Columns         []string
	OriginalColumns []string
	Name            string
	Alias           string
	Function        func(...interface{}) (interface{}, error)
	Parameters      []interface{}
}

type CassandraCompiler struct {
	Connection *compiler.Connection
	Session    *gocql.Session
}

type ExecutionOptions struct {
	Scan     bool
	Keyspace string
}

func NewCassandraCompiler(connection *compiler.Connection) *CassandraCompiler {
	return &CassandraCompiler{Connection: connection}
}

func (c *CassandraCompiler) Connect() error {
	if c.Session != nil {
		errors.New("nsQL cassandra transcompiler error: connection already exists")
	}

	if err := c.setSession(); err != nil {
		return err
	}

	return nil
}

func (c *CassandraCompiler) Disconnect() {
	if c.Session != nil {
		c.Session.Close()
	}
}

func (c *CassandraCompiler) Run(query string, options *compiler.Options) (interface{}, error) {
	parsed, err := parser.Parse(query)
	if err != nil {
		return "", err
	}

	if c.Session == nil {
		if err = c.setSession(); err != nil {
			return nil, err
		}
		defer c.Session.Close()
	}

	switch parsed.(type) {
	case *ast.DropTableStatement:
		statement, _ := parsed.(*ast.DropTableStatement)
		keyspace, table, err := c.getKeyspaceTable(statement.Table)
		if err != nil {
			return nil, err
		}
		execOpts := &ExecutionOptions{Scan: false, Keyspace: keyspace}
		_, _, err = c.execute(execOpts, "DROP TABLE "+keyspace+"."+table, nil)
		if err != nil {
			return nil, err
		}
		return nil, nil
	case *ast.CreateTableStatement:
		statement, _ := parsed.(*ast.CreateTableStatement)
		keyspace, table, err := c.getKeyspaceTable(statement.Table)
		if err != nil {
			return nil, err
		}
		execOpts := &ExecutionOptions{Scan: false, Keyspace: keyspace}
		fields := []string{}
		for _, f := range statement.Description.Fields {
			fields = append(fields, f.FieldName+" "+f.FieldType)
		}

		primaryKey := "PRIMARY KEY((" + strings.Join(statement.Description.PrimaryKey.Partitioning, ", ") + ")"
		if statement.Description.PrimaryKey.Clustering != nil {
			primaryKey += ", " + strings.Join(statement.Description.PrimaryKey.Clustering, ", ")
		}
		primaryKey += ")"

		tableDescription := "(" + strings.Join(fields, ", ") + ", " + primaryKey + ")"

		directiveList := []string{}
		for _, d := range statement.Directives {
			switch directive := d.(type) {
			case *ast.ClusteringOrder:
				order := []string{}
				for _, o := range directive.Order {
					field := o.FieldName
					if o.Ascending {
						field += " ASC"
					} else {
						field += " DESC"
					}
					order = append(order, field)
				}
				directiveList = append(directiveList, "CLUSTERING ORDER BY ("+strings.Join(order, ", ")+")")
			case *ast.CompactStorage:
				directiveList = append(directiveList, "COMPACT STORAGE")
			default:
				return nil, errors.New("nsQL cassandra transcompiler error: invalid directive")
			}
		}

		var directives string
		if len(directiveList) != 0 {
			directives = " WITH " + strings.Join(directiveList, " AND ")
		}

		_, _, err = c.execute(execOpts, "CREATE TABLE IF NOT EXISTS "+keyspace+"."+table+" "+
			tableDescription+directives, nil)
		if err != nil {
			return nil, err
		}
		return nil, nil
	case *ast.InsertStatement:
		statement, _ := parsed.(*ast.InsertStatement)
		keyspace, table, err := c.getKeyspaceTable(statement.Into)
		if err != nil {
			return nil, err
		}

		columns, _, err := c.getColumns(statement.Columns)
		if err != nil {
			return nil, err
		}

		values, stubs, err := c.getValues(statement.Values)
		if err != nil {
			return nil, err
		}

		execOpts := &ExecutionOptions{Scan: false,
			Keyspace: keyspace}
		_, _, err = c.execute(execOpts, "INSERT INTO "+keyspace+"."+table+" ("+columns+") VALUES ("+stubs+")", values...)
		if err != nil {
			return nil, err
		}
		return nil, nil
	case *ast.DeleteStatement:
		statement, _ := parsed.(*ast.DeleteStatement)
		keyspace, table, err := c.getKeyspaceTable(statement.From.Tables[0])
		if err != nil {
			return nil, err
		}

		filter, args, err := c.getWhere(statement.Where)
		if err != nil {
			return nil, err
		}

		if filter == "" {
			return nil, errors.New("nsQL cassandra transcompiler error: missing WHERE clause in DELETE " +
				"statement")
		}

		execOpts := &ExecutionOptions{Scan: false,
			Keyspace: keyspace}
		_, _, err = c.execute(execOpts, "DELETE FROM "+keyspace+"."+table+" WHERE "+filter, args...)
		if err != nil {
			return nil, err
		}

		return nil, nil
	case *ast.UpdateStatement:
		statement, _ := parsed.(*ast.UpdateStatement)
		keyspace, table, err := c.getKeyspaceTable(statement.Table)
		if err != nil {
			return nil, err
		}

		update, args1, err := c.getSet(statement.Update)
		if err != nil {
			return nil, err
		}

		filter, args2, err := c.getWhere(statement.Where)
		if err != nil {
			return nil, err
		}

		if filter == "" {
			return nil, errors.New("nsQL cassandra transcompiler error: missing WHERE clause in UPDATE " +
				"statement")
		}

		execOpts := &ExecutionOptions{Scan: false,
			Keyspace: keyspace}
		_, _, err = c.execute(execOpts, "UPDATE "+keyspace+"."+table+" SET "+update+" WHERE "+filter,
			append(args1, args2...)...)

		if err != nil {
			return nil, err
		}

		return nil, nil
	case *ast.SelectStatement:
		statement, _ := parsed.(*ast.SelectStatement)
		if len(statement.From.Tables) != 1 {
			return nil, errors.New("nsQL cassandra transcompiler error: invalid FROM in SELECT")
		}

		columns, transformers, err := c.getColumns(statement.Select.Expressions)
		if err != nil {
			return nil, err
		}

		keyspace, table, err := c.getKeyspaceTable(statement.From.Tables[0])
		if err != nil {
			return nil, err
		}

		filter, args, err := c.getWhere(statement.Where)
		if err != nil {
			return nil, err
		}

		if statement.GroupBy != nil || statement.OrderBy != nil {
			return nil, errors.New("nsQL cassandra transcompiler error: GROUP BY or ORDER BY not " +
				"supported in SELECT")
		}

		var data []map[string]interface{}

		var limit string
		if statement.Limit != "" {
			limit = " LIMIT " + statement.Limit
		}

		execOpts := &ExecutionOptions{Scan: true,
			Keyspace: keyspace}

		allowFiltering := ""
		if options.AllowFiltering {
			allowFiltering = " allow filtering"
		}

		var meta *gocql.KeyspaceMetadata

		if filter == "" {
			data, meta, err = c.execute(execOpts, "SELECT "+columns+" FROM "+keyspace+"."+table+
				limit+allowFiltering)
		} else {
			data, meta, err = c.execute(execOpts, "SELECT "+columns+" FROM "+keyspace+"."+table+" WHERE "+
				filter+limit+allowFiltering, args...)
		}

		if err != nil {
			return nil, err
		}

		result := make(map[string]interface{})
		var columnNames []interface{}
		var columnTypes []string
		var rows []interface{}

		if transformers[0].Columns[0] == "*" {
			transformer := transformers[0]
			transformers = nil
			if len(data) != 0 {
				for key, _ := range data[0] {
					newTransformer := &Transformer{Columns: []string{key},
						OriginalColumns: []string{key},
						Function:        transformer.Function,
						Alias:           key,
						Parameters:      transformer.Parameters}
					transformers = append(transformers, newTransformer)
				}
			}
		}

		for _, transformer := range transformers {
			columnNames = append(columnNames, transformer.Alias)
		}

		types := false
		for _, row := range data {
			var columnValues []interface{}
			for _, transformer := range transformers {
				parameters := []interface{}{}
				for _, col := range transformer.Columns {
					parameters = append(parameters, row[col])
				}

				parameters = append(parameters, transformer.Parameters...)

				for _, col := range transformer.OriginalColumns {
					parameters = append(parameters, meta.Tables[strings.ToLower(table)].
						Columns[strings.ToLower(col)].Validator)
				}

				value, err := transformer.Function(parameters...)
				if err != nil {
					return nil, err
				}

				if transformer.Name == TCOUNT {
					result["type"] = "int"
					result["value"] = fmt.Sprintf("%v", value)
					return result, nil
				}

				typeAlias, err := util.ToInternalType(reflect.TypeOf(value))
				if err != nil {
					return nil, err
				}

				if !types {
					columnTypes = append(columnTypes, typeAlias)
				}

				if options.ReturnTyped {
					columnValues = append(columnValues, value)
				} else {
					columnValues = append(columnValues, util.DataToString(value, typeAlias))
				}
			}
			types = true
			rows = append(rows, columnValues)
		}
		result["columns"] = columnNames
		result["types"] = columnTypes
		result["rows"] = rows

		return result, nil
	default:
		return nil, errors.New("nsQL cassandra transcompiler error: SELECT/UPDATE statements are not " +
			"supported")
	}
}

func (c *CassandraCompiler) setSession() error {
	port, err := strconv.Atoi(c.Connection.Port)
	if err != nil {
		return err
	}

	rp := new(gocql.SimpleRetryPolicy)
	rp.NumRetries = 5

	hostStrArr, err := getHostStrArr(c.Connection.Host)
	if err != nil {
		return err
	}

	cluster := gocql.NewCluster(hostStrArr...)
	cluster.Consistency = gocql.LocalQuorum
	cluster.Timeout = 3 * time.Second
	cluster.NumConns = 5
	cluster.RetryPolicy = rp
	cluster.ProtoVersion = CassandraProtoVersion
	cluster.PageSize = 500
	cluster.Port = port
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: c.Connection.Username,
		Password: c.Connection.Password,
	}

	if c.Connection.Version != "" {
		mlog.Debug("Setting CQL version to: %v", c.Connection.Version)
		cluster.CQLVersion = c.Connection.Version
	}

	if c.Session, err = cluster.CreateSession(); err != nil {
		errMsg := fmt.Sprintf("nsQL cassandra transcompiler error: unable to get session: %v", err)
		mlog.Error(errMsg)
		return errors.New(errMsg)
	}

	return nil
}

func getHostStrArr(hostStr string) ([]string, error) {
	hostStr = dktUtils.HostsToIps(hostStr)
	if hostStr == "" {
		return nil, fmt.Errorf("HostStr is empty")
	}

	return strings.Split(hostStr, ","), nil
}

func (c *CassandraCompiler) execute(options *ExecutionOptions,
	query string,
	args ...interface{}) ([]map[string]interface{}, *gocql.KeyspaceMetadata, error) {
	var data []map[string]interface{}
	metadata, err := c.Session.KeyspaceMetadata(options.Keyspace)
	if err != nil {
		return nil, nil, errors.New("nsQL cassandra transcompiler error: " + err.Error())
	}

	queryCall := c.Session.Query(query, args...)
	mlog.Debug("Query: %v", queryCall)
	if options.Scan {
		data, err = queryCall.Iter().SliceMap()
	} else {
		err = queryCall.Exec()
	}

	if err != nil {
		return nil, nil, errors.New("nsQL cassandra transcompiler error: " + err.Error())
	}

	return data, metadata, nil
}

func (c *CassandraCompiler) getKeyspaceTable(expression ast.Expression) (string, string, error) {
	tableName, ok := expression.(*ast.IdentifierExpression)
	if !ok {
		return "", "", errors.New("nsQL cassandra transcompiler error: invalid table name")
	}

	return tableName.Owner, tableName.Name, nil
}

func (c *CassandraCompiler) getSet(expressions []ast.Expression) (string, []interface{}, error) {
	var update []string
	var args []interface{}
	for _, expression := range expressions {
		switch expression.(type) {
		case *ast.ConditionalExpression:
			cond, _ := expression.(*ast.ConditionalExpression)

			left, lok := cond.GetLeft().(*ast.IdentifierExpression)
			right, rok := cond.GetRight().(*ast.LiteralExpression)

			if !lok || !rok {
				return "", nil, errors.New("nsQL cassandra transcompiler error: invalid expression " +
					"in SET")
			}

			operator := c.getOperator(cond.Operator)

			if operator != "=" {
				return "", nil, errors.New("nsQL cassandra transcompiler error: invalid expression " +
					"in SET")
			}

			update = append(update, left.Name+" "+operator+" ?")
			args = append(args, right.Value)

		default:
			return "", nil, errors.New("nsQL cassandra transcompiler error: invalid expression in SET")
		}
	}
	return strings.Join(update, ","), args, nil
}

func (c *CassandraCompiler) getWhere(expression ast.Expression) (string, []interface{}, error) {
	if expression == nil {
		return "", nil, nil
	}
	switch expression.(type) {
	case *ast.IdentifierExpression:
		identifierExpression, _ := expression.(*ast.IdentifierExpression)
		return identifierExpression.Name, []interface{}{}, nil
	case *ast.LiteralExpression:
		literalExpression, _ := expression.(*ast.LiteralExpression)
		return "?", []interface{}{literalExpression.Value}, nil
	case *ast.TemporalExpression:
		temporalExpression, _ := expression.(*ast.TemporalExpression)

		lFilter, lArgs, err := c.getWhere(temporalExpression.GetLeft())
		if err != nil {
			return "", nil, err
		}

		rFilter, rArgs, err := c.getWhere(temporalExpression.GetRight())
		if err != nil {
			return "", nil, err
		}

		filter := lFilter + " " + c.getOperator(temporalExpression.Operator) + " " + rFilter
		args := append(lArgs, rArgs...)

		return filter, args, nil
	case *ast.NumericExpression:
		numericExpression, _ := expression.(*ast.NumericExpression)

		lFilter, lArgs, err := c.getWhere(numericExpression.GetLeft())
		if err != nil {
			return "", nil, err
		}

		rFilter, rArgs, err := c.getWhere(numericExpression.GetRight())
		if err != nil {
			return "", nil, err
		}

		filter := lFilter + " " + c.getOperator(numericExpression.Operator) + " " + rFilter
		args := append(lArgs, rArgs...)

		return filter, args, nil
	case *ast.ColumnExpression:
		columnExpression, _ := expression.(*ast.ColumnExpression)

		lFilter, lArgs, err := c.getWhere(columnExpression.GetLeft())
		if err != nil {
			return "", nil, err
		}

		rFilter, rArgs, err := c.getWhere(columnExpression.GetRight())
		if err != nil {
			return "", nil, err
		}

		filter := lFilter + " " + c.getOperator(columnExpression.Operator) + " " + rFilter
		args := append(lArgs, rArgs...)

		return filter, args, nil
	case *ast.ConditionalExpression:
		conditionalExpression, _ := expression.(*ast.ConditionalExpression)

		lFilter, lArgs, err := c.getWhere(conditionalExpression.GetLeft())
		if err != nil {
			return "", nil, err
		}

		rFilter, rArgs, err := c.getWhere(conditionalExpression.GetRight())
		if err != nil {
			return "", nil, err
		}

		filter := lFilter + " " + c.getOperator(conditionalExpression.Operator) + " " + rFilter
		args := append(lArgs, rArgs...)

		return filter, args, nil
	case *ast.LogicalExpression:
		logicalExpression, _ := expression.(*ast.LogicalExpression)

		lFilter, lArgs, err := c.getWhere(logicalExpression.GetLeft())
		if err != nil {
			return "", nil, err
		}

		rFilter, rArgs, err := c.getWhere(logicalExpression.GetRight())
		if err != nil {
			return "", nil, err
		}

		filter := lFilter + " " + c.getOperator(logicalExpression.Operator) + " " + rFilter
		args := append(lArgs, rArgs...)

		return filter, args, nil
	default:
		return "", nil, errors.New("nsQL cassandra transcompiler: unsupported expression in WHERE")
	}
}

func (c *CassandraCompiler) getColumns(expressions []ast.Expression) (string, []*Transformer, error) {
	var columns []string
	var trans *Transformer
	var transformers []*Transformer
	mode := UNSET
	for _, expression := range expressions {
		switch expression.(type) {
		case *ast.IdentifierExpression:
			column, _ := expression.(*ast.IdentifierExpression)
			if column.Name == "*" {
				if mode == SINGULAR {
					return "", nil, errors.New("nsQL cassandra transcompiler error: invalid " +
						"column mix")
				}
				mode = GROUP
			} else {
				if mode == GROUP {
					return "", nil, errors.New("nsQL cassandra transcompiler error: invalid " +
						"column mix")
				}
				mode = SINGULAR
			}

			alias := ""
			if column.Alias != "" {
				alias = strings.ToLower(column.Alias)
				columns = append(columns, column.Name+" AS "+alias)
				trans = &Transformer{Columns: []string{alias}, OriginalColumns: []string{column.Name},
					Name: IDENTITY, Alias: alias, Parameters: []interface{}{}}
			} else {
				columns = append(columns, column.Name)
				trans = &Transformer{Columns: []string{strings.ToLower(column.Name)},
					OriginalColumns: []string{strings.ToLower(column.Name)},
					Name:            IDENTITY, Alias: column.Name, Parameters: []interface{}{}}
			}

			trans.Function = func(args ...interface{}) (interface{}, error) {
				if args[1] == "time" {
					if t, ok := args[0].(time.Time); ok {
						result := t.UnixNano()
						if result < 0 {
							result = 0
						}
						return result, nil
					}
					return nil, errors.New("nsQL cassandra transcompiler error: invalid " +
						"column type")
				}
				return args[0], nil
			}
			transformers = append(transformers, trans)
		case *ast.ToStringTransformer:
			function, _ := expression.(*ast.ToStringTransformer)
			column, _ := function.Parameters[0].(*ast.IdentifierExpression)

			if mode == GROUP {
				return "", nil, errors.New("nsQL cassandra transcompiler error: invalid " +
					"column mix")
			}
			mode = SINGULAR

			alias := ""
			if function.Alias != "" {
				alias = strings.ToLower(function.Alias)
				columns = append(columns, column.Name+" AS "+alias)
				trans = &Transformer{Columns: []string{alias}, OriginalColumns: []string{column.Name},
					Alias: alias, Name: function.Name}
			} else {
				columns = append(columns, column.Name)
				trans = &Transformer{Columns: []string{column.Name}, OriginalColumns: []string{column.Name}, Name: function.Name}
			}

			switch function.Name {
			case JSON_FETCH:
				field, _ := function.Parameters[1].(*ast.LiteralExpression)

				trans.Parameters = []interface{}{field.Value}

				trans.Function = func(args ...interface{}) (interface{}, error) {
					value, ok := args[0].(string)
					if !ok {
						return nil, errors.New("nsQL cassandra transcompiler error: " +
							"invalid first argument in JSON_FETCH, string required")
					}

					data := make(map[string]interface{})
					err := json.Unmarshal([]byte(value), &data)
					if err != nil {
						return nil, errors.New("nsQL cassandra transcompiler error: " +
							"malformed records received")
					}

					field, ok := args[1].(string)
					if !ok {
						return nil, errors.New("nsQL cassandra transcompiler error: " +
							"invalid second argument in JSON_FETCH, string required")
					}

					if _, ok := data[field]; ok {
						return data[field], nil
					}

					return "", nil
				}
				if trans.Alias == "" {
					trans.Alias = function.Name + "(" + field.Value + ")"
				}
				transformers = append(transformers, trans)
			case MAP_BLOB_JSON_FETCH:
				field, _ := function.Parameters[1].(*ast.LiteralExpression)
				subField, _ := function.Parameters[2].(*ast.LiteralExpression)

				trans.Parameters = []interface{}{field.Value, subField.Value}

				trans.Function = func(args ...interface{}) (interface{}, error) {
					value, ok := args[0].(map[string][]byte)
					if !ok {
						return nil, errors.New("nsQL cassandra transcompiler error: invalid " +
							"first argument in MAP_BLOB_JSON_FETCH, map[string][]byte" +
							"required")
					}

					field, ok := args[1].(string)
					if !ok {
						return nil, errors.New("nsQL cassandra transcompiler error: invalid " +
							"second argument in MAP_BLOB_JSON_FETCH, string required")
					}

					if _, ok := value[field]; !ok {
						return "", nil
					}

					var data interface{}
					err := json.Unmarshal(value[field], &data)
					if err != nil {
						return nil, errors.New("nsQL cassandra transcompiler error: " +
							"malformed records received")
					}

					subField, ok := args[2].(string)
					if !ok {
						return nil, errors.New("nsQL cassandra transcompiler error: invalid " +
							"third argument in MAP_BLOB_JSON_FETCH, string required")
					}

					converted, ok := data.(map[string]interface{})
					if ok {
						if _, ok := converted[subField]; ok {
							return converted[subField], nil
						}
						return "", nil
					}

					return data, nil
				}
				if trans.Alias == "" {
					trans.Alias = function.Name + "(" + field.Value + "," + subField.Value + ")"
				}
				transformers = append(transformers, trans)
			default:
				return "", nil, errors.New("nsQL cassandra transcompiler error: unknown function")
			}
		case *ast.ToNumericTransformer:
			function, _ := expression.(*ast.ToNumericTransformer)

			if mode == GROUP {
				return "", nil, errors.New("nsQL cassandra transcompiler error: invalid " +
					"column mix")
			}
			mode = SINGULAR

			switch function.Name {
			case SUBTRACT_TIMESTAMPS:
				column1, _ := function.Parameters[0].(*ast.IdentifierExpression)
				column2, _ := function.Parameters[1].(*ast.IdentifierExpression)

				if function.Alias != "" {
					alias1 := strings.ToLower(function.Alias) + "1"
					alias2 := strings.ToLower(function.Alias) + "2"
					columns = append(columns, column1.Name+" AS "+alias1)
					columns = append(columns, column2.Name+" AS "+alias2)
					trans = &Transformer{Columns: []string{alias1, alias2},
						OriginalColumns: []string{column1.Name, column2.Name},
						Alias:           function.Alias,
						Name:            function.Name}
				} else {
					columns = append(columns, column1.Name)
					columns = append(columns, column2.Name)
					trans = &Transformer{Columns: []string{column1.Name, column2.Name},
						OriginalColumns: []string{column1.Name, column2.Name},
						Alias:           function.Name + "(" + column1.Name + "," + column2.Name + ")",
						Name:            function.Name}
				}

				trans.Parameters = []interface{}{}

				trans.Function = func(args ...interface{}) (interface{}, error) {
					col1Type, _ := args[2].(string)
					col2Type, _ := args[3].(string)

					if col1Type == "time" || col2Type == "time" {
						return nil, errors.New("nsQL cassandra transcompiler error: invalid " +
							"first argument in SUBTRACT_TIMESTAMPS, timestamp required")
					}

					col1, ok := args[0].(time.Time)
					if !ok {
						return nil, errors.New("nsQL cassandra transcompiler error: invalid " +
							"first argument in SUBTRACT_TIMESTAMPS, timestamp required")
					}

					col2, ok := args[1].(time.Time)
					if !ok {
						return nil, errors.New("nsQL cassandra transcompiler error: invalid " +
							"second argument in SUBTRACT_TIMESTAMPS, timestamp required")
					}

					operand1 := col1.UnixNano()
					operand2 := col2.UnixNano()

					if operand1 <= 0 || operand2 <= 0 {
						return 0, nil
					}

					return operand1 - operand2, nil
				}
				transformers = append(transformers, trans)
			default:
				return "", nil, errors.New("nsQL cassandra transcompiler error: unknown function")
			}
		case *ast.ToNumericAggregator:
			function, _ := expression.(*ast.ToNumericAggregator)

			column, _ := function.Parameters[0].(*ast.IdentifierExpression)

			if function.Name == "mean" {
				function.Name = "avg"
			}

			if mode == GROUP {
				return "", nil, errors.New("nsQL cassandra transcompiler error: invalid " +
					"column mix")
			}
			mode = SINGULAR

			alias := ""
			if function.Alias != "" {
				alias = strings.ToLower(function.Alias)
				columns = append(columns, function.Name+"("+column.Name+") AS "+alias)
				trans = &Transformer{Columns: []string{alias},
					Name: function.Name, Parameters: []interface{}{}}
			} else {
				columns = append(columns, function.Name+"("+column.Name+")")
				trans = &Transformer{Columns: []string{"system." + strings.ToLower(function.Name) +
					"(" + column.Name + ")"},
					Name:       function.Name,
					Parameters: []interface{}{}}
			}

			trans.Function = func(args ...interface{}) (interface{}, error) {
				return args[0], nil
			}
			trans.Alias = trans.Columns[0]
			transformers = append(transformers, trans)
		case *ast.ToColumnAggregator:
			function, _ := expression.(*ast.ToColumnAggregator)

			column, _ := function.Parameters[0].(*ast.IdentifierExpression)

			if mode == GROUP {
				return "", nil, errors.New("nsQL cassandra transcompiler error: invalid " +
					"column mix")
			}
			mode = SINGULAR

			alias := ""
			if function.Alias != "" {
				alias = strings.ToLower(function.Alias)
				columns = append(columns, function.Name+"("+column.Name+") AS "+alias)
				trans = &Transformer{Columns: []string{alias},
					Name:       function.Name,
					Parameters: []interface{}{}}
			} else {
				columns = append(columns, function.Name+"("+column.Name+")")
				trans = &Transformer{Columns: []string{"system." + strings.ToLower(function.Name) +
					"(" + column.Name + ")"},
					Name:       function.Name,
					Parameters: []interface{}{}}
			}

			trans.Function = func(args ...interface{}) (interface{}, error) {
				return args[0], nil
			}
			trans.Alias = trans.Columns[0]
			transformers = append(transformers, trans)
		case *ast.TableAggregator:
			function, _ := expression.(*ast.TableAggregator)
			switch function.Name {
			case TCOUNT:
				columnName := "count"
				columns = append(columns, columnName+"(*)")
				trans := &Transformer{Columns: []string{columnName},
					Name: TCOUNT, Parameters: []interface{}{}}
				trans.Function = func(args ...interface{}) (interface{}, error) {
					return args[0], nil
				}
				trans.Alias = trans.Columns[0]
				transformers = append(transformers, trans)
			default:
				return "", nil, errors.New("nsQL cassandra transcompiler error: unknown table " +
					"aggregator")
			}
		default:
			return "", nil, errors.New("nsQL cassandra transcompiler error: invalid column name")
		}
	}

	return strings.Join(columns, ","), transformers, nil
}

func (c *CassandraCompiler) getValues(expressions []ast.Expression) ([]interface{}, string, error) {
	var values []interface{}
	var stubs []string
	for _, expression := range expressions {
		value, ok := expression.(*ast.LiteralExpression)
		if !ok {
			return nil, "", errors.New("nsQL cassandra transcompiler error: unknown literal")
		}

		if value.Token == parser.BINARY || value.Token == parser.COLLECTION {
			values = append(values, value.Original)
		} else {
			values = append(values, value.Value)
		}
	}

	for i := 0; i < len(values); i++ {
		stubs = append(stubs, "?")
	}

	return values, strings.Join(stubs, ","), nil
}

func (c *CassandraCompiler) getOperator(operator string) string {
	switch operator {
	case "==":
		return "="
	default:
		return operator
	}
}
