package main

import (
	"fmt"
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"sync"
	"time"
)

func main() {
	// Create a pool with go-redis (or redigo) which is the pool redisync will
	// use while communicating with Redis. This can also be any pool that
	// implements the `redis.Pool` interface.
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "192.168.2.112:6379",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	rs := redsync.New(pool)

	mutexname := "421"

	var gnum = 10
	var wg sync.WaitGroup
	wg.Add(gnum)
	for i := 0; i < gnum; i++ {
		go func() {
			defer wg.Done()
			// Obtain a new mutex by using the same name for all instances wanting the
			// same lock.
			mutex := rs.NewMutex(mutexname)
			fmt.Println("start getting lock")
			if err := mutex.Lock(); err != nil {
				panic(err)
			}
			fmt.Println("get the lock")
			time.Sleep(time.Second)
			if ok, err := mutex.Unlock(); !ok || err != nil {
				panic("unlock failed")
			}
			fmt.Println("unlock")
		}()
	}
	wg.Wait()
	// Obtain a lock for our given mutex. After this is successful, no one else
	// can obtain the same lock (the same mutex name) until we unlock it.

	// Do your work that requires the lock.

	// Release the lock so other processes or threads can obtain a lock.

}
