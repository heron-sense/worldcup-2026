CREATE table game_profile (
    id bigint unsigned primary key,
    open_id VARCHAR(64) not null default '',
    scene VARCHAR(8) not null default '',
    create_time bigint unsigned not null default 0,
    index idx_open_id(open_id, scene),
    scene_settings text
)