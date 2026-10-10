DO $$ BEGIN IF EXISTS(SELECT 1 FROM notification_dispatch) OR EXISTS(SELECT 1 FROM notifications) THEN RAISE EXCEPTION 'Refusing rollback with notification history';END IF;END $$;
DROP TABLE notifications,notification_dispatch;
DROP FUNCTION protect_notification();
