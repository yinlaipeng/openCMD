package db

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"yinlaipeng/openCMD/internal/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// sqliteDB 全局SQLite数据库连接
	sqliteDB *gorm.DB
	// sqliteOnce 确保初始化只执行一次
	sqliteOnce sync.Once
	// sqliteInitErr 初始化错误
	sqliteInitErr error
	// sqliteConfig 保存配置
	sqliteConfig *config.SqliteConfig
)

// InitSQLite 初始化SQLite数据库连接（单例模式）
func InitSQLite(config *config.SqliteConfig) error {
	// 保存配置
	sqliteConfig = config
	
	// 使用sync.Once确保初始化只执行一次
	sqliteOnce.Do(func() {
		sqliteInitErr = initSQLiteDB()
	})
	
	return sqliteInitErr
}

// initSQLiteDB 内部初始化SQLite数据库连接
func initSQLiteDB() error {
	if sqliteConfig == nil {
		return fmt.Errorf("sqlite config is not set")
	}
	
	// 创建数据库文件所在目录
	dbDir := filepath.Dir(sqliteConfig.FileName)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}
	
	// 配置GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// 连接数据库
	var err error
	sqliteDB, err = gorm.Open(sqlite.Open(sqliteConfig.FileName), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to sqlite: %w", err)
	}

	// 配置连接池
	sqlDB, err := sqliteDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Hour)

	return nil
}

// GetSQLiteDB 获取SQLite数据库连接（单例模式）
func GetSQLiteDB() (*gorm.DB, error) {
	// 如果还未初始化，返回错误
	if sqliteDB == nil {
		if sqliteInitErr != nil {
			return nil, sqliteInitErr
		}
		return nil, fmt.Errorf("sqlite database not initialized")
	}
	return sqliteDB, nil
}

// CloseSQLite 关闭SQLite数据库连接
func CloseSQLite() error {
	if sqliteDB != nil {
		sqlDB, err := sqliteDB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
