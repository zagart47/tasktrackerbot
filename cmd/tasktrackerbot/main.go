package main

import (
	"context"
	"tasktrackerbot/config"
	"tasktrackerbot/internal/service"
	"tasktrackerbot/internal/storage"
	"tasktrackerbot/internal/storage/cache"
	"tasktrackerbot/internal/storage/postgresql"
	"tasktrackerbot/internal/transport"
	"tasktrackerbot/internal/transport/handler"
	"tasktrackerbot/internal/usecase"
	"tasktrackerbot/pkg/migration"
)

func main() {
	// Получаем конфигурацию
	cfg := config.Configs

	// Инициализация пула postgres
	dbPool := postgresql.New(cfg.Postgres.DSN)

	// Инициализация стореджа
	storages := storage.NewStorages(dbPool)

	// Создание схемы в БД
	migration.Do(cfg.MigrationsPath, cfg.Postgres.DSN)

	// Инициализация кеша
	mc := cache.NewRedisClient(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Pwd)

	// Инициализация юзкейсов
	usecases := usecase.NewUsecases(storages, &mc)

	// Инициализация сервисов
	services := service.NewServices(usecases)

	// Заполняем кэш данными из БД
	services.Tasks.MakeTasksCache(context.Background())

	// Инициализация транспорта
	bot := transport.NewBotService(cfg.Bot.Token, services)

	// Инициализация хендлеров для бота
	handlers := handler.NewHandler(bot)
	handlers.InitHandlers()

	// Запуск бота
	bot.Start()
}
