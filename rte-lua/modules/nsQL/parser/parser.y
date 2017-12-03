%{package parser

import (
	"errors"
	"github.com/verizonlabs/northstar/rte-lua/modules/nsQL/parser/ast"
)
%}

%type<SelectStatement>         SelectStatement
%type<DeleteStatement>         DeleteStatement
%type<InsertStatement>         InsertStatement
%type<UpdateStatement>         UpdateStatement
%type<CreateTableStatement>    CreateTableStatement
%type<DropTableStatement>      DropTableStatement
%type<Directives>              Directives
%type<Properties>              Properties
%type<Property>                Property
%type<Order>                   Order
%type<CompoundOrder>           CompoundOrder
%type<Sorting>                 Sorting
%type<SimpleOrder>             SimpleOrder
%type<TableDescription>        TableDescription
%type<FieldDescriptions>       FieldDescriptions
%type<FieldDescription>        FieldDescription
%type<PrimaryKey>              PrimaryKey
%type<PartitioningKey>         PartitioningKey
%type<ClusteringColumns>       ClusteringColumns
%type<CompoundPartitioningKey> CompoundPartitioningKey
%type<Identifiers>             Identifiers
%type<SimplePartitioningKey>   SimplePartitioningKey
%type<Type>                    Type
%type<CompoundType>            CompoundType
%type<SimpleType>              SimpleType
%type<Select>                  Select
%type<Limit>                   Limit
%type<OrderBy>                 OrderBy
%type<GroupBy>                 GroupBy
%type<Having>                  Having
%type<Where>                   Where
%type<From>                    From
%type<Tables>                  Tables
%type<Join>                    Join
%type<On>                      On
%type<Table>                   Table
%type<Qualifier>               Qualifier
%type<Columns>                 Columns
%type<Column>                  Column
%type<Expressions>             Expressions
%type<Expression>              Expression
%type<LogicalExpression>       LogicalExpression
%type<ConditionalExpression>   ConditionalExpression
%type<OrdinaryExpression>      OrdinaryExpression
%type<StringExpression>        StringExpression
%type<TemporalExpression>      TemporalExpression
%type<NumericExpression>       NumericExpression
%type<ColumnExpression>        ColumnExpression
%type<SignedTimeInterval>      SignedTimeInterval
%type<TableAggregator>         TableAggregator
%type<ToColumnAggregator>      ToColumnAggregator
%type<ToNumericAggregator>     ToNumericAggregator
%type<ToNumericTransformer>    ToNumericTransformer
%type<ToTemporalTransformer>   ToTemporalTransformer
%type<ToStringTransformer>     ToStringTransformer
%type<TCount>                  TCount
%type<TCorr>                   TCorr
%type<TCov>                    TCov
%type<Min>                     Min
%type<Max>                     Max
%type<First>                   First
%type<Last>                    Last
%type<Count>                   Count
%type<Sum>                     Sum
%type<Mean>                    Mean
%type<Variance>                Variance
%type<Stdev>                   Stdev
%type<Corr>                    Corr
%type<Year>                    Year
%type<Month>                   Month
%type<Day>                     Day
%type<Hour>                    Hour
%type<Minute>                  Minute
%type<Second>                  Second
%type<Now>                     Now
%type<Map_Blob_Json_Fetch>     Map_Blob_Json_Fetch
%type<Json_Fetch>              Json_Fetch
%type<Subtract_Timestamps>     Subtract_Timestamps
%type<GenericParameter>        GenericParameter
%type<NumericParameter>        NumericParameter
%type<TemporalParameter>       TemporalParameter
%type<ColumnGroup>             ColumnGroup
%type<ColumnName>              ColumnName
%type<Identifier>              Identifier
%type<Literals>                Literals
%type<Literal>                 Literal
%type<Null>                    Null
%type<Number>                  Number
%type<Integer>                 Integer
%type<Float>                   Float
%type<String>                  String
%type<Boolean>                 Boolean
%type<Uuid>                    Uuid
%type<Timestamp>               Timestamp
%type<TimeInterval>            TimeInterval
%type<Binary>                  Binary
%type<Collection>              Collection
%type<InclusionComparator>     InclusionComparator
%type<IdentityComparator>      IdentityComparator
%type<RegularComparator>       RegularComparator
%type<EqualityComparator>      EqualityComparator
%type<RangeComparator>         RangeComparator

