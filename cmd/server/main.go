package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"taskmgmtsystem/internal/adaptors/persistance"
	"taskmgmtsystem/internal/config"
	"taskmgmtsystem/internal/interfaces/input/api/rest/handler"
	"taskmgmtsystem/internal/interfaces/input/api/rest/routes"
	"taskmgmtsystem/internal/usecase"
	"taskmgmtsystem/pkg/migrate"
)

func main() {

	database, err := persistance.NewDatabase()
	if err != nil {
		log.Fatalf("Error connecting to db %v", err)
	}
	fmt.Println("Connected to database")

	// fetch current cwd
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error fetching cwd %v", err)
	}

	//run migrations
	migrate := migrate.NewMigrate(
		database.GetDB(),
		cwd+"/migrations",
	)
	err = migrate.RunMigrations()
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// loadconfig
	configurations, err := config.LoadConfig()
	if err != nil {
		fmt.Println("failed to load config")
	}

	//repos
	userRepo := persistance.NewUserRepo(database)

	//services
	userService := usecase.NewUserService(userRepo)

	// handler
	userHandler := handler.NewUserHandler(configurations, userService)

	router := routes.InitRoutes(&userHandler)

	err = http.ListenAndServe(fmt.Sprintf(":%s", configurations.APP_PORT), router)
	if err != nil {
		fmt.Printf("failed to start server %v", err)
		os.Exit(1)
	}

}
