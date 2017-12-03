// Copyright 2017 Verizon. All rights reserved.
// See provided LICENSE file for use of this source code.

ALTER TABLE account.invocations ADD runtime text;
ALTER TABLE account.invocations ADD mainfn text;
ALTER TABLE account.invocations ADD url text;
ALTER TABLE account.invocations ADD code text;
ALTER TABLE account.invocations ADD timeout int;
ALTER TABLE account.invocations ADD memory bigint;
ALTER TABLE account.invocations ADD callback text;
ALTER TABLE account.invocations ADD args blob;
ALTER TABLE account.invocations ADD partition int;
ALTER TABLE account.invocations ADD updatedon timestamp;

ALTER TABLE account.snippets ADD memory bigint;
ALTER TABLE account.snippets ADD callback text;

ALTER TABLE account.templates ADD hash text;