%union {
	SelectStatement         ast.Expression
	DeleteStatement         ast.Expression
	InsertStatement         ast.Expression
	UpdateStatement         ast.Expression
	CreateTableStatement    ast.Expression
	DropTableStatement      ast.Expression
	Directives              []ast.Property
	Properties              []ast.Property
	Property                ast.Property
	Order                   []*ast.Order
	CompoundOrder           []*ast.Order
	Sorting                 []*ast.Order
	SimpleOrder             *ast.Order
	TableDescription        *ast.TableDescription
	FieldDescriptions       []*ast.FieldDescription
	FieldDescription        *ast.FieldDescription
	PrimaryKey              *ast.PrimaryKey
	PartitioningKey         []string
	CompoundPartitioningKey []string
	Identifiers             []string
	SimplePartitioningKey   []string
	ClusteringColumns       []string
	Type                    string
	CompoundType            string
	SimpleType              string
	Select                  *ast.Select
	Limit                   string
	OrderBy                 []ast.Expression
	GroupBy                 *ast.GroupBy
	Having                  ast.Expression
	Where                   ast.Expression
	From                    *ast.From
	Tables                  *ast.From
	Join                    string
	On                      ast.Expression
	Table                   ast.Expression
	Qualifier               string
	Columns                 []ast.Expression
	Column                  ast.Expression
	Expressions             []ast.Expression
	Expression              ast.Expression
	LogicalExpression       ast.Expression
	ConditionalExpression   ast.Expression
	OrdinaryExpression      ast.Expression
	StringExpression        ast.Expression
	TemporalExpression      ast.Expression
	NumericExpression       ast.Expression
	ColumnExpression        ast.Expression
	SignedTimeInterval      ast.Expression
	TableAggregator         *ast.TableAggregator
	ToColumnAggregator      *ast.ToColumnAggregator
	ToNumericAggregator     *ast.ToNumericAggregator
	ToNumericTransformer    *ast.ToNumericTransformer
	ToTemporalTransformer   *ast.ToTemporalTransformer
	ToStringTransformer     *ast.ToStringTransformer
	TCount                  *ast.TableAggregator
	TCorr                   *ast.TableAggregator
	TCov                    *ast.TableAggregator
	Min                     *ast.ToColumnAggregator
	Max                     *ast.ToColumnAggregator
	First                   *ast.ToColumnAggregator
	Last                    *ast.ToColumnAggregator
	Count                   *ast.ToNumericAggregator
	Sum                     *ast.ToNumericAggregator
	Mean                    *ast.ToNumericAggregator
	Variance                *ast.ToNumericAggregator
	Stdev                   *ast.ToNumericAggregator
	Corr                    *ast.ToNumericAggregator
	Year                    *ast.ToNumericTransformer
	Month                   *ast.ToNumericTransformer
	Day                     *ast.ToNumericTransformer
	Hour                    *ast.ToNumericTransformer
	Minute                  *ast.ToNumericTransformer
	Second                  *ast.ToNumericTransformer
	Now                     *ast.ToTemporalTransformer
	Map_Blob_Json_Fetch     *ast.ToStringTransformer
	Json_Fetch              *ast.ToStringTransformer
	Subtract_Timestamps     *ast.ToNumericTransformer
	GenericParameter        ast.Expression
	NumericParameter        ast.Expression
	TemporalParameter       ast.Expression
	ColumnGroup             *ast.IdentifierExpression
	ColumnName              *ast.IdentifierExpression
	Identifier              *Token
	Literals                []ast.Expression
	Literal                 *ast.LiteralExpression
	Null                    *ast.LiteralExpression
	Number                  *ast.LiteralExpression
	Integer                 *ast.LiteralExpression
	Float                   *ast.LiteralExpression
	String                  *ast.LiteralExpression
	Boolean                 *ast.LiteralExpression
	Uuid                    *ast.LiteralExpression
	Timestamp               *ast.LiteralExpression
	TimeInterval            *ast.LiteralExpression
	Binary                  *ast.LiteralExpression
	Collection              *ast.LiteralExpression
	InclusionComparator     string
	IdentityComparator      string
	RegularComparator       string
	EqualityComparator      string
	RangeComparator         string
}

%token<keyword>              FROM AS GROUP ORDER BY LIMIT INTO VALUES SET TABLE PRIMARY KEY CLUSTERING COMPACT STORAGE
%token<statement>            INSERT DELETE UPDATE SELECT CREATE DROP
%token<qualifier>            ALL DISTINCT EXISTS WITH ASC DESC
%token<join>                 JOIN INNER OUTER FULL LEFT RIGHT SEMI
%token<filter>               WHERE ON HAVING IF
%token<literal>              IDENTIFIER NULL INTEGER FLOAT STRING BOOLEAN UUID TIMESTAMP DATE TIME TIME_INTERVAL BINARY COLLECTION
%token<ordinary_operator>    PLUS_SIGN MINUS_SIGN ASTERISK SLASH PERCENT_SIGN AMPERSAND VERTICAL_BAR
%token<conditional_operator> IS EQUAL_SIGN EXCLAMATION_MARK LESS_THAN GREATER_THAN IN
%token<logical_operator>     AND OR NOT
%token<special_operator>     UNION INTERSECT BETWEEN
%token<separator>            COMMA PERIOD LEFT_PARANTHESIS RIGHT_PARANTHESIS SEMICOLON
%token<column_aggregator>    COUNT SUM MIN MAX FIRST LAST MEAN VARIANCE STDEV CORR
%token<column_transformer>   YEAR MONTH DAY HOUR MINUTE SECOND NOW MAP_BLOB_JSON_FETCH JSON_FETCH SUBTRACT_TIMESTAMPS
%token<table_aggregator>     TCOUNT TCORR TCOV
%token<data_types>           ASCII BIGINT BLOB BOOLEANTYPE COUNTER DECIMAL DOUBLE FLOATTYPE INET INT TEXT TIMESTAMPTYPE TIMEUUID UUIDTYPE VARCHAR VARINT LIST MAP
%token<unknown>              UNKNOWN

%left VERTICAL_BAR
%left AMPERSAND
%left PLUS_SIGN MINUS_SIGN
%left ASTERISK SLASH PERCENT_SIGN
%left UNARY_MINUS_SIGN UNARY_PLUS_SIGN
%left OR
%left AND
%left NOT
%left UNARY_NOT
%left UNION ALL
%left INTERSECT

%%
Statement
:	SelectStatement      SEMICOLON	{	getLexer(nsQLlex).Statement = $1	}
|	DeleteStatement      SEMICOLON	{	getLexer(nsQLlex).Statement = $1	}
|	InsertStatement      SEMICOLON  {	getLexer(nsQLlex).Statement = $1	}
|	UpdateStatement      SEMICOLON	{	getLexer(nsQLlex).Statement = $1	}
|	CreateTableStatement SEMICOLON  {       getLexer(nsQLlex).Statement = $1        }
|	DropTableStatement   SEMICOLON  {	getLexer(nsQLlex).Statement = $1        };

SelectStatement
:	LEFT_PARANTHESIS SelectStatement RIGHT_PARANTHESIS
	{	$$ = $2									}
|	SelectStatement UNION SelectStatement
	{	$$ = makeSelectExpression($1, $3, "union", getLexer(nsQLlex))		}
