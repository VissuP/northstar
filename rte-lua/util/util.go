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

package util

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/yuin/gopher-lua"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	TIME      = "time"
	BLOB      = "blob"
	MAP       = "map"
	ARRAY     = "array"
	BOOLEAN   = "bool"
	INTEGER   = "int"
	FLOAT     = "float"
	DOUBLE    = "double"
	STRING    = "string"
	UUID      = "gocql.UUID"
	TIMESTAMP = "time.Time"
	BYTEARRAY = "[]uint8"
	BIGINT    = "*big.Int"
	BIGFLOAT  = "*big.Float"
	IP        = "*net.IP"
	UNSIGNED  = "uint"
)

func FromLua(value interface{}) (interface{}, error) {
	switch converted := value.(type) {
	case lua.LBool:
		return bool(converted), nil
	case lua.LNumber:
		switch reflect.TypeOf(converted).Kind() {
		case reflect.Uint64:
			return uint64(converted), nil
		case reflect.Int64:
			return int64(converted), nil
		case reflect.Float64:
			return float64(converted), nil
		}
		return nil, errors.New("nsQL error: unknown Lua number type")
	case lua.LString:
		return string(converted), nil
	case *lua.LNilType:
		return nil, nil
	case *lua.LTable:
		maxn := converted.MaxN()
		if maxn == 0 {
			ret := make(map[interface{}]interface{})
			var convErr error
			converted.ForEach(func(k, v lua.LValue) {
				cKey, err := FromLua(k)
				if err != nil {
					convErr = err
					return
				}

				cValue, err := FromLua(v)
				if err != nil {
					convErr = err
					return
				}

				ret[cKey] = cValue
			})
			if convErr != nil {
				return nil, convErr
			}
			return ret, nil
		} else {
			ret := make([]interface{}, 0, maxn)
			for i := 1; i <= maxn; i++ {
				if sub, err := FromLua(converted.RawGetInt(i)); err != nil {
					return nil, err
				} else {
					ret = append(ret, sub)
				}
			}
			return ret, nil
		}
	}
	return nil, errors.New("nsQL error: unknown Lua type")
}

