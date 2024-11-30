package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	Error "github.com/ROHITHSAKTHIVEL/GoatRobotics/errors"
	"github.com/ROHITHSAKTHIVEL/GoatRobotics/models"
)

type Clients struct {
	Join    chan string      // Channel for joining clients
	Leave   chan string      // Channel for leaving clients
	Chat    chan models.Chat // Channel for chat messages
	Users   sync.Map         // Thread-safe map for connected users
	Message []*models.Chat   // Slice of messages for broadcasting
	Mu      sync.RWMutex     // Mutex for locking message slice
	Wg      sync.WaitGroup   // WaitGroup for managing goroutine completion
	done    chan struct{}    // Channel to signal shutdown
}

// Constructor for Clients
func NewClients() *Clients {
	return &Clients{
		Chat:  make(chan models.Chat, 100),
		Users: sync.Map{},
		Join:  make(chan string, 100),
		Leave: make(chan string, 100),
		done:  make(chan struct{}), // Initialize shutdown channel
	}
}

// IsClientExist checks if a client exists in the chat room (thread-safe)
func (c *Clients) IsClientExist(clientID string) bool {
	_, exists := c.Users.Load(clientID)

	return exists
}

// JoinClient handles the action of a client joining the chat room
func (c *Clients) JoinClient(clientID string) error {
	if c.IsClientExist(clientID) {
		log.Printf("Client %s already exists.\n", clientID)
		return Error.ClientAlreadyExist
	}

	c.Wg.Add(1)
	c.Join <- clientID
	return nil
}

// LeaveClient handles the action of a client leaving the chat room
func (c *Clients) LeaveClient(clientID string) error {
	if !c.IsClientExist(clientID) {
		log.Printf("Client %s does not exist.\n", clientID)
		return fmt.Errorf("client %s does not exist", clientID)
	}

	c.Wg.Add(1)
	c.Leave <- clientID
	return nil
}

// SendMessage handles sending a message from a client
func (c *Clients) SendMessage(message models.Chat) error {
	if !c.IsClientExist(message.ClientID) {
		log.Printf("Client %s does not exist. Message not sent.\n", message.ClientID)
		return fmt.Errorf("client %s does not exist", message.ClientID)
	}

	c.Wg.Add(1)
	c.Chat <- message
	return nil
}

// GetMessage retrieves all messages for a client
func (c *Clients) GetMessage(clientID string) (*models.MessagesResponse, error) {
	if !c.IsClientExist(clientID) {
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
		return nil, Error.RequestTimeout
	}
}

// GracefulShutdown shuts down all operations and releases resources
func (c *Clients) GracefulShutdown() {
	c.done <- struct{}{}
	close(c.done) // Signal shutdown
	c.Wg.Wait()   // Wait for all active operations to finish
	close(c.Join) // Close channels to release resources
	close(c.Leave)
	close(c.Chat)
	log.Println("All operations completed, shutdown successful.")
}

// Listener listens for client actions like join, leave, and send messages
func (c *Clients) Listener() {
	for {
		select {

		case <-c.done: // Explicit shutdown signal
			log.Println("Listener shutting down.")
			return

		case clientID := <-c.Join:

			if clientID != "" {
				c.Users.Store(clientID, true)
				log.Printf("Client %s has joined the chat.\n", clientID)
			}
			c.Wg.Done()

		case clientID := <-c.Leave:

			if clientID != "" {
				c.Users.Delete(clientID)
				log.Printf("Client %s has left the chat.\n", clientID)
			}
			c.Wg.Done()

		case message := <-c.Chat:

			message.SentTime = time.Now()
			c.Mu.Lock()
			c.Message = append(c.Message, &message)
			c.Mu.Unlock()
			log.Printf("Client %s sent a message: %s\n", message.ClientID, message.Message)
			c.Wg.Done()
		}
	}
}