|	SelectStatement UNION ALL SelectStatement
	{	$$ = makeSelectExpression($1, $4, "union all", getLexer(nsQLlex))	}
|	SelectStatement INTERSECT SelectStatement
	{	$$ = makeSelectExpression($1, $3, "intersect", getLexer(nsQLlex))	}
|	Select From Where GroupBy OrderBy Limit
	{	$$ = makeSelectStatement($1, $2, $3, $4, $5, $6, getLexer(nsQLlex))	};

DeleteStatement
:	DELETE From Where
	{	lexer := getLexer(nsQLlex)
		failOnNonTableName($2.Tables, lexer)
		failOnSubquery([]ast.Expression{$3}, lexer)
		$$ = makeDeleteStatement($2, $3, lexer)		};

InsertStatement
:	INSERT INTO Table LEFT_PARANTHESIS Columns RIGHT_PARANTHESIS VALUES LEFT_PARANTHESIS Literals RIGHT_PARANTHESIS
	{	lexer := getLexer(nsQLlex)
		failOnNonTableName([]ast.Expression{$3}, lexer)
		$$ = makeInsertStatement($3, $5, $9, lexer)				};

UpdateStatement
:	UPDATE Table SET Expressions Where
	{	lexer := getLexer(nsQLlex)
		failOnNonTableName([]ast.Expression{$2}, lexer)
		$$ = makeUpdateStatement($2, $4, $5, lexer)		};

CreateTableStatement
:	CREATE TABLE IF NOT EXISTS Table TableDescription Directives
	{	lexer := getLexer(nsQLlex)
		failOnNonTableName([]ast.Expression{$6}, lexer)
		$$ = makeCreateTableStatement($6, $7, $8, lexer)		};

DropTableStatement
:	DROP TABLE Table
	{	lexer := getLexer(nsQLlex)
		failOnNonTableName([]ast.Expression{$3}, lexer)
		$$ = makeDropTableStatement($3, lexer)				};

Directives
:	WITH Properties	{	$$ = $2		}
|			{	$$ = nil	};

Properties
:	Properties AND Property		{	$$ = append($1, $3)		}
|	Property			{	$$ = []ast.Property{$1}		};

Property
:	CLUSTERING ORDER BY Order	{	$$ = &ast.ClusteringOrder{Order: $4}		}
|	COMPACT STORAGE			{	$$ = &ast.CompactStorage{}			};

Order
:	CompoundOrder					{	$$ = $1				}
|	LEFT_PARANTHESIS SimpleOrder RIGHT_PARANTHESIS	{	$$ = []*ast.Order{$2}		};

CompoundOrder
:	LEFT_PARANTHESIS Sorting RIGHT_PARANTHESIS	{	$$ = $2		};

Sorting
:	Sorting COMMA SimpleOrder	{	$$ = append($1, $3)			}
|	SimpleOrder COMMA SimpleOrder	{	$$ = []*ast.Order{$1, $3}		};

SimpleOrder
:	Identifier ASC	{	$$ = &ast.Order{FieldName: $1.Value, Ascending: true}		}
|	Identifier DESC	{	$$ = &ast.Order{FieldName: $1.Value, Ascending: false}		};

TableDescription
:	LEFT_PARANTHESIS FieldDescriptions COMMA PrimaryKey RIGHT_PARANTHESIS
	{	$$ = &ast.TableDescription{Fields: $2, PrimaryKey: $4}		};

FieldDescriptions
:	FieldDescriptions COMMA FieldDescription
	{	$$ = append($1, $3)			}
|	FieldDescription
	{	$$ = []*ast.FieldDescription{$1}		};

FieldDescription
:	Identifier Type	{	$$ = &ast.FieldDescription{FieldName: $1.Value, FieldType: $2}		};

Type
:	CompoundType	{	$$ = $1		}
|	SimpleType	{	$$ = $1		};

CompoundType
:	LIST LESS_THAN SimpleType GREATER_THAN			{	$$ = "set<" + $3 + ">"			}
|	SET LESS_THAN SimpleType GREATER_THAN			{	$$ = "set<" + $3 + ">"			}
|	MAP LESS_THAN SimpleType COMMA SimpleType GREATER_THAN	{	$$ = "map<" + $3 + "," + $5 + ">"	};

SimpleType
:	ASCII		{	$$ = "ascii"		}
|	BIGINT		{	$$ = "bigint"		}
|	BLOB		{	$$ = "blob"		}
|	BOOLEANTYPE	{	$$ = "boolean"		}
|	COUNTER		{	$$ = "counter"		}
|	DECIMAL		{	$$ = "decimal"		}
|	DOUBLE		{	$$ = "double"		}
|	FLOATTYPE	{	$$ = "float"		}
|	INET		{	$$ = "inet"		}
|	INT		{	$$ = "int"		}
|	TEXT		{	$$ = "text"		}
|	TIMESTAMPTYPE	{	$$ = "timestamp"	}
|	TIMEUUID	{	$$ = "timeuuid"		}
|	UUIDTYPE	{	$$ = "uuid"		}
|	VARCHAR		{	$$ = "varchar"		}
|	VARINT		{	$$ = "varint"		};


PrimaryKey
:	PRIMARY KEY LEFT_PARANTHESIS PartitioningKey COMMA ClusteringColumns RIGHT_PARANTHESIS
	{	$$ = &ast.PrimaryKey{Partitioning: $4, Clustering: $6}		}
|	PRIMARY	KEY LEFT_PARANTHESIS PartitioningKey RIGHT_PARANTHESIS
	{	$$ = &ast.PrimaryKey{Partitioning: $4}				};

PartitioningKey
:	CompoundPartitioningKey		{	$$ = $1		}
|	SimplePartitioningKey		{	$$ = $1		};

