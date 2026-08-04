-- Migration 0004: Add items table for equipment
-- Up migration

CREATE TABLE items (
    id VARCHAR(64) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(32) NOT NULL,
    rarity VARCHAR(32),
    source VARCHAR(16) NOT NULL,
    edition ENUM('2014','2024') NOT NULL DEFAULT '2024',
    data JSON NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_items_type (type),
    INDEX idx_items_edition (edition),
    FULLTEXT INDEX idx_items_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;