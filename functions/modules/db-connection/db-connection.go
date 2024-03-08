package dbConnection

import (
	"database/sql"
	"log"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func GetDb() *bun.DB {
	dsn := "postgres://" + os.Getenv("PG_USER") + ":" + os.Getenv("PG_PASSWORD") + "@" + os.Getenv(("PG_HOST")) + ":5432/" + os.Getenv("PG_DATABASE")
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	log.Println(dsn)
	return bun.NewDB(sqldb, pgdialect.New())
}