CREATE FUNCTION borrowing_fine_amount(due_at timestamptz,until_at timestamptz) RETURNS bigint LANGUAGE sql IMMUTABLE AS $$ SELECT CASE WHEN due_at IS NULL OR until_at<=due_at THEN 0::bigint ELSE (ceil(extract(epoch FROM (until_at-due_at))/86400)*1000)::bigint END $$;
REVOKE ALL ON FUNCTION borrowing_fine_amount(timestamptz,timestamptz) FROM PUBLIC;
GRANT EXECUTE ON FUNCTION borrowing_fine_amount(timestamptz,timestamptz) TO elabtrack_runtime;
-- Additive Phase8 history; historical000001–000008 are unchanged.
ALTER TABLE borrowings ADD COLUMN completed_at timestamptz, ADD COLUMN fine_final_minor bigint;
DO $$ DECLARE v text; BEGIN
 FOR v IN SELECT conname FROM pg_constraint WHERE conrelid='borrowings'::regclass AND contype='c' AND pg_get_constraintdef(oid) ~ '(status|entry_path)' LOOP EXECUTE format('ALTER TABLE borrowings DROP CONSTRAINT %I',v); END LOOP;
END $$;
ALTER TABLE borrowings ADD CONSTRAINT borrowing_lifecycle CHECK(status IN ('PENDING','CHECKED_OUT','COMPLETED','DENIED','CANCELLED','EXPIRED'));
ALTER TABLE borrowings ADD CONSTRAINT borrowing_entry CHECK(entry_path IN ('REQUEST','DIRECT') AND ((entry_path='REQUEST' AND expires_at=created_at+interval '24 hours') OR (entry_path='DIRECT' AND expires_at IS NULL AND status IN ('CHECKED_OUT','COMPLETED'))));
ALTER TABLE borrowings ADD CONSTRAINT borrowing_issue_times CHECK((status IN ('CHECKED_OUT','COMPLETED') AND checked_out_at IS NOT NULL AND due_at IS NOT NULL AND due_at>checked_out_at AND terminal_at IS NULL) OR (status NOT IN ('CHECKED_OUT','COMPLETED') AND checked_out_at IS NULL AND due_at IS NULL));
ALTER TABLE borrowings ADD CONSTRAINT borrowing_terminal CHECK((status IN ('DENIED','CANCELLED','EXPIRED'))=(terminal_at IS NOT NULL));
ALTER TABLE borrowings ADD CONSTRAINT borrowing_denial CHECK((status='DENIED')=(length(btrim(denial_reason))>0));
ALTER TABLE borrowings ADD CONSTRAINT borrowing_completion CHECK((status='COMPLETED' AND completed_at IS NOT NULL AND completed_at>=checked_out_at AND fine_final_minor IS NOT NULL AND fine_final_minor>=0 AND fine_final_minor=borrowing_fine_amount(due_at,completed_at)) OR (status<>'COMPLETED' AND completed_at IS NULL AND fine_final_minor IS NULL));
ALTER TABLE borrowing_items ADD UNIQUE(id,borrowing_id,equipment_id);
CREATE TABLE return_lines (
 id uuid PRIMARY KEY,borrowing_id uuid NOT NULL,equipment_id uuid NOT NULL,item_id uuid NOT NULL,
 good bigint NOT NULL CHECK(good BETWEEN 0 AND 2147483647),damaged bigint NOT NULL CHECK(damaged BETWEEN 0 AND 2147483647),lost bigint NOT NULL CHECK(lost BETWEEN 0 AND 2147483647),
 actor_id uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,reason text NOT NULL CHECK(length(btrim(reason)) BETWEEN 1 AND 1000),occurred_at timestamptz NOT NULL,
 CHECK(good+damaged+lost BETWEEN 1 AND 2147483647),FOREIGN KEY(item_id,borrowing_id,equipment_id) REFERENCES borrowing_items(id,borrowing_id,equipment_id) ON DELETE RESTRICT,
 UNIQUE(id,equipment_id),UNIQUE(id,item_id,equipment_id),UNIQUE(id,borrowing_id,item_id,equipment_id)
);
CREATE TABLE replacement_obligations (
 id uuid PRIMARY KEY,return_line_id uuid NOT NULL,borrowing_id uuid NOT NULL,item_id uuid NOT NULL,equipment_id uuid NOT NULL,
 kind text NOT NULL CHECK(kind IN ('DAMAGED','LOST')),required bigint NOT NULL CHECK(required BETWEEN 1 AND 2147483647),
 FOREIGN KEY(return_line_id,borrowing_id,item_id,equipment_id) REFERENCES return_lines(id,borrowing_id,item_id,equipment_id) ON DELETE RESTRICT,
 UNIQUE(return_line_id,kind),UNIQUE(id,borrowing_id,equipment_id),UNIQUE(id,borrowing_id,item_id,equipment_id)
);
CREATE TABLE replacement_acceptances (
 id uuid PRIMARY KEY,obligation_id uuid NOT NULL,borrowing_id uuid NOT NULL,item_id uuid NOT NULL,equipment_id uuid NOT NULL,
 equivalent_confirmed boolean NOT NULL CHECK(equivalent_confirmed),equipment_name text NOT NULL CHECK(length(btrim(equipment_name)) BETWEEN 1 AND 200),
 quantity bigint NOT NULL CHECK(quantity BETWEEN 1 AND 2147483647),actor_id uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
 reason text NOT NULL CHECK(length(btrim(reason)) BETWEEN 1 AND 1000),occurred_at timestamptz NOT NULL,
 FOREIGN KEY(obligation_id,borrowing_id,item_id,equipment_id) REFERENCES replacement_obligations(id,borrowing_id,item_id,equipment_id) ON DELETE RESTRICT,UNIQUE(id,equipment_id),UNIQUE(id,item_id,equipment_id)
);
CREATE TABLE fine_clearances (
 id uuid PRIMARY KEY,borrowing_id uuid NOT NULL REFERENCES borrowings(id) ON DELETE RESTRICT,actor_id uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
 assessed_minor bigint NOT NULL CHECK(assessed_minor>0),cleared_minor bigint NOT NULL CHECK(cleared_minor>0 AND cleared_minor<=assessed_minor),
 method text NOT NULL CHECK(method IN ('PAID','WAIVED','OTHER_RESOLUTION')),note text NOT NULL CHECK(length(note)<=1000),occurred_at timestamptz NOT NULL
);
CREATE INDEX return_history ON return_lines(borrowing_id,occurred_at,id);
CREATE INDEX replacement_history ON replacement_obligations(borrowing_id,id);
CREATE INDEX replacement_equipment ON replacement_obligations(equipment_id,borrowing_id);
CREATE INDEX replacement_acceptance_history ON replacement_acceptances(obligation_id,occurred_at,id);
CREATE INDEX fine_clearance_history ON fine_clearances(borrowing_id,occurred_at,id);
CREATE INDEX borrowing_issued_due ON borrowings(due_at,id) WHERE status='CHECKED_OUT';
ALTER TABLE borrowing_events DROP CONSTRAINT borrowing_events_kind_check,DROP CONSTRAINT borrowing_events_borrowing_id_kind_key;
ALTER TABLE borrowing_events ADD CONSTRAINT borrowing_events_kind_check CHECK(kind IN ('SUBMITTED','CHECKED_OUT','DENIED','CANCELLED','EXPIRED','DIRECT','RETURN','REPLACEMENT','COMPLETED','FINE_CLEARED'));
CREATE UNIQUE INDEX borrowing_single_transition ON borrowing_events(borrowing_id,kind) WHERE kind NOT IN ('RETURN','REPLACEMENT','FINE_CLEARED');
ALTER TABLE inventory_movements ADD COLUMN return_line_id uuid,ADD COLUMN replacement_acceptance_id uuid;
ALTER TABLE inventory_movements ADD FOREIGN KEY(return_line_id,borrowing_item_id,equipment_id) REFERENCES return_lines(id,item_id,equipment_id) ON DELETE RESTRICT;
ALTER TABLE inventory_movements ADD FOREIGN KEY(replacement_acceptance_id,borrowing_item_id,equipment_id) REFERENCES replacement_acceptances(id,item_id,equipment_id) ON DELETE RESTRICT;
ALTER TABLE inventory_movements DROP CONSTRAINT inventory_movements_kind_check,DROP CONSTRAINT inventory_movement_source,DROP CONSTRAINT inventory_movement_vector;
ALTER TABLE inventory_movements ADD CONSTRAINT inventory_movements_kind_check CHECK(kind IN ('OPENING','ADD','REMOVE','RECONCILE','RESERVE','CHECKOUT','DIRECT','CANCELLED','DENIED','EXPIRED','RETURN','REPLACEMENT'));
ALTER TABLE inventory_movements ADD CONSTRAINT inventory_movement_source CHECK((kind IN ('RESERVE','CHECKOUT','DIRECT','CANCELLED','DENIED','EXPIRED','RETURN','REPLACEMENT'))=(borrowing_item_id IS NOT NULL) AND (kind='RETURN')=(return_line_id IS NOT NULL) AND (kind='REPLACEMENT')=(replacement_acceptance_id IS NOT NULL));
ALTER TABLE inventory_movements ADD CONSTRAINT inventory_movement_vector CHECK(
 (kind IN ('OPENING','ADD','REMOVE','RECONCILE') AND delta_reserved=0 AND delta_checked_out=0 AND delta_damaged_held=0)
 OR (kind='RESERVE' AND delta_available<0 AND delta_reserved=-delta_available AND delta_checked_out=0 AND delta_damaged_held=0 AND delta_total=0)
 OR (kind IN ('CANCELLED','DENIED','EXPIRED') AND delta_reserved<0 AND delta_available=-delta_reserved AND delta_checked_out=0 AND delta_damaged_held=0 AND delta_total=0)
 OR (kind='CHECKOUT' AND delta_reserved<0 AND delta_checked_out=-delta_reserved AND delta_available=0 AND delta_damaged_held=0 AND delta_total=0)
 OR (kind='DIRECT' AND delta_available<0 AND delta_checked_out=-delta_available AND delta_reserved=0 AND delta_damaged_held=0 AND delta_total=0)
 OR (kind='RETURN' AND delta_reserved=0 AND delta_available>=0 AND delta_checked_out<0 AND delta_damaged_held>=0 AND delta_total<=0 AND delta_available+delta_checked_out+delta_damaged_held=delta_total)
 OR (kind='REPLACEMENT' AND delta_reserved=0 AND delta_checked_out=0 AND delta_damaged_held=0 AND delta_available>0 AND delta_total=delta_available));
