package main

import (
	"app/internal/application"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// env
	godotenv.Load()

	// app
	// - config
	cfg := &application.ConfigDefaultInMemory{
		FileLoader: os.Getenv("PATH_FILE_LOADER_VEHICLES"),
		Addr:       os.Getenv("SERVER_ADDR"),
	}
	// - app
	app := application.NewDefaultInMemory(cfg)
	// - run
	if err := app.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
