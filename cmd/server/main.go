package main

import (
	"fmt"
	"log"
	"os"
	"taskmgmtsystem/internal/adaptors/persistance"
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

}
