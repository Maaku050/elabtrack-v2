CREATE TABLE borrowings (
 id uuid PRIMARY KEY,reference text NOT NULL UNIQUE,borrower_id uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
 borrower_name text NOT NULL,borrower_type text NOT NULL CHECK(borrower_type IN ('STUDENT','FACULTY')),student_id text NOT NULL DEFAULT '',
 acceptance_id uuid NOT NULL,status text NOT NULL CHECK(status IN ('PENDING','CHECKED_OUT','DENIED','CANCELLED','EXPIRED')),
 entry_path text NOT NULL CHECK(entry_path IN ('REQUEST','DIRECT')),
 created_at timestamptz NOT NULL,expires_at timestamptz,checked_out_at timestamptz,due_at timestamptz,terminal_at timestamptz,
 denial_reason text NOT NULL DEFAULT '' CHECK(length(denial_reason)<=1000),
 fine_currency text NOT NULL DEFAULT 'PHP' CHECK(fine_currency='PHP'),fine_daily_minor integer NOT NULL DEFAULT 1000 CHECK(fine_daily_minor=1000),fine_day_seconds integer NOT NULL DEFAULT 86400 CHECK(fine_day_seconds=86400),fine_rounding text NOT NULL DEFAULT 'CEILING_ELAPSED_24H' CHECK(fine_rounding='CEILING_ELAPSED_24H'),
 FOREIGN KEY(acceptance_id,borrower_id) REFERENCES terms_acceptances(id,user_id) ON DELETE RESTRICT,
 CHECK((entry_path='REQUEST' AND expires_at=created_at+interval '24 hours') OR (entry_path='DIRECT' AND expires_at IS NULL AND status='CHECKED_OUT')),
 CHECK((status='CHECKED_OUT' AND checked_out_at IS NOT NULL AND due_at IS NOT NULL AND due_at>checked_out_at AND terminal_at IS NULL) OR (status<>'CHECKED_OUT' AND checked_out_at IS NULL AND due_at IS NULL)),
 CHECK((status IN ('DENIED','CANCELLED','EXPIRED'))=(terminal_at IS NOT NULL)),
 CHECK((status='DENIED')=(length(btrim(denial_reason))>0)),UNIQUE(id,borrower_id)
);
CREATE TABLE borrowing_items (
 id uuid PRIMARY KEY,borrowing_id uuid NOT NULL REFERENCES borrowings(id) ON DELETE RESTRICT,
 equipment_id uuid NOT NULL REFERENCES equipment(id) ON DELETE RESTRICT,name text NOT NULL,
 quantity bigint NOT NULL CHECK(quantity BETWEEN 1 AND 2147483647),
 reserved_quantity bigint NOT NULL CHECK(reserved_quantity>=0),issued_quantity bigint NOT NULL CHECK(issued_quantity>=0),
 CHECK((reserved_quantity=quantity AND issued_quantity=0) OR (reserved_quantity=0 AND issued_quantity IN (0,quantity))),
 UNIQUE(borrowing_id,equipment_id),UNIQUE(id,equipment_id)
);
CREATE TABLE borrowing_events (
 id uuid PRIMARY KEY,borrowing_id uuid NOT NULL REFERENCES borrowings(id) ON DELETE RESTRICT,
 actor_id uuid REFERENCES users(id) ON DELETE RESTRICT,kind text NOT NULL CHECK(kind IN ('SUBMITTED','CHECKED_OUT','DENIED','CANCELLED','EXPIRED','DIRECT')),
 reason text NOT NULL DEFAULT '' CHECK(length(reason)<=1000),occurred_at timestamptz NOT NULL,
 CHECK((kind='EXPIRED')=(actor_id IS NULL)),UNIQUE(borrowing_id,kind)
);
CREATE TABLE borrowing_operation_receipts (
 actor_id uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,operation text NOT NULL,key uuid NOT NULL,payload_hash text NOT NULL CHECK(payload_hash ~ '^[a-f0-9]{64}$'),result jsonb NOT NULL,created_at timestamptz NOT NULL DEFAULT now(),PRIMARY KEY(actor_id,operation,key)
);
CREATE INDEX borrowing_owner_history ON borrowings(borrower_id,created_at DESC,id DESC);
CREATE INDEX borrowing_status_history ON borrowings(status,created_at DESC,id DESC);
CREATE INDEX borrowing_pending_expiry ON borrowings(expires_at,id) WHERE status='PENDING';
CREATE INDEX borrowing_equipment ON borrowing_items(equipment_id,borrowing_id);
CREATE INDEX borrowing_event_history ON borrowing_events(borrowing_id,occurred_at,id);
ALTER TABLE inventory_movements DROP CONSTRAINT inventory_movements_kind_check;
-- Remove only the Phase6 no-custody vector check; preserve conservation checks.
DO $$ DECLARE v text; BEGIN
 FOR v IN SELECT conname FROM pg_constraint WHERE conrelid='inventory_movements'::regclass AND contype='c' AND pg_get_constraintdef(oid) LIKE '%delta_reserved = 0%' LOOP EXECUTE format('ALTER TABLE inventory_movements DROP CONSTRAINT %I',v); END LOOP;
