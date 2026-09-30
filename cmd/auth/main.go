package main

import (
	"database/sql"
	"fmt"

	"github.com/dwikie/auth-manager/internal/api"
	"github.com/dwikie/auth-manager/internal/config"
	"github.com/dwikie/auth-manager/internal/db"
	"github.com/gin-gonic/gin"
)

type App struct {
	Config *config.AppConfig
	DB     *sql.DB
}

func main() {
	app, err := InitializeApps()
	if err != nil {
		fmt.Println("Error initializing the application:", err)
		return
	}

	defer app.DB.Close()

	server := gin.New()

	repo := api.RegisterRepositories(app.DB)
	service := api.RegisterServices(&repo)
	handler := api.InitializeHandler(&service)
	api.RegisterRoutes(server, &handler)

	server.Run(fmt.Sprintf("%s:%s", app.Config.Host, app.Config.Port))
}

func InitializeApps() (*App, error) {
	// Load the application configuration
	config := config.NewAppConfig()
	db, err := db.NewConnection(config.DatabaseUrl)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return nil, err
	}

	return &App{
		Config: config,
		DB:     db,
	}, err
}
