DO $$ BEGIN IF EXISTS(SELECT 1 FROM profile_images) THEN RAISE EXCEPTION 'Profile images or removal versions exist; preserve account audit and current presentation data'; END IF; END $$;
DROP TABLE profile_images;
