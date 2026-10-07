#!/bin/sh
# Fresh local Compose volumes only. Existing volumes are never auto-adopted.
set -eu
: "${MIGRATION_DB_PASSWORD:?}" "${RUNTIME_DB_PASSWORD:?}"
psql --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" --no-psqlrc --quiet --set ON_ERROR_STOP=1 <<'SQL'
\getenv db_name POSTGRES_DB
\getenv migration_password MIGRATION_DB_PASSWORD
\getenv runtime_password RUNTIME_DB_PASSWORD
CREATE ROLE elabtrack_migrator LOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION NOBYPASSRLS PASSWORD :'migration_password';
CREATE ROLE elabtrack_runtime LOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION NOBYPASSRLS PASSWORD :'runtime_password';
REVOKE ALL ON DATABASE :"db_name" FROM PUBLIC;
GRANT CONNECT ON DATABASE :"db_name" TO elabtrack_migrator, elabtrack_runtime;
ALTER SCHEMA public OWNER TO elabtrack_migrator;
REVOKE ALL ON SCHEMA public FROM PUBLIC;
GRANT USAGE ON SCHEMA public TO elabtrack_runtime;
-- Install the trusted extension with bootstrap credentials; migrator needs no
-- CREATE on the database or cluster privilege for the historical IF NOT EXISTS.
CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public;
ALTER DEFAULT PRIVILEGES FOR ROLE elabtrack_migrator IN SCHEMA public REVOKE ALL ON TABLES FROM PUBLIC;
-- No blanket runtime/default table grants: tracking stays private. Explicit
-- table grants belong to each future approved application migration.
SQL
