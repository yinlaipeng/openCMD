package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// 私有配置变量
var c *Config

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Logs    LogsConfig    `yaml:"logs"`
	MySQL   MysqlConfig   `yaml:"mysql"`
	Sqlite  SqliteConfig  `yaml:"sqlite"`
	Tencent TencentConfig `yaml:"tencent"`
	AWS     AWSConfig     `yaml:"aws"`
	Aliyun  AliyunConfig  `yaml:"aliyun"`
}

type ServerConfig struct {
	Port    int    `yaml:"port"`
	Host    string `yaml:"host"`
	Timeout string `yaml:"timeout"`
	DBType  string `yaml:"dbtype"`
}

type LogsConfig struct {
	Level   string `yaml:"level"`
	Format  string `yaml:"format"`
	Console bool   `yaml:"console"`
	// 是否输出到文件
	File       bool       `yaml:"file"`
	Dir        string     `yaml:"dir"`
	Lumberjack Lumberjack `yaml:"lumberjack"`
}

type Lumberjack struct {
	// 是否开启日志轮转
	Rotate     bool   `yaml:"rotate"`
	Filename   string `yaml:"filename"`
	MaxSize    int    `yaml:"maxsize"`
	MaxBackups int    `yaml:"maxbackups"`
	MaxAge     int    `yaml:"maxage"`
	Compress   bool   `yaml:"compress"`
}

type MysqlConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type SqliteConfig struct {
	FileName string `yaml:"filename"`
}

type TencentConfig struct {
	AppID  string `yaml:"appid"`
	AppKey string `yaml:"appkey"`
}

type AWSConfig struct {
	Region    string `yaml:"region"`
	AccessKey string `yaml:"accesskey"`
	SecretKey string `yaml:"secretkey"`
}

type AliyunConfig struct {
	Region    string `yaml:"region"`
	AccessKey string `yaml:"accesskey"`
	SecretKey string `yaml:"secretkey"`
}

func NewConfig() *Config {
	return &Config{}
}

// GetConfig 获取配置实例（返回拷贝，防止外部修改）
func GetConfig() *Config {
	if c == nil {
		return &Config{}
	}
	// 返回配置的深拷贝，防止外部修改
	copyConfig := *c
	return &copyConfig
}

// server默认值的函数
func (c *Config) ServerDefault() {
	if c.Server.Port == 0 {
		c.Server.Port = 8080
	}
	if c.Server.Host == "" {
		c.Server.Host = "0.0.0.0"
	}
	if c.Server.Timeout == "" {
		c.Server.Timeout = "5s"
	}
	if c.Server.DBType == "" {
		c.Server.DBType = "mysql"
	}
	if c.Server.DBType != "mysql" && c.Server.DBType != "sqlite" {
		c.Server.DBType = "mysql"
	}
}

// logs默认值的函数
func (c *Config) LogsDefault() {
	if c.Logs.Level == "" {
		c.Logs.Level = "info"
	}
	if c.Logs.Dir == "" {
		c.Logs.Dir = "logs/"
	}
	if c.Logs.Lumberjack.Filename == "" {
		c.Logs.Lumberjack.Filename = "openCMD.log"
	}
	if c.Logs.Format == "" {
		c.Logs.Format = "json"
	}
	if !c.Logs.Console {
		c.Logs.Console = true
	}
	// 是否开启日志轮转
	if !c.Logs.Lumberjack.Rotate {
		c.Logs.Lumberjack.Rotate = true
	}
	if !c.Logs.Lumberjack.Compress {
		c.Logs.Lumberjack.Compress = true
	}
	if c.Logs.Lumberjack.MaxSize == 0 {
		c.Logs.Lumberjack.MaxSize = 100
	}
	if c.Logs.Lumberjack.MaxBackups == 0 {
		c.Logs.Lumberjack.MaxBackups = 7
	}
	if c.Logs.Lumberjack.MaxAge == 0 {
		c.Logs.Lumberjack.MaxAge = 7
	}
}

// db默认值的函数
func (c *Config) mysqlDefault() {
	if c.Server.DBType == "mysql" {
		if c.MySQL.Host == "" {
			c.MySQL.Host = "localhost"
		}
		if c.MySQL.Port == 0 {
			c.MySQL.Port = 3306
		}
		if c.MySQL.User == "" {
			c.MySQL.User = "root"
		}
		if c.MySQL.Password == "" {
			c.MySQL.Password = "123456"
		}
		if c.MySQL.DBName == "" {
			c.MySQL.DBName = "openCMD"
		}
	} else if c.Server.DBType == "sqlite" {
		if c.Sqlite.FileName == "" {
			c.Sqlite.FileName = "openCMD.db"
		}
	}
}

// 增加一个配置默认值的函数
func (c *Config) SetDefault() {
	c.ServerDefault()
	c.mysqlDefault()
	c.LogsDefault()
}

// ValidateConfig 验证配置有效性
func (c *Config) ValidateConfig() error {
	// 验证服务器配置
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}
	if c.Server.Host == "" {
		return fmt.Errorf("server host cannot be empty")
	}

	// 验证数据库配置（如果需要连接数据库）
	if c.MySQL.Host == "" {
		return fmt.Errorf("db host cannot be empty")
	}
	if c.MySQL.Port == 0 {
		return fmt.Errorf("db port cannot be zero")
	}
	if c.MySQL.User == "" {
		return fmt.Errorf("db user cannot be empty")
	}
	if c.MySQL.Password == "" {
		return fmt.Errorf("db password cannot be empty")
	}
	if c.MySQL.DBName == "" {
		return fmt.Errorf("db name cannot be empty")
	}

	return nil
}

// LoadConfig 加载配置文件
func LoadConfig(path string) error {
	// 判断文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("config file not found: %s", path)
	}

	// 读取YAML文件
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// 初始化配置结构体
	config := NewConfig()

	// 解析YAML内容
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// 设置默认值（无论配置文件是否存在，都设置默认值）
	config.SetDefault()

	// 验证配置
	err = config.ValidateConfig()
	if err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// 更新全局配置
	c = config

	return nil
}
