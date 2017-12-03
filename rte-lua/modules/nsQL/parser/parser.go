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

import __yyfmt__ "fmt"

//line parser.y:2
import (
	"errors"
	"github.com/verizonlabs/northstar/rte-lua/modules/nsQL/parser/ast"
)

//line parser.y:112
type nsQLSymType struct {
	yys                     int
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

const FROM = 57346
const AS = 57347
const GROUP = 57348
const ORDER = 57349
const BY = 57350
const LIMIT = 57351
const INTO = 57352
const VALUES = 57353
const SET = 57354
const TABLE = 57355
const PRIMARY = 57356
const KEY = 57357
const CLUSTERING = 57358
const COMPACT = 57359
const STORAGE = 57360
const INSERT = 57361
const DELETE = 57362
const UPDATE = 57363
const SELECT = 57364
const CREATE = 57365
const DROP = 57366
const ALL = 57367
const DISTINCT = 57368
const EXISTS = 57369
const WITH = 57370
const ASC = 57371
const DESC = 57372
const JOIN = 57373
const INNER = 57374
const OUTER = 57375
const FULL = 57376
const LEFT = 57377
const RIGHT = 57378
const SEMI = 57379
const WHERE = 57380
const ON = 57381
const HAVING = 57382
const IF = 57383
const IDENTIFIER = 57384
const NULL = 57385
const INTEGER = 57386
const FLOAT = 57387
const STRING = 57388
const BOOLEAN = 57389
const UUID = 57390
const TIMESTAMP = 57391
const DATE = 57392
const TIME = 57393
const TIME_INTERVAL = 57394
const BINARY = 57395
const COLLECTION = 57396
const PLUS_SIGN = 57397
const MINUS_SIGN = 57398
const ASTERISK = 57399
const SLASH = 57400
const PERCENT_SIGN = 57401
const AMPERSAND = 57402
const VERTICAL_BAR = 57403
const IS = 57404
const EQUAL_SIGN = 57405
const EXCLAMATION_MARK = 57406
const LESS_THAN = 57407
const GREATER_THAN = 57408
const IN = 57409
const AND = 57410
const OR = 57411
const NOT = 57412
const UNION = 57413
const INTERSECT = 57414
const BETWEEN = 57415
const COMMA = 57416
const PERIOD = 57417
const LEFT_PARANTHESIS = 57418
const RIGHT_PARANTHESIS = 57419
const SEMICOLON = 57420
const COUNT = 57421
const SUM = 57422
const MIN = 57423
const MAX = 57424
const FIRST = 57425
const LAST = 57426
const MEAN = 57427
const VARIANCE = 57428
const STDEV = 57429
const CORR = 57430
const YEAR = 57431
const MONTH = 57432
const DAY = 57433
const HOUR = 57434
const MINUTE = 57435
const SECOND = 57436
const NOW = 57437
const MAP_BLOB_JSON_FETCH = 57438
const JSON_FETCH = 57439
const SUBTRACT_TIMESTAMPS = 57440
const TCOUNT = 57441
const TCORR = 57442
const TCOV = 57443
const ASCII = 57444
const BIGINT = 57445
const BLOB = 57446
const BOOLEANTYPE = 57447
const COUNTER = 57448
const DECIMAL = 57449
const DOUBLE = 57450
const FLOATTYPE = 57451
const INET = 57452
const INT = 57453
const TEXT = 57454
const TIMESTAMPTYPE = 57455
const TIMEUUID = 57456
const UUIDTYPE = 57457
const VARCHAR = 57458
const VARINT = 57459
const LIST = 57460
const MAP = 57461
const UNKNOWN = 57462
const UNARY_MINUS_SIGN = 57463
const UNARY_PLUS_SIGN = 57464
const UNARY_NOT = 57465

var nsQLToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"FROM",
	"AS",
	"GROUP",
	"ORDER",
	"BY",
	"LIMIT",
	"INTO",
	"VALUES",
	"SET",
	"TABLE",
	"PRIMARY",
	"KEY",
	"CLUSTERING",
	"COMPACT",
	"STORAGE",
	"INSERT",
	"DELETE",
	"UPDATE",
	"SELECT",
	"CREATE",
	"DROP",
	"ALL",
	"DISTINCT",
	"EXISTS",
	"WITH",
	"ASC",
	"DESC",
	"JOIN",
	"INNER",
	"OUTER",
	"FULL",
	"LEFT",
	"RIGHT",
	"SEMI",
	"WHERE",
	"ON",
	"HAVING",
	"IF",
	"IDENTIFIER",
	"NULL",
	"INTEGER",
	"FLOAT",
	"STRING",
	"BOOLEAN",
	"UUID",
	"TIMESTAMP",
	"DATE",
	"TIME",
	"TIME_INTERVAL",
	"BINARY",
	"COLLECTION",
	"PLUS_SIGN",
	"MINUS_SIGN",
	"ASTERISK",
	"SLASH",
	"PERCENT_SIGN",
	"AMPERSAND",
	"VERTICAL_BAR",
	"IS",
	"EQUAL_SIGN",
	"EXCLAMATION_MARK",
	"LESS_THAN",
	"GREATER_THAN",
	"IN",
	"AND",
	"OR",
	"NOT",
	"UNION",
	"INTERSECT",
	"BETWEEN",
	"COMMA",
	"PERIOD",
	"LEFT_PARANTHESIS",
	"RIGHT_PARANTHESIS",
	"SEMICOLON",
	"COUNT",
	"SUM",
	"MIN",
	"MAX",
	"FIRST",
	"LAST",
	"MEAN",
	"VARIANCE",
	"STDEV",
	"CORR",
	"YEAR",
	"MONTH",
	"DAY",
	"HOUR",
	"MINUTE",
	"SECOND",
	"NOW",
	"MAP_BLOB_JSON_FETCH",
	"JSON_FETCH",
	"SUBTRACT_TIMESTAMPS",
	"TCOUNT",
	"TCORR",
	"TCOV",
	"ASCII",
	"BIGINT",
	"BLOB",
	"BOOLEANTYPE",
	"COUNTER",
	"DECIMAL",
	"DOUBLE",
	"FLOATTYPE",
	"INET",
	"INT",
	"TEXT",
	"TIMESTAMPTYPE",
	"TIMEUUID",
	"UUIDTYPE",
	"VARCHAR",
	"VARINT",
	"LIST",
	"MAP",
	"UNKNOWN",
	"UNARY_MINUS_SIGN",
	"UNARY_PLUS_SIGN",
	"UNARY_NOT",
}
var nsQLStatenames = [...]string{}

const nsQLEofCode = 1
const nsQLErrCode = 2
const nsQLInitialStackSize = 16

