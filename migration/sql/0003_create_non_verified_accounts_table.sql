--liquibase formatted sql

--comment: Create non-verified accounts table
--changeset Septem:0003_create_non_verified_accounts_table

CREATE TABLE IF NOT EXISTS non_verified_accounts (
	id BIGSERIAL NOT NULL PRIMARY KEY,
	username VARCHAR(30) NOT NULL UNIQUE,
	email VARCHAR(50) NOT NULL UNIQUE,
	phone VARCHAR(20),
    password VARCHAR(255) NOT NULL,
	login_type VARCHAR(255) NOT NULL, 
    create_at TIMESTAMP NOT NULL,
    update_at TIMESTAMP NOT NULL
);

--rollback DROP TABLE non_verified_accounts