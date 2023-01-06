CREATE DATABASE "MMS_TEST";

\connect "MMS_TEST";

CREATE TABLE IF NOT EXISTS "mms_pairs" (
    id         serial PRIMARY KEY,
    pair       varchar(10) not null,
    time_stamp timestamp,
    mms_20     decimal null,
    mms_50     decimal null,
    mms_200    decimal null
);
