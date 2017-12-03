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
	"github.com/stretchr/testify/require"
	"github.com/verizonlabs/northstar/rte-lua/modules/nsQL/compiler"
	"testing"
)

const (
	CASSANDRA       = "cassandra"
	CASSANDRA_HOST  = "0.0.0.0"
	CASSANDRA_PORT  = "0"
	SPARK_HOST_PORT = "0.0.0.0:0"
)

var options = &compiler.Options{CassandraFetchLimit: 100}
var datasource = &compiler.DataSource{
	Protocol:   CASSANDRA,
	Connection: &compiler.Connection{Host: CASSANDRA_HOST, Port: CASSANDRA_PORT},
}
var transcompiler = NewSparkCompiler(SPARK_HOST_PORT, datasource)

func TestCompilationError01(t *testing.T) {
	query := "SELECT * FROM devicetxn.battery_history"
	_, err := transcompiler.compile(query, options)
	require.NotNil(t, err)
}

func TestAccessError01(t *testing.T) {
	query := "SELECT * FROM devicetxn.battery_history;"
	_, err := transcompiler.Run(query, options)
	require.NotNil(t, err)
}

func TestAccessError02(t *testing.T) {
	query := "SELECT * FROM devicetxn.battery_history;"
	_, err := transcompiler.Run(query, options)
	require.NotNil(t, err)
}

func TestLimitError03(t *testing.T) {
	query := "SELECT * FROM devicetxn.battery_history;"
	_, err := transcompiler.Run(query, nil)
	require.NotNil(t, err)
}

func TestSelect01(t *testing.T) {
	query := "SELECT TCOUNT() FROM devicetxn.battery_history;"
	_, err := transcompiler.compile(query, options)
	require.Nil(t, err)
}

func TestSelect02(t *testing.T) {
	query := "SELECT bh.imsi FROM devicetxn.battery_history as bh;"
	_, err := transcompiler.compile(query, options)
	require.Nil(t, err)
}

func TestSelect03(t *testing.T) {
	query := "SELECT bh1.imsi FROM devicetxn.battery_history as bh1 JOIN devicetxn.battery_history as bh2;"
	_, err := transcompiler.compile(query, options)
	require.Nil(t, err)
}

func TestSelect04(t *testing.T) {
	query := "SELECT DISTINCT imsi as id FROM devicetxn.battery_history WHERE battery_level > 0;"
	_, err := transcompiler.compile(query, options)
	require.Nil(t, err)
}

func TestSelect05(t *testing.T) {
	query := "SELECT count(imsi) as countid FROM devicetxn.battery_history GROUP BY imsi as id HAVING id > 0;"
	_, err := transcompiler.compile(query, options)
	require.Nil(t, err)
}

func TestSelect06(t *testing.T) {
	query := "SELECT imsi, battery_level FROM devicetxn.battery_history ORDER BY imsi LIMIT 5;"
	_, err := transcompiler.compile(query, options)
	require.Nil(t, err)
}

func TestSelect07(t *testing.T) {
	query := "(SELECT battery_level FROM devicetxn.battery_history) UNION (SELECT battery_level FROM " +
		"devicetxn.battery_history) UNION ALL (SELECT battery_level FROM devicetxn.battery_history) INTERSECT" +
		"(SELECT battery_level FROM devicetxn.battery_history);"
	_, err := transcompiler.compile(query, options)
	require.Nil(t, err)
}

func TestSelect08(t *testing.T) {
	query := "SELECT bh1.imsi FROM devicetxn.battery_history as bh1 OUTER JOIN devicetxn.battery_history as bh2 " +
		"ON bh1.imsi=bh1.imsi;"
	_, err := transcompiler.compile(query, options)
	require.Nil(t, err)
}

func TestSelect09(t *testing.T) {
	query := "SELECT bh1.imsi FROM devicetxn.battery_history as bh1 LEFT JOIN devicetxn.battery_history as bh2 " +
		"ON bh1.imsi=bh1.imsi;"
	_, err := transcompiler.compile(query, options)
	require.Nil(t, err)
}

func TestSelect10(t *testing.T) {
	query := "SELECT bh1.imsi FROM devicetxn.battery_history as bh1 LEFT SEMI JOIN devicetxn.battery_history as " +
		"bh2 ON bh1.imsi=bh1.imsi;"
	_, err := transcompiler.compile(query, options)
	require.Nil(t, err)
}

func TestSelect11(t *testing.T) {
	query := "SELECT bh1.imsi FROM (SELECT imsi FROM devicetxn.battery_history) as bh1;"
	_, err := transcompiler.compile(query, options)
	require.Nil(t, err)
}

func TestSelect12(t *testing.T) {
	query := "SELECT imsi FROM devicetxn.battery_history WHERE (msg_type = 'RetMsgLost' and event_time < " +
		"'2017-02-14 21:19:30' - 'INTERVAL 1 DAY') or msg_type = 'RetMsgInitial';"
	_, err := transcompiler.compile(query, options)
	require.Nil(t, err)
}

