CREATE TABLE notification_dispatch (
 event_key text PRIMARY KEY CHECK(length(event_key) BETWEEN 1 AND 160),
 borrowing_id uuid NOT NULL REFERENCES borrowings(id) ON DELETE RESTRICT,
 kind text NOT NULL CHECK(kind IN ('SUBMITTED','CHECKED_OUT','DIRECT','DENIED','CANCELLED','EXPIRED','RETURN','DAMAGED','LOST','REPLACEMENT_REQUIRED','REPLACEMENT','COMPLETED','FINE_CLEARED','DUE_SOON','OVERDUE','FINE_ASSESSED')),
 created_at timestamptz NOT NULL DEFAULT clock_timestamp(), UNIQUE(event_key,borrowing_id)
);
CREATE TABLE notifications (
 id uuid PRIMARY KEY,recipient_id uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
 event_key text NOT NULL,borrowing_id uuid NOT NULL,
 scope text NOT NULL CHECK(scope IN ('BORROWER','OPERATIONS')),
 kind text NOT NULL,title text NOT NULL CHECK(length(title) BETWEEN 1 AND 160),body text NOT NULL CHECK(length(body) BETWEEN 1 AND 1000),
 created_at timestamptz NOT NULL DEFAULT clock_timestamp(),read_at timestamptz,
 FOREIGN KEY(event_key,borrowing_id) REFERENCES notification_dispatch(event_key,borrowing_id) ON DELETE RESTRICT,
 UNIQUE(event_key,recipient_id),CHECK(read_at IS NULL OR read_at>=created_at)
);
CREATE INDEX notification_recipient_page ON notifications(recipient_id,created_at DESC,id DESC);
CREATE INDEX notification_unread ON notifications(recipient_id,scope) WHERE read_at IS NULL;
CREATE INDEX notification_borrowing ON notifications(borrowing_id,event_key);
CREATE TRIGGER notification_dispatch_immutable BEFORE UPDATE OR DELETE OR TRUNCATE ON notification_dispatch FOR EACH STATEMENT EXECUTE FUNCTION reject_inventory_history_change();
CREATE FUNCTION protect_notification() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN
 IF TG_OP<>'UPDATE' OR (to_jsonb(NEW)-'read_at') IS DISTINCT FROM (to_jsonb(OLD)-'read_at') THEN RAISE EXCEPTION 'Notification evidence is immutable';END IF;RETURN NEW;
END $$;
CREATE TRIGGER notification_history_guard BEFORE UPDATE OR DELETE ON notifications FOR EACH ROW EXECUTE FUNCTION protect_notification();
CREATE TRIGGER notification_no_truncate BEFORE TRUNCATE ON notifications FOR EACH STATEMENT EXECUTE FUNCTION reject_inventory_history_change();
REVOKE ALL ON FUNCTION protect_notification() FROM PUBLIC;
GRANT SELECT,INSERT ON notification_dispatch,notifications TO elabtrack_runtime;
GRANT UPDATE(read_at) ON notifications TO elabtrack_runtime;
