package model

// GameProfile 对应 game_profile 表，按 open_id + scene 存储场景状态。
type GameProfile struct {
	ID            int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	OpenID        string `gorm:"column:open_id;size:64;not null;index:idx_open_id,priority:1" json:"open_id"`
	Scene         string `gorm:"size:8;not null;index:idx_open_id,priority:2" json:"scene"`
	CreateTime    int64  `gorm:"column:create_time;not null;default:0" json:"create_time"`
	SceneSettings string `gorm:"column:scene_settings;type:text" json:"scene_settings"`
}

func (GameProfile) TableName() string {
	return "game_profile"
}
