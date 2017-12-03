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
	"github.com/stretchr/testify/require"
	"testing"
)

func TestColumnExpression01(t *testing.T) {
	parsed, err := Parse("select tbl.col1 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all tbl.col1 from ks.tbl;", parsed.ToString())
}

func TestColumnExpression02(t *testing.T) {
	parsed, err := Parse("select tbl.col1 + col2 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (tbl.col1 + col2) from ks.tbl;", parsed.ToString())
}

func TestColumnExpression03(t *testing.T) {
	parsed, err := Parse("select tbl.col1 - col2 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (tbl.col1 - col2) from ks.tbl;", parsed.ToString())
}

func TestColumnExpression04(t *testing.T) {
	parsed, err := Parse("select tbl.col1 * col2 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (tbl.col1 * col2) from ks.tbl;", parsed.ToString())
}

func TestColumnExpression05(t *testing.T) {
	parsed, err := Parse("select tbl.col1 / col2 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (tbl.col1 / col2) from ks.tbl;", parsed.ToString())
}

func TestColumnExpression06(t *testing.T) {
	parsed, err := Parse("select tbl.col1 % col2 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (tbl.col1 % col2) from ks.tbl;", parsed.ToString())
}

func TestColumnExpression07(t *testing.T) {
	parsed, err := Parse("select tbl.col1 & col2 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (tbl.col1 & col2) from ks.tbl;", parsed.ToString())
}

func TestColumnExpression08(t *testing.T) {
	parsed, err := Parse("select tbl.col1 | col2 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (tbl.col1 | col2) from ks.tbl;", parsed.ToString())
}

func TestNumericExpression01(t *testing.T) {
	parsed, err := Parse("select 17 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all 17 from ks.tbl;", parsed.ToString())
}

func TestNumericExpression02(t *testing.T) {
	parsed, err := Parse("select 17 + 34.39 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (17 + 34.39) from ks.tbl;", parsed.ToString())
}

func TestNumericExpression03(t *testing.T) {
	parsed, err := Parse("select 17 - 34.39 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (17 - 34.39) from ks.tbl;", parsed.ToString())
}

func TestNumericExpression04(t *testing.T) {
	parsed, err := Parse("select 17 * 34.39 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (17 * 34.39) from ks.tbl;", parsed.ToString())
}

func TestNumericExpression05(t *testing.T) {
	parsed, err := Parse("select 17 / 34.39 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (17 / 34.39) from ks.tbl;", parsed.ToString())
}

func TestNumericExpression06(t *testing.T) {
	parsed, err := Parse("select 17 % 34.39 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (17 % 34.39) from ks.tbl;", parsed.ToString())
}

func TestNumericExpression07(t *testing.T) {
	parsed, err := Parse("select 17 & 34.39 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (17 & 34.39) from ks.tbl;", parsed.ToString())
}

func TestNumericExpression08(t *testing.T) {
	parsed, err := Parse("select 17 | 34.39 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (17 | 34.39) from ks.tbl;", parsed.ToString())
}

func TestNumericExpression09(t *testing.T) {
	parsed, err := Parse("select tbl.col1 + 34.39 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (tbl.col1 + 34.39) from ks.tbl;", parsed.ToString())
}

func TestNumericExpression10(t *testing.T) {
	parsed, err := Parse("select tbl.col1 - 34.39 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (tbl.col1 - 34.39) from ks.tbl;", parsed.ToString())
}

func TestNumericExpression11(t *testing.T) {
	parsed, err := Parse("select tbl.col1 * 34.39 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (tbl.col1 * 34.39) from ks.tbl;", parsed.ToString())
}

func TestNumericExpression12(t *testing.T) {
	parsed, err := Parse("select tbl.col1 / 34.39 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (tbl.col1 / 34.39) from ks.tbl;", parsed.ToString())

}

func TestNumericExpression13(t *testing.T) {
	parsed, err := Parse("select tbl.col1 % 34.39 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (tbl.col1 % 34.39) from ks.tbl;", parsed.ToString())
}

func TestNumericExpression14(t *testing.T) {
	parsed, err := Parse("select tbl.col1 & 34.39 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (tbl.col1 & 34.39) from ks.tbl;", parsed.ToString())
}

func TestNumericExpression15(t *testing.T) {
	parsed, err := Parse("select tbl.col1 | 34.39 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (tbl.col1 | 34.39) from ks.tbl;", parsed.ToString())
}

