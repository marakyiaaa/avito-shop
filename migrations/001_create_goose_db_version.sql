CREATE TABLE IF NOT EXISTS goose_db_version (
    id serial NOT NULL,
    version_id bigint NOT NULL,
    is_applied boolean NOT NULL,
    tstamp timestamp NULL default now(),
    PRIMARY KEY(id)
);
