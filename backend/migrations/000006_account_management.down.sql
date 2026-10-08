DO $$ BEGIN
 IF EXISTS(SELECT 1 FROM borrower_profiles) OR EXISTS(SELECT 1 FROM account_activation_tokens) OR EXISTS(SELECT 1 FROM account_audit_events) OR EXISTS(SELECT 1 FROM account_operation_receipts) OR EXISTS(SELECT 1 FROM account_roster_batches) OR EXISTS(SELECT 1 FROM users WHERE activation_required OR activation_delivery_status<>'NOT_REQUIRED') THEN RAISE EXCEPTION 'Account management history exists; destructive rollback denied'; END IF;
END $$;
DROP TRIGGER classified_borrower_role ON users;
DROP TABLE account_roster_batches,account_audit_events,account_operation_receipts,account_activation_tokens,borrower_profiles;
DROP FUNCTION guard_borrower_profile_role(),reject_account_history_change(),guard_account_batch_confirmation();
ALTER TABLE users DROP COLUMN activation_required,DROP COLUMN activation_delivery_status;
