package main

import (
	"support-bot/internal/app"
)

const configPath = "config/config.yaml"

func main() {
	app.Run(configPath)
}