ClusteringColumns
:	Identifiers	{	$$ = $1				}
|	Identifier	{	$$ = []string{$1.Value}		};

CompoundPartitioningKey
:	LEFT_PARANTHESIS Identifiers RIGHT_PARANTHESIS	{	$$ = $2		};

Identifiers
:	Identifiers COMMA Identifier	{	$$ = append($1, $3.Value)		}
|	Identifier COMMA Identifier	{	$$ = []string{$1.Value, $3.Value}	};

SimplePartitioningKey
:	LEFT_PARANTHESIS Identifier RIGHT_PARANTHESIS	{	$$ = []string{$2.Value}		}
|	Identifier					{	$$ = []string{$1.Value}		};

Select
:	SELECT Qualifier Columns
	{	failOnSubquery($3, getLexer(nsQLlex))
		$$ = &ast.Select{Qualifier: $2, Expressions: $3}	}
|	SELECT TableAggregator
	{	$$ = &ast.Select{Expressions: []ast.Expression{$2}}	};


Limit
:				{	$$ = ""		}
|	LIMIT Integer		{	$$ = $2.Value	};

OrderBy
:				{	$$ = nil				}
|	ORDER BY Expressions	{	lexer := getLexer(nsQLlex)
					failOnSubquery($3, lexer)
					failOnNoColumnName($3, lexer)
					$$ = $3					};

GroupBy
:				{	$$ = nil				}
|	GROUP BY Columns Having	{	lexer := getLexer(nsQLlex)
					failOnSubquery($3, lexer)
					$$ = makeGroupBy($3, $4, lexer)		};

Having
:					{	$$ = nil	}
|	HAVING LogicalExpression	{	$$ = $2		};

Where
:					{	$$ = nil	}
|	WHERE LogicalExpression		{	$$ = $2		};

From
:	FROM Tables	{	finalizeFrom($2, getLexer(nsQLlex)); $$ = $2	};

Tables
:	Tables Join Table On
	{	$1.Tables = append($1.Tables, $3)
		$1.Joins = append($1.Joins, &ast.Join{Table: $3, Type: $2, On: $4})
		$$ = $1									}
|	Table
	{	$$ = &ast.From{Tables: []ast.Expression{$1}}				};

Join
:	JOIN			{	$$ = "inner"		}
|	INNER JOIN		{	$$ = "inner"		}
|	OUTER JOIN		{	$$ = "full_outer"	}
|	FULL OUTER JOIN		{	$$ = "full_outer"	}
|	FULL JOIN		{	$$ = "full_outer"	}
|	LEFT OUTER JOIN		{	$$ = "left_outer"	}
|	LEFT JOIN		{	$$ = "left_outer"	}
|	LEFT SEMI OUTER JOIN	{	$$ = "left_semi_outer"	}
|	LEFT SEMI JOIN		{	$$ = "left_semi_outer"	}
|	RIGHT OUTER JOIN	{	$$ = "right_outer"	}
|	RIGHT JOIN		{	$$ = "right_outer"	};

On
:				{	$$ = nil					}
|	ON LogicalExpression	{	lexer := getLexer(nsQLlex)
					failOnSubquery([]ast.Expression{$2}, lexer)
					$$ = $2						};

Table
:	Identifier PERIOD Identifier
	{	$$ = makeTableName($1.Value, $3.Value, getLexer(nsQLlex))	}
|	Identifier PERIOD Identifier AS Identifier
	{	table := makeTableName($1.Value, $3.Value, getLexer(nsQLlex))
		table.SetAlias($5.Value)
		$$ = table							}
|	LEFT_PARANTHESIS SelectStatement RIGHT_PARANTHESIS
	{	$$  = $2							}
|	LEFT_PARANTHESIS SelectStatement RIGHT_PARANTHESIS AS Identifier
	{	$2.SetAlias($5.Value)
		$$  = $2							};

Qualifier
:			{	$$ = "all"		}
|	ALL		{	$$ = "all"		}
|	DISTINCT	{	$$ = "distinct"		};

Columns
:	Columns COMMA Column	{	$$ = append($1, $3)		}
|	Column			{	$$ = []ast.Expression{$1}	};

Column
:	Expression			{	$$ = $1					}
|	Expression AS Identifier	{	$$ = $1; $$.SetAlias($3.Value)		}
|	ColumnGroup			{	$$ = $1					};

Expressions
:	Expressions COMMA Expression	{	$$ = append($1, $3)		}
|	Expression			{	$$ = []ast.Expression{$1}	}

Expression
:	LogicalExpression	{	$$ = $1		}
|	OrdinaryExpression	{	$$ = $1		};

LogicalExpression
:	LEFT_PARANTHESIS LogicalExpression RIGHT_PARANTHESIS
	{	$$ = $2	}
|	NOT LogicalExpression %prec UNARY_NOT
	{	$2.Negate(); $$ = $2						}
|	LogicalExpression OR LogicalExpression
	{	$$ = makeLogicalExpression($1, $3, "or", getLexer(nsQLlex))	}
|	LogicalExpression AND LogicalExpression
	{	$$ = makeLogicalExpression($1, $3, "and", getLexer(nsQLlex))	}
|	ConditionalExpression
	{	$$ = $1								};

ConditionalExpression
:	TemporalExpression RegularComparator TemporalExpression
	{	$$ = makeConditionalExpression($1, $3, $2, false, getLexer(nsQLlex))			}
|	TemporalExpression RegularComparator ColumnExpression
	{	lexer := getLexer(nsQLlex)
		failOnColumnExpression($3, lexer)
		$$ = makeConditionalExpression($1, $3, $2, false, lexer)				}
|	NumericExpression RegularComparator NumericExpression
	{	$$ = makeConditionalExpression($1, $3, $2, false, getLexer(nsQLlex))			}
|	NumericExpression RegularComparator ColumnExpression
	{	$$ = makeConditionalExpression($1, $3, $2, false, getLexer(nsQLlex))			}
