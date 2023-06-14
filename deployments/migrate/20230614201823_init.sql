-- +goose Up
CREATE TABLE IF NOT EXISTS users
(
    ID        VARCHAR(20) NOT NULL UNIQUE,
    Username  VARCHAR(20) NOT NULL,
    Password  VARCHAR(64) NOT NULL,
    Nickname  VARCHAR(20)          DEFAULT '',
    FirstName VARCHAR(20)          DEFAULT '',
    LastName  VARCHAR(20)          DEFAULT '',
    Email     VARCHAR(255)         DEFAULT '',
    Avatar    VARCHAR(300)         DEFAULT '',
    CreatedAt BIGINT      NOT NULL,
    UpdatedAt BIGINT      NOT NULL,
    Status    SMALLINT    NOT NULL DEFAULT 1,
    PRIMARY KEY (ID)
);

CREATE TABLE IF NOT EXISTS sessions
(
    ID              VARCHAR(20) NOT NULL UNIQUE,
    UserID          VARCHAR(20) NOT NULL,
    CreatedAt       BIGINT      NOT NULL,
    UpdatedAt       BIGINT      NOT NULL,
    ExpireAt        BIGINT      NOT NULL,
    IPAddress       VARCHAR(64)          DEFAULT '',
    IdpProvider     VARCHAR(20)          DEFAULT '',
    Platform        SMALLINT    NOT NULL DEFAULT 1,
    DeviceModel     VARCHAR(64)          DEFAULT '',
    DeviceName      VARCHAR(64)          DEFAULT '',
    DeviceOSVersion VARCHAR(64)          DEFAULT '',
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
