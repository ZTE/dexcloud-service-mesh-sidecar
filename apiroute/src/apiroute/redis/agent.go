package redis

import (
	"apiroute/managers/configmanager"
	"github.com/go-redis/redis"
	"strconv"
)

var (
	IP       = configmanager.CfgM.GetConfigInfo().Redis.Host
	port     = strconv.Itoa(configmanager.CfgM.GetConfigInfo().Redis.Port)
	poolSize = configmanager.CfgM.GetConfigInfo().Redis.Pool.MaxTotal
)

const (
	routerDB = 0
	tempDB   = 1
	mapperDB = 2
)

type Clients struct {
	RouterClient *redis.Client
	TempClient   *redis.Client
	MapperClient *redis.Client
}

func NewClients() *Clients {
	return &Clients{
		RouterClient: redis.NewClient(&redis.Options{
			Addr:     IP + ":" + port,
			PoolSize: poolSize,
			DB:       routerDB,
		}),

		TempClient: redis.NewClient(&redis.Options{
			Addr:     IP + ":" + port,
			PoolSize: poolSize,
			DB:       tempDB,
		}),

		MapperClient: redis.NewClient(&redis.Options{
			Addr:     IP + ":" + port,
			PoolSize: poolSize,
			DB:       mapperDB,
		}),
	}
}
