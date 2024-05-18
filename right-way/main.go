package main

import (
	"sync"

	"github.com/go-concurrence-treat/right-way/client"
	"github.com/go-concurrence-treat/right-way/redis"
)

var waitGroup sync.WaitGroup

func main() {

	adapter := redis.NewRedisAdapter()

	client_1 := client.NewClient(*adapter, "1")
	client_2 := client.NewClient(*adapter, "2")
	client_3 := client.NewClient(*adapter, "3")
	client_4 := client.NewClient(*adapter, "4")

	waitGroup.Add(4)

	go client_1.RequestResource()
	go client_2.RequestResource()
	go client_3.RequestResource()
	go client_4.RequestResource()

	waitGroup.Wait()
}
