--liquibase formatted sql

--comment: Update table `drytable` column `age` to NOT NULL
--changeset Septem:9995_update_age_column_to_not_null
UPDATE drytable SET age = 10 WHERE age IS NULL;
ALTER TABLE drytable ALTER COLUMN age SET NOT NULL;

--rollback ALTER TABLE drytable ALTER COLUMN age DROP NOT NULL;
