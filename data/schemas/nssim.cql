// Copyright 2017 Verizon. All rights reserved.
// See provided LICENSE file for use of this source code.

DROP KEYSPACE IF EXISTS nssim;
CREATE KEYSPACE IF NOT EXISTS nssim WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '3'};

CREATE TABLE IF NOT EXISTS nssim.sampledata (
    id uuid,
    rowid uuid,
    name varchar,
    createdtime timestamp,
    datevalue date,
    timevalue time,
    numvalue int,
    money double,
    floatvalue float,
    varintvalue varint,
    maxvalue bigint,
    data blob,
    ip inet,
    array set<text>,
    json varchar,
    mapdata map<text, blob>,
    PRIMARY KEY (rowid, id)
);
