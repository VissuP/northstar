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
	"encoding/hex"
	"github.com/satori/go.uuid"
	"regexp"
	"strings"
	"unicode"
)

func scan(l *nsQLLex) int {
	if l.Token, l.Position = scanAll(l.Code, l.Position); l.Token == nil {
		return 0
	} else {
		return l.Token.Type
	}
}

func scanAll(code string, position int) (*Token, int) {
	position = scanSpaces(code, position)
	if position == len(code) {
		return nil, 0
	}

	c := rune(code[position])
	if unicode.IsLetter(c) {
		return scanTextual(code, position)
	} else if unicode.IsDigit(c) {
		return scanNumeric(code, position)
	} else if c == '\'' || c == '`' {
		return scanLiteral(code, position)
	} else if c == '{' {
		return scanCollection(code, position)
	} else {
		return scanSymbolic(code, position)
	}
}

func scanSpaces(code string, position int) int {
	for position < len(code) {
		if !unicode.IsSpace(rune(code[position])) {
			break
		}
		position++
	}
	return position
}

func scanTextual(code string, position int) (*Token, int) {
	s := position
	position++
	for position < len(code) {
		c := rune(code[position])
		if c != '_' && !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			break
		}
		position++
	}
	textual := code[s:position]
	switch strings.ToLower(textual) {
	case "select":
		return &Token{Type: SELECT}, position
	case "create":
		return &Token{Type: CREATE}, position
	case "drop":
		return &Token{Type: DROP}, position
	case "delete":
		return &Token{Type: DELETE}, position
	case "insert":
		return &Token{Type: INSERT}, position
	case "update":
		return &Token{Type: UPDATE}, position
	case "from":
		return &Token{Type: FROM}, position
	case "where":
		return &Token{Type: WHERE}, position
	case "as":
		return &Token{Type: AS}, position
	case "all":
		return &Token{Type: ALL}, position
	case "distinct":
		return &Token{Type: DISTINCT}, position
	case "exists":
		return &Token{Type: EXISTS}, position
	case "with":
		return &Token{Type: WITH}, position
	case "asc":
		return &Token{Type: ASC}, position
	case "desc":
		return &Token{Type: DESC}, position
	case "join":
		return &Token{Type: JOIN}, position
	case "inner":
		return &Token{Type: INNER}, position
	case "outer":
		return &Token{Type: OUTER}, position
	case "full":
		return &Token{Type: FULL}, position
	case "left":
		return &Token{Type: LEFT}, position
	case "right":
		return &Token{Type: RIGHT}, position
	case "semi":
		return &Token{Type: SEMI}, position
	case "on":
		return &Token{Type: ON}, position
	case "group":
		return &Token{Type: GROUP}, position
	case "having":
		return &Token{Type: HAVING}, position
	case "if":
		return &Token{Type: IF}, position
	case "order":
		return &Token{Type: ORDER}, position
	case "by":
		return &Token{Type: BY}, position
	case "limit":
		return &Token{Type: LIMIT}, position
	case "into":
		return &Token{Type: INTO}, position
	case "values":
		return &Token{Type: VALUES}, position
	case "set":
		return &Token{Type: SET}, position
	case "table":
		return &Token{Type: TABLE}, position
	case "primary":
		return &Token{Type: PRIMARY}, position
	case "key":
		return &Token{Type: KEY}, position
	case "clustering":
		return &Token{Type: CLUSTERING}, position
	case "compact":
		return &Token{Type: COMPACT}, position
	case "storage":
		return &Token{Type: STORAGE}, position
	case "null":
		return &Token{Type: NULL, Value: textual}, position
	case "true":
		return &Token{Type: BOOLEAN, Value: textual}, position
	case "false":
		return &Token{Type: BOOLEAN, Value: textual}, position
	case "is":
		return &Token{Type: IS}, position
	case "in":
		return &Token{Type: IN}, position
	case "and":
		return &Token{Type: AND}, position
	case "or":
		return &Token{Type: OR}, position
	case "not":
		return &Token{Type: NOT}, position
	case "union":
		return &Token{Type: UNION}, position
	case "intersect":
		return &Token{Type: INTERSECT}, position
	case "between":
		return &Token{Type: BETWEEN}, position
	case "count":
		return &Token{Type: COUNT}, position
	case "sum":
		return &Token{Type: SUM}, position
	case "min":
		return &Token{Type: MIN}, position
	case "max":
		return &Token{Type: MAX}, position
	case "first":
		return &Token{Type: FIRST}, position
	case "last":
		return &Token{Type: LAST}, position
	case "mean", "avg":
		return &Token{Type: MEAN}, position
	case "variance":
		return &Token{Type: VARIANCE}, position
	case "stdev":
		return &Token{Type: STDEV}, position
	case "corr":
		return &Token{Type: CORR}, position
	case "now":
		return &Token{Type: NOW}, position
	case "year":
		return &Token{Type: YEAR}, position
	case "month":
		return &Token{Type: MONTH}, position
	case "day":
		return &Token{Type: DAY}, position
	case "hour":
		return &Token{Type: HOUR}, position
	case "minute":
		return &Token{Type: MINUTE}, position
	case "second":
		return &Token{Type: SECOND}, position
	case "tcount":
		return &Token{Type: TCOUNT}, position
	case "tcorr":
		return &Token{Type: TCORR}, position
	case "tcov":
		return &Token{Type: TCOV}, position
	case "ascii":
		return &Token{Type: ASCII}, position
	case "bigint":
		return &Token{Type: BIGINT}, position
	case "boolean":
		return &Token{Type: BOOLEANTYPE}, position
	case "counter":
		return &Token{Type: COUNTER}, position
	case "decimal":
		return &Token{Type: DECIMAL}, position
	case "double":
		return &Token{Type: DOUBLE}, position
	case "float":
		return &Token{Type: FLOATTYPE}, position
	case "inet":
		return &Token{Type: INET}, position
	case "int":
		return &Token{Type: INT}, position
	case "text":
		return &Token{Type: TEXT}, position
	case "timestamp":
		return &Token{Type: TIMESTAMPTYPE}, position
	case "timeuuid":
		return &Token{Type: TIMEUUID}, position
	case "uuid":
		return &Token{Type: UUIDTYPE}, position
	case "varchar":
		return &Token{Type: VARCHAR}, position
	case "varint":
		return &Token{Type: VARINT}, position
	case "list":
		return &Token{Type: LIST}, position
	case "map":
		return &Token{Type: MAP}, position
	case "map_blob_json_fetch":
		return &Token{Type: MAP_BLOB_JSON_FETCH}, position
	case "json_fetch":
		return &Token{Type: JSON_FETCH}, position
	case "subtract_timestamps":
		return &Token{Type: SUBTRACT_TIMESTAMPS}, position
	default:
		return &Token{Type: IDENTIFIER, Value: textual}, position
	}
}

