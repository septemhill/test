--liquibase formatted sql

--comment: Initialize tables 
--changeset Septem:0001_initialize_tables
CREATE TABLE IF NOT EXISTS accounts (
	id BIGSERIAL NOT NULL PRIMARY KEY,
	username VARCHAR(30) NOT NULL UNIQUE,
	email VARCHAR(50) NOT NULL UNIQUE,
	phone VARCHAR(20)
);

CREATE TABLE IF NOT EXISTS accounts_private (
	id BIGSERIAL NOT NULL PRIMARY KEY,
	email VARCHAR(50) NOT NULL,
	password VARCHAR(255) NOT NULL,
	login_type VARCHAR(255) NOT NULL, 
	CONSTRAINT fk_email FOREIGN KEY (email) REFERENCES accounts(email),
	CONSTRAINT login_type_check CHECK (login_type in ('NORMAL', 'OAUTH2_GOOGLE'))
);

CREATE TABLE IF NOT EXISTS perms (
	username VARCHAR(30) NOT NULL,
	perm int NOT NULL,
	PRIMARY KEY(username, perm),
	CONSTRAINT fk_username FOREIGN KEY (username) REFERENCES accounts(username)
);

CREATE TABLE IF NOT EXISTS articles (
	id BIGSERIAL NOT NULL PRIMARY KEY,
	author VARCHAR(30) NOT NULL,
	title VARCHAR(100) NOT NULL,
	content VARCHAR(10000),
	create_at TIMESTAMPTZ NOT NULL,
	update_at TIMESTAMPTZ NOT NULL,
	CONSTRAINT fk_author FOREIGN KEY (author) REFERENCES accounts(username)
);

CREATE TABLE IF NOT EXISTS articles_tag (
	art_id BIGSERIAL NOT NULL,
	tag_name VARCHAR(50) NOT NULL,
	PRIMARY KEY (art_id, tag_name),
	CONSTRAINT fk_art_id FOREIGN KEY (art_id) REFERENCES articles(id)
);

CREATE TABLE IF NOT EXISTS comments (
	id BIGSERIAL NOT NULL PRIMARY KEY,
	art_id BIGSERIAL NOT NULL, 
	author VARCHAR(30) NOT NULL,
	content VARCHAR(1000) NOT NULL,
	create_at TIMESTAMPTZ NOT NULL,
	update_at TIMESTAMPTZ NOT NULL,
	CONSTRAINT fk_article_id FOREIGN KEY (art_id) REFERENCES articles(id),
	CONSTRAINT fk_author FOREIGN KEY (author) REFERENCES accounts(username)
);