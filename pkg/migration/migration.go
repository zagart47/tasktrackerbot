package migration

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Do(filepath, dsn string) {
	// Создание схемы в БД
	m, err := migrate.New(filepath, dsn)
	if err != nil {
		log.Println(err)
	}
	if err := m.Up(); err != nil {
		log.Println(err)
	}
}
