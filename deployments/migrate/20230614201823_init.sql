-- +goose Up
CREATE TABLE IF NOT EXISTS users
(
    id         VARCHAR(20) NOT NULL UNIQUE,
    username   VARCHAR(20) NOT NULL,
    password   VARCHAR(64) NOT NULL,
    nickname   VARCHAR(20)          DEFAULT '',
    first_name VARCHAR(20)          DEFAULT '',
    last_name  VARCHAR(20)          DEFAULT '',
    email      VARCHAR(255)         DEFAULT '',
    avatar     VARCHAR(300)         DEFAULT '',
    created_at BIGINT      NOT NULL,
    updated_at BIGINT      NOT NULL,
    status     SMALLINT    NOT NULL DEFAULT 1,
    PRIMARY KEY (ID)
);

CREATE TABLE IF NOT EXISTS sessions
(
    id                VARCHAR(20) NOT NULL UNIQUE,
    user_id           VARCHAR(20) NOT NULL,
    created_at        BIGINT      NOT NULL,
    updated_at        BIGINT      NOT NULL,
    expire_at         BIGINT      NOT NULL,
    ip_address        VARCHAR(64)          DEFAULT '',
    idp_provider      VARCHAR(20)          DEFAULT '',
    platform          VARCHAR(20) NOT NULL DEFAULT 1,
    device_model      VARCHAR(64)          DEFAULT '',
    device_name       VARCHAR(64)          DEFAULT '',
    device_os_version VARCHAR(64)          DEFAULT '',
    PRIMARY KEY (ID)
);


-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS sessions;

-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