func ToLua(L *lua.LState, value interface{}) (lua.LValue, error) {
	if value == nil {
		return lua.LNil, nil
	}

	rType := reflect.TypeOf(value)
	rValue := reflect.ValueOf(value)
	switch rType.Kind() {
	case reflect.Bool:
		return lua.LBool(rValue.Bool()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return lua.LNumber(rValue.Uint()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return lua.LNumber(rValue.Int()), nil
	case reflect.Float32, reflect.Float64:
		return lua.LNumber(rValue.Float()), nil
	case reflect.String:
		return lua.LString(rValue.String()), nil
	case reflect.Slice:
		arr := L.CreateTable(rValue.Len(), 0)
		for i := 0; i < rValue.Len(); i++ {
			elem, err := ToLua(L, rValue.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			arr.Append(elem)
		}
		return arr, nil
	case reflect.Map:
		keys := rValue.MapKeys()
		tbl := L.CreateTable(0, len(keys))
		for _, key := range keys {
			k, err := ToLua(L, key.Interface())
			if err != nil {
				return nil, err
			}

			elem, err := ToLua(L, rValue.MapIndex(key).Interface())
			if err != nil {
				return nil, err
			}

			tbl.RawSetH(k, elem)
		}
		return tbl, nil
	case reflect.Struct:
		switch converted := value.(type) {
		case time.Time:
			return lua.LNumber(converted.Unix()), nil
		default:
			tbl := L.CreateTable(0, rValue.NumField())
			typ := rValue.Type()
			for i := 0; i < rValue.NumField(); i++ {
				k, err := ToLua(L, typ.Field(i).Name)
				if err != nil {
					return nil, err
				}

				elem, err := ToLua(L, rValue.Field(i).Interface())
				if err != nil {
					return nil, err
				}

				tbl.RawSetH(k, elem)
			}
			return tbl, nil
		}
	default:
		switch converted := value.(type) {
		case *big.Int:
			if converted == nil {
				return lua.LNumber(0), nil
			}
			return lua.LNumber(converted.Int64()), nil
		case *big.Float:
			if converted == nil {
				return lua.LNumber(0.0), nil
			}
			f, _ := converted.Float64()
			return lua.LNumber(f), nil
		case gocql.UUID:
			return lua.LString(converted.String()), nil
		}
	}

	return nil, errors.New("nsQL error: unknown Go type")
}

func ToInternalType(external reflect.Type) (string, error) {
	str := external.String()
	switch str {
	case BOOLEAN:
		return BOOLEAN, nil
	case STRING, UUID, IP:
		return STRING, nil
	case TIMESTAMP:
		return TIME, nil
	case BYTEARRAY:
		return BLOB, nil
	case BIGINT:
		return INTEGER, nil
	case BIGFLOAT:
		return DOUBLE, nil
	default:
		switch external.Kind() {
		case reflect.Slice:
			elemType, err := ToInternalType(external.Elem())
			if err != nil {
				return "", err
			}

			return "array[" + elemType + "]", nil
		case reflect.Map:
			keyType, err := ToInternalType(external.Key())
			if err != nil {
				return "", err
			}

			elemType, err := ToInternalType(external.Elem())
			if err != nil {
				return "", err
			}

			return "map[" + keyType + "]" + elemType, nil
		default:
			if strings.HasPrefix(str, UNSIGNED) || strings.HasPrefix(str, INTEGER) {
				return INTEGER, nil
			}

			if strings.HasPrefix(str, FLOAT) {
				return DOUBLE, nil
			}

		}
	}

	return "", errors.New("nsQL error: unknown Go type")
}

func DataToString(data interface{}, internalType string) string {
	if str, ok := data.(string); ok {
		return str
	}

	switch internalType {
	case TIME:
		var t time.Time
		var ok bool
		var unix float64

		if t, ok = data.(time.Time); !ok {
			if unix, ok = data.(float64); ok {
				t = time.Unix(int64(unix), 0).UTC()
			}
		}

		if t.Unix() < 0 {
			return ""
		}

		return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(),
			t.Second())
	case BLOB:
		var blob []byte
		var ok bool

		blob, ok = data.([]byte)
		if !ok {
			blobAsSlice, _ := data.([]interface{})
			for _, elem := range blobAsSlice {
				blobElem, _ := elem.(float64)
				blob = append(blob, uint8(blobElem))
			}
		}

		if len(blob) == 0 {
			return ""
		}
		return "0x" + hex.EncodeToString(blob)
	default:
		if strings.HasPrefix(internalType, ARRAY) {
			elemType := internalType[6 : len(internalType)-1]
			v := reflect.ValueOf(data)
			var array []string
			for i := 0; i < v.Len(); i++ {
				array = append(array, DataToString(v.Index(i).Interface(), elemType))
			}
			return fmt.Sprintf("%v", array)
		}

		if strings.HasPrefix(internalType, MAP) {
			end := 0
			nesting := 1
			rest := internalType[4:]
			for i, elem := range rest {
				if elem == ']' {
					nesting--
					if nesting == 0 {
						end = i
						break
					}
					continue
				}

				if elem == '[' {
					nesting++
					continue
				}
			}
			keyType := rest[:end]
			elemType := rest[end+1:]

			v := reflect.ValueOf(data)
			m := make(map[string]string)
			keys := v.MapKeys()
			for _, key := range keys {
				m[DataToString(key.Interface(), keyType)] =
					DataToString(v.MapIndex(key).Interface(), elemType)
			}
			return fmt.Sprintf("%v", m)
		}

		if internalType == "int" {
			switch converted := data.(type) {
			case float64:
				data = int64(converted)
			case *big.Int:
				if converted == nil {
					data = "0"
				}
			}
		}

		if internalType == "double" {
			if big, ok := data.(*big.Float); ok && big == nil {
				return "0"
			}
		}

		return fmt.Sprintf("%v", data)
	}
}

func StringToData(str, strType string) (interface{}, error) {
	switch strType {
	case BOOLEAN:
		if str == "true" {
			return true, nil
		}
		return false, nil
	case INTEGER:
		if str == "" {
			return 0, nil
		}

		if i, err := strconv.ParseInt(str, 10, 64); err != nil {
			return nil, err
		} else {
			return i, nil
		}
	case DOUBLE:
		if str == "" {
			return 0, nil
		}

		if d, err := strconv.ParseFloat(str, 64); err != nil {
			return nil, err
		} else {
			return d, nil
		}

	case TIME:
		if str == "" {
			return time.Time{}, nil
		}

		str = str[1 : len(str)-1]
		var t time.Time
		var err error
		if t, err = time.Parse("2006-01-02", str); err != nil {
			if t, err = time.Parse("2006-01-02 15:04:05.0", str); err != nil {
				return nil, err
			}
		}
		return t, nil
	case BLOB:
		if str == "" {
			return []byte{}, nil
		}

		if e, err := hex.DecodeString(str[2:]); err != nil {
			return nil, err
		} else {
			return e, nil
		}
	case STRING:
		if str == "" {
			return str, nil
		}
		return str[1 : len(str)-1], nil
	default:
		if strings.HasPrefix(strType, ARRAY) {
			return makeArrayFromString(str, strType)
		} else if strings.HasPrefix(strType, MAP) {
			return makeMapFromString(str, strType)
		} else {
			return nil, errors.New("nsQL error: " + strType + " is unknown")
		}
	}
}

func makeArrayFromString(array, arrayType string) ([]interface{}, error) {
	var result []interface{}
	var element string

	elements := getArrayElements(array)
	elementType := getArrayElementType(arrayType)

	for len(elements) != 0 {
		element, elements = getNextElement(elements, elementType)

		if value, err := StringToData(element, elementType); err != nil {
			return nil, err
		} else {
			result = append(result, value)
		}
	}

	return result, nil
}

func makeMapFromString(mmap, mmapType string) (map[interface{}]interface{}, error) {
	result := make(map[interface{}]interface{})
	var element string

	elements := getMapElements(mmap)
	keyType, valueType := getMapElementTypes(mmapType)

	for len(elements) != 0 {
		element, elements = getNextElement(elements, keyType)
		key, err := StringToData(element, keyType)
		if err != nil {
			return nil, err
		}
		element, elements = getNextElement(elements, valueType)
		value, err := StringToData(element, valueType)
		if err != nil {
			return nil, err
		}

		result[key] = value
	}

	return result, nil
}

func getArrayElements(array string) string {
	return strings.TrimFunc(array[1:len(array)-1], func(r rune) bool { return unicode.IsSpace(r) })
}

func getArrayElementType(arrayType string) string {
	return arrayType[6 : len(arrayType)-1]
}

func getMapElements(mmap string) string {
	return strings.TrimFunc(mmap[4:len(mmap)-1], func(r rune) bool { return unicode.IsSpace(r) })
}

func getMapElementTypes(mmapType string) (string, string) {
	return getNextList(mmapType[3:])
}

func getNextElement(str, strType string) (string, string) {
	str = strings.TrimLeftFunc(str, func(r rune) bool { return unicode.IsSpace(r) || r == ':' })

	if strings.HasPrefix(strType, ARRAY) {
		list, rest := getNextList(str)
		return "[" + list + "]", rest
	}

	if strings.HasPrefix(strType, MAP) {
		list, rest := getNextList(str)
		return "map[" + list + "]", rest
	}

	return getNextValue(str)
}

func getNextValue(str string) (string, string) {
	isString := false
	if str[0] == '"' {
		str = str[1:]
		isString = true
	}

	for i, v := range str {
		if isString && v == '"' {
			return "\"" + str[:i] + "\"", str[i+1:]
		}

		if !isString && (v == ':' || v == ' ') {
			return str[:i], str[i:]
		}
	}

	return str, ""
}

func getNextList(str string) (string, string) {
	var start, end, nesting int
	started := false
	for i, v := range str {
		if v == '[' {
			if !started {
				start = i + 1
				started = true
			}
			nesting++
		}

		if v == ']' {
			nesting--
			if nesting == 0 {
				end = i
				break
			}
		}
	}
	return str[start:end], str[end+1:]
}