func TestNumericExpression16(t *testing.T) {
	parsed, err := Parse("select 17 + col2 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (17 + col2) from ks.tbl;", parsed.ToString())
}

func TestNumericExpression17(t *testing.T) {
	parsed, err := Parse("select 17 - col2 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (17 - col2) from ks.tbl;", parsed.ToString())
}

func TestNumericExpression18(t *testing.T) {
	parsed, err := Parse("select 17 * col2 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (17 * col2) from ks.tbl;", parsed.ToString())
}

func TestNumericExpression19(t *testing.T) {
	parsed, err := Parse("select 17 / col2 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (17 / col2) from ks.tbl;", parsed.ToString())
}

func TestNumericExpression20(t *testing.T) {
	parsed, err := Parse("select 17 % col2 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (17 % col2) from ks.tbl;", parsed.ToString())
}

func TestNumericExpression21(t *testing.T) {
	parsed, err := Parse("select 17 & col2 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (17 & col2) from ks.tbl;", parsed.ToString())
}

func TestNumericExpression22(t *testing.T) {
	parsed, err := Parse("select 17 | col2 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (17 | col2) from ks.tbl;", parsed.ToString())
}

func TestTemporalExpression01(t *testing.T) {
	parsed, err := Parse("select '2016-06-15 00:00:00' from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all 2016-06-15 00:00:00 from ks.tbl;", parsed.ToString())
}

func TestTemporalExpression02(t *testing.T) {
	parsed, err := Parse("select tbl.col1 + 'interval 1 year 2 months' from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (tbl.col1 + interval 1 year 2 months) from ks.tbl;", parsed.ToString())
}

func TestTemporalExpression03(t *testing.T) {
	parsed, err := Parse("select tbl.col1 - 'interval 1 year 2 months' from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (tbl.col1 - interval 1 year 2 months) from ks.tbl;", parsed.ToString())
}

func TestTemporalExpression04(t *testing.T) {
	parsed, err := Parse("select 'interval 1 year 2 months' + tbl.col1 from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (interval 1 year 2 months + tbl.col1) from ks.tbl;", parsed.ToString())
}

func TestTemporalExpression05(t *testing.T) {
	parsed, err := Parse("select '2016-06-15 00:00:00' + 'interval 1 year 2 months' from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (2016-06-15 00:00:00 + interval 1 year 2 months) from ks.tbl;", parsed.ToString())
}

func TestTemporalExpression06(t *testing.T) {
	parsed, err := Parse("select '2016-06-15 00:00:00' - 'interval 1 year 2 months' from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select all (2016-06-15 00:00:00 - interval 1 year 2 months) from ks.tbl;", parsed.ToString())
}

func TestTemporalExpression07(t *testing.T) {
	_, err := Parse("select tbl.col1 + col2 + 'interval 1 year 2 months' from ks.tbl;")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: unexpected (tbl.col1 + col2)", err.Error())
}

func TestTemporalExpression08(t *testing.T) {
	_, err := Parse("select 'interval 1 year 2 months' + tbl.col1 + col2 from ks.tbl;")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: unexpected IDENTIFIER, expecting TIME_INTERVAL or PLUS_SIGN or MINUS_SIGN or LEFT_PARANTHESIS", err.Error())
}

func TestConditionalExpression01(t *testing.T) {
	parsed, err := Parse("select all * from ks.tbl where '2016-06-15 00:00:00' = col2 + 'interval 1 year 2 months';")
	require.Nil(t, err)
	require.Equal(t, "select all * from ks.tbl where (2016-06-15 00:00:00 == (col2 + interval 1 year 2 months));", parsed.ToString())
}

func TestConditionalExpression02(t *testing.T) {
	parsed, err := Parse("select all * from ks.tbl where 'interval 1 year 2 months' + tbl.col1 != col2;")
	require.Nil(t, err)
	require.Equal(t, "select all * from ks.tbl where ((interval 1 year 2 months + tbl.col1) != col2);", parsed.ToString())
}

func TestConditionalExpression03(t *testing.T) {
	parsed, err := Parse("select all * from ks.tbl where 17 + tbl.col1 <> col2 + 34.39;")
	require.Nil(t, err)
	require.Equal(t, "select all * from ks.tbl where ((17 + tbl.col1) != (col2 + 34.39));", parsed.ToString())
}

