// Copyright 2017 Verizon. All rights reserved.
// See provided LICENSE file for use of this source code.

CREATE KEYSPACE IF NOT EXISTS account WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '3'};

CREATE TABLE if not exists account.events (
    id              uuid,
    accountid       uuid,
    name            text,
    createdon       timestamp,
    PRIMARY KEY (accountid, id)
);

CREATE TABLE if not exists account.snippets (
    id               uuid,
    accountid        uuid,
    name             text,
    createdon        timestamp,
    updatedon        timestamp,
    runtime          text,
    mainfn           text,
    url              text,
    code             text,
    timeout          int,
    memory           bigint,
    callback         text,
    description      text,
    eventtype        text,
    eventid          text,
    PRIMARY KEY (accountid, id)
);

CREATE TABLE if not exists account.invocations (
    id              timeuuid,
    accountid       uuid,
    rteid           uuid,
    snippetid       uuid,
    partition       int,
    createdon       timestamp,
    startedon       timestamp,
    finishedon      timestamp,
    updatedon       timestamp,
    elapsedtime     double,
    runtime         text,
    mainfn          text,
    url             text,
    code            text,
    timeout         int,
    memory          bigint,
    callback        text,
    args            blob,
    stdout          text,
    result          text,
    status          text,
    errordescr      text,
    PRIMARY KEY (accountid, id)
) WITH default_time_to_live = 7889238 and CLUSTERING ORDER BY (id DESC);

create index if not exists on account.invocations(snippetid);

CREATE TABLE if not exists account.mappings (
    id              uuid,
    accountid       uuid,
    createdon       timestamp,
    eventid         uuid,
    snippetid       uuid,
    PRIMARY KEY (accountid, id)
);
create index if not exists on account.mappings(eventid);

CREATE TABLE if not exists account.access (
    id          uuid,
    createdon   timestamp,
    accountid   text,
    userid      text,
    permissions text,
    notebookid  text,
    PRIMARY KEY (id)
);

create index if not exists on account.access(accountid);
create index if not exists on account.access(userid);
create index if not exists on account.access(notebookid);

CREATE TABLE if not exists account.notebooks (
    id              uuid,
    createdon   	timestamp,
    version         timeuuid,
    data            text,
    PRIMARY KEY (id)
);

CREATE TABLE if not exists account.templates (
     id              uuid,
     createdon   	 timestamp,
     version         timeuuid,
     accountid       text,
     userid          text,
     name            text,
     description     text,
     type            text,
     data            text,
     published       int,
     hash            text,
     PRIMARY KEY (id)
);

create index if not exists on account.templates(accountid);
create index if not exists on account.templates(userid);
create index if not exists on account.templates(type);
create index if not exists on account.templates(published);
create index if not exists on account.templates(hash);

CREATE TYPE if not exists account.datasetColumn (
    name     text,
    datatype text
);

CREATE TYPE if not exists account.datasetTable (
    name     text,
    columns  map<text, frozen <datasetColumn>>
);

CREATE TABLE if not exists account.datasets (
    id               uuid,
    accountid        uuid,
    datasourceid     text,
    name             text,
    description      text,
    tables           map<text, frozen <datasetTable>>,
    createdon        timestamp,
    updatedon        timestamp,
    PRIMARY KEY (accountid, id)
);

create index if not exists on account.datasets(name);

CREATE TABLE if not exists account.datasources (
    id               uuid,
    accountid        uuid,
    name             text,
    description      text,
    protocol         text,
    host             text,
    port             int,
    options          map<text,text>,
    createdon        timestamp,
    updatedon        timestamp,
    PRIMARY KEY (accountid, id)
);
