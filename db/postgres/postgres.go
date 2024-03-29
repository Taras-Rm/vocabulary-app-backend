package postgres

import (
	"fmt"
	"vacabulary/config"

	"github.com/go-pg/pg/v10"
)

func NewPostgres(cfg config.PostgresConfig) *pg.DB {
	address := fmt.Sprintf(cfg.Host+":%s", cfg.Port)

	conn := pg.Connect(&pg.Options{
		Database: cfg.Database,
		Addr:     address,
		User:     cfg.User,
		Password: cfg.Password,
	})

	return conn
}

func DBConnectionString() string {
	conn := config.Config.Postgres
	addr := conn.Host + ":" + conn.Port
	return fmt.Sprintf("postgres://%s:%s@%s/%s", conn.User, conn.Password, addr, conn.Database)
}
