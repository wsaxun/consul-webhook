package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"consul-webhook/app"
	"consul-webhook/config"
	"consul-webhook/pkg"
)

func main() {
	// log日志输出
	log.SetOutput(os.Stdout)

	// 配置初始化
	config.InitConfig()

	// repo初始化
	pkg.InitRepo()

	engine := gin.New()
	port := config.GetConfig().App.Port
	server := http.Server{Addr: ":" + port, Handler: engine}

	// 设置gin运行模式
	envMap := map[string]string{
		"dev":    gin.DebugMode,
		"beta":   gin.TestMode,
		"online": gin.ReleaseMode,
	}
	gin.SetMode(envMap[config.GetConfig().App.Env])

	ch := make(chan os.Signal, 1)
	// 停止 SIGINT: 2(ctrl +c)  SIGKILL:9  SIGTERM:15
	// 只接受指定的信号
	signal.Notify(ch, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	go func() {
		select {
		case <-ch:
			log.Println("shutdown...")
			timeout := 5 * time.Second
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				log.Fatal("http server has exited")
			}
			log.Println("http server has exited")
		}
	}()

	// 全局中间件
	engine.Use(app.Recover)
	engine.Use(gin.Logger())

	// 注册路由
	app.Router(engine)

	// 启动
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("http server start failed: %v\n", err)
	}
}
