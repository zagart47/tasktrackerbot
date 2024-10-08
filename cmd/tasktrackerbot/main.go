package main

import (
	"log"
	"tasktrackerbot/config"
	"tasktrackerbot/internal/repository"
	"tasktrackerbot/internal/repository/postgresql"
	"tasktrackerbot/internal/service"
	"tasktrackerbot/internal/transport"
	"tasktrackerbot/internal/transport/handler"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Инициализация базы данных
	dbPool, err := postgresql.New(config.Configs.PostgreSQL.DSN)
	if err != nil {
		panic(err)
	}
	defer dbPool.Close()

	// Инициализация репозитория
	repos := repository.NewRepositories(dbPool)

	// Создание схемы в БД
	m, err := migrate.New("file://././internal/repository/postgresql/migrations",
		config.Configs.PostgreSQL.DSN)
	if err != nil {
		log.Println(err)
	}
	if err := m.Up(); err != nil {
		log.Println(err)
	}

	// Инициализация сервисов
	services := service.NewServices(repos)

	// Инициализация транспорта
	bot := transport.NewBotService(services)

	// Инициализация хендлеров для бота
	handlers := handler.NewHandler(bot)
	handlers.InitHandlers()

	// Запуск бота
	bot.Start()
}
