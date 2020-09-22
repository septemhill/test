--liquibase formatted sql

--comment: Create non-verified accounts table
--changeset Septem:0003_add_active_column_in_accounts

ALTER TABLE accounts ADD COLUMN active BOOLEAN;
UPDATE accounts SET active = true;
ALTER TABLE accounts ALTER COLUMN active SET NOT NULL;

--rollback ALTER TABLE accounts DROP active;