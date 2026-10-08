DO $$ BEGIN
 IF EXISTS(SELECT 1 FROM equipment) OR EXISTS(SELECT 1 FROM equipment_categories) OR EXISTS(SELECT 1 FROM inventory_operation_receipts) OR EXISTS(SELECT 1 FROM inventory_audit_events) THEN RAISE EXCEPTION 'Inventory data/history exists; destructive rollback denied'; END IF;
END $$;
ALTER TABLE equipment DROP CONSTRAINT equipment_current_image;
DROP TABLE inventory_operation_receipts,inventory_audit_events,inventory_movements,equipment_images,equipment,equipment_categories;
DROP FUNCTION reject_inventory_history_change();
