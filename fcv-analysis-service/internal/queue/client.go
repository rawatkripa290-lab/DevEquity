package queue

import (
	"github.com/hibiken/asynq"
)

func NewClient(addr, password string) *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{Addr: addr, Password: password})
}

func NewServer(addr, password string) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{Addr: addr, Password: password},
		asynq.Config{
			Concurrency: 5,
			Queues:      map[string]int{"critical": 3, "default": 2},
		},
	)
}
