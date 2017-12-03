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

package nsQL

import (
	"errors"
	"os"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/rte-lua/modules/nsQL/compiler"
	"github.com/verizonlabs/northstar/rte-lua/modules/nsQL/compiler/cassandra"
	"github.com/verizonlabs/northstar/rte-lua/modules/nsQL/compiler/spark"
	"github.com/verizonlabs/northstar/rte-lua/modules/nsQL/constants"
)

func getCompiler(processing *compiler.Processing) (compiler.Compiler, error) {
	if processing.DataSource.Protocol == "" || processing.DataSource.Connection.Host == "" ||
		processing.DataSource.Connection.Port == "" {
		return nil, errors.New("nsQL error: data source protocol, host and port must be defined")
	}

	switch processing.Backend {
	case constants.NATIVE:
		switch processing.DataSource.Protocol {
		case constants.CASSANDRA:
			return cassandra.NewCassandraCompiler(processing.DataSource.Connection), nil
		default:
			return nil, errors.New("nsQL error: " + processing.DataSource.Protocol + " is not a supported " +
				"data source protocol")
		}
	case constants.SPARK:
		if spark, err := getSpark(processing.DataSource); err != nil {
			return nil, err
		} else {
			return spark, nil
		}
	default:
		return nil, errors.New("nsQL error: " + processing.Backend + " is not a supported data " +
			"processing backend")
	}
}

func getSpark(dataSource *compiler.DataSource) (*spark.SparkCompiler, error) {
	hostPort := os.Getenv("DPE_SPARK_HOST_PORT")
	if hostPort == "" {
		return nil, errors.New("nsQL error: unable to find spark host port")
	}

	mlog.Info("Spark host port: %s", hostPort)
	return spark.NewSparkCompiler(hostPort, dataSource), nil
}
