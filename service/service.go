package service

import (
	"context"
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
	Users   map[string]bool
	Message []*models.Chat
	Mu      sync.RWMutex
	Wg      sync.WaitGroup
}

// Constructor for Clients
func NewClients() *Clients {
	return &Clients{
		Chat:  make(chan models.Chat, 100),
		Join:  make(chan string),
		Leave: make(chan string),
		Users: make(map[string]bool),
	}
}

// Check if a client exists (thread-safe)
func (c *Clients) isClientExist(client string) bool {

	c.Mu.RLock()
	defer c.Mu.RUnlock()
	
	_, exists := c.Users[client]
	
	return exists
}

// Handle client joining
func (c *Clients) JoinClient(clientID string) error {
	if c.isClientExist(clientID) {
		log.Printf("Client %s already exists.\n", clientID)
		return fmt.Errorf("client %s already exists", clientID)
	}

	c.Join <- clientID
	return nil
}

// Handle client leaving
func (c *Clients) LeaveClient(clientID string) error {
	if !c.isClientExist(clientID) {
		log.Printf("Client %s does not exist.\n", clientID)
		return fmt.Errorf("client %s does not exist", clientID)
	}
	c.Leave <- clientID
	return nil
}


// Handle sending a message
func (c *Clients) SendMessage(message models.Chat) error {
	if !c.isClientExist(message.ClientID) {
		log.Printf("Client %s does not exist. Message not sent.\n", message.ClientID)
		return fmt.Errorf("client %s does not exist", message.ClientID)
	}
	c.Chat <- message
	return nil
}

// Retrieve messages for a client
func (c *Clients) GetMessage(clientID string) (*models.MessagesResponse, error) {
	if !c.isClientExist(clientID) {
		log.Printf("Client %s does not exist.\n", clientID)
		return nil, fmt.Errorf("client %s does not exist", clientID)
	}

	// Lock the messages for thread-safe access
	c.Mu.RLock()
	messages := append([]*models.Chat{}, c.Message...)
	c.Mu.RUnlock()

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	responseChannel := make(chan models.MessagesResponse)

	go func() {
		response := models.MessagesResponse{
			Messages:     messages,
			ID:           clientID,
			ResponseTime: time.Now(),
		}
		if len(messages) == 0 {
			response.Message = "No new messages"
		}
		responseChannel <- response
	}()

	select {
	case response := <-responseChannel:
		return &response, nil
	case <-ctx.Done():
		log.Printf("Timeout while retrieving messages for client %s.\n", clientID)
		return nil, fmt.Errorf("request timed out")
	}
}

// Listener for handling join, leave, and chat events
func (c *Clients) Listener() {
	for {
		select {
		case clientID := <-c.Join:
			c.Mu.Lock()
			c.Users[clientID] = true
			c.Mu.Unlock()
			fmt.Printf("Client %s has joined the chat.\n", clientID)

		case clientID := <-c.Leave:
			c.Mu.Lock()
			delete(c.Users, clientID)
			c.Mu.Unlock()
			fmt.Printf("Client %s has left the chat.\n", clientID)

		case message := <-c.Chat:
			message.SentTime = time.Now()
			c.Mu.Lock()
			c.Message = append(c.Message, &message)
			c.Mu.Unlock()
			fmt.Printf("Client %s sent a message: %s\n", message.ClientID, message.Message)
		}
	}
}
