-- Migration 0005: Add ETL control table for 5etools delta tracking
-- Up migration

CREATE TABLE etl_control (
    file_path VARCHAR(512) NOT NULL PRIMARY KEY,
    file_hash CHAR(64) NOT NULL,
    last_processed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    entities_inserted INT DEFAULT 0,
    entities_updated INT DEFAULT 0,
    status ENUM('pending','success','failed') DEFAULT 'pending',
    error_log TEXT,
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;