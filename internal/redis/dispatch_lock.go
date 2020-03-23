package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"

	"demo/pkg/enum"
)

const dispatchLockPrefix = "dispatch_lock:%d"

type DispatchLock struct {
	client *redis.Client
}

func NewDispatchLock(client *redis.Client) *DispatchLock {
	return &DispatchLock{client: client}
}

// Lock
func (s *DispatchLock) Lock(id int, expiration time.Duration) (bool, error) {
	key := fmt.Sprintf(dispatchLockPrefix, id)
	return s.client.SetNX(key, time.Now().Format(enum.DateTimeFormat), expiration).Result()
}
