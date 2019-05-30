//go:generate go-bindata -pkg model -ignore .*\.go .
//go:generate go fmt .
package model

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/go_bindata"
	_ "github.com/lib/pq"
	"log"
)

var DB *sql.DB
var psql squirrel.StatementBuilderType

func init() {
	var err error
	DB, err = sql.Open("postgres", DatabaseConnection)
	if err != nil {
		panic(err)
	}
	source, err := bindata.WithInstance(bindata.Resource(AssetNames(), func(name string) ([]byte, error) {
		return Asset(name)
	}))
	if err != nil {
		panic(err)
	}

	driver, err := postgres.WithInstance(DB, &postgres.Config{
		MigrationsTable: "grabber_schema_migrations",
	})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithInstance("go-bindata", source, "postgres", driver)
	if err != nil {
		panic(err)
	}

	err = m.Up()
	if err == nil {
		log.Println("Migrated.")
	} else if err == migrate.ErrNoChange {
		log.Println("Migration not required.")
	} else {
		panic(err)
	}

	psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(DB)
}
