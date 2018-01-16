package redis

import (
	"github.com/go-redis/redis"
)

const (
	count = 100
)

func FilterKeys(match string, client *redis.Client) ([]string, error) {
	var (
		cursor uint64
		out    []string
	)

	for {
		var keys []string
		var err error
		keys, cursor, err = client.Scan(cursor, match, count).Result()
		if err != nil {
			return nil, err
		}

		for _, v := range keys {
			out = append(out, v)
		}

		if cursor == 0 {
			break
		}
	}

	return out, nil
}
