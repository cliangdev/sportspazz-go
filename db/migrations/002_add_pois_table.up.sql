CREATE TABLE IF NOT EXISTS pois (
     internal_id BIGINT AUTO_INCREMENT PRIMARY KEY,
     id VARCHAR(36) NOT NULL,
     created_on DATETIME(3) NOT NULL,
     updated_on DATETIME(3) NOT NULL,
     created_by VARCHAR(36) NOT NULL,
     updated_by VARCHAR(36) NOT NULL,
     name VARCHAR(150) NOT NULL,
     sport_type VARCHAR(50) NOT NULL,
     description VARCHAR(8192),
     note VARCHAR(4096),
     address VARCHAR(150),
     city_id VARCHAR(255) NOT NULL,
     thumbnail_url VARCHAR(500),
     website_url VARCHAR(100),
     INDEX idx_1 (city_id, sport_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