|	ColumnExpression InclusionComparator LEFT_PARANTHESIS SelectStatement RIGHT_PARANTHESIS
	{	lexer := getLexer(nsQLlex)
		failOnNonColumnName($1, lexer)
		failOnNonSelectStatement($4, lexer)
		$$ = makeConditionalExpression($1, $4, $2, true, lexer)					}
|	ColumnExpression RegularComparator TemporalExpression
	{	lexer := getLexer(nsQLlex)
		failOnColumnExpression($1, lexer)
		$$ = makeConditionalExpression($1, $3, $2, false, lexer)				}
|	ColumnExpression RegularComparator NumericExpression
	{	$$ = makeConditionalExpression($1, $3, $2, false, getLexer(nsQLlex))			}
|	ColumnExpression RegularComparator ColumnExpression
	{	$$ = makeConditionalExpression($1, $3, $2, false, getLexer(nsQLlex))			}
|	ColumnExpression IdentityComparator Null
	{	lexer := getLexer(nsQLlex)
		failOnNonColumnName($1, lexer)
		$$ = makeConditionalExpression($1, $3, $2, false, lexer)				}
|	ColumnExpression RegularComparator StringExpression
	{	lexer := getLexer(nsQLlex)
		failOnNonColumnName($1, lexer)
		$$ = makeConditionalExpression($1, $3, $2, false, lexer)				}
|	ColumnExpression EqualityComparator Boolean
	{	lexer := getLexer(nsQLlex)
		failOnNonColumnName($1, lexer)
		$$ = makeConditionalExpression($1, $3, $2, false, lexer)				}
|	ColumnExpression EqualityComparator Uuid
	{	lexer := getLexer(nsQLlex)
		failOnNonColumnName($1, lexer)
		$$ = makeConditionalExpression($1, $3, $2, false, lexer)				}
|	StringExpression RegularComparator ColumnExpression
	{	lexer := getLexer(nsQLlex)
		failOnNonColumnName($3, lexer)
		$$ = makeConditionalExpression($1, $3, $2, false, lexer)				}
|	Boolean EqualityComparator ColumnExpression
	{	lexer := getLexer(nsQLlex)
		failOnNonColumnName($3, lexer)
		$$ = makeConditionalExpression($1, $3, $2, false, lexer)				}
|	Uuid EqualityComparator ColumnExpression
	{	lexer := getLexer(nsQLlex)
		failOnNonColumnName($3, lexer)
		$$ = makeConditionalExpression($1, $3, $2, false, lexer)				};

OrdinaryExpression
:	StringExpression	{	$$ = $1		}
|	TemporalExpression	{	$$ = $1		}
|	NumericExpression	{	$$ = $1		}
|	ColumnExpression	{	$$ = $1		};

StringExpression
:	LEFT_PARANTHESIS StringExpression RIGHT_PARANTHESIS
	{	$$ = $2						}
|	String
	{	$$ = $1						}
|	ToStringTransformer
	{	$$ = $1						};

TemporalExpression
:	LEFT_PARANTHESIS TemporalExpression RIGHT_PARANTHESIS
	{	$$ = $2								}
|	TemporalExpression PLUS_SIGN SignedTimeInterval
	{	$$ = makeTemporalExpression($1, $3, "+", getLexer(nsQLlex))	}
|	TemporalExpression MINUS_SIGN SignedTimeInterval
	{	$$ = makeTemporalExpression($1, $3, "-", getLexer(nsQLlex))	}
|	SignedTimeInterval PLUS_SIGN TemporalExpression
	{	$$ = makeTemporalExpression($1, $3, "+", getLexer(nsQLlex))	}
|	ColumnExpression PLUS_SIGN SignedTimeInterval
	{	lexer := getLexer(nsQLlex)
		failOnColumnExpression($1, lexer)
		$$ = makeTemporalExpression($1, $3, "+", lexer)			}
|	ColumnExpression MINUS_SIGN SignedTimeInterval
	{	lexer := getLexer(nsQLlex)
		failOnColumnExpression($1, lexer)
		$$ = makeTemporalExpression($1, $3, "-", lexer)			}
|	SignedTimeInterval PLUS_SIGN ColumnExpression
	{	lexer := getLexer(nsQLlex)
		failOnColumnExpression($3, lexer)
		$$ = makeTemporalExpression($1, $3, "+", lexer)			}
|	Timestamp
	{	$$ = $1								}
|	ToTemporalTransformer
	{	$$ = $1								};

NumericExpression
:	LEFT_PARANTHESIS NumericExpression RIGHT_PARANTHESIS
	{	$$ = $2								}
|	PLUS_SIGN NumericExpression %prec UNARY_PLUS_SIGN
	{	$$ = $2								}
|	MINUS_SIGN NumericExpression %prec UNARY_MINUS_SIGN
	{	$$ = makeNumericExpression(nil, $2, "-", getLexer(nsQLlex))	}
|	NumericExpression PLUS_SIGN NumericExpression
	{	$$ = makeNumericExpression($1, $3, "+", getLexer(nsQLlex))	}
|	NumericExpression MINUS_SIGN NumericExpression
	{	$$ = makeNumericExpression($1, $3, "-", getLexer(nsQLlex))	}
|	NumericExpression ASTERISK NumericExpression
	{	$$ = makeNumericExpression($1, $3, "*", getLexer(nsQLlex))	}
|	NumericExpression SLASH NumericExpression
	{	$$ = makeNumericExpression($1, $3, "/", getLexer(nsQLlex))	}
|	NumericExpression PERCENT_SIGN NumericExpression
	{	$$ = makeNumericExpression($1, $3, "%", getLexer(nsQLlex))	}
|	NumericExpression AMPERSAND NumericExpression
	{	$$ = makeNumericExpression($1, $3, "&", getLexer(nsQLlex))	}