func TestConditionalExpression04(t *testing.T) {
	parsed, err := Parse("select all * from ks.tbl where 17 + tbl.col1 < col2 + tbl.col3;")
	require.Nil(t, err)
	require.Equal(t, "select all * from ks.tbl where ((17 + tbl.col1) < (col2 + tbl.col3));", parsed.ToString())
}

func TestConditionalExpression05(t *testing.T) {
	parsed, err := Parse("select all * from ks.tbl where tbl.col1 in (select all col1 from ks.tbl);")
	require.Nil(t, err)
	require.Equal(t, "select all * from ks.tbl where (tbl.col1 in select all col1 from ks.tbl;);", parsed.ToString())
}

func TestConditionalExpression06(t *testing.T) {
	parsed, err := Parse("select all * from ks.tbl where tbl.col1 not in (select all col2 as col1 from ks.tbl);")
	require.Nil(t, err)
	require.Equal(t, "select all * from ks.tbl where (tbl.col1 not in select all col2 from ks.tbl;);", parsed.ToString())
}

func TestConditionalExpression07(t *testing.T) {
	parsed, err := Parse("select all * from ks.tbl where col2 <= 'interval 1 year 2 months' + tbl.col1;")
	require.Nil(t, err)
	require.Equal(t, "select all * from ks.tbl where (col2 <= (interval 1 year 2 months + tbl.col1));", parsed.ToString())
}

func TestConditionalExpression08(t *testing.T) {
	parsed, err := Parse("select all * from ks.tbl where col2 + tbl.col3 > 17 + tbl.col1;")
	require.Nil(t, err)
	require.Equal(t, "select all * from ks.tbl where ((col2 + tbl.col3) > (17 + tbl.col1));", parsed.ToString())
}

func TestConditionalExpression09(t *testing.T) {
	parsed, err := Parse("select all * from ks.tbl where col2 + tbl.col3 >= tbl.col1;")
	require.Nil(t, err)
	require.Equal(t, "select all * from ks.tbl where ((col2 + tbl.col3) >= tbl.col1);", parsed.ToString())
}

func TestConditionalExpression10(t *testing.T) {
	parsed, err := Parse("select all * from ks.tbl where tbl.col1 is null;")
	require.Nil(t, err)
	require.Equal(t, "select all * from ks.tbl where (tbl.col1 is null);", parsed.ToString())
}

func TestConditionalExpression11(t *testing.T) {
	parsed, err := Parse("select all * from ks.tbl where tbl.col1 is not null;")
	require.Nil(t, err)
	require.Equal(t, "select all * from ks.tbl where (tbl.col1 is not null);", parsed.ToString())
}

func TestConditionalExpression12(t *testing.T) {
	parsed, err := Parse("select all * from ks.tbl where tbl.col1 = 'some text';")
	require.Nil(t, err)
	require.Equal(t, "select all * from ks.tbl where (tbl.col1 == some text);", parsed.ToString())
}

func TestConditionalExpression13(t *testing.T) {
	parsed, err := Parse("select all * from ks.tbl where tbl.col1 != true;")
	require.Nil(t, err)
	require.Equal(t, "select all * from ks.tbl where (tbl.col1 != true);", parsed.ToString())
}

func TestConditionalExpression14(t *testing.T) {
	parsed, err := Parse("select all * from ks.tbl where tbl.col1 <> '123e4567-e89b-12d3-a456-426655440000';")
	require.Nil(t, err)
	require.Equal(t, "select all * from ks.tbl where (tbl.col1 != 123e4567-e89b-12d3-a456-426655440000);", parsed.ToString())
}

func TestConditionalExpression15(t *testing.T) {
	parsed, err := Parse("select all * from ks.tbl where 'some text' < tbl.col1;")
	require.Nil(t, err)
	require.Equal(t, "select all * from ks.tbl where (some text < tbl.col1);", parsed.ToString())
}

func TestConditionalExpression17(t *testing.T) {
	parsed, err := Parse("select all * from ks.tbl where false = tbl.col1;")
	require.Nil(t, err)
	require.Equal(t, "select all * from ks.tbl where (false == tbl.col1);", parsed.ToString())
}

func TestConditionalExpression18(t *testing.T) {
	parsed, err := Parse("select all * from ks.tbl where '123e4567-e89b-12d3-a456-426655440000' != tbl.col1;")
	require.Nil(t, err)
	require.Equal(t, "select all * from ks.tbl where (123e4567-e89b-12d3-a456-426655440000 != tbl.col1);", parsed.ToString())
}

