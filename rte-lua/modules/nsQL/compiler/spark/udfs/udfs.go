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

package udfs

const (
	MAP_BLOB_JSON_FETCH = "var map_blob_json_fetch = udf { (fields: Map[String, Array[Byte]], capability: String, field: " +
		"String) => if(fields.contains(capability)) {var byteBuffer = java.nio.ByteBuffer.wrap(fields(capability)); " +
		"var charBuffer = java.nio.charset.StandardCharsets.UTF_8.decode(byteBuffer); var jsonAsOption = " +
		"scala.util.parsing.json.JSON.parseFull(charBuffer.toString()); if(jsonAsOption != None) {" +
		"var jsonAsMap = jsonAsOption.get.asInstanceOf[Map[String, Any]]; if(jsonAsMap.contains(field)) " +
		"{ jsonAsMap(field).toString() } else { \\\"\\\" } } else { \\\"\\\" } } else { \\\"\\\" } }"
	JSON_FETCH = "var json_fetch = udf { (stats: String, field: String) => var jsonAsOption = " +
		"scala.util.parsing.json.JSON.parseFull(stats); if(jsonAsOption != None) { var jsonAsMap = " +
		"jsonAsOption.get.asInstanceOf[Map[String, Any]]; if(jsonAsMap.contains(field)) { " +
		"jsonAsMap(field).toString()} else {\\\"\\\"} } else { \\\"\\\" } }"
	SUBTRACT_TIMESTAMPS = "var subtract_timestamps = udf { (operand1: java.sql.Timestamp, operand2: " +
		"java.sql.Timestamp) => if(operand1 == null || operand2 == null){ 0 }else{ val time1 = " +
		"operand1.getTime()*1000000; val time2 = operand2.getTime()*1000000; time1 - time2 } }"
)
