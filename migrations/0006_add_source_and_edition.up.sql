-- Migration 0003: Add source and edition columns to content tables
-- Up migration

-- Classes
ALTER TABLE classes ADD COLUMN source VARCHAR(16) NOT NULL DEFAULT 'XPHB';
ALTER TABLE classes ADD COLUMN edition ENUM('2014','2024') NOT NULL DEFAULT '2024';
CREATE INDEX idx_classes_edition ON classes(edition);

-- Subclasses
ALTER TABLE subclasses ADD COLUMN source VARCHAR(16) NOT NULL DEFAULT 'XPHB';
ALTER TABLE subclasses ADD COLUMN edition ENUM('2014','2024') NOT NULL DEFAULT '2024';
CREATE INDEX idx_subclasses_edition ON subclasses(edition);

-- Species
ALTER TABLE species ADD COLUMN source VARCHAR(16) NOT NULL DEFAULT 'XPHB';
ALTER TABLE species ADD COLUMN edition ENUM('2014','2024') NOT NULL DEFAULT '2024';
CREATE INDEX idx_species_edition ON species(edition);

-- Backgrounds
ALTER TABLE backgrounds ADD COLUMN source VARCHAR(16) NOT NULL DEFAULT 'XPHB';
ALTER TABLE backgrounds ADD COLUMN edition ENUM('2014','2024') NOT NULL DEFAULT '2024';
CREATE INDEX idx_backgrounds_edition ON backgrounds(edition);

-- Spells
ALTER TABLE spells ADD COLUMN source VARCHAR(16) NOT NULL DEFAULT 'PHB';
ALTER TABLE spells ADD COLUMN edition ENUM('2014','2024') NOT NULL DEFAULT '2014';
CREATE INDEX idx_spells_edition ON spells(edition);

-- Feats
ALTER TABLE feats ADD COLUMN source VARCHAR(16) NOT NULL DEFAULT 'XPHB';
ALTER TABLE feats ADD COLUMN edition ENUM('2014','2024') NOT NULL DEFAULT '2024';
CREATE INDEX idx_feats_edition ON feats(edition);

-- Features
ALTER TABLE features ADD COLUMN source VARCHAR(16) NOT NULL DEFAULT 'PHB';
ALTER TABLE features ADD COLUMN edition ENUM('2014','2024') NOT NULL DEFAULT '2014';
CREATE INDEX idx_features_edition ON features(edition);