func TestConditionalExpression19(t *testing.T) {
	_, err := Parse("select all * from ks.tbl where 'interval 1 year 2 months' + tbl.col1 != col2 + tbl.col3;")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: unexpected (col2 + tbl.col3)", err.Error())
}

func TestConditionalExpression20(t *testing.T) {
	_, err := Parse("select all * from ks.tbl where tbl.col1 + 1 in (select all col1 from ks.tbl);")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: unexpected IN", err.Error())
}

func TestConditionalExpression21(t *testing.T) {
	_, err := Parse("select all * from ks.tbl where tbl.col1 not in ((select all col1 from ks.tbl) union (select all col1 from ks.tbl));")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: IDENTIFIER expected", err.Error())
}

func TestConditionalExpression22(t *testing.T) {
	_, err := Parse("select all * from ks.tbl where tbl.col1 in (select all col2 from ks.tbl);")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: syntax error: unexpected column mismatch after IN", err.Error())
}

func TestConditionalExpression23(t *testing.T) {
	_, err := Parse("select all * from ks.tbl where tbl.col1 + col2 <> '123e4567-e89b-12d3-a456-426655440000';")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: IDENTIFIER expected", err.Error())
}

func TestConditionalExpression24(t *testing.T) {
	_, err := Parse("select all * from ks.tbl where tbl.col1 + col2 is null;")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: IDENTIFIER expected", err.Error())
}

func TestConditionalExpression25(t *testing.T) {
	_, err := Parse("select all * from ks.tbl where tbl.col1 + col2 is not null;")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: IDENTIFIER expected", err.Error())
}

func TestConditionalExpression26(t *testing.T) {
	_, err := Parse("select all * from ks.tbl where tbl.col1 + col2 > 'some text';")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: IDENTIFIER expected", err.Error())
}

func TestConditionalExpression27(t *testing.T) {
	_, err := Parse("select all * from ks.tbl where tbl.col1 + col2 != true;")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: IDENTIFIER expected", err.Error())
}

func TestConditionalExpression28(t *testing.T) {
	_, err := Parse("select all * from ks.tbl where tbl.col1 + col2 <> '123e4567-e89b-12d3-a456-426655440000';")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: IDENTIFIER expected", err.Error())
}

func TestConditionalExpression29(t *testing.T) {
	_, err := Parse("select all * from ks.tbl where 'some text' < tbl.col1 + col2;")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: IDENTIFIER expected", err.Error())
}

func TestConditionalExpression30(t *testing.T) {
	_, err := Parse("select all * from ks.tbl where true = tbl.col1 + col2;")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: IDENTIFIER expected", err.Error())
}

func TestConditionalExpression31(t *testing.T) {
	_, err := Parse("select all * from ks.tbl where '123e4567-e89b-12d3-a456-426655440000' != tbl.col1 + col2;")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: IDENTIFIER expected", err.Error())
}

func TestLogicalExpression01(t *testing.T) {
	parsed, err := Parse("select all * from ks.tbl where '2016-06-15 00:00:00' = col2 + 'interval 1 year 2 months' or 17 + tbl.col1 < col2 + tbl.col3;")
	require.Nil(t, err)
	require.Equal(t, "select all * from ks.tbl where ((2016-06-15 00:00:00 == (col2 + interval 1 year 2 months)) or ((17 + tbl.col1) < (col2 + tbl.col3)));", parsed.ToString())
}

func TestLogicalExpression02(t *testing.T) {
	parsed, err := Parse("select all * from ks.tbl where 17 + tbl.col1 < col2 + tbl.col3 and '2016-06-15 00:00:00' = col2 + 'interval 1 year 2 months';")
	require.Nil(t, err)
	require.Equal(t, "select all * from ks.tbl where (((17 + tbl.col1) < (col2 + tbl.col3)) and (2016-06-15 00:00:00 == (col2 + interval 1 year 2 months)));", parsed.ToString())
}

func TestLogicalExpression03(t *testing.T) {
	parsed, err := Parse("select all * from ks.tbl where not (17 + tbl.col1 < col2 + tbl.col3 and '2016-06-15 00:00:00' = col2 + 'interval 1 year 2 months');")
	require.Nil(t, err)
	require.Equal(t, "select all * from ks.tbl where (((17 + tbl.col1) >= (col2 + tbl.col3)) or (2016-06-15 00:00:00 != (col2 + interval 1 year 2 months)));", parsed.ToString())
}

