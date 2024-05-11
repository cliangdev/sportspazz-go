CREATE TABLE IF NOT EXISTS users (
     internal_id BIGINT AUTO_INCREMENT PRIMARY KEY,
     id VARCHAR(36),
     created_on DATETIME(3) NOT NULL,
     updated_on DATETIME(3) NOT NULL,
     email VARCHAR(320) NOT NULL,
     UNIQUE(id),
     UNIQUE(email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
