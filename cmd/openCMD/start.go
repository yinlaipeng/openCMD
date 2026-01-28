package openCMD

import (
	"os"
	"yinlaipeng/openCMD/internal/config"
	"yinlaipeng/openCMD/internal/db"
	"yinlaipeng/openCMD/internal/log"

	"github.com/spf13/cobra"
)

var configPath string

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the openCMD server",
	Long:  `Start the openCMD server with the specified configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 执行启动逻辑

		// 1. 验证配置文件是否存在
		if configPath == "" {
			log.Error("config path is empty")
			os.Exit(1)
		}

		// 2. 加载配置
		err := config.LoadConfig(configPath)
		if err != nil {
			log.Errorf("load config failed: %v", err)
			os.Exit(1)
		}

		// 3. 加载日志
		if err := log.InitLogger(); err != nil {
			log.Errorf("init logger failed: %v", err)
			os.Exit(1)
		}

		// 4. 初始化数据库
		db.InitMySQL()
	},
}
