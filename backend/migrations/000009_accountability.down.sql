DO $$ BEGIN IF EXISTS(SELECT 1 FROM return_lines) OR EXISTS(SELECT 1 FROM replacement_obligations) OR EXISTS(SELECT 1 FROM replacement_acceptances) OR EXISTS(SELECT 1 FROM fine_clearances) OR EXISTS(SELECT 1 FROM borrowings WHERE status='COMPLETED') THEN RAISE EXCEPTION 'Refusing rollback with accountability history'; END IF; END $$;
DROP INDEX inventory_return_effect,inventory_replacement_effect,inventory_borrowing_effect;
ALTER TABLE inventory_movements DROP CONSTRAINT inventory_movements_kind_check,DROP CONSTRAINT inventory_movement_source,DROP CONSTRAINT inventory_movement_vector;
ALTER TABLE inventory_movements DROP COLUMN return_line_id,DROP COLUMN replacement_acceptance_id;
ALTER TABLE inventory_movements ADD CONSTRAINT inventory_movements_kind_check CHECK(kind IN ('OPENING','ADD','REMOVE','RECONCILE','RESERVE','CHECKOUT','DIRECT','CANCELLED','DENIED','EXPIRED'));
ALTER TABLE inventory_movements ADD CONSTRAINT inventory_movement_source CHECK((kind IN ('RESERVE','CHECKOUT','DIRECT','CANCELLED','DENIED','EXPIRED'))=(borrowing_item_id IS NOT NULL));
ALTER TABLE inventory_movements ADD CONSTRAINT inventory_movement_vector CHECK(
 (kind IN ('OPENING','ADD','REMOVE','RECONCILE') AND delta_reserved=0 AND delta_checked_out=0 AND delta_damaged_held=0)
 OR (kind='RESERVE' AND delta_available<0 AND delta_reserved=-delta_available AND delta_checked_out=0 AND delta_damaged_held=0 AND delta_total=0)
 OR (kind IN ('CANCELLED','DENIED','EXPIRED') AND delta_reserved<0 AND delta_available=-delta_reserved AND delta_checked_out=0 AND delta_damaged_held=0 AND delta_total=0)
 OR (kind='CHECKOUT' AND delta_reserved<0 AND delta_checked_out=-delta_reserved AND delta_available=0 AND delta_damaged_held=0 AND delta_total=0)
 OR (kind='DIRECT' AND delta_available<0 AND delta_checked_out=-delta_available AND delta_reserved=0 AND delta_damaged_held=0 AND delta_total=0));
CREATE UNIQUE INDEX inventory_borrowing_effect ON inventory_movements(borrowing_item_id,kind) WHERE borrowing_item_id IS NOT NULL;
CREATE OR REPLACE FUNCTION protect_borrowing_state() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN
 IF TG_OP<>'UPDATE' THEN RAISE EXCEPTION 'Borrowing history cannot be removed'; END IF;
 IF OLD.status<>'PENDING' OR NEW.status NOT IN ('CHECKED_OUT','DENIED','CANCELLED','EXPIRED') OR (to_jsonb(NEW)-ARRAY['status','checked_out_at','due_at','terminal_at','denial_reason']) IS DISTINCT FROM (to_jsonb(OLD)-ARRAY['status','checked_out_at','due_at','terminal_at','denial_reason']) THEN RAISE EXCEPTION 'Illegal borrowing transition'; END IF;
 RETURN NEW;
END $$;
DROP INDEX borrowing_single_transition;
ALTER TABLE borrowing_events DROP CONSTRAINT borrowing_events_kind_check;
ALTER TABLE borrowing_events ADD CONSTRAINT borrowing_events_kind_check CHECK(kind IN ('SUBMITTED','CHECKED_OUT','DENIED','CANCELLED','EXPIRED','DIRECT')),ADD UNIQUE(borrowing_id,kind);
DROP TABLE fine_clearances,replacement_acceptances,replacement_obligations,return_lines;
DROP FUNCTION enforce_accountability_insert();
REVOKE UPDATE(damaged_held) ON equipment FROM elabtrack_runtime;
REVOKE UPDATE(completed_at,fine_final_minor) ON borrowings FROM elabtrack_runtime;
ALTER TABLE borrowings DROP CONSTRAINT borrowing_lifecycle,DROP CONSTRAINT borrowing_entry,DROP CONSTRAINT borrowing_issue_times,DROP CONSTRAINT borrowing_terminal,DROP CONSTRAINT borrowing_denial,DROP CONSTRAINT borrowing_completion;
ALTER TABLE borrowings DROP COLUMN completed_at,DROP COLUMN fine_final_minor;
ALTER TABLE borrowings ADD CHECK(status IN ('PENDING','CHECKED_OUT','DENIED','CANCELLED','EXPIRED')),ADD CHECK(entry_path IN ('REQUEST','DIRECT')),
ADD CHECK((entry_path='REQUEST' AND expires_at=created_at+interval '24 hours') OR (entry_path='DIRECT' AND expires_at IS NULL AND status='CHECKED_OUT')),
ADD CHECK((status='CHECKED_OUT' AND checked_out_at IS NOT NULL AND due_at IS NOT NULL AND due_at>checked_out_at AND terminal_at IS NULL) OR (status<>'CHECKED_OUT' AND checked_out_at IS NULL AND due_at IS NULL)),
ADD CHECK((status IN ('DENIED','CANCELLED','EXPIRED'))=(terminal_at IS NOT NULL)),ADD CHECK((status='DENIED')=(length(btrim(denial_reason))>0));
DROP INDEX borrowing_issued_due;
ALTER TABLE borrowing_items DROP CONSTRAINT borrowing_items_id_borrowing_id_equipment_id_key;

DROP FUNCTION borrowing_fine_amount(timestamptz,timestamptz);