func TestToTemporalTransformer01(t *testing.T) {
	parsed, err := Parse("select distinct now() from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select distinct now() from ks.tbl;", parsed.ToString())
}

func TestToTemporalTransformer02(t *testing.T) {
	_, err := Parse("select distinct now(1) from ks.tbl;")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: unexpected INTEGER, expecting RIGHT_PARANTHESIS", err.Error())
}

func TestToNumericTransformer01(t *testing.T) {
	parsed, err := Parse("select distinct year(tbl.col1) from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select distinct year(tbl.col1) from ks.tbl;", parsed.ToString())
}

func TestToNumericTransformer02(t *testing.T) {
	parsed, err := Parse("select distinct month(min(tbl.col1 + 'interval 1 year 2 months')) from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select distinct month(min((tbl.col1 + interval 1 year 2 months))) from ks.tbl;", parsed.ToString())
}

func TestToNumericTransformer03(t *testing.T) {
	parsed, err := Parse("select distinct day('2016-06-15 00:00:00') from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select distinct day(2016-06-15 00:00:00) from ks.tbl;", parsed.ToString())
}

func TestToNumericTransformer04(t *testing.T) {
	parsed, err := Parse("select distinct hour(now()) from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select distinct hour(now()) from ks.tbl;", parsed.ToString())
}

func TestToNumericTransformer05(t *testing.T) {
	parsed, err := Parse("select distinct minute('2016-06-15 00:00:00' - 'interval 1 year 2 months') from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select distinct minute((2016-06-15 00:00:00 - interval 1 year 2 months)) from ks.tbl;", parsed.ToString())
}

func TestToNumericTransformer06(t *testing.T) {
	_, err := Parse("select distinct second(tbl.col1 + col2) from ks.tbl;")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: unexpected (tbl.col1 + col2)", err.Error())
}

func TestToNumericAggregator01(t *testing.T) {
	parsed, err := Parse("select distinct count(tbl.col1) from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select distinct count(tbl.col1) from ks.tbl;", parsed.ToString())
}

func TestToNumericAggregator02(t *testing.T) {
	parsed, err := Parse("select distinct count(tbl.col1 + col2) from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select distinct count((tbl.col1 + col2)) from ks.tbl;", parsed.ToString())
}

func TestToNumericAggregator03(t *testing.T) {
	parsed, err := Parse("select distinct count(17) from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select distinct count(17) from ks.tbl;", parsed.ToString())
}

func TestToNumericAggregator04(t *testing.T) {
	parsed, err := Parse("select distinct count(year('2016-06-15 00:00:00')) from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select distinct count(year(2016-06-15 00:00:00)) from ks.tbl;", parsed.ToString())
}

func TestToNumericAggregator05(t *testing.T) {
	parsed, err := Parse("select distinct count(34.39 + year('2016-06-15 00:00:00') + tbl.col1) from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select distinct count(((34.39 + year(2016-06-15 00:00:00)) + tbl.col1)) from ks.tbl;", parsed.ToString())
}

func TestToNumericAggregator06(t *testing.T) {
	parsed, err := Parse("select distinct count('2016-06-15 00:00:00') from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select distinct count(2016-06-15 00:00:00) from ks.tbl;", parsed.ToString())
}

func TestToNumericAggregator07(t *testing.T) {
	parsed, err := Parse("select distinct count(now()) from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select distinct count(now()) from ks.tbl;", parsed.ToString())
}

func TestToNumericAggregator08(t *testing.T) {
	parsed, err := Parse("select distinct count(now() + 'interval 1 year 2 months') from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select distinct count((now() + interval 1 year 2 months)) from ks.tbl;", parsed.ToString())
}

func TestToNumericAggregator09(t *testing.T) {
	parsed, err := Parse("select distinct count('2016-06-15 00:00:00' = col2 + 'interval 1 year 2 months' or 17 + tbl.col1 < col2 + tbl.col3) from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select distinct count(((2016-06-15 00:00:00 == (col2 + interval 1 year 2 months)) or ((17 + tbl.col1) < (col2 + tbl.col3)))) from ks.tbl;", parsed.ToString())
}

func TestToNumericAggregator10(t *testing.T) {
	_, err := Parse("select distinct count(min(now())) from ks.tbl;")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: unexpected aggregation function in count", err.Error())
}

