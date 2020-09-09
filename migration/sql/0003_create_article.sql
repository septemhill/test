--liquibase formatted sql

--comment: Create table `article`, and add sample data
--changeset Septem:0003_create_article_table_add_sample_data
CREATE TABLE article(
	aid BIGSERIAL NOT NULL PRIMARY KEY,
	author VARCHAR NOT NULL,
	title VARCHAR NOT NULL,
	content VARCHAR
);

INSERT INTO article VALUES 
	(DEFAULT, 'Septem Li', 'No news is good news', 'Balabababa'),
	(DEFAULT, 'Septem Li', 'Today''s weather is good', 'Hmmmm .... really?');

--rollback DROP TABLE IF EXISTS article
