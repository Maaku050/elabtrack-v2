-- Explicit foundation setup, executed as schema owner after migrate-up.
-- No runtime role membership, grant option, schema CREATE or tracking access.
BEGIN;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE public.users, public.refresh_tokens TO elabtrack_runtime;
REVOKE ALL ON TABLE public.schema_migrations FROM elabtrack_runtime, PUBLIC;
COMMIT;
-- Future approved migrations grant DML on named application tables in their own
-- transaction. Grant sequences/functions only if that table/workflow needs them.
-- Current UUID IDs require no sequences; citext uses its extension's public
-- type/function permissions. No custom SECURITY DEFINER function is introduced.
