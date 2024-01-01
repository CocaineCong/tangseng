-- Licensed to the Apache Software Foundation (ASF) under one
-- or more contributor license agreements.  See the NOTICE file
-- distributed with this work for additional information
-- regarding copyright ownership.  The ASF licenses this file
-- to you under the Apache License, Version 2.0 (the
-- "License"); you may not use this file except in compliance
-- with the License.  You may obtain a copy of the License at
--
--   http://www.apache.org/licenses/LICENSE-2.0
--
-- Unless required by applicable law or agreed to in writing,
-- software distributed under the License is distributed on an
-- "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
-- KIND, either express or implied.  See the License for the
-- specific language governing permissions and limitations
-- under the License.

CREATE DATABASE TEST;

USE TEST;

CREATE TABLE `sr_on_mac` (
     `c0` int(11) NULL COMMENT "",
     `c1` date NULL COMMENT "",
     `c2` datetime NULL COMMENT "",
     `c3` varchar(65533) NULL COMMENT ""
) ENGINE=OLAP
DUPLICATE KEY(`c0`)
PARTITION BY RANGE (c1) (
  START ("2022-02-01") END ("2022-02-10") EVERY (INTERVAL 1 DAY)
)
DISTRIBUTED BY HASH(`c0`) BUCKETS 1
PROPERTIES (
"replication_num" = "1",
"in_memory" = "false",
"storage_format" = "DEFAULT"
);

CREATE TABLE `test_upload` (
     `doc_id` int(11) NULL COMMENT "",
     `url` varchar(255) NULL COMMENT "",
     `title` varchar(65533) NULL COMMENT "",
     `desc` varchar(65533) NULL COMMENT "",
     `score` FLOAT COMMENT ""
) ENGINE=OLAP
     DUPLICATE KEY(`doc_id`)
DISTRIBUTED BY HASH(`doc_id`) BUCKETS 1
PROPERTIES (
"replication_num" = "1",
"in_memory" = "false",
"storage_format" = "DEFAULT"
);

insert into sr_on_mac values (1, '2022-02-01', '2022-02-01 10:47:57', '111');
insert into sr_on_mac values (2, '2022-02-02', '2022-02-02 10:47:57', '222');
insert into sr_on_mac values (3, '2022-02-03', '2022-02-03 10:47:57', '333');
