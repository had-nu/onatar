-- Migration 0003: Add source and edition columns to content tables
-- Down migration

-- Features
ALTER TABLE features DROP INDEX idx_features_edition;
ALTER TABLE features DROP COLUMN edition;
ALTER TABLE features DROP COLUMN source;

-- Feats
ALTER TABLE feats DROP INDEX idx_feats_edition;
ALTER TABLE feats DROP COLUMN edition;
ALTER TABLE feats DROP COLUMN source;

-- Spells
ALTER TABLE spells DROP INDEX idx_spells_edition;
ALTER TABLE spells DROP COLUMN edition;
ALTER TABLE spells DROP COLUMN source;

-- Backgrounds
ALTER TABLE backgrounds DROP INDEX idx_backgrounds_edition;
ALTER TABLE backgrounds DROP COLUMN edition;
ALTER TABLE backgrounds DROP COLUMN source;

-- Species
ALTER TABLE species DROP INDEX idx_species_edition;
ALTER TABLE species DROP COLUMN edition;
ALTER TABLE species DROP COLUMN source;

-- Subclasses
ALTER TABLE subclasses DROP INDEX idx_subclasses_edition;
ALTER TABLE subclasses DROP COLUMN edition;
ALTER TABLE subclasses DROP COLUMN source;

-- Classes
ALTER TABLE classes DROP INDEX idx_classes_edition;
ALTER TABLE classes DROP COLUMN edition;
ALTER TABLE classes DROP COLUMN source;