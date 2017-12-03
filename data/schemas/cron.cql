// Copyright 2017 Verizon. All rights reserved.
// See provided LICENSE file for use of this source code.

CREATE KEYSPACE IF NOT EXISTS cron WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '3'};

CREATE TABLE if not exists cron.jobs (
    datacenter       text,
    accountid        uuid,
    id               uuid,
    name             text,
    snippetid        uuid,
    schedule         text,
    disabled         boolean,
    updatedon        timestamp,
    description      text,
    PRIMARY KEY (datacenter, accountid, id)
);
