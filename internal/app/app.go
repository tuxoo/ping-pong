package app

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"ping-pong/internal/client"
	"ping-pong/internal/config"
	"ping-pong/internal/server"
	"syscall"
)

func Run(configPath string) {
	cfg, err := config.InitConfig(configPath)
	if err != nil {
		logrus.Fatalf("error occured durinng the initialization configs: %s", err.Error())
	}

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		logrus.Println("pong")
	})
	httpServer := server.NewHTTPServer(cfg)

	go func() {
		if err := httpServer.Run(cfg.ServerConfig.Port); err != nil {
			logrus.Errorf("error occured while running http server: %s", err.Error())
		}
	}()

	httpClient := client.NewHTTPClient(cfg)
	str := httpClient.GetString("http://localhost:8000/api/v1/ping")
	logrus.Print(str)

	logrus.Printf("Ping Pong application has started on [%s] port", cfg.ServerConfig.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Ping Pong application shutting down")

	if err := httpServer.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on http server shutting down: %s", err.Error())
	}
}