func scanNumeric(code string, position int) (*Token, int) {
	s := position
	position++
	for position < len(code) {
		c := rune(code[position])
		if c != '.' && !unicode.IsDigit(c) {
			break
		}
		position++
	}
	numeric := code[s:position]
	if isInteger(numeric) {
		return &Token{Type: INTEGER, Value: numeric}, position
	}
	if isFloat(numeric) {
		return &Token{Type: FLOAT, Value: numeric}, position
	}
	return &Token{Type: UNKNOWN}, position
}

func scanCollection(code string, position int) (*Token, int) {
	position++
	var token, key *Token
	collection := make(map[interface{}]interface{})
	isSet := -1
	for position < len(code) {
		position = scanSpaces(code, position)
		if position == len(code) {
			return &Token{Type: UNKNOWN}, position
		}

		c := rune(code[position])
		if unicode.IsDigit(c) {
			token, position = scanNumeric(code, position)
			if token.Type == UNKNOWN {
				return token, position
			}
		} else if c == '\'' {
			token, position = scanLiteral(code, position)
			if token.Type == UNKNOWN {
				return token, position
			}
		} else if c == ':' {
			if token == nil {
				return &Token{Type: UNKNOWN}, position
			}

			if isSet == -1 || isSet == 0 {
				key = token
				token = nil
				isSet = 0
			} else {
				return &Token{Type: UNKNOWN}, position
			}

			position++
		} else if c == ',' || c == '}' {
			if token == nil {
				return &Token{Type: UNKNOWN}, position
			}

			if isSet == -1 || isSet == 1 {
				isSet = 1
				if token.Type != BINARY && token.Type != COLLECTION {
					collection[token.Value] = ""
				} else {
					collection[token.Original] = ""
				}
				token = nil
			} else {
				if key == nil {
					return &Token{Type: UNKNOWN}, position
				} else {
					if (key.Type != BINARY && key.Type != COLLECTION) &&
						(token.Type != BINARY && token.Type != COLLECTION) {
						collection[key.Value] = token.Value
					} else if (key.Type != BINARY && key.Type != COLLECTION) &&
						(token.Type == BINARY || token.Type == COLLECTION) {
						collection[key.Value] = token.Original
					} else if (key.Type == BINARY || key.Type == COLLECTION) &&
						(token.Type != BINARY && token.Type != COLLECTION) {
						collection[key.Original] = token.Value
					} else if (key.Type == BINARY || key.Type == COLLECTION) &&
						(token.Type == BINARY || token.Type == COLLECTION) {
						collection[key.Original] = token.Original
					}

					key, token = nil, nil
				}
			}
			position++

			if c == '}' {
				if isSet == 0 {
					return &Token{Type: COLLECTION, Original: collection}, position
				} else {
					set := []interface{}{}
					for elem, _ := range collection {
						set = append(set, elem)
					}
					return &Token{Type: COLLECTION, Original: set}, position
				}
			}
		} else {
			return &Token{Type: UNKNOWN}, position
		}
	}
	return &Token{Type: UNKNOWN}, position
}

