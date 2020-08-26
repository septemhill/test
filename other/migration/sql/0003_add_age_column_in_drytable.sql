--liquibase formatted sql

--comment: Add one column for table `drytable` 
--changeset Septem:0003_add_age_column_for_drytable
ALTER TABLE drytable ADD COLUMN age INTEGER;

--rollback ALTER TABLE drytable DROP COLUMN IF EXISTS age; 