//line parser.y:946
type Token struct {
	Type     int
	Value    string
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

//line yacctab:1
var nsQLExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const nsQLPrivate = 57344

const nsQLLast = 1016

var nsQLAct = [...]int{

	152, 550, 485, 77, 76, 75, 80, 533, 84, 460,
	71, 479, 64, 30, 467, 355, 161, 61, 98, 179,
	23, 468, 354, 74, 2, 60, 49, 30, 22, 30,
	17, 18, 487, 24, 21, 30, 67, 16, 87, 20,
	29, 19, 45, 47, 37, 38, 353, 562, 11, 10,
	12, 15, 13, 14, 53, 521, 57, 163, 62, 17,
	18, 151, 362, 147, 517, 448, 465, 52, 213, 54,
	32, 144, 464, 463, 150, 59, 186, 187, 162, 172,
	151, 171, 177, 278, 279, 181, 182, 183, 184, 454,
	275, 79, 175, 150, 199, 210, 212, 560, 272, 271,
	561, 453, 17, 18, 277, 8, 215, 219, 164, 137,
	138, 139, 140, 559, 548, 271, 558, 549, 42, 43,
	44, 72, 489, 490, 491, 492, 493, 494, 495, 496,
	497, 498, 499, 500, 501, 502, 503, 504, 486, 488,
	429, 428, 119, 242, 243, 291, 292, 189, 190, 191,
	192, 193, 194, 195, 30, 181, 182, 183, 184, 545,
	176, 67, 544, 427, 426, 178, 290, 67, 265, 273,
	199, 269, 148, 424, 216, 220, 169, 168, 423, 151,
	151, 266, 267, 62, 264, 270, 258, 422, 259, 62,
	173, 148, 150, 150, 248, 199, 543, 17, 18, 542,
	421, 420, 276, 48, 322, 321, 419, 507, 317, 285,
	506, 417, 295, 300, 302, 304, 306, 308, 310, 312,
	319, 316, 416, 415, 325, 328, 330, 332, 334, 336,
	338, 32, 166, 414, 341, 388, 342, 344, 185, 188,
	197, 348, 358, 359, 360, 361, 413, 67, 269, 357,
	357, 357, 357, 357, 364, 364, 364, 364, 364, 364,
	350, 351, 390, 391, 30, 539, 32, 275, 352, 62,
	15, 15, 380, 46, 241, 370, 529, 528, 289, 293,
	456, 440, 313, 389, 371, 372, 373, 374, 240, 239,
	148, 148, 323, 326, 365, 366, 367, 368, 369, 238,
	31, 400, 401, 402, 392, 237, 236, 284, 235, 349,
	406, 401, 402, 234, 185, 188, 197, 233, 232, 314,
	348, 215, 219, 231, 8, 8, 230, 229, 228, 409,
	227, 226, 225, 32, 224, 343, 223, 222, 408, 185,
	188, 197, 410, 119, 221, 160, 287, 288, 143, 142,
	141, 50, 363, 363, 363, 363, 363, 363, 435, 200,
	201, 202, 203, 204, 205, 206, 560, 407, 247, 325,
	328, 170, 137, 138, 139, 140, 430, 431, 176, 216,
	220, 274, 349, 216, 220, 73, 56, 260, 559, 532,
	186, 187, 166, 151, 482, 436, 433, 432, 181, 182,
	183, 184, 162, 462, 260, 425, 150, 418, 405, 508,
	176, 412, 272, 411, 376, 449, 450, 441, 442, 443,
	444, 445, 446, 447, 375, 166, 18, 441, 442, 169,
	168, 400, 340, 165, 451, 169, 149, 32, 339, 552,
	405, 461, 357, 531, 95, 96, 151, 119, 455, 282,
	287, 288, 281, 530, 174, 149, 181, 182, 211, 150,
	470, 472, 471, 473, 452, 469, 505, 281, 214, 218,
	516, 286, 181, 182, 183, 184, 137, 138, 139, 140,
	515, 514, 283, 461, 393, 394, 395, 396, 397, 398,
	120, 280, 512, 213, 32, 323, 326, 186, 187, 349,
	395, 396, 397, 121, 148, 181, 182, 183, 184, 268,
	122, 470, 472, 471, 473, 91, 469, 523, 524, 525,
	519, 320, 518, 395, 396, 397, 398, 399, 58, 535,
	540, 202, 203, 204, 32, 541, 191, 192, 193, 50,
	551, 186, 187, 513, 535, 382, 535, 89, 90, 557,
	556, 553, 437, 554, 149, 149, 386, 148, 385, 387,
	563, 564, 403, 404, 395, 396, 397, 398, 399, 384,
	257, 32, 256, 383, 294, 299, 301, 303, 305, 307,
	309, 311, 250, 315, 274, 249, 324, 327, 329, 331,
	333, 335, 337, 252, 458, 251, 189, 190, 191, 192,
	193, 194, 195, 347, 181, 182, 183, 184, 546, 547,
	263, 356, 356, 356, 356, 356, 489, 490, 491, 492,
	493, 494, 495, 496, 497, 498, 499, 500, 501, 502,
	503, 504, 200, 201, 202, 203, 204, 205, 206, 209,
	181, 182, 183, 184, 207, 480, 481, 208, 345, 346,
	202, 203, 204, 205, 274, 393, 394, 395, 396, 397,
	398, 399, 510, 345, 346, 202, 203, 204, 205, 206,
	189, 190, 191, 192, 193, 194, 195, 274, 34, 522,
	33, 438, 347, 214, 218, 274, 189, 190, 191, 192,
	193, 194, 273, 403, 404, 395, 396, 397, 398, 399,
	254, 55, 253, 28, 347, 32, 255, 121, 123, 91,
	89, 90, 95, 96, 122, 119, 520, 25, 82, 83,
	66, 345, 346, 202, 203, 204, 205, 206, 27, 378,
	379, 324, 327, 69, 393, 394, 395, 396, 397, 68,
	246, 509, 124, 125, 137, 138, 139, 140, 126, 127,
	128, 129, 130, 131, 132, 133, 134, 135, 120, 117,
	118, 136, 345, 346, 202, 203, 204, 245, 149, 200,
	201, 202, 203, 204, 205, 206, 209, 181, 182, 183,
	184, 207, 146, 262, 208, 32, 261, 121, 123, 91,
	89, 90, 95, 96, 122, 119, 167, 26, 82, 83,
	1, 180, 198, 196, 356, 189, 190, 191, 192, 193,
	194, 195, 475, 69, 189, 190, 191, 192, 193, 68,
	474, 149, 124, 125, 137, 138, 139, 140, 126, 127,
	128, 129, 130, 131, 132, 133, 134, 135, 120, 117,
	118, 136, 32, 94, 121, 123, 91, 99, 466, 95,
	96, 122, 119, 63, 112, 82, 83, 393, 394, 395,
	396, 397, 398, 399, 200, 201, 202, 203, 204, 205,
	206, 93, 92, 97, 111, 110, 318, 109, 108, 124,
	125, 137, 138, 139, 140, 126, 127, 128, 129, 130,
	131, 132, 133, 134, 135, 120, 117, 118, 136, 32,
	107, 121, 123, 106, 105, 104, 103, 102, 122, 119,
	101, 100, 82, 83, 320, 121, 123, 91, 89, 90,
	95, 96, 122, 116, 476, 477, 154, 155, 156, 157,
	158, 159, 115, 217, 114, 113, 124, 125, 137, 138,
	139, 140, 126, 127, 128, 129, 130, 131, 132, 133,
	134, 135, 41, 40, 32, 136, 121, 123, 39, 78,
	81, 86, 85, 122, 88, 36, 65, 297, 298, 70,
	35, 381, 153, 51, 434, 145, 244, 377, 9, 484,
	483, 538, 537, 555, 536, 511, 459, 439, 296, 534,
	527, 124, 125, 137, 138, 139, 140, 126, 127, 128,
	129, 130, 131, 132, 133, 134, 135, 526, 478, 457,
	136, 7, 6, 5, 4, 3,
}
var nsQLPact = [...]int{

	29, -1000, -41, -37, -39, -44, -50, -58, 249, 793,
	793, 693, 224, 667, 665, 19, -1000, 248, 249, -1000,
	-1000, -1000, -1000, -1000, 126, 501, 224, 501, 224, 689,
	311, 249, -1000, 487, 224, 663, -1000, -1000, -1000, -1000,
	-1000, -1000, 274, 273, 272, 354, 249, -1000, -1000, 776,
	743, 895, -1000, -1000, 269, 743, 492, 31, 363, -1000,
	351, -1000, 791, -1000, 361, -1000, -1000, 296, 743, 743,
	-1000, 409, 442, 541, 714, 393, 393, -1000, -1000, 438,
	-1000, -1000, 857, 857, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, 268, 261, -1000,
	260, -1000, -1000, -1000, 258, 256, 255, 254, 252, 251,
	250, 247, 242, 241, 237, 232, 230, 229, 223, 213,
	212, 197, 492, 492, 354, 760, 732, 361, 442, 541,
	714, 409, 293, 224, -1000, 554, 551, 562, 669, 539,
	663, 313, -1000, 781, 778, 583, 663, 492, 743, 743,
	452, 108, 22, 335, 92, 577, 13, -1000, 28, -1000,
	-1000, -1000, 428, 386, 419, 395, 90, 90, 912, 912,
	912, 912, 912, 912, 912, 912, 206, 800, 478, 500,
	857, 857, 912, 912, 912, 912, 912, -1000, 371, 362,
	28, 401, 28, 395, -1000, -1000, -1000, 857, -1000, -1000,
	-1000, 492, 492, 191, 743, 912, 912, 912, 912, 912,
	395, 395, 395, 395, 395, 395, 492, 743, 743, 743,
	743, -1000, 350, 340, 720, 722, 663, 492, 506, -1000,
	-1000, 542, -1000, 538, -1000, 525, 528, -1000, 158, -1000,
	743, 492, 492, 224, -1000, -1000, 367, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 802, 28, 28, 28,
	-1000, -1000, -1000, -1000, 486, 638, 395, 291, 291, -1000,
	90, 90, 90, -1000, 750, 666, 912, 912, 912, 479,
	474, 479, 474, -1000, -1000, -1000, -1000, -1000, -1000, 759,
	707, 631, 593, 249, 486, 750, 809, -1000, 800, -1000,
	-1000, -1000, -1000, -1000, 479, 474, -1000, 479, 474, -1000,
	-1000, -1000, -1000, -1000, -1000, 759, 707, 631, 593, -1000,
	-1000, 802, 802, -1000, 466, 912, 912, 615, 608, 190,
	339, 337, -1000, 169, -1000, 156, 750, 666, 146, 145,
	134, 333, 129, 486, 638, 124, 123, 110, 101, 96,
	331, 87, 86, 64, 63, 492, 492, -1000, 459, 743,
	318, -1000, 743, -1000, -1000, 521, -1000, -1000, 670, -1000,
	-1000, -1000, 205, 28, 28, 28, 28, 28, 28, 28,
	600, -1000, -1000, 291, 291, 21, 507, 291, -12, 38,
	304, 469, 469, -1000, -1000, -1000, -1000, -1000, 912, -1000,
	-1000, -1000, -1000, -1000, -1000, 492, -1000, -1000, -1000, -1000,
	24, 12, -1000, 330, -1000, 743, 361, -1000, 204, 566,
	492, 443, 443, -1000, -1000, -1000, 679, 429, -1000, 329,
	-4, -5, -11, -1000, -1000, 361, 871, -1000, 629, 320,
	-1000, 20, 469, -1000, -1000, -1000, 133, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 341, -1000,
	734, 644, 529, -1000, -1000, -1000, 416, 415, 405, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -13, -1000, 871, 629, 708,
	-1000, -22, -1000, 664, 514, 514, 514, -1000, -1000, -1000,
	201, -1000, 200, 387, 377, 315, -1000, -1000, 492, 189,
	-1000, -1000, 514, 122, 85, 579, 40, -1000, -1000, 492,
	-1000, 373, -1000, 492, -1000, 492, -1000, -1000, 492, -1000,
	39, 23, -1000, -1000, -1000, -30, 314, 292, -1000, 492,
	492, -1000, -1000, -1000, -1000,
}
var nsQLPgo = [...]int{

	0, 24, 1015, 1014, 1013, 1012, 1011, 1009, 1008, 11,
	1007, 990, 989, 7, 987, 986, 9, 985, 984, 983,
	982, 1, 981, 980, 979, 2, 978, 977, 976, 975,
	974, 26, 717, 973, 972, 971, 40, 970, 25, 17,
	16, 22, 12, 969, 966, 10, 121, 385, 23, 91,
	965, 964, 962, 961, 960, 959, 958, 953, 952, 935,
	934, 932, 923, 911, 910, 907, 906, 905, 904, 903,
	900, 878, 877, 875, 874, 873, 872, 871, 854, 46,
	15, 62, 853, 38, 0, 848, 14, 21, 8, 18,
	847, 3, 5, 4, 6, 843, 820, 812, 803, 802,
	165, 19, 801, 800,
}
var nsQLR1 = [...]int{

	0, 103, 103, 103, 103, 103, 103, 1, 1, 1,
	1, 1, 2, 3, 4, 5, 6, 7, 7, 8,
	8, 9, 9, 10, 10, 11, 12, 12, 13, 13,
	14, 15, 15, 16, 23, 23, 24, 24, 24, 25,
	25, 25, 25, 25, 25, 25, 25, 25, 25, 25,
	25, 25, 25, 25, 25, 17, 17, 18, 18, 19,
	19, 20, 21, 21, 22, 22, 26, 26, 27, 27,
	28, 28, 29, 29, 30, 30, 31, 31, 32, 33,
	33, 34, 34, 34, 34, 34, 34, 34, 34, 34,
	34, 34, 35, 35, 36, 36, 36, 36, 37, 37,
	37, 38, 38, 39, 39, 39, 40, 40, 41, 41,
	42, 42, 42, 42, 42, 43, 43, 43, 43, 43,
	43, 43, 43, 43, 43, 43, 43, 43, 43, 43,
	44, 44, 44, 44, 45, 45, 45, 46, 46, 46,
	46, 46, 46, 46, 46, 46, 47, 47, 47, 47,
	47, 47, 47, 47, 47, 47, 47, 47, 47, 47,
	47, 47, 47, 47, 47, 47, 47, 47, 47, 47,
	47, 47, 47, 48, 48, 48, 48, 48, 48, 48,
	48, 48, 48, 48, 48, 49, 49, 49, 49, 50,
	50, 50, 51, 51, 51, 51, 52, 52, 52, 52,
	52, 52, 53, 53, 53, 53, 53, 53, 53, 54,
	55, 55, 56, 57, 58, 59, 60, 61, 62, 63,
	64, 65, 66, 67, 68, 69, 70, 71, 72, 73,
	74, 75, 76, 77, 78, 79, 80, 80, 81, 81,
	82, 82, 83, 83, 84, 85, 85, 86, 86, 86,
	86, 86, 86, 86, 86, 87, 88, 88, 89, 89,
	90, 91, 92, 93, 94, 94, 95, 96, 97, 98,
	98, 99, 99, 100, 100, 101, 101, 101, 102, 102,
	102, 102,
}
var nsQLR2 = [...]int{

	0, 2, 2, 2, 2, 2, 2, 3, 3, 4,
	3, 6, 3, 10, 5, 8, 3, 2, 0, 3,
	1, 4, 2, 1, 3, 3, 3, 3, 2, 2,
	5, 3, 1, 2, 1, 1, 4, 4, 6, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 7, 5, 1, 1, 1,
	1, 3, 3, 3, 3, 1, 3, 2, 0, 2,
	0, 3, 0, 4, 0, 2, 0, 2, 2, 4,
	1, 1, 2, 2, 3, 2, 3, 2, 4, 3,
	3, 2, 0, 2, 3, 5, 3, 5, 0, 1,
	1, 3, 1, 1, 3, 1, 3, 1, 1, 1,
	3, 2, 3, 3, 1, 3, 3, 3, 3, 5,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	1, 1, 1, 1, 3, 1, 1, 3, 3, 3,
	3, 3, 3, 3, 1, 1, 3, 2, 2, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	1, 1, 1, 3, 2, 2, 3, 3, 3, 3,
	3, 3, 3, 1, 1, 3, 2, 2, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 3, 6, 6, 4, 4, 4, 4, 4,
	4, 4, 4, 4, 6, 4, 4, 4, 4, 4,
	4, 3, 8, 6, 6, 1, 1, 1, 1, 1,
	1, 3, 1, 3, 1, 3, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	2, 1, 2, 1, 1, 1, 2, 2, 1, 2,
	1, 2,
}
var nsQLChk = [...]int{

	-1000, -103, -1, -2, -3, -4, -5, -6, 76, -26,
	20, 19, 21, 23, 24, 22, 78, 71, 72, 78,
	78, 78, 78, 78, -1, -32, 4, -32, 10, -36,
	-84, 76, 42, 13, 13, -37, -50, 25, 26, -56,
	-57, -58, 99, 100, 101, -1, 25, -1, 77, -31,
	38, -33, -36, -31, -36, 12, 75, -1, 41, -36,
	-38, -39, -41, -82, -42, -44, 57, -84, 76, 70,
	-43, -45, -46, -47, -48, -92, -93, -91, -55, -49,
	-94, -54, 55, 56, -88, -52, -53, -83, -51, 47,
	48, 46, -76, -77, -95, 49, 50, -75, -89, -90,
	-63, -64, -65, -66, -67, -68, -69, -70, -71, -72,
	-73, -74, -78, -59, -60, -61, -62, 96, 97, 52,
	95, 44, 51, 45, 79, 80, 85, 86, 87, 88,
	89, 90, 91, 92, 93, 94, 98, 81, 82, 83,
	84, 76, 76, 76, -1, -29, 6, -42, -46, -47,
	-48, -45, -84, -34, 31, 32, 33, 34, 35, 36,
	76, -40, -41, -84, 77, 70, 74, 5, 69, 68,
	75, -42, -45, -46, -47, -48, -49, -42, -100, -101,
	-102, 63, 64, 65, 66, -100, 55, 56, -100, 55,
	56, 57, 58, 59, 60, 61, -98, -100, -99, -101,
	55, 56, 57, 58, 59, 60, 61, 67, 70, 62,
	-101, 65, -101, 55, -47, -48, -49, 76, -47, -48,
	-49, 76, 76, 76, 76, 76, 76, 76, 76, 76,
	76, 76, 76, 76, 76, 76, 76, 76, 76, 76,
	76, 77, -84, -84, -28, 7, 8, 75, -36, 31,
	31, 33, 31, 33, 31, 37, 33, 31, -38, -31,
	74, 5, 5, 27, -39, -84, -42, -42, 57, -84,
	77, 77, 77, 77, 77, 77, -48, 76, 55, 56,
	63, 66, 63, 63, -46, -48, 76, 55, 56, -49,
	76, 55, 56, -49, -47, -48, 76, 55, 56, -47,
	-48, -47, -48, -47, -48, -47, -48, -47, -48, -47,
	-48, -47, -48, 76, -46, -47, -48, -45, 76, -87,
	43, -92, -93, -49, -47, -48, -49, -47, -48, -47,
	-48, -47, -48, -47, -48, -47, -48, -47, -48, 67,
	70, -48, -48, -46, -48, 55, 56, -47, -48, -49,
	-83, -83, 77, -79, -41, -80, -47, -48, -80, -80,
	-80, -80, -81, -46, -48, -81, -81, -81, -81, -81,
	-83, -79, -79, -79, -79, 74, 74, -27, 9, 8,
	-38, -35, 39, 31, 31, 33, 31, 31, 77, -41,
	-84, -84, -36, 55, 56, 57, 58, 59, 60, 61,
	-48, -48, -48, 55, 56, -46, -48, 76, -1, -45,
	-48, 74, 74, 77, 77, 77, 77, 77, 74, 77,
	77, 77, 77, 77, 77, 74, 77, 77, 77, 77,
	-84, -84, -89, -40, -30, 40, -42, 31, 11, -14,
	76, -48, -48, -48, -48, -48, -48, -48, 77, -91,
	-91, -80, -83, 77, 77, -42, 76, -7, 28, -15,
	-16, -84, 74, 77, 77, 77, -85, -86, -87, -88,
	-91, -92, -93, -94, -96, -97, 53, 54, -8, -9,
	16, 17, 74, -23, -24, -25, 118, 12, 119, 102,
	103, 104, 105, 106, 107, 108, 109, 110, 111, 112,
	113, 114, 115, 116, 117, -91, 77, 74, 68, 7,
	18, -17, -16, 14, 65, 65, 65, 77, -86, -9,
	8, 77, 15, -25, -25, -25, -10, -11, 76, 76,
	66, 66, 74, -13, -12, -84, -18, -20, -22, 76,
	-84, -25, 77, 74, 77, 74, 29, 30, 74, 77,
	-21, -84, 66, -13, -13, -19, -21, -84, 77, 74,
	74, 77, 77, -84, -84,
}
var nsQLDef = [...]int{

	0, -2, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 98, 1, 0, 0, 2,
	3, 4, 5, 6, 0, 76, 0, 76, 0, 0,
	0, 0, 244, 0, 0, 0, 67, 99, 100, 189,
	190, 191, 0, 0, 0, 8, 0, 10, 7, 72,
	0, 78, 80, 12, 0, 0, 0, 0, 0, 16,
	66, 102, 103, 105, 108, 109, 240, 242, 0, 0,
	114, 130, 131, 132, 133, 0, 0, 135, 136, 0,
	144, 145, 0, 0, 170, 171, 172, 183, 184, 262,
	263, 261, 210, 211, 188, 264, 265, 209, 256, 257,
	196, 197, 198, 199, 200, 201, 202, 203, 204, 205,
	206, 207, 208, 192, 193, 194, 195, 0, 0, 266,
	0, 258, 259, 260, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 9, 70, 0, 77, 0, 0,
	0, 0, 242, 0, 81, 0, 0, 0, 0, 0,
	0, 76, 107, 94, 96, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 111, 0, 273,
	274, 275, 0, 278, 280, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 273,
	0, 0, 0, 0, 0, 0, 0, 269, 0, 271,
	0, 0, 0, 0, 147, 174, 186, 0, 148, 175,
	187, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 212, 0, 0, 68, 0, 0, 0, 92, 82,
	83, 0, 85, 0, 87, 0, 0, 91, 0, 14,
	0, 0, 0, 0, 101, 104, 112, 113, 241, 243,
	110, 134, 137, 146, 173, 185, 127, 0, 0, 0,
	276, 277, 279, 281, 115, 116, 0, 0, 0, 138,
	0, 0, 0, 139, 117, 118, 0, 0, 0, 149,
	156, 150, 157, 151, 158, 152, 159, 153, 160, 154,
	161, 155, 162, 0, 120, 121, 122, 124, 0, 123,
	255, 125, 126, 141, 163, 176, 142, 164, 177, 165,
	178, 166, 179, 167, 180, 168, 181, 169, 182, 270,
	272, 128, 129, 140, 143, 0, 0, 0, 0, 0,
	0, 0, 231, 0, 235, 0, 236, 237, 0, 0,
	0, 0, 0, 238, 239, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 11, 0, 0,
	74, 79, 0, 84, 86, 0, 89, 90, 0, 106,
	95, 97, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 174, 175, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 219, 220, 221, 222, 223, 0, 225,
	226, 227, 228, 229, 230, 0, 215, 216, 217, 218,
	0, 0, 69, 71, 73, 0, 93, 88, 0, 18,
	0, 176, 177, 178, 179, 180, 181, 182, 119, 0,
	0, 0, 0, 213, 214, 75, 0, 15, 0, 0,
	32, 0, 0, 233, 224, 234, 0, 246, 247, 248,
	249, 250, 251, 252, 253, 254, 267, 268, 17, 20,
	0, 0, 0, 33, 34, 35, 0, 0, 0, 39,
	40, 41, 42, 43, 44, 45, 46, 47, 48, 49,
	50, 51, 52, 53, 54, 0, 13, 0, 0, 0,
	22, 0, 31, 0, 0, 0, 0, 232, 245, 19,
	0, 30, 0, 0, 0, 0, 21, 23, 0, 0,
	36, 37, 0, 0, 0, 0, 0, 57, 58, 0,
	65, 0, 24, 0, 25, 0, 28, 29, 0, 56,
	0, 0, 38, 27, 26, 0, 59, 60, 61, 0,
	0, 64, 55, 62, 63,
}
var nsQLTok1 = [...]int{

	1,
}
var nsQLTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 49, 50, 51,
	52, 53, 54, 55, 56, 57, 58, 59, 60, 61,
	62, 63, 64, 65, 66, 67, 68, 69, 70, 71,
	72, 73, 74, 75, 76, 77, 78, 79, 80, 81,
	82, 83, 84, 85, 86, 87, 88, 89, 90, 91,
	92, 93, 94, 95, 96, 97, 98, 99, 100, 101,
	102, 103, 104, 105, 106, 107, 108, 109, 110, 111,
	112, 113, 114, 115, 116, 117, 118, 119, 120, 121,
	122, 123,
}
var nsQLTok3 = [...]int{
	0,
}

var nsQLErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	nsQLDebug        = 0
	nsQLErrorVerbose = true
)

type nsQLLexer interface {
	Lex(lval *nsQLSymType) int
	Error(s string)
}

type nsQLParser interface {
	Parse(nsQLLexer) int
	Lookahead() int
}

type nsQLParserImpl struct {
	lval  nsQLSymType
	stack [nsQLInitialStackSize]nsQLSymType
	char  int
}

func (p *nsQLParserImpl) Lookahead() int {
	return p.char
}

func nsQLNewParser() nsQLParser {
	return &nsQLParserImpl{}
}

const nsQLFlag = -1000

func nsQLTokname(c int) string {
	if c >= 1 && c-1 < len(nsQLToknames) {
		if nsQLToknames[c-1] != "" {
			return nsQLToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func nsQLStatname(s int) string {
	if s >= 0 && s < len(nsQLStatenames) {
		if nsQLStatenames[s] != "" {
			return nsQLStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func nsQLErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !nsQLErrorVerbose {
		return "syntax error"
	}

	for _, e := range nsQLErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + nsQLTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := nsQLPact[state]
	for tok := TOKSTART; tok-1 < len(nsQLToknames); tok++ {
		if n := base + tok; n >= 0 && n < nsQLLast && nsQLChk[nsQLAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if nsQLDef[state] == -2 {
		i := 0
		for nsQLExca[i] != -1 || nsQLExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; nsQLExca[i] >= 0; i += 2 {
			tok := nsQLExca[i]
			if tok < TOKSTART || nsQLExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if nsQLExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += nsQLTokname(tok)
	}
	return res
}

func nsQLlex1(lex nsQLLexer, lval *nsQLSymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = nsQLTok1[0]
		goto out
	}
	if char < len(nsQLTok1) {
		token = nsQLTok1[char]
		goto out
	}
	if char >= nsQLPrivate {
		if char < nsQLPrivate+len(nsQLTok2) {
			token = nsQLTok2[char-nsQLPrivate]
			goto out
		}
	}
	for i := 0; i < len(nsQLTok3); i += 2 {
		token = nsQLTok3[i+0]
		if token == char {
			token = nsQLTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = nsQLTok2[1] /* unknown char */
	}
	if nsQLDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", nsQLTokname(token), uint(char))
	}
	return char, token
}

func nsQLParse(nsQLlex nsQLLexer) int {
	return nsQLNewParser().Parse(nsQLlex)
}

func (nsQLrcvr *nsQLParserImpl) Parse(nsQLlex nsQLLexer) int {
	var nsQLn int
	var nsQLVAL nsQLSymType
	var nsQLDollar []nsQLSymType
	_ = nsQLDollar // silence set and not used
	nsQLS := nsQLrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	nsQLstate := 0
	nsQLrcvr.char = -1
	nsQLtoken := -1 // nsQLrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		nsQLstate = -1
		nsQLrcvr.char = -1
		nsQLtoken = -1
	}()
	nsQLp := -1
	goto nsQLstack

ret0:
	return 0

ret1:
	return 1

nsQLstack:
	/* put a state and value onto the stack */
	if nsQLDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", nsQLTokname(nsQLtoken), nsQLStatname(nsQLstate))
	}

	nsQLp++
	if nsQLp >= len(nsQLS) {
		nyys := make([]nsQLSymType, len(nsQLS)*2)
		copy(nyys, nsQLS)
		nsQLS = nyys
	}
	nsQLS[nsQLp] = nsQLVAL
	nsQLS[nsQLp].yys = nsQLstate

nsQLnewstate:
	nsQLn = nsQLPact[nsQLstate]
	if nsQLn <= nsQLFlag {
		goto nsQLdefault /* simple state */
	}
	if nsQLrcvr.char < 0 {
		nsQLrcvr.char, nsQLtoken = nsQLlex1(nsQLlex, &nsQLrcvr.lval)
	}
	nsQLn += nsQLtoken
	if nsQLn < 0 || nsQLn >= nsQLLast {
		goto nsQLdefault
	}
	nsQLn = nsQLAct[nsQLn]
	if nsQLChk[nsQLn] == nsQLtoken { /* valid shift */
		nsQLrcvr.char = -1
		nsQLtoken = -1
		nsQLVAL = nsQLrcvr.lval
		nsQLstate = nsQLn
		if Errflag > 0 {
			Errflag--
		}
		goto nsQLstack
	}

nsQLdefault:
	/* default state action */
	nsQLn = nsQLDef[nsQLstate]
	if nsQLn == -2 {
		if nsQLrcvr.char < 0 {
			nsQLrcvr.char, nsQLtoken = nsQLlex1(nsQLlex, &nsQLrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if nsQLExca[xi+0] == -1 && nsQLExca[xi+1] == nsQLstate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			nsQLn = nsQLExca[xi+0]
			if nsQLn < 0 || nsQLn == nsQLtoken {
				break
			}
		}
		nsQLn = nsQLExca[xi+1]
		if nsQLn < 0 {
			goto ret0
		}
	}
	if nsQLn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			nsQLlex.Error(nsQLErrorMessage(nsQLstate, nsQLtoken))
			Nerrs++
			if nsQLDebug >= 1 {
				__yyfmt__.Printf("%s", nsQLStatname(nsQLstate))
				__yyfmt__.Printf(" saw %s\n", nsQLTokname(nsQLtoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for nsQLp >= 0 {
				nsQLn = nsQLPact[nsQLS[nsQLp].yys] + nsQLErrCode
				if nsQLn >= 0 && nsQLn < nsQLLast {
					nsQLstate = nsQLAct[nsQLn] /* simulate a shift of "error" */
					if nsQLChk[nsQLstate] == nsQLErrCode {
						goto nsQLstack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if nsQLDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", nsQLS[nsQLp].yys)
				}
				nsQLp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if nsQLDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", nsQLTokname(nsQLtoken))
			}
			if nsQLtoken == nsQLEofCode {
				goto ret1
			}
			nsQLrcvr.char = -1
			nsQLtoken = -1
			goto nsQLnewstate /* try again in the same state */
		}
	}

	/* reduction by production nsQLn */
	if nsQLDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", nsQLn, nsQLStatname(nsQLstate))
	}

	nsQLnt := nsQLn
	nsQLpt := nsQLp
	_ = nsQLpt // guard against "declared and not used"

	nsQLp -= nsQLR2[nsQLn]
	// nsQLp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if nsQLp+1 >= len(nsQLS) {
		nyys := make([]nsQLSymType, len(nsQLS)*2)
		copy(nyys, nsQLS)
		nsQLS = nyys
	}
	nsQLVAL = nsQLS[nsQLp+1]

	/* consult goto table to find next state */
	nsQLn = nsQLR1[nsQLn]
	nsQLg := nsQLPgo[nsQLn]
	nsQLj := nsQLg + nsQLS[nsQLp].yys + 1

	if nsQLj >= nsQLLast {
		nsQLstate = nsQLAct[nsQLg]
	} else {
		nsQLstate = nsQLAct[nsQLj]
		if nsQLChk[nsQLstate] != -nsQLn {
			nsQLstate = nsQLAct[nsQLg]
		}
	}
	// dummy call; replaced with literal code
	switch nsQLnt {

	case 1:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:248
		{
			getLexer(nsQLlex).Statement = nsQLDollar[1].SelectStatement
		}
	case 2:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:249
		{
			getLexer(nsQLlex).Statement = nsQLDollar[1].DeleteStatement
		}
	case 3:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:250
		{
			getLexer(nsQLlex).Statement = nsQLDollar[1].InsertStatement
		}
	case 4:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:251
		{
			getLexer(nsQLlex).Statement = nsQLDollar[1].UpdateStatement
		}
	case 5:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:252
		{
			getLexer(nsQLlex).Statement = nsQLDollar[1].CreateTableStatement
		}
	case 6:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:253
		{
			getLexer(nsQLlex).Statement = nsQLDollar[1].DropTableStatement
		}
	case 7:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:257
		{
			nsQLVAL.SelectStatement = nsQLDollar[2].SelectStatement
		}
	case 8:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:259
		{
			nsQLVAL.SelectStatement = makeSelectExpression(nsQLDollar[1].SelectStatement, nsQLDollar[3].SelectStatement, "union", getLexer(nsQLlex))
		}
	case 9:
		nsQLDollar = nsQLS[nsQLpt-4 : nsQLpt+1]
		//line parser.y:261
		{
			nsQLVAL.SelectStatement = makeSelectExpression(nsQLDollar[1].SelectStatement, nsQLDollar[4].SelectStatement, "union all", getLexer(nsQLlex))
		}
	case 10:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:263
		{
			nsQLVAL.SelectStatement = makeSelectExpression(nsQLDollar[1].SelectStatement, nsQLDollar[3].SelectStatement, "intersect", getLexer(nsQLlex))
		}
	case 11:
		nsQLDollar = nsQLS[nsQLpt-6 : nsQLpt+1]
		//line parser.y:265
		{
			nsQLVAL.SelectStatement = makeSelectStatement(nsQLDollar[1].Select, nsQLDollar[2].From, nsQLDollar[3].Where, nsQLDollar[4].GroupBy, nsQLDollar[5].OrderBy, nsQLDollar[6].Limit, getLexer(nsQLlex))
		}
	case 12:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:269
		{
			lexer := getLexer(nsQLlex)
			failOnNonTableName(nsQLDollar[2].From.Tables, lexer)
			failOnSubquery([]ast.Expression{nsQLDollar[3].Where}, lexer)
			nsQLVAL.DeleteStatement = makeDeleteStatement(nsQLDollar[2].From, nsQLDollar[3].Where, lexer)
		}
	case 13:
		nsQLDollar = nsQLS[nsQLpt-10 : nsQLpt+1]
		//line parser.y:276
		{
			lexer := getLexer(nsQLlex)
			failOnNonTableName([]ast.Expression{nsQLDollar[3].Table}, lexer)
			nsQLVAL.InsertStatement = makeInsertStatement(nsQLDollar[3].Table, nsQLDollar[5].Columns, nsQLDollar[9].Literals, lexer)
		}
	case 14:
		nsQLDollar = nsQLS[nsQLpt-5 : nsQLpt+1]
		//line parser.y:282
		{
			lexer := getLexer(nsQLlex)
			failOnNonTableName([]ast.Expression{nsQLDollar[2].Table}, lexer)
			nsQLVAL.UpdateStatement = makeUpdateStatement(nsQLDollar[2].Table, nsQLDollar[4].Expressions, nsQLDollar[5].Where, lexer)
		}
	case 15:
		nsQLDollar = nsQLS[nsQLpt-8 : nsQLpt+1]
		//line parser.y:288
		{
			lexer := getLexer(nsQLlex)
			failOnNonTableName([]ast.Expression{nsQLDollar[6].Table}, lexer)
			nsQLVAL.CreateTableStatement = makeCreateTableStatement(nsQLDollar[6].Table, nsQLDollar[7].TableDescription, nsQLDollar[8].Directives, lexer)
		}
	case 16:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:294
		{
			lexer := getLexer(nsQLlex)
			failOnNonTableName([]ast.Expression{nsQLDollar[3].Table}, lexer)
			nsQLVAL.DropTableStatement = makeDropTableStatement(nsQLDollar[3].Table, lexer)
		}
	case 17:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:299
		{
			nsQLVAL.Directives = nsQLDollar[2].Properties
		}
	case 18:
		nsQLDollar = nsQLS[nsQLpt-0 : nsQLpt+1]
		//line parser.y:300
		{
			nsQLVAL.Directives = nil
		}
	case 19:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:303
		{
			nsQLVAL.Properties = append(nsQLDollar[1].Properties, nsQLDollar[3].Property)
		}
	case 20:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:304
		{
			nsQLVAL.Properties = []ast.Property{nsQLDollar[1].Property}
		}
	case 21:
		nsQLDollar = nsQLS[nsQLpt-4 : nsQLpt+1]
		//line parser.y:307
		{
			nsQLVAL.Property = &ast.ClusteringOrder{Order: nsQLDollar[4].Order}
		}
	case 22:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:308
		{
			nsQLVAL.Property = &ast.CompactStorage{}
		}
	case 23:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:311
		{
			nsQLVAL.Order = nsQLDollar[1].CompoundOrder
		}
	case 24:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:312
		{
			nsQLVAL.Order = []*ast.Order{nsQLDollar[2].SimpleOrder}
		}
	case 25:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:315
		{
			nsQLVAL.CompoundOrder = nsQLDollar[2].Sorting
		}
	case 26:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:318
		{
			nsQLVAL.Sorting = append(nsQLDollar[1].Sorting, nsQLDollar[3].SimpleOrder)
		}
	case 27:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:319
		{
			nsQLVAL.Sorting = []*ast.Order{nsQLDollar[1].SimpleOrder, nsQLDollar[3].SimpleOrder}
		}
	case 28:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:322
		{
			nsQLVAL.SimpleOrder = &ast.Order{FieldName: nsQLDollar[1].Identifier.Value, Ascending: true}
		}
	case 29:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:323
		{
			nsQLVAL.SimpleOrder = &ast.Order{FieldName: nsQLDollar[1].Identifier.Value, Ascending: false}
		}
	case 30:
		nsQLDollar = nsQLS[nsQLpt-5 : nsQLpt+1]
		//line parser.y:327
		{
			nsQLVAL.TableDescription = &ast.TableDescription{Fields: nsQLDollar[2].FieldDescriptions, PrimaryKey: nsQLDollar[4].PrimaryKey}
		}
	case 31:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:331
		{
			nsQLVAL.FieldDescriptions = append(nsQLDollar[1].FieldDescriptions, nsQLDollar[3].FieldDescription)
		}
	case 32:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:333
		{
			nsQLVAL.FieldDescriptions = []*ast.FieldDescription{nsQLDollar[1].FieldDescription}
		}
	case 33:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:336
		{
			nsQLVAL.FieldDescription = &ast.FieldDescription{FieldName: nsQLDollar[1].Identifier.Value, FieldType: nsQLDollar[2].Type}
		}
	case 34:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:339
		{
			nsQLVAL.Type = nsQLDollar[1].CompoundType
		}
	case 35:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:340
		{
			nsQLVAL.Type = nsQLDollar[1].SimpleType
		}
	case 36:
		nsQLDollar = nsQLS[nsQLpt-4 : nsQLpt+1]
		//line parser.y:343
		{
			nsQLVAL.CompoundType = "set<" + nsQLDollar[3].SimpleType + ">"
		}
	case 37:
		nsQLDollar = nsQLS[nsQLpt-4 : nsQLpt+1]
		//line parser.y:344
		{
			nsQLVAL.CompoundType = "set<" + nsQLDollar[3].SimpleType + ">"
		}
	case 38:
		nsQLDollar = nsQLS[nsQLpt-6 : nsQLpt+1]
		//line parser.y:345
		{
			nsQLVAL.CompoundType = "map<" + nsQLDollar[3].SimpleType + "," + nsQLDollar[5].SimpleType + ">"
		}
	case 39:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:348
		{
			nsQLVAL.SimpleType = "ascii"
		}
	case 40:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:349
		{
			nsQLVAL.SimpleType = "bigint"
		}
	case 41:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:350
		{
			nsQLVAL.SimpleType = "blob"
		}
	case 42:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:351
		{
			nsQLVAL.SimpleType = "boolean"
		}
	case 43:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:352
		{
			nsQLVAL.SimpleType = "counter"
		}
	case 44:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:353
		{
			nsQLVAL.SimpleType = "decimal"
		}
	case 45:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:354
		{
			nsQLVAL.SimpleType = "double"
		}
	case 46:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:355
		{
			nsQLVAL.SimpleType = "float"
		}
	case 47:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:356
		{
			nsQLVAL.SimpleType = "inet"
		}
	case 48:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:357
		{
			nsQLVAL.SimpleType = "int"
		}
	case 49:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:358
		{
			nsQLVAL.SimpleType = "text"
		}
	case 50:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:359
		{
			nsQLVAL.SimpleType = "timestamp"
		}
	case 51:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:360
		{
			nsQLVAL.SimpleType = "timeuuid"
		}
	case 52:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:361
		{
			nsQLVAL.SimpleType = "uuid"
		}
	case 53:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:362
		{
			nsQLVAL.SimpleType = "varchar"
		}
	case 54:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:363
		{
			nsQLVAL.SimpleType = "varint"
		}
	case 55:
		nsQLDollar = nsQLS[nsQLpt-7 : nsQLpt+1]
		//line parser.y:368
		{
			nsQLVAL.PrimaryKey = &ast.PrimaryKey{Partitioning: nsQLDollar[4].PartitioningKey, Clustering: nsQLDollar[6].ClusteringColumns}
		}
	case 56:
		nsQLDollar = nsQLS[nsQLpt-5 : nsQLpt+1]
		//line parser.y:370
		{
			nsQLVAL.PrimaryKey = &ast.PrimaryKey{Partitioning: nsQLDollar[4].PartitioningKey}
		}
	case 57:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:373
		{
			nsQLVAL.PartitioningKey = nsQLDollar[1].CompoundPartitioningKey
		}
	case 58:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:374
		{
			nsQLVAL.PartitioningKey = nsQLDollar[1].SimplePartitioningKey
		}
	case 59:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:377
		{
			nsQLVAL.ClusteringColumns = nsQLDollar[1].Identifiers
		}
	case 60:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:378
		{
			nsQLVAL.ClusteringColumns = []string{nsQLDollar[1].Identifier.Value}
		}
	case 61:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:381
		{
			nsQLVAL.CompoundPartitioningKey = nsQLDollar[2].Identifiers
		}
	case 62:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:384
		{
			nsQLVAL.Identifiers = append(nsQLDollar[1].Identifiers, nsQLDollar[3].Identifier.Value)
		}
	case 63:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:385
		{
			nsQLVAL.Identifiers = []string{nsQLDollar[1].Identifier.Value, nsQLDollar[3].Identifier.Value}
		}
	case 64:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:388
		{
			nsQLVAL.SimplePartitioningKey = []string{nsQLDollar[2].Identifier.Value}
		}
	case 65:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:389
		{
			nsQLVAL.SimplePartitioningKey = []string{nsQLDollar[1].Identifier.Value}
		}
	case 66:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:393
		{
			failOnSubquery(nsQLDollar[3].Columns, getLexer(nsQLlex))
			nsQLVAL.Select = &ast.Select{Qualifier: nsQLDollar[2].Qualifier, Expressions: nsQLDollar[3].Columns}
		}
	case 67:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:396
		{
			nsQLVAL.Select = &ast.Select{Expressions: []ast.Expression{nsQLDollar[2].TableAggregator}}
		}
	case 68:
		nsQLDollar = nsQLS[nsQLpt-0 : nsQLpt+1]
		//line parser.y:400
		{
			nsQLVAL.Limit = ""
		}
	case 69:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:401
		{
			nsQLVAL.Limit = nsQLDollar[2].Integer.Value
		}
	case 70:
		nsQLDollar = nsQLS[nsQLpt-0 : nsQLpt+1]
		//line parser.y:404
		{
			nsQLVAL.OrderBy = nil
		}
	case 71:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:405
		{
			lexer := getLexer(nsQLlex)
			failOnSubquery(nsQLDollar[3].Expressions, lexer)
			failOnNoColumnName(nsQLDollar[3].Expressions, lexer)
			nsQLVAL.OrderBy = nsQLDollar[3].Expressions
		}
	case 72:
		nsQLDollar = nsQLS[nsQLpt-0 : nsQLpt+1]
		//line parser.y:411
		{
			nsQLVAL.GroupBy = nil
		}
	case 73:
		nsQLDollar = nsQLS[nsQLpt-4 : nsQLpt+1]
		//line parser.y:412
		{
			lexer := getLexer(nsQLlex)
			failOnSubquery(nsQLDollar[3].Columns, lexer)
			nsQLVAL.GroupBy = makeGroupBy(nsQLDollar[3].Columns, nsQLDollar[4].Having, lexer)
		}
	case 74:
		nsQLDollar = nsQLS[nsQLpt-0 : nsQLpt+1]
		//line parser.y:417
		{
			nsQLVAL.Having = nil
		}
	case 75:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:418
		{
			nsQLVAL.Having = nsQLDollar[2].LogicalExpression
		}
	case 76:
		nsQLDollar = nsQLS[nsQLpt-0 : nsQLpt+1]
		//line parser.y:421
		{
			nsQLVAL.Where = nil
		}
	case 77:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:422
		{
			nsQLVAL.Where = nsQLDollar[2].LogicalExpression
		}
	case 78:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:425
		{
			finalizeFrom(nsQLDollar[2].Tables, getLexer(nsQLlex))
			nsQLVAL.From = nsQLDollar[2].Tables
		}
	case 79:
		nsQLDollar = nsQLS[nsQLpt-4 : nsQLpt+1]
		//line parser.y:429
		{
			nsQLDollar[1].Tables.Tables = append(nsQLDollar[1].Tables.Tables, nsQLDollar[3].Table)
			nsQLDollar[1].Tables.Joins = append(nsQLDollar[1].Tables.Joins, &ast.Join{Table: nsQLDollar[3].Table, Type: nsQLDollar[2].Join, On: nsQLDollar[4].On})
			nsQLVAL.Tables = nsQLDollar[1].Tables
		}
	case 80:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:433
		{
			nsQLVAL.Tables = &ast.From{Tables: []ast.Expression{nsQLDollar[1].Table}}
		}
	case 81:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:436
		{
			nsQLVAL.Join = "inner"
		}
	case 82:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:437
		{
			nsQLVAL.Join = "inner"
		}
	case 83:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:438
		{
			nsQLVAL.Join = "full_outer"
		}
	case 84:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:439
		{
			nsQLVAL.Join = "full_outer"
		}
	case 85:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:440
		{
			nsQLVAL.Join = "full_outer"
		}
	case 86:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:441
		{
			nsQLVAL.Join = "left_outer"
		}
	case 87:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:442
		{
			nsQLVAL.Join = "left_outer"
		}
	case 88:
		nsQLDollar = nsQLS[nsQLpt-4 : nsQLpt+1]
		//line parser.y:443
		{
			nsQLVAL.Join = "left_semi_outer"
		}
	case 89:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:444
		{
			nsQLVAL.Join = "left_semi_outer"
		}
	case 90:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:445
		{
			nsQLVAL.Join = "right_outer"
		}
	case 91:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:446
		{
			nsQLVAL.Join = "right_outer"
		}
	case 92:
		nsQLDollar = nsQLS[nsQLpt-0 : nsQLpt+1]
		//line parser.y:449
		{
			nsQLVAL.On = nil
		}
	case 93:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:450
		{
			lexer := getLexer(nsQLlex)
			failOnSubquery([]ast.Expression{nsQLDollar[2].LogicalExpression}, lexer)
			nsQLVAL.On = nsQLDollar[2].LogicalExpression
		}
	case 94:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:456
		{
			nsQLVAL.Table = makeTableName(nsQLDollar[1].Identifier.Value, nsQLDollar[3].Identifier.Value, getLexer(nsQLlex))
		}
	case 95:
		nsQLDollar = nsQLS[nsQLpt-5 : nsQLpt+1]
		//line parser.y:458
		{
			table := makeTableName(nsQLDollar[1].Identifier.Value, nsQLDollar[3].Identifier.Value, getLexer(nsQLlex))
			table.SetAlias(nsQLDollar[5].Identifier.Value)
			nsQLVAL.Table = table
		}
	case 96:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:462
		{
			nsQLVAL.Table = nsQLDollar[2].SelectStatement
		}
	case 97:
		nsQLDollar = nsQLS[nsQLpt-5 : nsQLpt+1]
		//line parser.y:464
		{
			nsQLDollar[2].SelectStatement.SetAlias(nsQLDollar[5].Identifier.Value)
			nsQLVAL.Table = nsQLDollar[2].SelectStatement
		}
	case 98:
		nsQLDollar = nsQLS[nsQLpt-0 : nsQLpt+1]
		//line parser.y:468
		{
			nsQLVAL.Qualifier = "all"
		}
	case 99:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:469
		{
			nsQLVAL.Qualifier = "all"
		}
	case 100:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:470
		{
			nsQLVAL.Qualifier = "distinct"
		}
	case 101:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:473
		{
			nsQLVAL.Columns = append(nsQLDollar[1].Columns, nsQLDollar[3].Column)
		}
	case 102:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:474
		{
			nsQLVAL.Columns = []ast.Expression{nsQLDollar[1].Column}
		}
	case 103:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:477
		{
			nsQLVAL.Column = nsQLDollar[1].Expression
		}
	case 104:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:478
		{
			nsQLVAL.Column = nsQLDollar[1].Expression
			nsQLVAL.Column.SetAlias(nsQLDollar[3].Identifier.Value)
		}
	case 105:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:479
		{
			nsQLVAL.Column = nsQLDollar[1].ColumnGroup
		}
	case 106:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:482
		{
			nsQLVAL.Expressions = append(nsQLDollar[1].Expressions, nsQLDollar[3].Expression)
		}
	case 107:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:483
		{
			nsQLVAL.Expressions = []ast.Expression{nsQLDollar[1].Expression}
		}
	case 108:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:486
		{
			nsQLVAL.Expression = nsQLDollar[1].LogicalExpression
		}
	case 109:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:487
		{
			nsQLVAL.Expression = nsQLDollar[1].OrdinaryExpression
		}
	case 110:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:491
		{
			nsQLVAL.LogicalExpression = nsQLDollar[2].LogicalExpression
		}
	case 111:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:493
		{
			nsQLDollar[2].LogicalExpression.Negate()
			nsQLVAL.LogicalExpression = nsQLDollar[2].LogicalExpression
		}
	case 112:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:495
		{
			nsQLVAL.LogicalExpression = makeLogicalExpression(nsQLDollar[1].LogicalExpression, nsQLDollar[3].LogicalExpression, "or", getLexer(nsQLlex))
		}
	case 113:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:497
		{
			nsQLVAL.LogicalExpression = makeLogicalExpression(nsQLDollar[1].LogicalExpression, nsQLDollar[3].LogicalExpression, "and", getLexer(nsQLlex))
		}
	case 114:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:499
		{
			nsQLVAL.LogicalExpression = nsQLDollar[1].ConditionalExpression
		}
	case 115:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:503
		{
			nsQLVAL.ConditionalExpression = makeConditionalExpression(nsQLDollar[1].TemporalExpression, nsQLDollar[3].TemporalExpression, nsQLDollar[2].RegularComparator, false, getLexer(nsQLlex))
		}
	case 116:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:505
		{
			lexer := getLexer(nsQLlex)
			failOnColumnExpression(nsQLDollar[3].ColumnExpression, lexer)
			nsQLVAL.ConditionalExpression = makeConditionalExpression(nsQLDollar[1].TemporalExpression, nsQLDollar[3].ColumnExpression, nsQLDollar[2].RegularComparator, false, lexer)
		}
	case 117:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:509
		{
			nsQLVAL.ConditionalExpression = makeConditionalExpression(nsQLDollar[1].NumericExpression, nsQLDollar[3].NumericExpression, nsQLDollar[2].RegularComparator, false, getLexer(nsQLlex))
		}
	case 118:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:511
		{
			nsQLVAL.ConditionalExpression = makeConditionalExpression(nsQLDollar[1].NumericExpression, nsQLDollar[3].ColumnExpression, nsQLDollar[2].RegularComparator, false, getLexer(nsQLlex))
		}
	case 119:
		nsQLDollar = nsQLS[nsQLpt-5 : nsQLpt+1]
		//line parser.y:513
		{
			lexer := getLexer(nsQLlex)
			failOnNonColumnName(nsQLDollar[1].ColumnExpression, lexer)
			failOnNonSelectStatement(nsQLDollar[4].SelectStatement, lexer)
			nsQLVAL.ConditionalExpression = makeConditionalExpression(nsQLDollar[1].ColumnExpression, nsQLDollar[4].SelectStatement, nsQLDollar[2].InclusionComparator, true, lexer)
		}
	case 120:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:518
		{
			lexer := getLexer(nsQLlex)
			failOnColumnExpression(nsQLDollar[1].ColumnExpression, lexer)
			nsQLVAL.ConditionalExpression = makeConditionalExpression(nsQLDollar[1].ColumnExpression, nsQLDollar[3].TemporalExpression, nsQLDollar[2].RegularComparator, false, lexer)
		}
	case 121:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:522
		{
			nsQLVAL.ConditionalExpression = makeConditionalExpression(nsQLDollar[1].ColumnExpression, nsQLDollar[3].NumericExpression, nsQLDollar[2].RegularComparator, false, getLexer(nsQLlex))
		}
	case 122:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:524
		{
			nsQLVAL.ConditionalExpression = makeConditionalExpression(nsQLDollar[1].ColumnExpression, nsQLDollar[3].ColumnExpression, nsQLDollar[2].RegularComparator, false, getLexer(nsQLlex))
		}
	case 123:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:526
		{
			lexer := getLexer(nsQLlex)
			failOnNonColumnName(nsQLDollar[1].ColumnExpression, lexer)
			nsQLVAL.ConditionalExpression = makeConditionalExpression(nsQLDollar[1].ColumnExpression, nsQLDollar[3].Null, nsQLDollar[2].IdentityComparator, false, lexer)
		}
	case 124:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:530
		{
			lexer := getLexer(nsQLlex)
			failOnNonColumnName(nsQLDollar[1].ColumnExpression, lexer)
			nsQLVAL.ConditionalExpression = makeConditionalExpression(nsQLDollar[1].ColumnExpression, nsQLDollar[3].StringExpression, nsQLDollar[2].RegularComparator, false, lexer)
		}
	case 125:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:534
		{
			lexer := getLexer(nsQLlex)
			failOnNonColumnName(nsQLDollar[1].ColumnExpression, lexer)
			nsQLVAL.ConditionalExpression = makeConditionalExpression(nsQLDollar[1].ColumnExpression, nsQLDollar[3].Boolean, nsQLDollar[2].EqualityComparator, false, lexer)
		}
	case 126:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:538
		{
			lexer := getLexer(nsQLlex)
			failOnNonColumnName(nsQLDollar[1].ColumnExpression, lexer)
			nsQLVAL.ConditionalExpression = makeConditionalExpression(nsQLDollar[1].ColumnExpression, nsQLDollar[3].Uuid, nsQLDollar[2].EqualityComparator, false, lexer)
		}
	case 127:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:542
		{
			lexer := getLexer(nsQLlex)
			failOnNonColumnName(nsQLDollar[3].ColumnExpression, lexer)
			nsQLVAL.ConditionalExpression = makeConditionalExpression(nsQLDollar[1].StringExpression, nsQLDollar[3].ColumnExpression, nsQLDollar[2].RegularComparator, false, lexer)
		}
	case 128:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:546
		{
			lexer := getLexer(nsQLlex)
			failOnNonColumnName(nsQLDollar[3].ColumnExpression, lexer)
			nsQLVAL.ConditionalExpression = makeConditionalExpression(nsQLDollar[1].Boolean, nsQLDollar[3].ColumnExpression, nsQLDollar[2].EqualityComparator, false, lexer)
		}
	case 129:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:550
		{
			lexer := getLexer(nsQLlex)
			failOnNonColumnName(nsQLDollar[3].ColumnExpression, lexer)
			nsQLVAL.ConditionalExpression = makeConditionalExpression(nsQLDollar[1].Uuid, nsQLDollar[3].ColumnExpression, nsQLDollar[2].EqualityComparator, false, lexer)
		}
	case 130:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:555
		{
			nsQLVAL.OrdinaryExpression = nsQLDollar[1].StringExpression
		}
	case 131:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:556
		{
			nsQLVAL.OrdinaryExpression = nsQLDollar[1].TemporalExpression
		}
	case 132:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:557
		{
			nsQLVAL.OrdinaryExpression = nsQLDollar[1].NumericExpression
		}
	case 133:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:558
		{
			nsQLVAL.OrdinaryExpression = nsQLDollar[1].ColumnExpression
		}
	case 134:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:562
		{
			nsQLVAL.StringExpression = nsQLDollar[2].StringExpression
		}
	case 135:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:564
		{
			nsQLVAL.StringExpression = nsQLDollar[1].String
		}
	case 136:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:566
		{
			nsQLVAL.StringExpression = nsQLDollar[1].ToStringTransformer
		}
	case 137:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:570
		{
			nsQLVAL.TemporalExpression = nsQLDollar[2].TemporalExpression
		}
	case 138:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:572
		{
			nsQLVAL.TemporalExpression = makeTemporalExpression(nsQLDollar[1].TemporalExpression, nsQLDollar[3].SignedTimeInterval, "+", getLexer(nsQLlex))
		}
	case 139:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:574
		{
			nsQLVAL.TemporalExpression = makeTemporalExpression(nsQLDollar[1].TemporalExpression, nsQLDollar[3].SignedTimeInterval, "-", getLexer(nsQLlex))
		}
	case 140:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:576
		{
			nsQLVAL.TemporalExpression = makeTemporalExpression(nsQLDollar[1].SignedTimeInterval, nsQLDollar[3].TemporalExpression, "+", getLexer(nsQLlex))
		}
	case 141:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:578
		{
			lexer := getLexer(nsQLlex)
			failOnColumnExpression(nsQLDollar[1].ColumnExpression, lexer)
			nsQLVAL.TemporalExpression = makeTemporalExpression(nsQLDollar[1].ColumnExpression, nsQLDollar[3].SignedTimeInterval, "+", lexer)
		}
	case 142:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:582
		{
			lexer := getLexer(nsQLlex)
			failOnColumnExpression(nsQLDollar[1].ColumnExpression, lexer)
			nsQLVAL.TemporalExpression = makeTemporalExpression(nsQLDollar[1].ColumnExpression, nsQLDollar[3].SignedTimeInterval, "-", lexer)
		}
	case 143:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:586
		{
			lexer := getLexer(nsQLlex)
			failOnColumnExpression(nsQLDollar[3].ColumnExpression, lexer)
			nsQLVAL.TemporalExpression = makeTemporalExpression(nsQLDollar[1].SignedTimeInterval, nsQLDollar[3].ColumnExpression, "+", lexer)
		}
	case 144:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:590
		{
			nsQLVAL.TemporalExpression = nsQLDollar[1].Timestamp
		}
	case 145:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:592
		{
			nsQLVAL.TemporalExpression = nsQLDollar[1].ToTemporalTransformer
		}
	case 146:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:596
		{
			nsQLVAL.NumericExpression = nsQLDollar[2].NumericExpression
		}
	case 147:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:598
		{
			nsQLVAL.NumericExpression = nsQLDollar[2].NumericExpression
		}
	case 148:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:600
		{
			nsQLVAL.NumericExpression = makeNumericExpression(nil, nsQLDollar[2].NumericExpression, "-", getLexer(nsQLlex))
		}
	case 149:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:602
		{
			nsQLVAL.NumericExpression = makeNumericExpression(nsQLDollar[1].NumericExpression, nsQLDollar[3].NumericExpression, "+", getLexer(nsQLlex))
		}
	case 150:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:604
		{
			nsQLVAL.NumericExpression = makeNumericExpression(nsQLDollar[1].NumericExpression, nsQLDollar[3].NumericExpression, "-", getLexer(nsQLlex))
		}
	case 151:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:606
		{
			nsQLVAL.NumericExpression = makeNumericExpression(nsQLDollar[1].NumericExpression, nsQLDollar[3].NumericExpression, "*", getLexer(nsQLlex))
		}
	case 152:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:608
		{
			nsQLVAL.NumericExpression = makeNumericExpression(nsQLDollar[1].NumericExpression, nsQLDollar[3].NumericExpression, "/", getLexer(nsQLlex))
		}
	case 153:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:610
		{
			nsQLVAL.NumericExpression = makeNumericExpression(nsQLDollar[1].NumericExpression, nsQLDollar[3].NumericExpression, "%", getLexer(nsQLlex))
		}
	case 154:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:612
		{
			nsQLVAL.NumericExpression = makeNumericExpression(nsQLDollar[1].NumericExpression, nsQLDollar[3].NumericExpression, "&", getLexer(nsQLlex))
		}
	case 155:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:614
		{
			nsQLVAL.NumericExpression = makeNumericExpression(nsQLDollar[1].NumericExpression, nsQLDollar[3].NumericExpression, "|", getLexer(nsQLlex))
		}
	case 156:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:616
		{
			nsQLVAL.NumericExpression = makeNumericExpression(nsQLDollar[1].NumericExpression, nsQLDollar[3].ColumnExpression, "+", getLexer(nsQLlex))
		}
	case 157:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:618
		{
			nsQLVAL.NumericExpression = makeNumericExpression(nsQLDollar[1].NumericExpression, nsQLDollar[3].ColumnExpression, "-", getLexer(nsQLlex))
		}
	case 158:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:620
		{
			nsQLVAL.NumericExpression = makeNumericExpression(nsQLDollar[1].NumericExpression, nsQLDollar[3].ColumnExpression, "*", getLexer(nsQLlex))
		}
	case 159:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:622
		{
			nsQLVAL.NumericExpression = makeNumericExpression(nsQLDollar[1].NumericExpression, nsQLDollar[3].ColumnExpression, "/", getLexer(nsQLlex))
		}
	case 160:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:624
		{
			nsQLVAL.NumericExpression = makeNumericExpression(nsQLDollar[1].NumericExpression, nsQLDollar[3].ColumnExpression, "%", getLexer(nsQLlex))
		}
	case 161:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:626
		{
			nsQLVAL.NumericExpression = makeNumericExpression(nsQLDollar[1].NumericExpression, nsQLDollar[3].ColumnExpression, "&", getLexer(nsQLlex))
		}
	case 162:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:628
		{
			nsQLVAL.NumericExpression = makeNumericExpression(nsQLDollar[1].NumericExpression, nsQLDollar[3].ColumnExpression, "|", getLexer(nsQLlex))
		}
	case 163:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:630
		{
			nsQLVAL.NumericExpression = makeNumericExpression(nsQLDollar[1].ColumnExpression, nsQLDollar[3].NumericExpression, "+", getLexer(nsQLlex))
		}
	case 164:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:632
		{
			nsQLVAL.NumericExpression = makeNumericExpression(nsQLDollar[1].ColumnExpression, nsQLDollar[3].NumericExpression, "-", getLexer(nsQLlex))
		}
	case 165:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:634
		{
			nsQLVAL.NumericExpression = makeNumericExpression(nsQLDollar[1].ColumnExpression, nsQLDollar[3].NumericExpression, "*", getLexer(nsQLlex))
		}
	case 166:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:636
		{
			nsQLVAL.NumericExpression = makeNumericExpression(nsQLDollar[1].ColumnExpression, nsQLDollar[3].NumericExpression, "/", getLexer(nsQLlex))
		}
	case 167:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:638
		{
			nsQLVAL.NumericExpression = makeNumericExpression(nsQLDollar[1].ColumnExpression, nsQLDollar[3].NumericExpression, "%", getLexer(nsQLlex))
		}
	case 168:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:640
		{
			nsQLVAL.NumericExpression = makeNumericExpression(nsQLDollar[1].ColumnExpression, nsQLDollar[3].NumericExpression, "&", getLexer(nsQLlex))
		}
	case 169:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:642
		{
			nsQLVAL.NumericExpression = makeNumericExpression(nsQLDollar[1].ColumnExpression, nsQLDollar[3].NumericExpression, "|", getLexer(nsQLlex))
		}
	case 170:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:644
		{
			nsQLVAL.NumericExpression = nsQLDollar[1].Number
		}
	case 171:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:646
		{
			nsQLVAL.NumericExpression = nsQLDollar[1].ToNumericAggregator
		}
	case 172:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:648
		{
			nsQLVAL.NumericExpression = nsQLDollar[1].ToNumericTransformer
		}
	case 173:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:652
		{
			nsQLVAL.ColumnExpression = nsQLDollar[2].ColumnExpression
		}
	case 174:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:654
		{
			nsQLVAL.ColumnExpression = makeColumnExpression(nil, nsQLDollar[2].ColumnExpression, "+", getLexer(nsQLlex))
		}
	case 175:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:656
		{
			nsQLVAL.ColumnExpression = makeColumnExpression(nil, nsQLDollar[2].ColumnExpression, "-", getLexer(nsQLlex))
		}
	case 176:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:658
		{
			nsQLVAL.ColumnExpression = makeColumnExpression(nsQLDollar[1].ColumnExpression, nsQLDollar[3].ColumnExpression, "+", getLexer(nsQLlex))
		}
	case 177:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:660
		{
			nsQLVAL.ColumnExpression = makeColumnExpression(nsQLDollar[1].ColumnExpression, nsQLDollar[3].ColumnExpression, "-", getLexer(nsQLlex))
		}
	case 178:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:662
		{
			nsQLVAL.ColumnExpression = makeColumnExpression(nsQLDollar[1].ColumnExpression, nsQLDollar[3].ColumnExpression, "*", getLexer(nsQLlex))
		}
	case 179:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:664
		{
			nsQLVAL.ColumnExpression = makeColumnExpression(nsQLDollar[1].ColumnExpression, nsQLDollar[3].ColumnExpression, "/", getLexer(nsQLlex))
		}
	case 180:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:666
		{
			nsQLVAL.ColumnExpression = makeColumnExpression(nsQLDollar[1].ColumnExpression, nsQLDollar[3].ColumnExpression, "%", getLexer(nsQLlex))
		}
	case 181:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:668
		{
			nsQLVAL.ColumnExpression = makeColumnExpression(nsQLDollar[1].ColumnExpression, nsQLDollar[3].ColumnExpression, "&", getLexer(nsQLlex))
		}
	case 182:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:670
		{
			nsQLVAL.ColumnExpression = makeColumnExpression(nsQLDollar[1].ColumnExpression, nsQLDollar[3].ColumnExpression, "|", getLexer(nsQLlex))
		}
	case 183:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:672
		{
			nsQLVAL.ColumnExpression = nsQLDollar[1].ColumnName
		}
	case 184:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:674
		{
			nsQLVAL.ColumnExpression = nsQLDollar[1].ToColumnAggregator
		}
	case 185:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:678
		{
			nsQLVAL.SignedTimeInterval = nsQLDollar[2].SignedTimeInterval
		}
	case 186:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:680
		{
			nsQLVAL.SignedTimeInterval = makeSignedLiteralExpression(nil, nsQLDollar[2].SignedTimeInterval, "+", getLexer(nsQLlex))
		}
	case 187:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:682
		{
			nsQLVAL.SignedTimeInterval = makeSignedLiteralExpression(nil, nsQLDollar[2].SignedTimeInterval, "-", getLexer(nsQLlex))
		}
	case 188:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:684
		{
			nsQLVAL.SignedTimeInterval = nsQLDollar[1].TimeInterval
		}
	case 189:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:687
		{
			nsQLVAL.TableAggregator = nsQLDollar[1].TCount
		}
	case 190:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:688
		{
			nsQLVAL.TableAggregator = nsQLDollar[1].TCorr
		}
	case 191:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:689
		{
			nsQLVAL.TableAggregator = nsQLDollar[1].TCov
		}
	case 192:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:692
		{
			nsQLVAL.ToColumnAggregator = nsQLDollar[1].Min
		}
	case 193:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:693
		{
			nsQLVAL.ToColumnAggregator = nsQLDollar[1].Max
		}
	case 194:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:694
		{
			nsQLVAL.ToColumnAggregator = nsQLDollar[1].First
		}
	case 195:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:695
		{
			nsQLVAL.ToColumnAggregator = nsQLDollar[1].Last
		}
	case 196:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:698
		{
			nsQLVAL.ToNumericAggregator = nsQLDollar[1].Count
		}
	case 197:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:699
		{
			nsQLVAL.ToNumericAggregator = nsQLDollar[1].Sum
		}
	case 198:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:700
		{
			nsQLVAL.ToNumericAggregator = nsQLDollar[1].Mean
		}
	case 199:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:701
		{
			nsQLVAL.ToNumericAggregator = nsQLDollar[1].Variance
		}
	case 200:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:702
		{
			nsQLVAL.ToNumericAggregator = nsQLDollar[1].Stdev
		}
	case 201:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:703
		{
			nsQLVAL.ToNumericAggregator = nsQLDollar[1].Corr
		}
	case 202:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:706
		{
			nsQLVAL.ToNumericTransformer = nsQLDollar[1].Year
		}
	case 203:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:707
		{
			nsQLVAL.ToNumericTransformer = nsQLDollar[1].Month
		}
	case 204:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:708
		{
			nsQLVAL.ToNumericTransformer = nsQLDollar[1].Day
		}
	case 205:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:709
		{
			nsQLVAL.ToNumericTransformer = nsQLDollar[1].Hour
		}
	case 206:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:710
		{
			nsQLVAL.ToNumericTransformer = nsQLDollar[1].Minute
		}
	case 207:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:711
		{
			nsQLVAL.ToNumericTransformer = nsQLDollar[1].Second
		}
	case 208:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:712
		{
			nsQLVAL.ToNumericTransformer = nsQLDollar[1].Subtract_Timestamps
		}
	case 209:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:715
		{
			nsQLVAL.ToTemporalTransformer = nsQLDollar[1].Now
		}
	case 210:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:718
		{
			nsQLVAL.ToStringTransformer = nsQLDollar[1].Map_Blob_Json_Fetch
		}
	case 211:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:719
		{
			nsQLVAL.ToStringTransformer = nsQLDollar[1].Json_Fetch
		}
	case 212:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:723
		{
			nsQLVAL.TCount = makeTableAggregator("tcount", nil, getLexer(nsQLlex))
		}
	case 213:
		nsQLDollar = nsQLS[nsQLpt-6 : nsQLpt+1]
		//line parser.y:727
		{
			nsQLVAL.TCorr = makeTableAggregator("tcorr", []string{nsQLDollar[3].Identifier.Value, nsQLDollar[5].Identifier.Value}, getLexer(nsQLlex))
		}
	case 214:
		nsQLDollar = nsQLS[nsQLpt-6 : nsQLpt+1]
		//line parser.y:731
		{
			nsQLVAL.TCov = makeTableAggregator("tcov", []string{nsQLDollar[3].Identifier.Value, nsQLDollar[5].Identifier.Value}, getLexer(nsQLlex))
		}
	case 215:
		nsQLDollar = nsQLS[nsQLpt-4 : nsQLpt+1]
		//line parser.y:735
		{
			nsQLVAL.Min = makeToColumnAggregator("min", []ast.Expression{nsQLDollar[3].GenericParameter}, getLexer(nsQLlex))
		}
	case 216:
		nsQLDollar = nsQLS[nsQLpt-4 : nsQLpt+1]
		//line parser.y:739
		{
			nsQLVAL.Max = makeToColumnAggregator("max", []ast.Expression{nsQLDollar[3].GenericParameter}, getLexer(nsQLlex))
		}
	case 217:
		nsQLDollar = nsQLS[nsQLpt-4 : nsQLpt+1]
		//line parser.y:743
		{
			nsQLVAL.First = makeToColumnAggregator("first", []ast.Expression{nsQLDollar[3].GenericParameter}, getLexer(nsQLlex))
		}
	case 218:
		nsQLDollar = nsQLS[nsQLpt-4 : nsQLpt+1]
		//line parser.y:747
		{
			nsQLVAL.Last = makeToColumnAggregator("last", []ast.Expression{nsQLDollar[3].GenericParameter}, getLexer(nsQLlex))
		}
	case 219:
		nsQLDollar = nsQLS[nsQLpt-4 : nsQLpt+1]
		//line parser.y:751
		{
			nsQLVAL.Count = makeToNumericAggregator("count", []ast.Expression{nsQLDollar[3].GenericParameter}, getLexer(nsQLlex))
		}
	case 220:
		nsQLDollar = nsQLS[nsQLpt-4 : nsQLpt+1]
		//line parser.y:755
		{
			nsQLVAL.Sum = makeToNumericAggregator("sum", []ast.Expression{nsQLDollar[3].NumericParameter}, getLexer(nsQLlex))
		}
	case 221:
		nsQLDollar = nsQLS[nsQLpt-4 : nsQLpt+1]
		//line parser.y:759
		{
			nsQLVAL.Mean = makeToNumericAggregator("mean", []ast.Expression{nsQLDollar[3].NumericParameter}, getLexer(nsQLlex))
		}
	case 222:
		nsQLDollar = nsQLS[nsQLpt-4 : nsQLpt+1]
		//line parser.y:763
		{
			nsQLVAL.Variance = makeToNumericAggregator("variance", []ast.Expression{nsQLDollar[3].NumericParameter}, getLexer(nsQLlex))
		}
	case 223:
		nsQLDollar = nsQLS[nsQLpt-4 : nsQLpt+1]
		//line parser.y:767
		{
			nsQLVAL.Stdev = makeToNumericAggregator("stdev", []ast.Expression{nsQLDollar[3].NumericParameter}, getLexer(nsQLlex))
		}
	case 224:
		nsQLDollar = nsQLS[nsQLpt-6 : nsQLpt+1]
		//line parser.y:771
		{
			nsQLVAL.Corr = makeToNumericAggregator("corr", []ast.Expression{nsQLDollar[3].NumericParameter, nsQLDollar[5].NumericParameter}, getLexer(nsQLlex))
		}
	case 225:
		nsQLDollar = nsQLS[nsQLpt-4 : nsQLpt+1]
		//line parser.y:775
		{
			nsQLVAL.Year = makeToNumericTransformer("year", []ast.Expression{nsQLDollar[3].TemporalParameter}, getLexer(nsQLlex))
		}
	case 226:
		nsQLDollar = nsQLS[nsQLpt-4 : nsQLpt+1]
		//line parser.y:779
		{
			nsQLVAL.Month = makeToNumericTransformer("month", []ast.Expression{nsQLDollar[3].TemporalParameter}, getLexer(nsQLlex))
		}
	case 227:
		nsQLDollar = nsQLS[nsQLpt-4 : nsQLpt+1]
		//line parser.y:783
		{
			nsQLVAL.Day = makeToNumericTransformer("day", []ast.Expression{nsQLDollar[3].TemporalParameter}, getLexer(nsQLlex))
		}
	case 228:
		nsQLDollar = nsQLS[nsQLpt-4 : nsQLpt+1]
		//line parser.y:787
		{
			nsQLVAL.Hour = makeToNumericTransformer("hour", []ast.Expression{nsQLDollar[3].TemporalParameter}, getLexer(nsQLlex))
		}
	case 229:
		nsQLDollar = nsQLS[nsQLpt-4 : nsQLpt+1]
		//line parser.y:791
		{
			nsQLVAL.Minute = makeToNumericTransformer("minute", []ast.Expression{nsQLDollar[3].TemporalParameter}, getLexer(nsQLlex))
		}
	case 230:
		nsQLDollar = nsQLS[nsQLpt-4 : nsQLpt+1]
		//line parser.y:795
		{
			nsQLVAL.Second = makeToNumericTransformer("second", []ast.Expression{nsQLDollar[3].TemporalParameter}, getLexer(nsQLlex))
		}
	case 231:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:799
		{
			nsQLVAL.Now = makeToTemporalTransformer("now", nil, getLexer(nsQLlex))
		}
	case 232:
		nsQLDollar = nsQLS[nsQLpt-8 : nsQLpt+1]
		//line parser.y:803
		{
			nsQLVAL.Map_Blob_Json_Fetch = makeToStringTransformer("map_blob_json_fetch", []ast.Expression{nsQLDollar[3].ColumnName, nsQLDollar[5].String, nsQLDollar[7].String}, getLexer(nsQLlex))
		}
	case 233:
		nsQLDollar = nsQLS[nsQLpt-6 : nsQLpt+1]
		//line parser.y:807
		{
			nsQLVAL.Json_Fetch = makeToStringTransformer("json_fetch", []ast.Expression{nsQLDollar[3].ColumnName, nsQLDollar[5].String}, getLexer(nsQLlex))
		}
	case 234:
		nsQLDollar = nsQLS[nsQLpt-6 : nsQLpt+1]
		//line parser.y:811
		{
			nsQLVAL.Subtract_Timestamps = makeToNumericTransformer("subtract_timestamps", []ast.Expression{nsQLDollar[3].ColumnName, nsQLDollar[5].ColumnName}, getLexer(nsQLlex))
		}
	case 235:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:815
		{
			lexer := getLexer(nsQLlex)
			failOnSubquery([]ast.Expression{nsQLDollar[1].Expression}, lexer)
			nsQLVAL.GenericParameter = nsQLDollar[1].Expression
		}
	case 236:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:820
		{
			nsQLVAL.NumericParameter = nsQLDollar[1].NumericExpression
		}
	case 237:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:821
		{
			nsQLVAL.NumericParameter = nsQLDollar[1].ColumnExpression
		}
	case 238:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:824
		{
			nsQLVAL.TemporalParameter = nsQLDollar[1].TemporalExpression
		}
	case 239:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:825
		{
			failOnColumnExpression(nsQLDollar[1].ColumnExpression, getLexer(nsQLlex))
			nsQLVAL.TemporalParameter = nsQLDollar[1].ColumnExpression
		}
	case 240:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:828
		{
			nsQLVAL.ColumnGroup = makeColumnName("", "*", getLexer(nsQLlex))
		}
	case 241:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:829
		{
			nsQLVAL.ColumnGroup = makeColumnName(nsQLDollar[1].Identifier.Value, "*", getLexer(nsQLlex))
		}
	case 242:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:832
		{
			nsQLVAL.ColumnName = makeColumnName("", nsQLDollar[1].Identifier.Value, getLexer(nsQLlex))
		}
	case 243:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:833
		{
			nsQLVAL.ColumnName = makeColumnName(nsQLDollar[1].Identifier.Value, nsQLDollar[3].Identifier.Value, getLexer(nsQLlex))
		}
	case 244:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:836
		{
			nsQLVAL.Identifier = getLexer(nsQLlex).Token
		}
	case 245:
		nsQLDollar = nsQLS[nsQLpt-3 : nsQLpt+1]
		//line parser.y:839
		{
			nsQLVAL.Literals = append(nsQLDollar[1].Literals, nsQLDollar[3].Literal)
		}
	case 246:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:840
		{
			nsQLVAL.Literals = []ast.Expression{nsQLDollar[1].Literal}
		}
	case 247:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:843
		{
			nsQLVAL.Literal = nsQLDollar[1].Null
		}
	case 248:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:844
		{
			nsQLVAL.Literal = nsQLDollar[1].Number
		}
	case 249:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:845
		{
			nsQLVAL.Literal = nsQLDollar[1].String
		}
	case 250:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:846
		{
			nsQLVAL.Literal = nsQLDollar[1].Boolean
		}
	case 251:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:847
		{
			nsQLVAL.Literal = nsQLDollar[1].Uuid
		}
	case 252:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:848
		{
			nsQLVAL.Literal = nsQLDollar[1].Timestamp
		}
	case 253:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:849
		{
			nsQLVAL.Literal = nsQLDollar[1].Binary
		}
	case 254:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:850
		{
			nsQLVAL.Literal = nsQLDollar[1].Collection
		}
	case 255:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:854
		{
			lexer := getLexer(nsQLlex)
			token := lexer.Token
			nsQLVAL.Null = makeLiteralExpression(token.Type, token.Value, token.Original, lexer)
		}
	case 256:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:859
		{
			nsQLVAL.Number = nsQLDollar[1].Integer
		}
	case 257:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:860
		{
			nsQLVAL.Number = nsQLDollar[1].Float
		}
	case 258:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:864
		{
			lexer := getLexer(nsQLlex)
			token := lexer.Token
			nsQLVAL.Integer = makeLiteralExpression(token.Type, token.Value, token.Original, lexer)
		}
	case 259:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:868
		{
			lexer := getLexer(nsQLlex)
			token := lexer.Token
			nsQLVAL.Integer = makeLiteralExpression(token.Type, token.Value, token.Original, lexer)
		}
	case 260:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:874
		{
			lexer := getLexer(nsQLlex)
			token := lexer.Token
			nsQLVAL.Float = makeLiteralExpression(token.Type, token.Value, token.Original, lexer)
		}
	case 261:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:880
		{
			lexer := getLexer(nsQLlex)
			token := lexer.Token
			nsQLVAL.String = makeLiteralExpression(token.Type, token.Value, token.Original, lexer)
		}
	case 262:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:886
		{
			lexer := getLexer(nsQLlex)
			token := lexer.Token
			nsQLVAL.Boolean = makeLiteralExpression(token.Type, token.Value, token.Original, lexer)
		}
	case 263:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:892
		{
			lexer := getLexer(nsQLlex)
			token := lexer.Token
			nsQLVAL.Uuid = makeLiteralExpression(token.Type, token.Value, token.Original, lexer)
		}
	case 264:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:898
		{
			lexer := getLexer(nsQLlex)
			token := lexer.Token
			nsQLVAL.Timestamp = makeLiteralExpression(token.Type, token.Value, token.Original, lexer)
		}
	case 265:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:902
		{
			lexer := getLexer(nsQLlex)
			token := lexer.Token
			nsQLVAL.Timestamp = makeLiteralExpression(token.Type, token.Value, token.Original, lexer)
		}
	case 266:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:908
		{
			lexer := getLexer(nsQLlex)
			token := lexer.Token
			nsQLVAL.TimeInterval = makeLiteralExpression(token.Type, token.Value, token.Original, lexer)
		}
	case 267:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:914
		{
			lexer := getLexer(nsQLlex)
			token := lexer.Token
			nsQLVAL.Binary = makeLiteralExpression(token.Type, token.Value, token.Original, lexer)
		}
	case 268:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:920
		{
			lexer := getLexer(nsQLlex)
			token := lexer.Token
			nsQLVAL.Collection = makeLiteralExpression(token.Type, token.Value, token.Original, lexer)
		}
	case 269:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:925
		{
			nsQLVAL.InclusionComparator = "in"
		}
	case 270:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:926
		{
			nsQLVAL.InclusionComparator = "not in"
		}
	case 271:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:929
		{
			nsQLVAL.IdentityComparator = "is"
		}
	case 272:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:930
		{
			nsQLVAL.IdentityComparator = "is not"
		}
	case 273:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:933
		{
			nsQLVAL.RegularComparator = nsQLDollar[1].EqualityComparator
		}
	case 274:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:934
		{
			nsQLVAL.RegularComparator = nsQLDollar[1].RangeComparator
		}
	case 275:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:937
		{
			nsQLVAL.EqualityComparator = "=="
		}
	case 276:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:938
		{
			nsQLVAL.EqualityComparator = "!="
		}
	case 277:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:939
		{
			nsQLVAL.EqualityComparator = "!="
		}
	case 278:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:942
		{
			nsQLVAL.RangeComparator = "<"
		}
	case 279:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:943
		{
			nsQLVAL.RangeComparator = "<="
		}
	case 280:
		nsQLDollar = nsQLS[nsQLpt-1 : nsQLpt+1]
		//line parser.y:944
		{
			nsQLVAL.RangeComparator = ">"
		}
	case 281:
		nsQLDollar = nsQLS[nsQLpt-2 : nsQLpt+1]
		//line parser.y:945
		{
			nsQLVAL.RangeComparator = ">="
		}
	}
	goto nsQLstack /* stack new state and value */
}