DROP INDEX inventory_borrowing_effect;
CREATE UNIQUE INDEX inventory_borrowing_effect ON inventory_movements(borrowing_item_id,kind) WHERE borrowing_item_id IS NOT NULL AND kind NOT IN ('RETURN','REPLACEMENT');
CREATE UNIQUE INDEX inventory_return_effect ON inventory_movements(return_line_id) WHERE return_line_id IS NOT NULL;
CREATE UNIQUE INDEX inventory_replacement_effect ON inventory_movements(replacement_acceptance_id) WHERE replacement_acceptance_id IS NOT NULL;
CREATE OR REPLACE FUNCTION protect_borrowing_state() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN
 IF TG_OP<>'UPDATE' THEN RAISE EXCEPTION 'Borrowing history cannot be removed'; END IF;
 IF OLD.status='PENDING' THEN
  IF NEW.status NOT IN ('CHECKED_OUT','DENIED','CANCELLED','EXPIRED') OR (to_jsonb(NEW)-ARRAY['status','checked_out_at','due_at','terminal_at','denial_reason']) IS DISTINCT FROM (to_jsonb(OLD)-ARRAY['status','checked_out_at','due_at','terminal_at','denial_reason']) THEN RAISE EXCEPTION 'Illegal borrowing transition'; END IF;
 ELSIF OLD.status='CHECKED_OUT' THEN
  IF NEW.status<>'COMPLETED' OR (to_jsonb(NEW)-ARRAY['status','completed_at','fine_final_minor']) IS DISTINCT FROM (to_jsonb(OLD)-ARRAY['status','completed_at','fine_final_minor']) THEN RAISE EXCEPTION 'Illegal completion'; END IF;
  IF EXISTS(SELECT 1 FROM borrowing_items i WHERE i.borrowing_id=OLD.id AND i.issued_quantity<>(SELECT COALESCE(sum(good+damaged+lost),0) FROM return_lines l WHERE l.item_id=i.id)) OR EXISTS(SELECT 1 FROM replacement_obligations o WHERE o.borrowing_id=OLD.id AND o.required<>(SELECT COALESCE(sum(quantity),0) FROM replacement_acceptances a WHERE a.obligation_id=o.id)) THEN RAISE EXCEPTION 'Unresolved borrowing'; END IF;
 ELSE RAISE EXCEPTION 'Immutable terminal borrowing'; END IF;
 RETURN NEW;
