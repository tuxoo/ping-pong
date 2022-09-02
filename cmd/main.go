package main

import (
	"ping-pong/internal/app"
)

const (
	configPath = "config/config"
)

func main() {
	app.Run(configPath)
}
