package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/Thrun12/golang-assignment/internal/config"
)

func main() {
	var (
		direction string
		steps     int
	)

	flag.StringVar(&direction, "direction", "up", "Migration direction: up or down")
	flag.IntVar(&steps, "steps", 0, "Number of migrations to apply (0 for all)")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Create migrator
	m, err := migrate.New(
		fmt.Sprintf("file://%s", cfg.MigrationPath),
		cfg.DatabaseURL,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create migrator: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		sourceErr, dbErr := m.Close()
		if sourceErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to close source: %v\n", sourceErr)
		}
		if dbErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to close database: %v\n", dbErr)
		}
	}()

	// Run migrations
	switch direction {
	case "up":
		if steps > 0 {
			err = m.Steps(steps)
		} else {
			err = m.Up()
		}
	case "down":
		if steps > 0 {
			err = m.Steps(-steps)
		} else {
			err = m.Down()
		}
	default:
		fmt.Fprintf(os.Stderr, "Invalid direction: %s (must be 'up' or 'down')\n", direction)
		os.Exit(1)
	}

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		fmt.Fprintf(os.Stderr, "Migration failed: %v\n", err)
		os.Exit(1)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("No migrations to apply")
	} else {
		fmt.Printf("Successfully applied migrations (%s)\n", direction)
	}
}
