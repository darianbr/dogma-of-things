CREATE TABLE IF NOT EXISTS device_states (
    device_id    TEXT PRIMARY KEY,
    device_name  TEXT NOT NULL,
    device_type  TEXT NOT NULL,
    protocol     TEXT NOT NULL,
    room_id      TEXT,
    manufacturer TEXT,
    model        TEXT,
    external_id  TEXT,
    capabilities TEXT,  -- JSON array
    state        TEXT,  -- JSON object
    status       TEXT NOT NULL DEFAULT 'unknown',
    retired      INTEGER NOT NULL DEFAULT 0,
    registered_at DATETIME NOT NULL,
    updated_at   DATETIME NOT NULL
);
