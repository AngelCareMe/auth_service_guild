CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    battletag TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS blizzard_users (
    id TEXT PRIMARY KEY,
    battletag TEXT NOT NULL UNIQUE,
    FOREIGN KEY (battletag) REFERENCES users(battletag) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS jwt_tokens (
    user_id TEXT PRIMARY KEY,
    battletag TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    expiry TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS blizzard_tokens (
    user_id TEXT PRIMARY KEY,
    blizzard_id TEXT NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    expiry TIMESTAMP NOT NULL,
    token_type TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);