|	NumericExpression VERTICAL_BAR NumericExpression
	{	$$ = makeNumericExpression($1, $3, "|", getLexer(nsQLlex))	}
|	NumericExpression PLUS_SIGN ColumnExpression
	{	$$ = makeNumericExpression($1, $3, "+", getLexer(nsQLlex))	}
|	NumericExpression MINUS_SIGN ColumnExpression
	{	$$ = makeNumericExpression($1, $3, "-", getLexer(nsQLlex))	}
|	NumericExpression ASTERISK ColumnExpression
	{	$$ = makeNumericExpression($1, $3, "*", getLexer(nsQLlex))	}
|	NumericExpression SLASH ColumnExpression
	{	$$ = makeNumericExpression($1, $3, "/", getLexer(nsQLlex))	}
|	NumericExpression PERCENT_SIGN ColumnExpression
	{	$$ = makeNumericExpression($1, $3, "%", getLexer(nsQLlex))	}
|	NumericExpression AMPERSAND ColumnExpression
	{	$$ = makeNumericExpression($1, $3, "&", getLexer(nsQLlex))	}
|	NumericExpression VERTICAL_BAR ColumnExpression
	{	$$ = makeNumericExpression($1, $3, "|", getLexer(nsQLlex))	}
|	ColumnExpression PLUS_SIGN NumericExpression
	{	$$ = makeNumericExpression($1, $3, "+", getLexer(nsQLlex))	}
|	ColumnExpression MINUS_SIGN NumericExpression
	{	$$ = makeNumericExpression($1, $3, "-", getLexer(nsQLlex))	}
|	ColumnExpression ASTERISK NumericExpression
	{	$$ = makeNumericExpression($1, $3, "*", getLexer(nsQLlex))	}
|	ColumnExpression SLASH NumericExpression
	{	$$ = makeNumericExpression($1, $3, "/", getLexer(nsQLlex))	}
|	ColumnExpression PERCENT_SIGN NumericExpression
	{	$$ = makeNumericExpression($1, $3, "%", getLexer(nsQLlex))	}
|	ColumnExpression AMPERSAND NumericExpression
	{	$$ = makeNumericExpression($1, $3, "&", getLexer(nsQLlex))	}
|	ColumnExpression VERTICAL_BAR NumericExpression
	{	$$ = makeNumericExpression($1, $3, "|", getLexer(nsQLlex))	}
|	Number
	{	$$ = $1								}
|	ToNumericAggregator
	{	$$ = $1								}
|	ToNumericTransformer
	{	$$ = $1								};

ColumnExpression
:	LEFT_PARANTHESIS ColumnExpression RIGHT_PARANTHESIS
	{	$$ = $2								}
|	PLUS_SIGN ColumnExpression %prec UNARY_PLUS_SIGN
	{	$$ = makeColumnExpression(nil, $2, "+", getLexer(nsQLlex))	}
|	MINUS_SIGN ColumnExpression %prec UNARY_MINUS_SIGN
	{	$$ = makeColumnExpression(nil, $2, "-", getLexer(nsQLlex))	}
|	ColumnExpression PLUS_SIGN ColumnExpression
	{	$$ = makeColumnExpression($1, $3, "+", getLexer(nsQLlex))	}
|	ColumnExpression MINUS_SIGN ColumnExpression
	{	$$ = makeColumnExpression($1, $3, "-", getLexer(nsQLlex))	}
|	ColumnExpression ASTERISK ColumnExpression
	{	$$ = makeColumnExpression($1, $3, "*", getLexer(nsQLlex))	}
|	ColumnExpression SLASH ColumnExpression
	{	$$ = makeColumnExpression($1, $3, "/", getLexer(nsQLlex))	}
|	ColumnExpression PERCENT_SIGN ColumnExpression
	{	$$ = makeColumnExpression($1, $3, "%", getLexer(nsQLlex))	}
|	ColumnExpression AMPERSAND ColumnExpression
	{	$$ = makeColumnExpression($1, $3, "&", getLexer(nsQLlex))	}
|	ColumnExpression VERTICAL_BAR ColumnExpression
	{	$$ = makeColumnExpression($1, $3, "|", getLexer(nsQLlex))	}
|	ColumnName
	{	$$ = $1								}
|	ToColumnAggregator
	{	$$ = $1								};

SignedTimeInterval
:	LEFT_PARANTHESIS SignedTimeInterval RIGHT_PARANTHESIS
	{	$$ = $2									}
|	PLUS_SIGN SignedTimeInterval %prec UNARY_PLUS_SIGN
	{	$$ = makeSignedLiteralExpression(nil, $2, "+", getLexer(nsQLlex))	}
|	MINUS_SIGN SignedTimeInterval %prec UNARY_MINUS_SIGN
	{	$$ = makeSignedLiteralExpression(nil, $2, "-", getLexer(nsQLlex))	}
|	TimeInterval
	{	$$ = $1									};

TableAggregator
:	TCount	{	$$ = $1		}
|	TCorr	{	$$ = $1		}
|	TCov	{	$$ = $1		};

ToColumnAggregator
:	Min	{	$$ = $1		}
|	Max	{	$$ = $1		}
|	First	{	$$ = $1		}
|	Last	{	$$ = $1		};

ToNumericAggregator
:	Count		{	$$ = $1		}
|	Sum		{	$$ = $1		}
|	Mean		{	$$ = $1		}
|	Variance	{	$$ = $1		}
|	Stdev		{	$$ = $1		}
|	Corr		{	$$ = $1		};

ToNumericTransformer
:	Year			{	$$ = $1		}
|	Month			{	$$ = $1		}
|	Day			{	$$ = $1		}
|	Hour			{	$$ = $1		}
|	Minute			{	$$ = $1		}
|	Second			{	$$ = $1		}
|	Subtract_Timestamps	{	$$ = $1		};

