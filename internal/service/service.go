package service

import (
	"tasktrackerbot/internal/repository"
	"tasktrackerbot/internal/usecase"
)

type Services struct {
	usecases usecase.Usecases
	Tasks    Tasks
}

func NewServices(repo repository.Repositories) Services {
	usecases := usecase.NewUsecases(repo)
	services := NewTaskService(usecases)
	return Services{
		Tasks: services,
	}
}
