-- Migration 0003: Auth tables + characters table + campaigns + campaign_members
-- Down migration

ALTER TABLE characters DROP FOREIGN KEY fk_characters_campaign;
ALTER TABLE characters DROP INDEX idx_characters_campaign;
ALTER TABLE characters DROP INDEX idx_characters_user;
ALTER TABLE characters DROP FOREIGN KEY fk_characters_user;

DROP TABLE IF EXISTS campaign_members;
DROP TABLE IF EXISTS campaigns;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS characters;
DROP TABLE IF EXISTS users;