END $$;
DO $$ DECLARE t text;BEGIN FOREACH t IN ARRAY ARRAY['return_lines','replacement_obligations','replacement_acceptances','fine_clearances'] LOOP
 EXECUTE format('CREATE TRIGGER accountability_immutable BEFORE UPDATE OR DELETE OR TRUNCATE ON %I FOR EACH STATEMENT EXECUTE FUNCTION reject_inventory_history_change()',t);
 EXECUTE format('GRANT SELECT,INSERT ON %I TO elabtrack_runtime',t);
END LOOP;END $$;
GRANT UPDATE(completed_at,fine_final_minor) ON borrowings TO elabtrack_runtime;
GRANT UPDATE(damaged_held,total_tracked) ON equipment TO elabtrack_runtime;
-- The aggregate lock independently protects cross-row bounds for SQL writers.
CREATE FUNCTION enforce_accountability_insert() RETURNS trigger LANGUAGE plpgsql AS $$ DECLARE b borrowings;bound bigint;used bigint;l return_lines;BEGIN
 SELECT * INTO STRICT b FROM borrowings WHERE id=NEW.borrowing_id FOR UPDATE;
 IF TG_TABLE_NAME='return_lines' THEN
  IF b.status<>'CHECKED_OUT' THEN RAISE EXCEPTION 'Return requires issued loan';END IF;
  SELECT issued_quantity INTO STRICT bound FROM borrowing_items WHERE id=NEW.item_id;
  SELECT COALESCE(sum(good+damaged+lost),0) INTO used FROM return_lines WHERE item_id=NEW.item_id;
  IF used+NEW.good+NEW.damaged+NEW.lost>bound THEN RAISE EXCEPTION 'Return exceeds custody';END IF;
 ELSIF TG_TABLE_NAME='replacement_obligations' THEN
  SELECT * INTO STRICT l FROM return_lines WHERE id=NEW.return_line_id;
  IF NEW.required<>(CASE WHEN NEW.kind='DAMAGED' THEN l.damaged ELSE l.lost END) THEN RAISE EXCEPTION 'Obligation must equal incident';END IF;
 ELSIF TG_TABLE_NAME='replacement_acceptances' THEN
  IF b.status<>'CHECKED_OUT' THEN RAISE EXCEPTION 'Replacement requires open loan';END IF;
  SELECT required INTO STRICT bound FROM replacement_obligations WHERE id=NEW.obligation_id;
  SELECT COALESCE(sum(quantity),0) INTO used FROM replacement_acceptances WHERE obligation_id=NEW.obligation_id;
  IF used+NEW.quantity>bound THEN RAISE EXCEPTION 'Replacement exceeds obligation';END IF;
 ELSE
  IF b.status NOT IN ('CHECKED_OUT','COMPLETED') THEN RAISE EXCEPTION 'Fine requires issued loan';END IF;
  bound:=COALESCE(b.fine_final_minor,borrowing_fine_amount(b.due_at,NEW.occurred_at));
  SELECT COALESCE(sum(cleared_minor),0) INTO used FROM fine_clearances WHERE borrowing_id=NEW.borrowing_id;
  IF NEW.assessed_minor<>bound OR NEW.cleared_minor<>bound-used OR NOT EXISTS(SELECT 1 FROM users WHERE id=NEW.actor_id AND role='ADMIN' AND is_active AND NOT activation_required) THEN RAISE EXCEPTION 'Invalid Admin full clearance';END IF;
 END IF;
 IF TG_TABLE_NAME IN ('return_lines','replacement_acceptances') THEN
  IF NOT EXISTS(SELECT 1 FROM users WHERE id=NEW.actor_id AND role IN ('STAFF','ADMIN') AND is_active AND NOT activation_required) THEN RAISE EXCEPTION 'Invalid operational actor';END IF;
 END IF;
 RETURN NEW;
END $$;
DO $$ DECLARE t text;BEGIN FOREACH t IN ARRAY ARRAY['return_lines','replacement_obligations','replacement_acceptances','fine_clearances'] LOOP
 EXECUTE format('CREATE TRIGGER accountability_bounds BEFORE INSERT ON %I FOR EACH ROW EXECUTE FUNCTION enforce_accountability_insert()',t);
END LOOP;END $$;
REVOKE ALL ON FUNCTION enforce_accountability_insert() FROM PUBLIC;
