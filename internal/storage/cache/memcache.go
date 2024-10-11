package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"tasktrackerbot/internal/entity"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	redis redis.Client
}

type Cacher interface {
	Set(task entity.Task) error
	Get(key string) (entity.Task, error)
	GetAll(ctx context.Context) ([]entity.Task, error)
	Delete(key string) error
}

func NewRedisClient(addr, port, pwd string) Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", addr, port),
		Password: pwd,
		DB:       0,
	})
	return Redis{
		redis: *rdb,
	}
}

func (r *Redis) Set(task entity.Task) error {
	ctx := context.Background()
	taskJson, err := json.Marshal(task)
	if err != nil {
		return err
	}
	err = r.redis.Set(ctx, strconv.FormatInt(task.ID, 10), taskJson, 0).Err()
	return err
}

func (r *Redis) Get(key string) (entity.Task, error) {
	ctx := context.Background()
	taskJson, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		return entity.Task{}, err
	}
	var task entity.Task
	err = json.Unmarshal([]byte(taskJson), &task)
	return task, nil
}

func (r *Redis) GetAll(ctx context.Context) ([]entity.Task, error) {
	keys, err := r.redis.Keys(ctx, "*").Result()
	if err != nil {
		return nil, err
	}
	var tasks []entity.Task
	for _, key := range keys {
		taskData, err := r.redis.Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		var task entity.Task
		err = json.Unmarshal([]byte(taskData), &task)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *Redis) Delete(key string) error {
	ctx := context.Background()
	err := r.redis.Del(ctx, key).Err()
	return err
}