func TestToNumericAggregator11(t *testing.T) {
	_, err := Parse("select distinct count(tbl.col1 + col2 + min(now())) from ks.tbl;")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: unexpected aggregation function in count", err.Error())
}

func TestToNumericAggregator12(t *testing.T) {
	_, err := Parse("select distinct count(count(now())) from ks.tbl;")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: unexpected aggregation function in count", err.Error())
}

func TestToNumericAggregator13(t *testing.T) {
	_, err := Parse("select distinct count(17 + year(now()) + tbl.col1 - count(now())) from ks.tbl;")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: unexpected aggregation function in count", err.Error())
}

func TestToNumericAggregator14(t *testing.T) {
	parsed, err := Parse("select distinct sum(17) from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select distinct sum(17) from ks.tbl;", parsed.ToString())
}

func TestToNumericAggregator15(t *testing.T) {
	parsed, err := Parse("select distinct mean(year(now())) from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select distinct mean(year(now())) from ks.tbl;", parsed.ToString())
}

func TestToNumericAggregator16(t *testing.T) {
	parsed, err := Parse("select distinct variance(17 + year(now()) + tbl.col1) from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select distinct variance(((17 + year(now())) + tbl.col1)) from ks.tbl;", parsed.ToString())
}

func TestToNumericAggregator17(t *testing.T) {
	parsed, err := Parse("select distinct stdev(tbl.col1) from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select distinct stdev(tbl.col1) from ks.tbl;", parsed.ToString())
}

func TestToNumericAggregator18(t *testing.T) {
	parsed, err := Parse("select distinct corr(tbl.col1 + col2, 17 + year(now())) from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select distinct corr((tbl.col1 + col2), (17 + year(now()))) from ks.tbl;", parsed.ToString())
}

func TestToNumericAggregator19(t *testing.T) {
	_, err := Parse("select distinct sum(17 + sum(34.39) + tbl.col1) from ks.tbl;")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: unexpected aggregation function in sum", err.Error())
}

func TestToNumericAggregator20(t *testing.T) {
	_, err := Parse("select distinct mean(17 + sum(34.39) + year(now()) + tbl.col1) from ks.tbl;")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: unexpected aggregation function in mean", err.Error())
}

func TestToNumericAggregator21(t *testing.T) {
	_, err := Parse("select distinct variance(min(tbl.col1)) from ks.tbl;")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: unexpected aggregation function in variance", err.Error())
}

func TestToNumericAggregator22(t *testing.T) {
	_, err := Parse("select distinct stdev(min(tbl.col1) + col2) from ks.tbl;")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: unexpected aggregation function in stdev", err.Error())
}

func TestToColumnAggregator01(t *testing.T) {
	parsed, err := Parse("select distinct min(tbl.col1) from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select distinct min(tbl.col1) from ks.tbl;", parsed.ToString())
}

func TestToColumnAggregator02(t *testing.T) {
	parsed, err := Parse("select distinct max(tbl.col1 + col2) from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select distinct max((tbl.col1 + col2)) from ks.tbl;", parsed.ToString())
}

func TestToColumnAggregator03(t *testing.T) {
	parsed, err := Parse("select distinct first(17) from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select distinct first(17) from ks.tbl;", parsed.ToString())
}

func TestToColumnAggregator04(t *testing.T) {
	parsed, err := Parse("select distinct first(year(now())) from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select distinct first(year(now())) from ks.tbl;", parsed.ToString())
}

func TestToColumnAggregator05(t *testing.T) {
	parsed, err := Parse("select distinct last(17 + year(now()) + col2) from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select distinct last(((17 + year(now())) + col2)) from ks.tbl;", parsed.ToString())
}

func TestToColumnAggregator06(t *testing.T) {
	parsed, err := Parse("select distinct min('2016-06-15 00:00:00') from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select distinct min(2016-06-15 00:00:00) from ks.tbl;", parsed.ToString())
}

func TestToColumnAggregator07(t *testing.T) {
	parsed, err := Parse("select distinct max(now()) from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select distinct max(now()) from ks.tbl;", parsed.ToString())
}

func TestToColumnAggregator08(t *testing.T) {
	parsed, err := Parse("select distinct first(now() + 'interval 1 year 2 months') from ks.tbl;")
	require.Nil(t, err)
	require.Equal(t, "select distinct first((now() + interval 1 year 2 months)) from ks.tbl;", parsed.ToString())
}

