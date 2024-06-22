CREATE TABLE IF NOT EXISTS pois (
    internal_id BIGSERIAL PRIMARY KEY,
    id VARCHAR(36) NOT NULL,
    created_on TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(36) NOT NULL,
    updated_by VARCHAR(36) NOT NULL,
    name VARCHAR(150) NOT NULL,
    sport_type VARCHAR(50) NOT NULL,
    description VARCHAR(8192),
    note VARCHAR(4096),
    address VARCHAR(150),
    city_id VARCHAR(255) NOT NULL,
    thumbnail_url VARCHAR(500),
    website VARCHAR(100)
);

CREATE INDEX idx_1 ON pois (city_id, sport_type)
