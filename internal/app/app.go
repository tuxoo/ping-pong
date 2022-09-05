package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"ping-pong/internal/client"
	"ping-pong/internal/config"
	"ping-pong/internal/server"
	"syscall"
	"time"
)

func Run(configPath string) {
	dir, _ := os.Getwd()
	logDir := fmt.Sprintf("%s/logs", dir)

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.Mkdir(logDir, 0755)
		if err != nil {
			return
		}
	}

	logFile := fmt.Sprintf("%s/pp.log", logDir)

	log, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE, 0755)
	logrus.SetOutput(log)

	cfg, err := config.InitConfig(configPath)
	if err != nil {
		logrus.Fatalf("error occured durinng the initialization configs: %s", err.Error())
	}

	httpClient := client.NewHTTPClient(cfg)

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		go func() {
			time.Sleep(1 * time.Second)
			url := fmt.Sprintf("%s:%s/ping", cfg.AppConfig.PingPongServiceUrl, cfg.AppConfig.PingPongServicePort)
			resp := httpClient.GetString(url)
			logrus.Info(resp)
		}()

		body, err := json.Marshal("pong")
		if err != nil {
			// Handle Error
		}

		if _, err := w.Write(body); err != nil {
			// Handle Error
		}
	})

	httpServer := server.NewHTTPServer(cfg)

	go func() {
		if err := httpServer.Run(cfg.ServerConfig.Port); err != nil {
			logrus.Errorf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Printf("Ping Pong application has started on [%s] port", cfg.ServerConfig.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Ping Pong application shutting down")

	if err := httpServer.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on http server shutting down: %s", err.Error())
	}
}
