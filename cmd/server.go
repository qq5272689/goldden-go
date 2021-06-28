package cmd

import (
	"gitee.com/goldden-go/goldden-go/pkg/db"
	"gitee.com/goldden-go/goldden-go/pkg/server/http_server"
	"gitee.com/goldden-go/goldden-go/pkg/service"
	"gitee.com/goldden-go/goldden-go/pkg/utils/jwt"
	"gitee.com/goldden-go/goldden-go/pkg/utils/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动服务",
	Long:  `启动服务`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := serverInit(cmd)
		if err != nil {
			logger.Error("初始化服务失败！！！", zap.Error(err))
			return err
		}
		return s.ListenAndServe()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	serverCmd.Flags().BoolP("migrate", "", false, "数据库migrate")
}

func serverInit(cmd *cobra.Command) (s *http_server.HttpServer, err error) {
	if err = db.OpenDB("goldden_go", viper.GetString("mysql.dsn")); err != nil {
		return nil, err
	}
	if migrate, _ := cmd.Flags().GetBool("migrate"); migrate {
		if err = db.SetupDatabase(db.DB); err != nil {
			return nil, err
		}
	}
	if err = service.GetUserServiceDB(db.DB).InitSuperAdmin(); err != nil {
		return nil, err
	}
	s = http_server.NewHttpServer(viper.GetString("env"), viper.GetString("listen"))
	gj, err := jwt.NewGolddenJwt(viper.GetInt("jwt.exp"), viper.GetString("jwt.publicKey"), viper.GetString("jwt.privateKey"))
	if err != nil {
		return nil, err
	}
	s.AddMiddleware(gj.GinJwtMiddleware, db.GormMiddleware())
	return
}
