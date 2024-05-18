package main

import (
	"sync"

	"github.com/go-concurrence-treat/wrong-way/client"
	"github.com/go-concurrence-treat/wrong-way/redis"
)

var waitGroup sync.WaitGroup

func main() {
	adapter := redis.NewRedisAdapter()
	service := redis.NewRedisLockingService(adapter)

	client_1 := client.NewClient(*service, "1")
	client_2 := client.NewClient(*service, "2")
	client_3 := client.NewClient(*service, "3")
	client_4 := client.NewClient(*service, "4")

	waitGroup.Add(4)

	go client_1.RequestResource()
	go client_2.RequestResource()
	go client_3.RequestResource()
	go client_4.RequestResource()

	waitGroup.Wait()
}