ToTemporalTransformer
:	Now	{	$$ = $1		};

ToStringTransformer
:	Map_Blob_Json_Fetch	{	$$ = $1		}
|	Json_Fetch	{	$$ = $1		};

TCount
:	TCOUNT LEFT_PARANTHESIS RIGHT_PARANTHESIS
	{	$$ = makeTableAggregator("tcount", nil, getLexer(nsQLlex))				};

TCorr
:	TCORR LEFT_PARANTHESIS Identifier COMMA Identifier RIGHT_PARANTHESIS
	{	$$ = makeTableAggregator("tcorr", []string{$3.Value, $5.Value}, getLexer(nsQLlex))	};

TCov
:	TCOV LEFT_PARANTHESIS Identifier COMMA Identifier RIGHT_PARANTHESIS
	{	$$ = makeTableAggregator("tcov", []string{$3.Value, $5.Value}, getLexer(nsQLlex))	};

Min
:	MIN LEFT_PARANTHESIS GenericParameter RIGHT_PARANTHESIS
	{	$$ = makeToColumnAggregator("min", []ast.Expression{$3}, getLexer(nsQLlex))		};

Max
:	MAX LEFT_PARANTHESIS GenericParameter RIGHT_PARANTHESIS
	{	$$ = makeToColumnAggregator("max", []ast.Expression{$3}, getLexer(nsQLlex))		};

First
:	FIRST LEFT_PARANTHESIS GenericParameter RIGHT_PARANTHESIS
	{	$$ = makeToColumnAggregator("first", []ast.Expression{$3}, getLexer(nsQLlex))		};

Last
:	LAST LEFT_PARANTHESIS GenericParameter RIGHT_PARANTHESIS
	{	$$ = makeToColumnAggregator("last", []ast.Expression{$3}, getLexer(nsQLlex))		};

Count
:	COUNT LEFT_PARANTHESIS GenericParameter RIGHT_PARANTHESIS
	{	$$ = makeToNumericAggregator("count", []ast.Expression{$3}, getLexer(nsQLlex))		};

Sum
:	SUM LEFT_PARANTHESIS NumericParameter RIGHT_PARANTHESIS
	{	$$ = makeToNumericAggregator("sum", []ast.Expression{$3}, getLexer(nsQLlex))		};

Mean
:	MEAN LEFT_PARANTHESIS NumericParameter RIGHT_PARANTHESIS
	{	$$ = makeToNumericAggregator("mean", []ast.Expression{$3}, getLexer(nsQLlex))		};

Variance
:	VARIANCE LEFT_PARANTHESIS NumericParameter RIGHT_PARANTHESIS
	{	$$ = makeToNumericAggregator("variance", []ast.Expression{$3}, getLexer(nsQLlex))	};

Stdev
:	STDEV LEFT_PARANTHESIS NumericParameter RIGHT_PARANTHESIS
	{	$$ = makeToNumericAggregator("stdev", []ast.Expression{$3}, getLexer(nsQLlex))		};

Corr
:	CORR LEFT_PARANTHESIS NumericParameter COMMA NumericParameter RIGHT_PARANTHESIS
	{	$$ = makeToNumericAggregator("corr", []ast.Expression{$3, $5}, getLexer(nsQLlex))	};

Year
:	YEAR LEFT_PARANTHESIS TemporalParameter RIGHT_PARANTHESIS
	{	$$ = makeToNumericTransformer("year", []ast.Expression{$3}, getLexer(nsQLlex))		};

Month
:	MONTH LEFT_PARANTHESIS TemporalParameter RIGHT_PARANTHESIS
	{	$$ = makeToNumericTransformer("month", []ast.Expression{$3}, getLexer(nsQLlex))		};

Day
:	DAY LEFT_PARANTHESIS TemporalParameter RIGHT_PARANTHESIS
	{	$$ = makeToNumericTransformer("day", []ast.Expression{$3}, getLexer(nsQLlex))		};

Hour
:	HOUR LEFT_PARANTHESIS TemporalParameter RIGHT_PARANTHESIS
	{	$$ = makeToNumericTransformer("hour", []ast.Expression{$3}, getLexer(nsQLlex))		};

Minute
:	MINUTE LEFT_PARANTHESIS TemporalParameter RIGHT_PARANTHESIS
	{	$$ = makeToNumericTransformer("minute", []ast.Expression{$3}, getLexer(nsQLlex))	};

Second
:	SECOND LEFT_PARANTHESIS TemporalParameter RIGHT_PARANTHESIS
	{	$$ = makeToNumericTransformer("second", []ast.Expression{$3}, getLexer(nsQLlex))	};

Now
:	NOW LEFT_PARANTHESIS RIGHT_PARANTHESIS
	{	$$ = makeToTemporalTransformer("now", nil, getLexer(nsQLlex))				};

Map_Blob_Json_Fetch
:	MAP_BLOB_JSON_FETCH LEFT_PARANTHESIS ColumnName COMMA String COMMA String RIGHT_PARANTHESIS
	{	$$ = makeToStringTransformer("map_blob_json_fetch", []ast.Expression{$3, $5, $7}, getLexer(nsQLlex))};

Json_Fetch
:	JSON_FETCH LEFT_PARANTHESIS ColumnName COMMA String RIGHT_PARANTHESIS
	{	$$ = makeToStringTransformer("json_fetch", []ast.Expression{$3, $5}, getLexer(nsQLlex))	}

Subtract_Timestamps
:	SUBTRACT_TIMESTAMPS LEFT_PARANTHESIS ColumnName COMMA ColumnName RIGHT_PARANTHESIS
	{	$$ = makeToNumericTransformer("subtract_timestamps", []ast.Expression{$3, $5}, getLexer(nsQLlex))}