func scanLiteral(code string, position int) (*Token, int) {
	delimiter := rune(code[position])
	position++
	s := position
	for position < len(code) {
		c := rune(code[position])
		if c == delimiter && rune(code[position-1]) != '\\' {
			break
		}
		position++
	}

	if position == len(code) {
		return &Token{Type: UNKNOWN}, position
	}

	literal := code[s:position]
	position++

	if delimiter == '`' {
		if !isTextual(literal) {
			return &Token{Type: UNKNOWN}, position
		}
		return &Token{Type: IDENTIFIER, Value: literal}, position
	}

	if isUUID(literal) {
		return &Token{Type: UUID, Value: literal}, position
	}

	if isTimestamp(literal) {
		return &Token{Type: TIMESTAMP, Value: literal}, position
	}

	if isDate(literal) {
		return &Token{Type: DATE, Value: literal}, position
	}

	if isTime(literal) {
		return &Token{Type: TIME, Value: literal}, position
	}

	if isInterval(strings.ToUpper(literal)) {
		return &Token{Type: TIME_INTERVAL, Value: literal}, position
	}

	if isBinary(literal) {
		if decoded, err := hex.DecodeString(literal[2:]); err != nil {
			return &Token{Type: UNKNOWN}, position
		} else {
			return &Token{Type: BINARY, Original: decoded}, position
		}
	}

	return &Token{Type: STRING, Value: literal}, position
}

func scanSymbolic(code string, position int) (*Token, int) {
	var symbol int
	switch code[position] {
	case '<':
		symbol = LESS_THAN
	case '>':
		symbol = GREATER_THAN
	case '!':
		symbol = EXCLAMATION_MARK
	case '=':
		symbol = EQUAL_SIGN
	case '+':
		symbol = PLUS_SIGN
	case '-':
		symbol = MINUS_SIGN
	case '*':
		symbol = ASTERISK
	case '/':
		symbol = SLASH
	case '%':
		symbol = PERCENT_SIGN
	case '&':
		symbol = AMPERSAND
	case '|':
		symbol = VERTICAL_BAR
	case ',':
		symbol = COMMA
	case '.':
		symbol = PERIOD
	case '(':
		symbol = LEFT_PARANTHESIS
	case ')':
		symbol = RIGHT_PARANTHESIS
	case ';':
		symbol = SEMICOLON
	default:
		symbol = UNKNOWN
	}
	position++
	return &Token{Type: symbol}, position
}

func isInteger(number string) bool {
	r := regexp.MustCompile(`^-?[0-9]+$`)
	return r.MatchString(number)
}

func isFloat(number string) bool {
	r := regexp.MustCompile(`^-?[0-9]+[.][0-9]+$`)
	return r.MatchString(number)
}

func isUUID(id string) bool {
	_, err := uuid.FromString(id)
	if err != nil {
		return false
	}
	return true
}

func isTimestamp(timestamp string) bool {
	r := regexp.MustCompile(`^[0-9]{4}-(0[0-9]|1[0-2])-([0-2][0-9]|3[0-1])\s+([0-1][0-9]|2[0-3])(:[0-5][0-9]){2}$`)
	return r.MatchString(timestamp)
}

func isDate(date string) bool {
	r := regexp.MustCompile(`^[0-9]{4}-(0[0-9]|1[0-2])-([0-2][0-9]|3[0-1])$`)
	return r.MatchString(date)
}

func isBinary(binary string) bool {
	r := regexp.MustCompile(`^0(x|X)(([0-9]|(a|A)|(b|B)|(c|C)|(d|D)|(e|E)|(f|F)))+$`)
	return r.MatchString(binary)
}

func isTime(time string) bool {
	r := regexp.MustCompile(`^([0-1][0-9]|2[0-3])(:[0-5][0-9]){2}$`)
	return r.MatchString(time)
}

func isInterval(interval string) bool {
	r := regexp.MustCompile(`^INTERVAL(\s[0-9]+\s(YEAR|MONTH|WEEK|DAY|HOUR|MINUTE|SECOND)S*)+$`)
	return r.MatchString(interval)
}

func isTextual(textual string) bool {
	for i, _ := range textual {
		c := rune(textual[i])
		if c != '_' && !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
