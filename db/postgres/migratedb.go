package postgres

import (
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateDB() {
	migration, err := migrate.New("file://migrations/postgres", DBConnectionString()+"?sslmode=disable")
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