GenericParameter
:	Expression
	{	lexer := getLexer(nsQLlex)
		failOnSubquery([]ast.Expression{$1}, lexer)
		$$ = $1						};

NumericParameter
:	NumericExpression	{	$$ = $1		}
|	ColumnExpression	{	$$ = $1		};

TemporalParameter
:	TemporalExpression	{	$$ = $1								}
|	ColumnExpression	{	failOnColumnExpression($1, getLexer(nsQLlex)); $$ = $1		};

ColumnGroup
:	ASTERISK			{	$$ = makeColumnName("", "*", getLexer(nsQLlex))			}
|	Identifier PERIOD ASTERISK	{	$$ = makeColumnName($1.Value, "*", getLexer(nsQLlex))		};

ColumnName
:	Identifier			{	$$ = makeColumnName("", $1.Value, getLexer(nsQLlex))		}
|	Identifier PERIOD Identifier	{	$$ = makeColumnName($1.Value, $3.Value, getLexer(nsQLlex))	};

Identifier
:	IDENTIFIER	{	$$ = getLexer(nsQLlex).Token	};

Literals
:	Literals COMMA Literal	{	$$ = append($1, $3)		}
|	Literal			{	$$ = []ast.Expression{$1}	};

Literal
:	Null		{	$$ = $1		}
|	Number		{	$$ = $1		}
|	String		{	$$ = $1		}
|	Boolean		{	$$ = $1		}
|	Uuid		{	$$ = $1		}
|	Timestamp	{	$$ = $1		}
|	Binary		{	$$ = $1		}
|	Collection	{	$$ = $1		};

Null
:	NULL
	{	lexer := getLexer(nsQLlex)
		token := lexer.Token
		$$ = makeLiteralExpression(token.Type, token.Value, token.Original, lexer)	};

Number
:	Integer	{	$$ = $1		}
|	Float	{	$$ = $1		};

Integer
:	INTEGER
	{	lexer := getLexer(nsQLlex)
		token := lexer.Token
		$$ = makeLiteralExpression(token.Type, token.Value, token.Original, lexer)	}
|	TIME
	{	lexer := getLexer(nsQLlex)
		token := lexer.Token
		$$ = makeLiteralExpression(token.Type, token.Value, token.Original, lexer)	};

Float
:	FLOAT
	{	lexer := getLexer(nsQLlex)
		token := lexer.Token
		$$ = makeLiteralExpression(token.Type, token.Value, token.Original, lexer)	};

String
:	STRING
	{	lexer := getLexer(nsQLlex)
		token := lexer.Token
		$$ = makeLiteralExpression(token.Type, token.Value, token.Original, lexer)	};

Boolean
:	BOOLEAN
	{	lexer := getLexer(nsQLlex)
		token := lexer.Token
		$$ = makeLiteralExpression(token.Type, token.Value, token.Original, lexer)	};

Uuid
:	UUID
	{	lexer := getLexer(nsQLlex)
		token := lexer.Token
		$$ = makeLiteralExpression(token.Type, token.Value, token.Original, lexer)	};

Timestamp
:	TIMESTAMP
	{	lexer := getLexer(nsQLlex)
		token := lexer.Token
		$$ = makeLiteralExpression(token.Type, token.Value, token.Original, lexer)	}
|	DATE
	{	lexer := getLexer(nsQLlex)
		token := lexer.Token
		$$ = makeLiteralExpression(token.Type, token.Value, token.Original, lexer)};

TimeInterval
:	TIME_INTERVAL
	{	lexer := getLexer(nsQLlex)
		token := lexer.Token
		$$ = makeLiteralExpression(token.Type, token.Value, token.Original, lexer)	};

Binary
:	BINARY
	{	lexer := getLexer(nsQLlex)
		token := lexer.Token
		$$ = makeLiteralExpression(token.Type, token.Value, token.Original, lexer)	};

Collection
:	COLLECTION
	{	lexer := getLexer(nsQLlex)
		token := lexer.Token
		$$ = makeLiteralExpression(token.Type, token.Value, token.Original, lexer)	};

InclusionComparator
:	IN	{	$$ = "in"	}
|	NOT IN	{	$$ = "not in"	};

IdentityComparator
:	IS	{	$$ = "is"	}
|	IS NOT	{	$$ = "is not"	};

RegularComparator
:	EqualityComparator	{	$$ = $1		}
|	RangeComparator		{	$$ = $1		};

EqualityComparator
:	EQUAL_SIGN			{	$$ = "=="	}
|	EXCLAMATION_MARK EQUAL_SIGN	{	$$ = "!="	}
|	LESS_THAN GREATER_THAN		{	$$ = "!="	};

RangeComparator
:	LESS_THAN		{	$$ = "<"	}
|	LESS_THAN EQUAL_SIGN	{	$$ = "<="	}
|	GREATER_THAN		{	$$ = ">"	}
|	GREATER_THAN EQUAL_SIGN	{	$$ = ">="	};
%%

type Token struct {
	Type  int
	Value string
	Original interface{}
}

type nsQLLex struct {
	Code      string
	Position  int
	Token     *Token
	Statement ast.Expression
}

func (l *nsQLLex) Lex(lval *nsQLSymType) int {
	return scan(l)
}

func (l *nsQLLex) Error(s string) {
	panic(errors.New("nsQL " + s))
}

func getLexer(l nsQLLexer) *nsQLLex {
	lexer, _ := l.(*nsQLLex)
	return lexer
}

func Parse(code string) (statement ast.Expression, err error) {
	defer func() {
		if r := recover(); r != nil {
			statement = nil
			e, _ := r.(error)
			err = e
		}
	}()
	if len(code) == 0 {
		err = errors.New("nsql syntax error: empty query")
		return
	}
	lexer := &nsQLLex{Code: code}
	nsQLParse(lexer)
	statement = lexer.Statement
	return
}
