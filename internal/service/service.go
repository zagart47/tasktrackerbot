package service

import "tasktrackerbot/internal/usecase"

type Services struct {
	Tasks Tasks
}

func NewServices(usecases usecase.Usecases) Services {
	taskService := NewTaskService(usecases)
	return Services{
		Tasks: taskService,
	}
}
