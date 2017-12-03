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
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/verizonlabs/northstar/pkg/mlog"
	dktUtils "github.com/verizonlabs/northstar/pkg/utils"
	"github.com/verizonlabs/northstar/data/config"
)

type DB struct {
	hostStrArr []string
	keyspace   string
	cluster    *gocql.ClusterConfig
	Session    *gocql.Session
}

func GetHostStrArr() []string {
	hostStr := os.Getenv("CASSANDRA_HOST")
	hostStr = dktUtils.HostsToIps(hostStr)
	if hostStr == "" {
		hostStr = os.Getenv("HOST")
		mlog.Error("No host mentioned in CASSANDRA_HOST.. Defaulting to local host")
	}
	return strings.Split(hostStr, ",")
}

func NewDB(keyspace string) *DB {
	db := new(DB)
	hostStr := os.Getenv("CASSANDRA_HOST")
	if hostStr == "" {
		hostStr = os.Getenv("HOST")
		mlog.Error("No host mentioned in CASSANDRA_HOST.. Defaulting to local host")
	}

	portStr := os.Getenv("CASSANDRA_NATIVE_TRANSPORT_PORT")
	if portStr == "" {
		portStr = "9042"
		mlog.Error("No cassandra native transport port mentioned.. Defaulting to 9042")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		mlog.Error("Error parsing native transport port from environment", err)
	}
	rp := new(gocql.SimpleRetryPolicy)

	rp.NumRetries = 5

	db.hostStrArr = GetHostStrArr()
	db.keyspace = keyspace
	db.cluster = gocql.NewCluster(db.hostStrArr...)
	db.cluster.Keyspace = keyspace
	db.cluster.Consistency = gocql.LocalQuorum
	db.cluster.Timeout = 3 * time.Second
	db.cluster.NumConns = 5
	db.cluster.RetryPolicy = rp
	db.cluster.ProtoVersion = config.CassandraProtoVersion
	db.cluster.PageSize = 500
	db.cluster.Port = port
	if config.CassandraDatacenter != "" {
		db.cluster.HostFilter = gocql.DataCentreHostFilter(config.CassandraDatacenter)
	}
	if config.CassandraCQLVersion != "" {
		db.cluster.CQLVersion = config.CassandraCQLVersion
	}
	if config.CassandraAuthEnabled == true {
		db.cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: config.CassandraUsername,
			Password: config.CassandraPassword,
		}
	}
	db.Session = db.GetSession()
	return db
}

func (db *DB) GetSession() *gocql.Session {

	if db.Session == nil || db.Session.Closed() {
		var err error
		db.cluster.Hosts = GetHostStrArr()
		db.Session, err = db.cluster.CreateSession()
		if err != nil {
			mlog.Alarm("Error creating database session", db.cluster, err)
		}
	}
	return db.Session
}

// GetSession() function is depracated.. Use GetSessionWithError() in dblib.go when you create new data service modules
// to return the error
func (db *DB) GetSessionWithError() (*gocql.Session, error) {

	var err error
	if db.Session == nil || db.Session.Closed() {

		db.Session, err = db.cluster.CreateSession()
		if err != nil {
			mlog.Alarm("Error creating database session", db.cluster, err)
		}
	}
	return db.Session, err
}
