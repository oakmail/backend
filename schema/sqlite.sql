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
	account       UNSIGNED BIG INT
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

CREATE TABLE tokens (
	id            CHARACTER(20) PRIMARY KEY DEFAULT (uuid()),
	date_created  DATETIME,
	date_modified DATETIME,
	owner         UNSIGNED BIG INT,
	expiry_date   DATETIME,
	type          TEXT,
	perms         TEXT,
	application   UNSIGNED BIG INT
);
