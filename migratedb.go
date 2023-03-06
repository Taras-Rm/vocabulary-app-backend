package main

import (
	"fmt"
	"os"
	"vacabulary/db/postgres"

	"github.com/golang-migrate/migrate/v4"
)

func migrateDB() {
	migration, err := migrate.New("file://migrations/postgres", postgres.DBConnectionString()+"?sslmode=disable")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = migration.Up()
	if err != nil && err != migrate.ErrNoChange {
		fmt.Println(err.Error())
		os.Exit(1)
	} else if err == migrate.ErrNoChange {
		fmt.Println("DB schema is up to date")
	} else {
		fmt.Println("Success migrated DB schema")
	}
}
