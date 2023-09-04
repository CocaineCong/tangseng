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