END $$;
ALTER TABLE inventory_movements ALTER COLUMN actor_id DROP NOT NULL;
ALTER TABLE inventory_movements ADD COLUMN borrowing_item_id uuid;
ALTER TABLE inventory_movements ADD FOREIGN KEY(borrowing_item_id,equipment_id) REFERENCES borrowing_items(id,equipment_id) ON DELETE RESTRICT;
ALTER TABLE inventory_movements ADD CONSTRAINT inventory_movements_kind_check CHECK(kind IN ('OPENING','ADD','REMOVE','RECONCILE','RESERVE','CHECKOUT','DIRECT','CANCELLED','DENIED','EXPIRED'));
ALTER TABLE inventory_movements ADD CONSTRAINT inventory_movement_source CHECK((kind IN ('RESERVE','CHECKOUT','DIRECT','CANCELLED','DENIED','EXPIRED'))=(borrowing_item_id IS NOT NULL));
ALTER TABLE inventory_movements ADD CONSTRAINT inventory_movement_actor CHECK((kind='EXPIRED')=(actor_id IS NULL));
ALTER TABLE inventory_movements ADD CONSTRAINT inventory_movement_vector CHECK(
 (kind IN ('OPENING','ADD','REMOVE','RECONCILE') AND delta_reserved=0 AND delta_checked_out=0 AND delta_damaged_held=0)
 OR (kind='RESERVE' AND delta_available<0 AND delta_reserved=-delta_available AND delta_checked_out=0 AND delta_damaged_held=0 AND delta_total=0)
 OR (kind IN ('CANCELLED','DENIED','EXPIRED') AND delta_reserved<0 AND delta_available=-delta_reserved AND delta_checked_out=0 AND delta_damaged_held=0 AND delta_total=0)
 OR (kind='CHECKOUT' AND delta_reserved<0 AND delta_checked_out=-delta_reserved AND delta_available=0 AND delta_damaged_held=0 AND delta_total=0)
 OR (kind='DIRECT' AND delta_available<0 AND delta_checked_out=-delta_available AND delta_reserved=0 AND delta_damaged_held=0 AND delta_total=0));
CREATE UNIQUE INDEX inventory_borrowing_effect ON inventory_movements(borrowing_item_id,kind) WHERE borrowing_item_id IS NOT NULL;
CREATE FUNCTION protect_borrowing_state() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN
 IF TG_OP<>'UPDATE' THEN RAISE EXCEPTION 'Borrowing history cannot be removed'; END IF;
 IF OLD.status<>'PENDING' OR NEW.status NOT IN ('CHECKED_OUT','DENIED','CANCELLED','EXPIRED') OR (to_jsonb(NEW)-ARRAY['status','checked_out_at','due_at','terminal_at','denial_reason']) IS DISTINCT FROM (to_jsonb(OLD)-ARRAY['status','checked_out_at','due_at','terminal_at','denial_reason']) THEN RAISE EXCEPTION 'Illegal borrowing transition'; END IF;
 RETURN NEW;
END $$;
CREATE TRIGGER borrowing_state_guard BEFORE UPDATE OR DELETE ON borrowings FOR EACH ROW EXECUTE FUNCTION protect_borrowing_state();
CREATE TRIGGER borrowing_no_truncate BEFORE TRUNCATE ON borrowings FOR EACH STATEMENT EXECUTE FUNCTION reject_inventory_history_change();
CREATE FUNCTION protect_borrowing_item() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN
 IF TG_OP<>'UPDATE' OR OLD.reserved_quantity=0 OR NEW.reserved_quantity<>0 OR (to_jsonb(NEW)-ARRAY['reserved_quantity','issued_quantity']) IS DISTINCT FROM (to_jsonb(OLD)-ARRAY['reserved_quantity','issued_quantity']) THEN RAISE EXCEPTION 'Immutable borrowing item'; END IF; RETURN NEW;
END $$;
CREATE TRIGGER borrowing_item_guard BEFORE UPDATE OR DELETE ON borrowing_items FOR EACH ROW EXECUTE FUNCTION protect_borrowing_item();
CREATE TRIGGER borrowing_items_no_truncate BEFORE TRUNCATE ON borrowing_items FOR EACH STATEMENT EXECUTE FUNCTION reject_inventory_history_change();
CREATE TRIGGER borrowing_events_immutable BEFORE UPDATE OR DELETE OR TRUNCATE ON borrowing_events FOR EACH STATEMENT EXECUTE FUNCTION reject_inventory_history_change();
CREATE TRIGGER borrowing_receipts_immutable BEFORE UPDATE OR DELETE OR TRUNCATE ON borrowing_operation_receipts FOR EACH STATEMENT EXECUTE FUNCTION reject_inventory_history_change();
REVOKE ALL ON FUNCTION protect_borrowing_state(),protect_borrowing_item() FROM PUBLIC;
GRANT SELECT,INSERT ON borrowings,borrowing_items,borrowing_events,borrowing_operation_receipts TO elabtrack_runtime;
GRANT UPDATE(status,checked_out_at,due_at,terminal_at,denial_reason) ON borrowings TO elabtrack_runtime;
GRANT UPDATE(reserved_quantity,issued_quantity) ON borrowing_items TO elabtrack_runtime;
GRANT UPDATE(reserved,checked_out) ON equipment TO elabtrack_runtime;
