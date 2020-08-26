--liquibase formatted sql

--comment: Create table `drytable`, and add sample data
--changeset Septem:0001_create_drytable_table_add_sample_data
CREATE TABLE drytable (
	id BIGSERIAL NOT NULL PRIMARY KEY,
	fname CHAR(25),
	lname CHAR(25)
);

INSERT INTO drytable VALUES 
(DEFAULT, 'Septem', 'Li'),
(DEFAULT, 'Nicole', 'Chen');

--rollback DROP TABLE IF EXISTS drytable;
