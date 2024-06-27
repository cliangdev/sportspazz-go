ALTER TABLE pois ADD COLUMN google_place_id VARCHAR(255);

ALTER TABLE pois ADD CONSTRAINT unique_google_place_id UNIQUE (google_place_id);
