package service

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ROHITHSAKTHIVEL/GoatRobotics/models"
)

type Clients struct {
	Join    chan string
	Leave   chan string
	Chat    chan models.Chat
	Users   map[string]interface{}
	Message []*models.Chat
	Mu      sync.Mutex
	Wg      sync.WaitGroup
}

func NewClients() *Clients {
	return &Clients{
		Chat:  make(chan models.Chat, 100),
		Join:  make(chan string),
		Leave: make(chan string),
		Users: make(map[string]interface{}),
	}

}

func clientExist(user map[string]interface{}, client string) bool {
	for clientID := range user {
		if clientID == client {
			return true
		}
	}
	return false
}

func (c *Clients) JoinClient(clientID string) {

	if !clientExist(c.Users, clientID) {
		log.Println("User Alreadt Exist")
		return
	}
	c.Join <- clientID
}

func (c *Clients) LeaveClient(clientID string) {

	if !clientExist(c.Users, clientID) {
		log.Println("User Alreadt Exist")
		return
	}
	c.Leave <- clientID
}

func (c *Clients) SendMessage(message models.Chat) {

	for clientID := range c.Users {
		if clientID == message.ClientID {
			c.Chat <- message
			break
		}
	}

	if !clientExist(c.Users, message.ClientID) {
		log.Println("User Alreadt Exist")
		return
	}

}

func (c *Clients) Listener() {
	for {
		select {
		case clientID := <-c.Join:
			c.Mu.Lock()
			c.Users[clientID] = true
			c.Mu.Unlock()
			fmt.Printf("%s has joined the chat.\n", clientID)

		case clientID := <-c.Leave:
			c.Mu.Lock()
			delete(c.Users, clientID)
			c.Mu.Unlock()
			fmt.Printf("%s has left the chat.\n", clientID)

		case message := <-c.Chat:
			message.SentTime = time.Now()
			c.Mu.Lock()
			fmt.Printf("%s has sent the chat %s.\n", message.ClientID, message.Message)
			c.Message = append(c.Message, &message)
			c.Mu.Unlock()

		default:
			time.After(1 * time.Millisecond)
		}
	}
}
