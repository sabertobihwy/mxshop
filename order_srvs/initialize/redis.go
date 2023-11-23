package initialize

import (
	"fmt"
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"

	"mxshop_srvs/order_srvs/global"
)

func InitializeRedis(host string, port int) {
	client := goredislib.NewClient(&goredislib.Options{
		Addr: fmt.Sprintf("%s:%d", host, port),
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	global.Redsync = redsync.New(pool)
}
