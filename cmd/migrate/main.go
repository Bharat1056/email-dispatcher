package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"                     // core migration engine
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // postgres driver for migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"       // used for file system migrations

	"github.com/user/queue/internal/config"
)

const (
	migrationsPath = "file://docs/migrations"
)

func main() {
	cfg := config.LoadConfig()

	dbURL := cfg.DBURL
	// golang-migrate expects postgres:// scheme, but Neon/Postgres often provide postgresql://
	if len(dbURL) >= 13 && dbURL[:13] == "postgresql://" {
		dbURL = "postgres://" + dbURL[13:]
	}

	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		log.Fatalf("could not create migrate instance: %v", err)
	}

	if len(os.Args) < 2 {
		usage()
		return
	}

	cmd := os.Args[1]
	switch cmd {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("could not run up migrations: %v", err)
		}
		fmt.Println("Migrations applied successfully!")
	case "down":
		if err := m.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("could not run down migrations: %v", err)
		}
		fmt.Println("Last migration reverted successfully!")
	case "force":
		if len(os.Args) < 3 {
			log.Fatal("force command requires a version number")
		}
		version, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("invalid version number: %v", err)
		}
		if err := m.Force(version); err != nil {
			log.Fatalf("could not force version: %v", err)
		}
		fmt.Printf("Forced version to %d\n", version)
	default:
		usage()
	}
}

func usage() {
	fmt.Println("Usage: go run cmd/migrate/main.go [up|down|force <version>]")
}
