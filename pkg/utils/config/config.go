package config

import (
	"github.com/mitchellh/go-homedir"
	"github.com/qq5272689/goldden-go/pkg/utils/base_dir"
	"github.com/qq5272689/goldden-go/pkg/utils/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"path"
)

func init() {
	// 16为密码加密
	viper.SetDefault("golddengo.password.key", "KY9ciRr1Q7sOgjVV")
	// mysql连接url
	viper.SetDefault("mysql.dsn", "goldden_go:goldden_go123@tcp(192.168.8.154:12301)/goldden_go?charset=utf8&parseTime=True&loc=Local")
	//监听地址
	viper.SetDefault("listen", ":8080")
}

func InitConfig(cfgFile string) error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		logger.Debug("use default cfgfile")
		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		viper.AddConfigPath(home)
		bd := base_dir.GetBaseDir()
		viper.AddConfigPath(bd)
		viper.AddConfigPath(path.Join(bd, "conf"))
		viper.AddConfigPath(path.Join(bd, "etc"))
		viper.SetConfigName("base-service")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		logger.Info("Using config file:" + viper.ConfigFileUsed())
	} else {
		logger.Warn("read in config", zap.Error(err))
	}
	return nil
}
