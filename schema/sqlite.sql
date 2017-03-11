CREATE TABLE accounts (
	id                 UNSIGNED BIG INT PRIMARY KEY DEFAULT (next_id()),
	date_created       DATETIME,
	date_modified      DATETIME,
	type               TEXT,
	main_address       TEXT,
	identity           TEXT,
	password           TEXT,
	subscription       TEXT,
	blocked            BOOLEAN,
	alt_email          TEXT,
	alt_email_verified DATETIME
);

CREATE TABLE addresses (
	id            TEXT PRIMARY KEY,
	date_created  DATETIME,
	date_modified DATETIME,
	account       UNSIGNED BIG INT,
	public_key    UNSIGNED BIG INT
);

CREATE TABLE applications (
	id            UNSIGNED BIG INT PRIMARY KEY DEFAULT (next_id()),
	date_created  DATETIME,
	date_modified DATETIME,
	owner         UNSIGNED BIG INT,
	secret        TEXT,
	callback      TEXT,
	name          TEXT,
	email         TEXT,
	home_page     TEXT,
	description	  TEXT
);

CREATE TABLE public_keys (
	id                  UNSIGNED BIG INT PRIMARY KEY,
	date_created        DATETIME,
	date_modified       DATETIME,
	owner               UNSIGNED BIG INT,
	algorithm           INT,
	length              INT,
	body                BLOB,
	key_id_string       TEXT,
	key_id_short_string TEXT,
	master_key          UNSIGNED BIG INT
);

CREATE TABLE public_key_identities (
	id             UNSIGNED BIG INT PRIMARY KEY DEFAULT (next_id()),
	date_created   DATETIME,
	date_modified  DATETIME,
	name           TEXT,
	self_signature UNSIGNED BIG INT,
	signatures     JSON
);

CREATE TABLE public_key_signatures (
	id                     UNSIGNED BIG INT PRIMARY KEY DEFAULT (next_id()),
	type                   INT,
	algorithm              INT,
	hash                   INT,
	creation_time          DATETIME,
	sig_lifetime_secs      INT,
	key_lifetime_secs      INT,
	issuer_key_id          UNSIGNED BIG INT,
	is_primary_id          BOOLEAN,
	revocation_reason      INT,
	revocation_reason_text TEXT
);

CREATE TABLE resources (
	id            CHARACTER(20) PRIMARY KEY DEFAULT (uuid()),
	date_created  DATETIME,
	date_modified DATETIME,
	owner         UNSIGNED BIG INT,
	meta          JSON,
	tags          JSON,
	file          TEXT,
	upload_token  TEXT
);

CREATE TABLE tokens (
	id             CHARACTER(20) PRIMARY KEY DEFAULT (uuid()),
	date_created   DATETIME,
	date_modified  DATETIME,
	owner          UNSIGNED BIG INT,
	expiry_date    DATETIME,
	type           TEXT,
	perms          TEXT,
	application    UNSIGNED BIG INT,
	reference_type TEXT,
	reference_id   UNSIGNED BIG INT
);