func TestToColumnAggregator09(t *testing.T) {
	_, err := Parse("select distinct last(min(now())) from ks.tbl;")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: unexpected aggregation function in last", err.Error())
}

func TestToColumnAggregator10(t *testing.T) {
	_, err := Parse("select distinct min(tbl.col1 + col2 + min(tbl.col3)) from ks.tbl;")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: unexpected aggregation function in min", err.Error())
}

func TestToColumnAggregator11(t *testing.T) {
	_, err := Parse("select distinct max(count(now())) from ks.tbl;")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: unexpected aggregation function in max", err.Error())
}

func TestToColumnAggregator12(t *testing.T) {
	_, err := Parse("select distinct first(34.39 + year(now()) + tbl.col1 + sum(col2)) from ks.tbl;")
	require.NotNil(t, err)
	require.Equal(t, "nsQL syntax error: unexpected aggregation function in first", err.Error())
}

func TestMapBlobFetch01(t *testing.T) {
	_, err := Parse("select map_blob_json_fetch(tbl.col1, 'capability', 'field') from ks.tbl;")
	require.Nil(t, err)
}

func TestJsonFetch01(t *testing.T) {
	_, err := Parse("select json_fetch(tbl.col1, 'field') from ks.tbl;")
	require.Nil(t, err)
}

func TestEscapeSequence01(t *testing.T) {
	_, err := Parse("select `group` from ks.tbl;")
	require.Nil(t, err)
}

func TestDelete01(t *testing.T) {
	_, err := Parse("delete from ks.tbl where col1 = '\\'\\'';")
	require.Nil(t, err)
}

func TestInsert01(t *testing.T) {
	_, err := Parse("insert into ks.tbl (col1, col2) values (1, 'sometest\\'s');")
	require.Nil(t, err)
}

func TestUpdate01(t *testing.T) {
	_, err := Parse("update ks.tbl set col1 = 1 where col1 = 2;")
	require.Nil(t, err)
}

func TestDateAndTime01(t *testing.T) {
	_, err := Parse("select * from ks.tbl where col1 = '2017-05-05' and col2 = '01:02:03';")
	require.Nil(t, err)
}

func TestNumericParseError01(t *testing.T) {
	_, err := Parse("select * from ks.tbl where col1 = 2..1;")
	require.NotNil(t, err)
}

func TestBinaryField01(t *testing.T) {
	_, err := Parse("insert into ks.tbl (col1) values ('0x00000A');")
	require.Nil(t, err)
}

func TestSymbolError01(t *testing.T) {
	_, err := Parse("select * from ks.tbl where col1 = ^;")
	require.NotNil(t, err)
}

func TestUnterminatedLiteralError01(t *testing.T) {
	_, err := Parse("select * from ks.tbl where col1 = '")
	require.NotNil(t, err)
}

func TestJoin01(t *testing.T) {
	_, err := Parse("select * from ks.tbl1 inner join ks.tbl2 on tbl1.col1 = tbl2.col1;")
	require.Nil(t, err)
}

func TestJoin02(t *testing.T) {
	_, err := Parse("select * from ks.tbl1 full outer join ks.tbl2 on tbl1.col1 = tbl2.col1;")
	require.Nil(t, err)
}

func TestJoin03(t *testing.T) {
	_, err := Parse("select * from ks.tbl1 left semi join ks.tbl2 on tbl1.col1 = tbl2.col1;")
	require.Nil(t, err)
}

func TestJoin04(t *testing.T) {
	_, err := Parse("select * from ks.tbl1 right join ks.tbl2 on tbl1.col1 = tbl2.col1;")
	require.Nil(t, err)
}

func TestGroupByOrderLimit01(t *testing.T) {
	_, err := Parse("select count(col2) as some_count from ks.tbl1 group by col1 having some_count < 5 order by col1 limit 5;")
	require.Nil(t, err)
}

func TestIntersection01(t *testing.T) {
	_, err := Parse("(select * from ks.tbl1) intersect (select * from ks.tbl2);")
	require.Nil(t, err)
}

func TestTableAggregate01(t *testing.T) {
	_, err := Parse("select TCOUNT() from ks.tbl1;")
	require.Nil(t, err)
}

func TestTableAggregate02(t *testing.T) {
	_, err := Parse("select TCOV(col1, col2) from ks.tbl1;")
	require.Nil(t, err)
}

