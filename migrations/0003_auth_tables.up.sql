-- Migration 0003: Auth tables + characters table + campaigns + campaign_members
-- Up migration

-- Users table (GitHub OAuth)
CREATE TABLE users (
    id              VARCHAR(64)  NOT NULL PRIMARY KEY,
    github_id       BIGINT       NOT NULL UNIQUE,
    login           VARCHAR(128) NOT NULL,
    name            VARCHAR(256),
    avatar_url      VARCHAR(512),
    email           VARCHAR(256),
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Characters table (cloud, owned by user)
CREATE TABLE characters (
    id              VARCHAR(64)  NOT NULL PRIMARY KEY,
    user_id         VARCHAR(64)  NOT NULL,
    campaign_id     VARCHAR(64)  NULL,
    name            VARCHAR(128) NOT NULL,
    is_npc          BOOLEAN      NOT NULL DEFAULT FALSE,
    draft           JSON         NOT NULL,
    sheet           JSON         NULL,
    live            JSON         NULL,
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_characters_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_characters_user ON characters(user_id);

-- Sessions table
CREATE TABLE sessions (
    id              VARCHAR(128) NOT NULL PRIMARY KEY,
    user_id         VARCHAR(64)  NOT NULL,
    expires_at      TIMESTAMP    NOT NULL,
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_sessions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_sessions_user ON sessions(user_id);
CREATE INDEX idx_sessions_expires ON sessions(expires_at);

-- Campaigns table (owned by DM)
CREATE TABLE campaigns (
    id              VARCHAR(64)  NOT NULL PRIMARY KEY,
    owner_id        VARCHAR(64)  NOT NULL,
    name            VARCHAR(128) NOT NULL,
    description     TEXT,
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_campaigns_owner FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_campaigns_owner ON campaigns(owner_id);

-- Campaign members table
CREATE TABLE campaign_members (
    campaign_id     VARCHAR(64)  NOT NULL,
    user_id         VARCHAR(64)  NOT NULL,
    role            ENUM('dm', 'player') NOT NULL DEFAULT 'player',
    joined_at       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (campaign_id, user_id),
    CONSTRAINT fk_members_campaign FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE,
    CONSTRAINT fk_members_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Add campaign_id FK to characters (after campaigns table exists)
ALTER TABLE characters
    ADD CONSTRAINT fk_characters_campaign FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE SET NULL;
CREATE INDEX idx_characters_campaign ON characters(campaign_id);