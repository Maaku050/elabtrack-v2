DO $$ BEGIN IF EXISTS(SELECT 1 FROM borrowings) OR EXISTS(SELECT 1 FROM borrowing_operation_receipts) OR EXISTS(SELECT 1 FROM inventory_movements WHERE borrowing_item_id IS NOT NULL) THEN RAISE EXCEPTION 'Borrowing history exists; destructive rollback denied'; END IF; END $$;
REVOKE UPDATE(reserved,checked_out) ON equipment FROM elabtrack_runtime;
ALTER TABLE inventory_movements DROP CONSTRAINT inventory_movement_vector,DROP CONSTRAINT inventory_movement_source,DROP CONSTRAINT inventory_movement_actor,DROP CONSTRAINT inventory_movements_kind_check;
ALTER TABLE inventory_movements DROP COLUMN borrowing_item_id;
ALTER TABLE inventory_movements ALTER COLUMN actor_id SET NOT NULL;
ALTER TABLE inventory_movements ADD CONSTRAINT inventory_movements_kind_check CHECK(kind IN ('OPENING','ADD','REMOVE','RECONCILE'));
ALTER TABLE inventory_movements ADD CHECK(delta_reserved=0 AND delta_checked_out=0 AND delta_damaged_held=0);
DROP TABLE borrowing_operation_receipts,borrowing_events,borrowing_items,borrowings;
DROP FUNCTION protect_borrowing_state(),protect_borrowing_item();
