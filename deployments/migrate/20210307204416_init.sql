-- +goose Up
CREATE TABLE IF NOT EXISTS profiles
(
    id             BIGINT(20) UNSIGNED NOT NULL COMMENT 'unique identity number',
    account        VARCHAR(16)         NOT NULL COMMENT 'user login identity account',
    password       VARCHAR(16)         NOT NULL COMMENT 'user login identity password',
    nickname       VARCHAR(20)                  DEFAULT '' COMMENT 'user nickname',
    first_name     VARCHAR(20)                  DEFAULT '' COMMENT 'user first name',
    last_name      VARCHAR(20)                  DEFAULT '' COMMENT 'user last name',
    email          VARCHAR(254)                 DEFAULT '' COMMENT 'user email address',
    avatar         VARCHAR(300)                 DEFAULT '' COMMENT 'means user profile picture URL',
    created_at     DATETIME            NOT NULL COMMENT 'this account create time',
    updated_at     DATETIME            NOT NULL COMMENT 'this account update time',
    suspend_status TINYINT(1) UNSIGNED NOT NULL DEFAULT 1 COMMENT 'this account is suspend',
    PRIMARY KEY (`id`),
    INDEX (`id`, `account`)
)
    COMMENT 'define identity profile schema';

CREATE TABLE IF NOT EXISTS access
(
    user_id    BIGINT(20) UNSIGNED NOT NULL COMMENT 'user identity number',
    role       TINYINT(1) UNSIGNED NOT NULL DEFAULT 1 COMMENT 'user default role',
    created_at DATETIME            NOT NULL COMMENT 'role create time',
    PRIMARY KEY (`user_id`),
    INDEX (`user_id`)
)
    COMMENT 'define user access schema';
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS profiles;
DROP TABLE IF EXISTS access;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
