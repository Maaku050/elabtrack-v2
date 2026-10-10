-- Optional current account photo; immutable account audit/receipts retain actions.
-- No return-evidence storage, account provisioning or role changes.
CREATE TABLE profile_images (
 account_id uuid PRIMARY KEY REFERENCES users(id) ON DELETE RESTRICT,
 image_id uuid UNIQUE,version bigint NOT NULL CHECK(version>0),
 png bytea,sha256 text,width integer,height integer,
 updated_at timestamptz NOT NULL DEFAULT clock_timestamp(),
 CHECK ((image_id IS NULL AND png IS NULL AND sha256 IS NULL AND width IS NULL AND height IS NULL) OR
 (image_id IS NOT NULL AND png IS NOT NULL AND sha256 IS NOT NULL AND width IS NOT NULL AND height IS NOT NULL AND octet_length(png) BETWEEN 1 AND 524288 AND sha256~'^[a-f0-9]{64}$' AND width BETWEEN 1 AND 2048 AND height BETWEEN 1 AND 2048 AND width::bigint*height<=4194304))
);
REVOKE ALL ON profile_images FROM PUBLIC;
GRANT SELECT,INSERT ON profile_images TO elabtrack_runtime;
GRANT UPDATE(image_id,version,png,sha256,width,height,updated_at) ON profile_images TO elabtrack_runtime;
