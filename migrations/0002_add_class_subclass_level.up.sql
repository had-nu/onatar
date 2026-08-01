ALTER TABLE classes
    ADD COLUMN subclass_level INT NOT NULL DEFAULT 0 AFTER spellcaster;
