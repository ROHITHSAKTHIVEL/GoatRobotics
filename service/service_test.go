package service_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	apihandlers "github.com/ROHITHSAKTHIVEL/GoatRobotics/apiHandlers"
	"github.com/ROHITHSAKTHIVEL/GoatRobotics/models"
	"github.com/ROHITHSAKTHIVEL/GoatRobotics/service"
)

var (
	JoinClientConcurrencyCount = 100
)

func TestPing(t *testing.T) {
	// Create a new instance of the Ping handler
	handler := http.HandlerFunc(apihandlers.Ping)

	// Create a new HTTP request to the ping endpoint
	req := httptest.NewRequest("GET", "/ping", nil)

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the handler with the request and recorder
	handler.ServeHTTP(rr, req)

	// Check that the status code is 200 OK
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %v, but got %v", http.StatusOK, rr.Code)
	}

	// Check that the response body is correct
	expected := map[string]string{"message": "Pinged Successfully"}
	var actual map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&actual); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if actual["message"] != expected["message"] {
		t.Errorf("Expected message %v, but got %v", expected["message"], actual["message"])
	}
}

// func TestJoinClient(t *testing.T){
// 	client:=service.NewClients()

// 	concurrentcalls:=100
// 	for i:=0; i<= concurrentcalls;i++{
//       go func(){

//		  }()
//		}
//	}
func TestJoinClientConcurrency(t *testing.T) {
	// Initialize the Clients object
	clients := service.NewClients()

	// Start the Listener in a separate goroutine to handle join, leave, and chat operations
	go clients.Listener()

	// Use a WaitGroup to synchronize goroutines
	var wg sync.WaitGroup

	id := make(chan string, JoinClientConcurrencyCount)
	wg.Add(2 * JoinClientConcurrencyCount) // Add to wait group count for both operations

	go func() {
		for i := 0; i < JoinClientConcurrencyCount; i++ {
			go func(i int) {
				defer wg.Done() // Decrement when done
				clientID := GenerateClientID(i)
				
				 
				if err := clients.JoinClient(clientID); err != nil {
					t.Errorf("Failed to join client %s: %v", clientID, err)
				}
				
				id <- clientID
			}(i)
		}
	}()

	go func() {
		for i := 0; i < JoinClientConcurrencyCount; i++ {
			go func() {
				defer wg.Done() // Decrement when done
				clientID := <-id
				
				if err := clients.LeaveClient(clientID); err != nil {
					t.Errorf("Failed to leave client %s: %v", clientID, err)
				}
				
			}()
		}
	}()

	meschan := make(chan models.Chat)
	go func() {
		for i := 0; i < JoinClientConcurrencyCount; i++ {
			go func() {
				defer wg.Done() // Decrement when done
				clientID := <-id
				cha := models.Chat{
					ClientID: clientID,
					Message:  "sample",
					SentTime: time.Now(),
				}
				
				
				if err := clients.SendMessage(cha); err != nil {
					t.Errorf("Failed to send measage %s: %v", clientID, err)
				}
				
				meschan <- cha
			}()
		}
	}()

	go func() {
		for i := 0; i < JoinClientConcurrencyCount; i++ {
			go func() {
				defer wg.Done() // Decrement when done
				val := <-meschan	
				
				msg, err := clients.GetMessage(val.ClientID)
				if err != nil {
					t.Errorf("Failed to get message %s: %v", val.ClientID, err)
				}

				fmt.Println(msg)
			}()
		}
	}()

	wg.Wait()

}

func GenerateClientID(index int) string {
	// Format the client ID to include the index and current timestamp (for uniqueness)
	// Example: "client-1-1631719375" where 1631719375 is the Unix timestamp
	timestamp := time.Now().Unix()
	clientID := fmt.Sprintf("client-%d-%d", index, timestamp)
	return clientID
}
