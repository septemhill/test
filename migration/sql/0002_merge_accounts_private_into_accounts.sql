--liquibase formatted sql

--comment: Merge private information into `accounts` table
--changeset Septem:0002_merge_accounts_private_into_accounts

ALTER TABLE accounts ADD COLUMN password VARCHAR(255);
ALTER TABLE accounts ADD COLUMN login_type VARCHAR(255);
ALTER TABLE accounts ADD COLUMN create_at TIMESTAMP;
ALTER TABLE accounts ADD COLUMN update_at TIMESTAMP;

UPDATE accounts SET password = accounts_private.password FROM accounts_private WHERE accounts.email = accounts_private.email;
UPDATE accounts SET login_type = accounts_private.login_type FROM accounts_private WHERE accounts.email = accounts_private.email;
UPDATE accounts SET create_at = now(), update_at = now();

ALTER TABLE accounts ALTER COLUMN password SET NOT NULL;
ALTER TABLE accounts ALTER COLUMN login_type SET NOT NULL;
ALTER TABLE accounts ADD CONSTRAINT login_type_check CHECK (login_type IN ('NORMAL', 'OAUTH2_GOOGLE'));
ALTER TABLE accounts ALTER COLUMN create_at SET NOT NULL;
ALTER TABLE accounts ALTER COLUMN update_at SET NOT NULL;

DROP TABLE accounts_private;

--rollback CREATE TABLE IF NOT EXISTS accounts_private (
--rollback 	id BIGSERIAL NOT NULL PRIMARY KEY,
--rollback 	email VARCHAR(50) NOT NULL,
--rollback 	password VARCHAR(255) NOT NULL,
--rollback 	login_type VARCHAR(255) NOT NULL, 
--rollback 	CONSTRAINT fk_email FOREIGN KEY (email) REFERENCES accounts(email),
--rollback 	CONSTRAINT login_type_check CHECK (login_type in ('NORMAL', 'OAUTH2_GOOGLE'))
--rollback );
--rollback 
--rollback UPDATE accounts_private SET email = accounts.email, password = accounts.password, login_type = accounts.login_type FROM accounts;
--rollback ALTER TABLE accounts DROP COLUMN password;
--rollback ALTER TABLE accounts DROP COLUMN login_type;
--rollback ALTER TABLE accounts DROP COLUMN create_at;
--rollback ALTER TABLE accounts DROP COLUMN update_at;
