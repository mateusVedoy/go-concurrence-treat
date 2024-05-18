package client

import (
	"fmt"
	"time"

	"github.com/go-concurrence-treat/right-way/redis"
)

type Client struct {
	service redis.RedisAdapter
	ID      string
}

func NewClient(
	service redis.RedisAdapter,
	ID string,
) *Client {

	return &Client{
		service: service,
		ID:      ID,
	}
}

func (client *Client) RequestResource() {
	client.accessResource()
	client.print(
		fmt.Sprintf("\nThe client %s is in possession of the resource\n", client.ID),
	)
	delay()
	client.releaseResource()
}

func (client *Client) moment() {
	fmt.Print("\nTS: " + time.Now().Format("2006-01-02T15:04:05"))
}

func (client *Client) print(message string) {
	client.moment()
	fmt.Print(message)
}

func (client *Client) accessResource() {

	ATTEMPTS := 1

	for ATTEMPTS <= 3 {

		tryAccess := client.service.TryAccess()

		if tryAccess != nil {
			client.print(
				fmt.Sprintf("\nAttempt %d .Resource not available for client %s. Reason: %s\n", ATTEMPTS, client.ID, tryAccess.Error()),
			)
			ATTEMPTS++
			delay()
			continue
		}

		tryLock := client.service.Lock()

		if tryLock != nil {
			client.print(
				fmt.Sprintf("\nAttempt %d. Resource not locked for client %s. Reason: %s\n", ATTEMPTS, client.ID, tryLock.Error()),
			)
			ATTEMPTS++
			delay()
			continue
		}
		client.print(
			fmt.Sprintf("\nAttempt %d. Client %s locked resource successfully\n", ATTEMPTS, client.ID),
		)
		break
	}
}

func (client *Client) releaseResource() {

	unlockErr := client.service.Unlock()

	if unlockErr != nil {
		client.print(
			fmt.Sprintf("\nResource not unlock for client %s. Reason: %s\n", client.ID, unlockErr.Error()),
		)
	}
	client.print(
		fmt.Sprintf("\nResource released for client %s\n", client.ID),
	)
}

func delay() {
	time.Sleep(time.Second * 5)
}
