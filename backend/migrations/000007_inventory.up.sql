CREATE TABLE equipment_categories (
 id uuid PRIMARY KEY, name text NOT NULL CHECK(length(btrim(name)) BETWEEN 1 AND 100),
 is_active boolean NOT NULL DEFAULT true, version bigint NOT NULL DEFAULT 1 CHECK(version>0),
 created_at timestamptz NOT NULL DEFAULT now(), updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX equipment_categories_name_key ON equipment_categories(lower(name));
CREATE TABLE equipment (
 id uuid PRIMARY KEY,name text NOT NULL CHECK(length(btrim(name)) BETWEEN 1 AND 160),
 description text NOT NULL DEFAULT '' CHECK(length(description)<=2000), category_id uuid REFERENCES equipment_categories(id) ON DELETE RESTRICT,
 status text NOT NULL CHECK(status IN ('ACTIVE','INACTIVE','ARCHIVED')),
 available bigint NOT NULL DEFAULT 0 CHECK(available BETWEEN 0 AND 2147483647),
 reserved bigint NOT NULL DEFAULT 0 CHECK(reserved BETWEEN 0 AND 2147483647),
 checked_out bigint NOT NULL DEFAULT 0 CHECK(checked_out BETWEEN 0 AND 2147483647),
 damaged_held bigint NOT NULL DEFAULT 0 CHECK(damaged_held BETWEEN 0 AND 2147483647),
 total_tracked bigint NOT NULL DEFAULT 0 CHECK(total_tracked BETWEEN 0 AND 2147483647),
 stock_sequence bigint NOT NULL DEFAULT 0 CHECK(stock_sequence>=0),metadata_version bigint NOT NULL DEFAULT 1 CHECK(metadata_version>0),
 image_id uuid,created_at timestamptz NOT NULL DEFAULT now(),updated_at timestamptz NOT NULL DEFAULT now(),
 CHECK(total_tracked=available+reserved+checked_out+damaged_held), CHECK(status<>'ARCHIVED' OR (reserved=0 AND checked_out=0))
);
CREATE TABLE equipment_images (
 id uuid PRIMARY KEY,equipment_id uuid NOT NULL REFERENCES equipment(id) ON DELETE RESTRICT,
 png bytea NOT NULL CHECK(octet_length(png) BETWEEN 1 AND 524288),sha256 text NOT NULL CHECK(sha256 ~ '^[a-f0-9]{64}$'),
 width integer NOT NULL CHECK(width BETWEEN 1 AND 2048),height integer NOT NULL CHECK(height BETWEEN 1 AND 2048),
 actor_id uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,created_at timestamptz NOT NULL DEFAULT now(),
 UNIQUE(id,equipment_id),CHECK(width::bigint*height<=4194304)
);
ALTER TABLE equipment ADD CONSTRAINT equipment_current_image FOREIGN KEY(image_id,id) REFERENCES equipment_images(id,equipment_id) ON DELETE RESTRICT;
CREATE TABLE inventory_movements (
 id uuid PRIMARY KEY,equipment_id uuid NOT NULL REFERENCES equipment(id) ON DELETE RESTRICT,
 sequence bigint NOT NULL CHECK(sequence>0),actor_id uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
 kind text NOT NULL CHECK(kind IN ('OPENING','ADD','REMOVE','RECONCILE')),
 delta_available bigint NOT NULL,delta_reserved bigint NOT NULL DEFAULT 0,delta_checked_out bigint NOT NULL DEFAULT 0,delta_damaged_held bigint NOT NULL DEFAULT 0,delta_total bigint NOT NULL,
 after_available bigint NOT NULL CHECK(after_available BETWEEN 0 AND 2147483647),after_reserved bigint NOT NULL CHECK(after_reserved BETWEEN 0 AND 2147483647),after_checked_out bigint NOT NULL CHECK(after_checked_out BETWEEN 0 AND 2147483647),after_damaged_held bigint NOT NULL CHECK(after_damaged_held BETWEEN 0 AND 2147483647),after_total bigint NOT NULL CHECK(after_total BETWEEN 0 AND 2147483647),
 reason text NOT NULL CHECK(length(btrim(reason)) BETWEEN 1 AND 1000),created_at timestamptz NOT NULL DEFAULT now(),
 UNIQUE(equipment_id,sequence),CHECK(delta_total=delta_available+delta_reserved+delta_checked_out+delta_damaged_held),CHECK(after_total=after_available+after_reserved+after_checked_out+after_damaged_held),CHECK(delta_reserved=0 AND delta_checked_out=0 AND delta_damaged_held=0)
);
CREATE TABLE inventory_audit_events (
 id uuid PRIMARY KEY,actor_id uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
 equipment_id uuid REFERENCES equipment(id) ON DELETE RESTRICT,category_id uuid REFERENCES equipment_categories(id) ON DELETE RESTRICT,
 action text NOT NULL CHECK(length(action) BETWEEN 1 AND 80),before_data jsonb,after_data jsonb NOT NULL,created_at timestamptz NOT NULL DEFAULT now(),
 CHECK(num_nonnulls(equipment_id,category_id)=1)
);
CREATE TABLE inventory_operation_receipts (
 actor_id uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,operation text NOT NULL,key uuid NOT NULL,
 payload_hash text NOT NULL,result jsonb NOT NULL,created_at timestamptz NOT NULL DEFAULT now(),PRIMARY KEY(actor_id,operation,key)
);
CREATE INDEX equipment_catalog ON equipment(status,lower(name),id);
CREATE INDEX equipment_category ON equipment(category_id,status);
CREATE INDEX inventory_audit_equipment ON inventory_audit_events(equipment_id,created_at,id);
CREATE FUNCTION reject_inventory_history_change() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'Inventory history is immutable'; END $$;
CREATE TRIGGER inventory_movements_immutable BEFORE UPDATE OR DELETE OR TRUNCATE ON inventory_movements FOR EACH STATEMENT EXECUTE FUNCTION reject_inventory_history_change();
CREATE TRIGGER inventory_audit_immutable BEFORE UPDATE OR DELETE OR TRUNCATE ON inventory_audit_events FOR EACH STATEMENT EXECUTE FUNCTION reject_inventory_history_change();
CREATE TRIGGER inventory_receipts_immutable BEFORE UPDATE OR DELETE OR TRUNCATE ON inventory_operation_receipts FOR EACH STATEMENT EXECUTE FUNCTION reject_inventory_history_change();
CREATE TRIGGER equipment_images_immutable BEFORE UPDATE OR DELETE OR TRUNCATE ON equipment_images FOR EACH STATEMENT EXECUTE FUNCTION reject_inventory_history_change();
REVOKE ALL ON FUNCTION reject_inventory_history_change() FROM PUBLIC;
GRANT SELECT,INSERT ON equipment_categories,equipment,equipment_images,inventory_movements,inventory_audit_events,inventory_operation_receipts TO elabtrack_runtime;
GRANT UPDATE(name,is_active,version,updated_at) ON equipment_categories TO elabtrack_runtime;
GRANT UPDATE(name,description,category_id,status,available,total_tracked,stock_sequence,metadata_version,image_id,updated_at) ON equipment TO elabtrack_runtime;
