package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"yinlaipeng/openCMD/internal/config"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger
var globalFields logrus.Fields

// InitLogger 初始化日志
func InitLogger() error {
	logConfig := config.GetConfig().Logs
	// 构建完整的日志文件路径
	logFilePath := filepath.Join(logConfig.Dir, logConfig.Lumberjack.Filename)

	// 创建日志目录
	// 使用filepath.Dir(logFilePath)而不是直接使用logConfig.Dir
	// 这样即使logConfig.Lumberjack.Filename中包含子目录路径，也能正确创建完整的目录结构
	logDir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// 初始化logger
	Logger = logrus.New()

	// 设置日志级别
	// 把日志级别转换为 logrus 级别
	switch logConfig.Level {
	case "debug":
		Logger.SetLevel(logrus.DebugLevel)
	case "info":
		Logger.SetLevel(logrus.InfoLevel)
	case "warn":
		Logger.SetLevel(logrus.WarnLevel)
	case "error":
		Logger.SetLevel(logrus.ErrorLevel)
	default:
		Logger.SetLevel(logrus.InfoLevel)
	}

	// 设置日志格式
	if logConfig.Format == "json" {
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// 配置输出
	var writers []io.Writer

	// 控制台输出
	if logConfig.Console {
		writers = append(writers, os.Stdout)
	}

	// 文件输出
	if logConfig.File {
		if logConfig.Lumberjack.Rotate {
			// 使用lumberjack实现日志轮转
			lumberjackLogger := &lumberjack.Logger{
				Filename:   logFilePath,                     // 完整的日志文件路径
				MaxSize:    logConfig.Lumberjack.MaxSize,    // 日志文件最大大小
				MaxBackups: logConfig.Lumberjack.MaxBackups, // 保留的最大备份日志文件数
				MaxAge:     logConfig.Lumberjack.MaxAge,     // 保留的最大日志文件天数
				Compress:   logConfig.Lumberjack.Compress,   // 是否压缩旧的日志文件
			}
			writers = append(writers, lumberjackLogger)
		} else {
			// 普通文件输出
			file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return fmt.Errorf("failed to open log file: %w", err)
			}
			writers = append(writers, file)
		}
	}

	// 设置多输出
	if len(writers) > 0 {
		// logrus 没有 NewMultiWriter，使用 io.MultiWriter 代替
		Logger.SetOutput(io.MultiWriter(writers...))
	}

	return nil
}

// GetLogger 获取日志实例
func GetLogger() *logrus.Logger {
	if Logger == nil {
		// 如果未初始化，返回默认日志实例
		defaultLogger := logrus.New()
		defaultLogger.SetLevel(logrus.InfoLevel)
		defaultLogger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
		return defaultLogger
	}
	return Logger
}

// WithFields 添加字段到日志条目
func WithFields(fields logrus.Fields) *logrus.Entry {
	// 如果fields为nil，创建一个空的Fields
	if fields == nil {
		fields = make(logrus.Fields)
	}

	// 合并全局字段和本地字段
	if len(globalFields) > 0 {
		for k, v := range globalFields {
			if _, exists := fields[k]; !exists {
				fields[k] = v
			}
		}
	}
	return GetLogger().WithFields(fields)
}

// SetGlobalFields 设置全局字段
func SetGlobalFields(fields logrus.Fields) {
	globalFields = fields
}

// AddGlobalField 添加单个全局字段
func AddGlobalField(key string, value interface{}) {
	if globalFields == nil {
		globalFields = make(logrus.Fields)
	}
	globalFields[key] = value
}

// 以下是便捷方法

func Debug(args ...interface{}) {
	WithFields(nil).Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	WithFields(nil).Debugf(format, args...)
}

func Info(args ...interface{}) {
	WithFields(nil).Info(args...)
}

func Infof(format string, args ...interface{}) {
	WithFields(nil).Infof(format, args...)
}

func Warn(args ...interface{}) {
	WithFields(nil).Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	WithFields(nil).Warnf(format, args...)
}

func Error(args ...interface{}) {
	WithFields(nil).Error(args...)
}

func Errorf(format string, args ...interface{}) {
	WithFields(nil).Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	WithFields(nil).Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	WithFields(nil).Fatalf(format, args...)
}

// 以下是带字段的便捷方法

func DebugWithFields(fields logrus.Fields, args ...interface{}) {
	WithFields(fields).Debug(args...)
}

func DebugfWithFields(fields logrus.Fields, format string, args ...interface{}) {
	WithFields(fields).Debugf(format, args...)
}

func InfoWithFields(fields logrus.Fields, args ...interface{}) {
	WithFields(fields).Info(args...)
}

func InfofWithFields(fields logrus.Fields, format string, args ...interface{}) {
	WithFields(fields).Infof(format, args...)
}

func WarnWithFields(fields logrus.Fields, args ...interface{}) {
	WithFields(fields).Warn(args...)
}

func WarnfWithFields(fields logrus.Fields, format string, args ...interface{}) {
	WithFields(fields).Warnf(format, args...)
}

func ErrorWithFields(fields logrus.Fields, args ...interface{}) {
	WithFields(fields).Error(args...)
}

func ErrorfWithFields(fields logrus.Fields, format string, args ...interface{}) {
	WithFields(fields).Errorf(format, args...)
}

func FatalWithFields(fields logrus.Fields, args ...interface{}) {
	WithFields(fields).Fatal(args...)
}

func FatalfWithFields(fields logrus.Fields, format string, args ...interface{}) {
	WithFields(fields).Fatalf(format, args...)
}
