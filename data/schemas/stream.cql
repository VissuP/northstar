// Copyright 2017 Verizon. All rights reserved.
// See provided LICENSE file for use of this source code.

CREATE KEYSPACE IF NOT EXISTS stream WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '3'};

CREATE TABLE stream.jobs (
    id           uuid,
    accountid    uuid,
    invocationId uuid,
    memory       bigint,
    source       blob,
    functions    blob,
    createdon    timestamp,
    updatedon    timestamp,
    status       text,
    errordescr   text,
    description  text,
    PRIMARY KEY (accountid, id)
);
