-- 世界杯场景数据表（PostgreSQL）
CREATE TABLE IF NOT EXISTS game_profile (
    id BIGSERIAL PRIMARY KEY,
    open_id VARCHAR(64) NOT NULL DEFAULT '',
    scene VARCHAR(8) NOT NULL DEFAULT '',
    create_time BIGINT NOT NULL DEFAULT 0,
    scene_settings TEXT
);

CREATE INDEX IF NOT EXISTS idx_open_id ON game_profile (open_id, scene);
