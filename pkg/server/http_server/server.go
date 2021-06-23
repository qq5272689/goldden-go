package http_server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/qq5272689/goldden-go/pkg/server/http_server/handlers"
	"github.com/qq5272689/goldden-go/pkg/utils/gin_middleware"
	ghttp "github.com/qq5272689/goldden-go/pkg/utils/http"
	"github.com/qq5272689/goldden-go/pkg/utils/logger"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ginZapRecoveryErrResponse struct {
}

func (ger ginZapRecoveryErrResponse) SetErr(err error) interface{} {
	return ghttp.CommonErrResult(err)
}

type HttpServer struct {
	g *gin.Engine
	//viper.GetString("listen")
	Addr string
}

func (hs *HttpServer) Server() *gin.Engine {
	return hs.g
}

// @title GOLDDEN-GO接口
// @version 1.0
// @description GOLDDEN-GO接口
// @BasePath "/api"
func (hs *HttpServer) router() {
	v1 := hs.g.Group("/api/v1")
	//用户相关
	v1.GET("/user/:userid", handlers.GetUser)
	v1.GET("/user", handlers.SearchUser)
	v1.PUT("/user", handlers.UpdateUser)
	v1.POST("/user", handlers.CreateUser)
	v1.DELETE("/user", handlers.DeleteUser)

	//登录相关
	v1.GET("/verify", handlers.Verify)
	v1.GET("/logout", handlers.LogOut)
	v1.GET("/login/local", handlers.LoginLocal)
}

func (hs *HttpServer) listenAndServe() {
	logger.Info("start listenAndServe", zap.String("listen addr", hs.Addr))
	srv := &http.Server{
		Addr:    hs.Addr,
		Handler: hs.g,
	}
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("listen fail", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Debug("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown ", zap.Error(err))
	}

	logger.Debug("Server exiting")
	//return http.ListenAndServe(hs.Addr, hs.g)
}

func (hs *HttpServer) ListenAndServe() {
	hs.g.Use(gin_middleware.GinZapLogger(logger.GetLogger()), gin_middleware.GinZapRecovery(logger.GetLogger(), ginZapRecoveryErrResponse{}))
	env := viper.GetString("env")
	if env == "dev" || env == "local" {
		hs.g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	hs.router()
	hs.listenAndServe()
}
