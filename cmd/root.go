package cmd

import (
	"github.com/qq5272689/goldden-go/pkg/utils/config"
	"github.com/qq5272689/goldden-go/pkg/utils/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var cfgFile string
var env string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goldden",
	Short: "",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(rootInit)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/base-service.yaml)")
	rootCmd.PersistentFlags().StringVar(&env, "env", "local", "env name (default is local")
	viper.BindPFlag("env", rootCmd.PersistentFlags().Lookup("env"))
}

func rootInit() {
	logger.JsonLoggerInit(env)
	logger.Debug("cfg:" + cfgFile)
	if err := config.InitConfig(cfgFile, "base-service"); err != nil {
		logger.GetLogger().Fatal("InitConfig Fail!!!", zap.Error(err))
	}
	logger.Debug("config:", zap.Any("all", viper.ConfigFileUsed()))
}