func TestSelect13(t *testing.T) {
	query := "SELECT TCORR(battery_level, current_voltage) FROM devicetxn.battery_history;"
	_, err := transcompiler.compile(query, options)
	require.Nil(t, err)
}

func TestSelect14(t *testing.T) {
	query := "SELECT TCOV(battery_level, current_voltage) FROM devicetxn.battery_history;"
	_, err := transcompiler.compile(query, options)
	require.Nil(t, err)
}

func TestSelect15(t *testing.T) {
	query := "SELECT MIN(battery_level) FROM devicetxn.battery_history;"
	_, err := transcompiler.compile(query, options)
	require.Nil(t, err)
}

func TestSelect16(t *testing.T) {
	query := "SELECT DAY(event_time), JSON_FETCH(msg_type, 'lost'), MAP_BLOB_JSON_FETCH(msg_type, 'lost', 'time'), " +
		"NOW() FROM devicetxn.battery_history;"
	_, err := transcompiler.compile(query, options)
	require.Nil(t, err)
}

func TestSelect17(t *testing.T) {
	query := "SELECT original_voltage - current_voltage, battery_level - 10, event_time + -'INTERVAL 1 DAY'" +
		"FROM devicetxn.battery_history;"
	_, err := transcompiler.compile(query, options)
	require.Nil(t, err)
}

func TestSelect18(t *testing.T) {
	query := "SELECT battery_level FROM devicetxn.battery_history WHERE imsi IN (SELECT imsi FROM " +
		"devicetxn.battery_history) and imsi NOT IN (SELECT imsi FROM devicetxn.battery_history) and imsi IS " +
		"NULL and imsi IS NOT NULL and imsi != 5;"
	_, err := transcompiler.compile(query, options)
	require.Nil(t, err)
}

func TestSelect19(t *testing.T) {
	query := "SELECT battery_level FROM devicetxn.battery_history WHERE imsi IN (SELECT imsi FROM " +
		"devicetxn.battery_history) or imsi NOT IN (SELECT imsi FROM devicetxn.battery_history) and imsi IS " +
		"NULL and imsi IS NOT NULL or imsi != 5;"
	_, err := transcompiler.compile(query, options)
	require.Nil(t, err)
}

func TestSelect20(t *testing.T) {
	query := "SELECT battery_level FROM devicetxn.battery_history WHERE imsi IN (SELECT imsi FROM " +
		"devicetxn.battery_history) or imsi IS NULL;"
	_, err := transcompiler.compile(query, options)
	require.Nil(t, err)
}

func TestSelect21(t *testing.T) {
	query := "SELECT battery_level FROM devicetxn.battery_history WHERE imsi IS NULL or imsi IN (SELECT imsi FROM " +
		"devicetxn.battery_history);"
	_, err := transcompiler.compile(query, options)
	require.Nil(t, err)
}

func TestSelect22(t *testing.T) {
	query := "SELECT battery_level FROM devicetxn.battery_history WHERE imsi IS NULL and imsi IN (SELECT imsi FROM " +
		"devicetxn.battery_history);"
	_, err := transcompiler.compile(query, options)
	require.Nil(t, err)
}

func TestSelect23(t *testing.T) {
	query := "SELECT battery_level FROM devicetxn.battery_history WHERE imsi IS NULL and imsi IN (SELECT imsi FROM " +
		"devicetxn.battery_history);"
	_, err := transcompiler.compile(query, &compiler.Options{CassandraFetchLimit: 100})
	require.Nil(t, err)
}

func TestSelect24(t *testing.T) {
	query := "SELECT battery_level FROM devicetxn.battery_history WHERE imsi IS NULL and imsi IN (SELECT imsi FROM " +
		"devicetxn.battery_history);"
	_, err := transcompiler.compile(query, &compiler.Options{})
	require.Nil(t, err)
}

func TestSelect25(t *testing.T) {
	query := "SELECT battery_level FROM devicetxn.battery_history WHERE imsi IS NULL and imsi IN (SELECT imsi FROM " +
		"devicetxn.battery_history);"
	_, err := transcompiler.compile(query, &compiler.Options{})
	require.Nil(t, err)
}

func TestSelect26(t *testing.T) {
	query := "SELECT battery_level FROM devicetxn.battery_history WHERE imsi IS NULL and imsi IN (SELECT imsi FROM " +
		"devicetxn.battery_history);"
	_, err := transcompiler.compile(query, &compiler.Options{CassandraFetchLimit: 100000000})
	require.Nil(t, err)
}

func TestSelect27(t *testing.T) {
	query := "SELECT battery_level FROM devicetxn.battery_history WHERE imsi IS NULL and imsi IN (SELECT imsi FROM " +
		"devicetxn.battery_history);"
	_, err := transcompiler.compile(query, &compiler.Options{})
	require.Nil(t, err)
}

func TestSelectError01(t *testing.T) {
	query := "INSERT INTO devicetxn.battery_history (imsi) VALUES (5);"
	_, err := transcompiler.compile(query, options)
	require.NotNil(t, err)
}
