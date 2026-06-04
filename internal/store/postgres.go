package store

import (
	"errors"
	"fmt"
	"time"

	"game-server/internal/config"
	"game-server/internal/model"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// PostgresStore PostgreSQL 存储（世界杯场景数据）
type PostgresStore struct {
	db *gorm.DB
}

// NewPostgresStore 创建 PostgreSQL 连接
func NewPostgresStore(cfg *config.PostgresConfig) (*PostgresStore, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("连接 PostgreSQL 失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层 sql.DB 失败: %w", err)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}

	if err := db.AutoMigrate(&model.GameProfile{}); err != nil {
		return nil, fmt.Errorf("数据库迁移失败: %w", err)
	}

	zap.L().Info("PostgreSQL 初始化成功")
	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// GetGameProfile 按 open_id + scene 查询档案
func (s *PostgresStore) GetGameProfile(openID, scene string) (*model.GameProfile, error) {
	var profile model.GameProfile
	err := s.db.Where("open_id = ? AND scene = ?", openID, scene).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// GetOrCreateGameProfile 不存在时创建默认 scene_settings
func (s *PostgresStore) GetOrCreateGameProfile(openID, scene string) (*model.GameProfile, error) {
	profile, err := s.GetGameProfile(openID, scene)
	if err == nil {
		return profile, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	defaults, err := model.DefaultSceneSettings().JSON()
	if err != nil {
		return nil, err
	}

	profile = &model.GameProfile{
		OpenID:        openID,
		Scene:         scene,
		CreateTime:    time.Now().Unix(),
		SceneSettings: defaults,
	}
	if err := s.db.Create(profile).Error; err != nil {
		return nil, err
	}
	return profile, nil
}

// SaveSceneSettings 更新 scene_settings 字段
func (s *PostgresStore) SaveSceneSettings(profileID int64, settings string) error {
	return s.db.Model(&model.GameProfile{}).Where("id = ?", profileID).
		Update("scene_settings", settings).Error
}