func TestTableAggregate03(t *testing.T) {
	_, err := Parse("select TCORR(col1, col2) from ks.tbl1;")
	require.Nil(t, err)
}

func TestCollection01(t *testing.T) {
	_, err := Parse("insert into ks.tbl (col1) values ({ ")
	require.NotNil(t, err)
}

func TestCollection02(t *testing.T) {
	_, err := Parse("insert into ks.tbl (col1) values ({1});")
	require.Nil(t, err)
}

func TestCollection03(t *testing.T) {
	_, err := Parse("insert into ks.tbl (col1) values ({'1'});")
	require.Nil(t, err)
}

func TestCollection04(t *testing.T) {
	_, err := Parse("insert into ks.tbl (col1) values ({1..1});")
	require.NotNil(t, err)
}

func TestCollection05(t *testing.T) {
	_, err := Parse("insert into ks.tbl (col1) values ({'1});")
	require.NotNil(t, err)
}

func TestCollection06(t *testing.T) {
	_, err := Parse("insert into ks.tbl (col1) values ({^});")
	require.NotNil(t, err)
}

func TestCollection07(t *testing.T) {
	_, err := Parse("insert into ks.tbl (col1) values ({")
	require.NotNil(t, err)
}

func TestCollection08(t *testing.T) {
	_, err := Parse("insert into ks.tbl (col1) values ({,});")
	require.NotNil(t, err)
}

func TestCollection09(t *testing.T) {
	_, err := Parse("insert into ks.tbl (col1) values ({1,});")
	require.NotNil(t, err)
}

func TestCollection10(t *testing.T) {
	_, err := Parse("insert into ks.tbl (col1) values ({:1,1:1});")
	require.NotNil(t, err)
}

func TestCollection11(t *testing.T) {
	_, err := Parse("insert into ks.tbl (col1) values ({1:1,1});")
	require.NotNil(t, err)
}

func TestCollection12(t *testing.T) {
	_, err := Parse("insert into ks.tbl (col1) values ({1, 1:1});")
	require.NotNil(t, err)
}

func TestCollection13(t *testing.T) {
	_, err := Parse("insert into ks.tbl (col1) values ({1:1, 2:10});")
	require.Nil(t, err)
}

func TestTableCreation01(t *testing.T) {
	_, err := Parse("create table if not exists ks.tbl (col1 text, primary key(col1));")
	require.Nil(t, err)
}

func TestTableCreation02(t *testing.T) {
	_, err := Parse("create table if not exists ks.tbl (col1 text, col2 text, primary key(col1));")
	require.Nil(t, err)
}

func TestTableCreation03(t *testing.T) {
	_, err := Parse("create table if not exists ks.tbl (col1 text, col2 text, primary key(col1, col2));")
	require.Nil(t, err)
}

func TestTableCreation04(t *testing.T) {
	_, err := Parse("create table if not exists ks.tbl (col1 text, col2 text, primary key((col1, col2)));")
	require.Nil(t, err)
}

func TestTableCreation05(t *testing.T) {
	_, err := Parse("create table if not exists ks.tbl (col1 text, col2 text, col3 text, " +
		"primary key((col1, col2), col3));")
	require.Nil(t, err)
}

func TestTableCreation06(t *testing.T) {
	_, err := Parse("create table if not exists ks.tbl (col1 text, col2 text, col3 text, col4 text, " +
		"primary key((col1, col2), col3));")
	require.Nil(t, err)
}

func TestTableCreation07(t *testing.T) {
	_, err := Parse("create table if not exists ks.tbl (col1 text, col2 text, primary key(col1)) " +
		"with clustering order by (col2 asc);")
	require.Nil(t, err)
}

func TestTableCreation08(t *testing.T) {
	_, err := Parse("create table if not exists ks.tbl (col1 text, col2 text, col3 text, " +
		"primary key(col1, col2, col3)) with clustering order by (col2 asc, col3 desc) and compact storage;")
	require.Nil(t, err)
}

func TestTableCreation09(t *testing.T) {
	_, err := Parse("create table if not exists ks.tbl (col1 int, col2 set<text>, col3 list<double>, col4 map<int,float>, primary key(col1));")
	require.Nil(t, err)
}

func TestTableDrop01(t *testing.T) {
	_, err := Parse("drop table ks.tbl;")
	require.Nil(t, err)
}
