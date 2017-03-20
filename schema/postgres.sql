CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE SEQUENCE table_id_seq START 1;

CREATE OR REPLACE FUNCTION next_id(OUT result bigint) AS $$
DECLARE
    our_epoch bigint := 1487595908000;
    seq_id bigint;
    now_millis bigint;
    shard_id int := 0;
BEGIN
    SELECT nextval('table_id_seq') %% 1024 INTO seq_id;
    SELECT FLOOR(EXTRACT(EPOCH FROM clock_timestamp()) * 1000) INTO now_millis;
    result := (now_millis - our_epoch) << 23;
    result := result | (shard_id <<10);
    result := result | (seq_id);
END;
    $$ LANGUAGE PLPGSQL;


CREATE TABLE accounts (
	id                 BIGINT PRIMARY KEY DEFAULT (next_id()),
	date_created       TIMESTAMP WITH TIME ZONE,
	date_modified      TIMESTAMP WITH TIME ZONE,
	type               TEXT,
	main_address       TEXT,
	identity           TEXT,
	password           TEXT,
	subscription       TEXT,
	blocked            BOOLEAN,
	alt_email          TEXT,
	alt_email_verified TIMESTAMP WITH TIME ZONE
);

CREATE TABLE addresses (
	id            TEXT PRIMARY KEY,
	styled_id     TEXT,
	date_created  TIMESTAMP WITH TIME ZONE,
	date_modified TIMESTAMP WITH TIME ZONE,
	account       BIGINT,
	public_key    BIGINT
);

CREATE TABLE applications (
	id            BIGINT PRIMARY KEY DEFAULT (next_id()),
	date_created  TIMESTAMP WITH TIME ZONE,
	date_modified TIMESTAMP WITH TIME ZONE,
	owner         BIGINT,
	secret        TEXT,
	callback      TEXT,
	name          TEXT,
	email         TEXT,
	home_page     TEXT,
	description	  TEXT
);

CREATE TABLE public_keys (
	id                  BIGINT PRIMARY KEY,
	date_created        TIMESTAMP WITH TIME ZONE,
	date_modified       TIMESTAMP WITH TIME ZONE,
	owner               BIGINT,
	algorithm           INT,
	length              INT,
	body                BYTEA,
	key_id_string       TEXT,
	key_id_short_string TEXT,
	master_key          BIGINT
);

CREATE TABLE public_key_identities (
	id             BIGINT PRIMARY KEY DEFAULT (next_id()),
	date_created   TIMESTAMP WITH TIME ZONE,
	date_modified  TIMESTAMP WITH TIME ZONE,
	name           TEXT,
	self_signature BIGINT
);

CREATE TABLE public_key_signatures (
	id                     BIGINT PRIMARY KEY DEFAULT (next_id()),
	identity               BIGINT,
	type                   INT,
	algorithm              INT,
	hash                   INT,
	creation_time          TIMESTAMP WITH TIME ZONE,
	sig_lifetime_secs      INT,
	key_lifetime_secs      INT,
	issuer_key_id          BIGINT,
	is_primary_id          BOOLEAN,
	revocation_reason      INT,
	revocation_reason_text TEXT
);

CREATE TABLE resources (
	id            BIGINT PRIMARY KEY DEFAULT (next_id()),
	date_created  TIMESTAMP WITH TIME ZONE,
	date_modified TIMESTAMP WITH TIME ZONE,
	owner         BIGINT,
	meta          JSON,
	tags          text[],
	file          TEXT,
	upload_token  TEXT
);

CREATE TABLE tokens (
	id             UUID PRIMARY KEY DEFAULT (gen_random_uuid()),
	date_created   TIMESTAMP WITH TIME ZONE,
	date_modified  TIMESTAMP WITH TIME ZONE,
	owner          BIGINT,
	expiry_date    TIMESTAMP WITH TIME ZONE,
	type           TEXT,
	perms          TEXT,
	application    BIGINT,
	reference_type TEXT,
	reference_id   BIGINT
);
