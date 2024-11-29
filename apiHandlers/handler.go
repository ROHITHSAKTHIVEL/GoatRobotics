package apihandlers

import (
	"encoding/json"
	"net/http"

	"github.com/ROHITHSAKTHIVEL/GoatRobotics/models"
	"github.com/ROHITHSAKTHIVEL/GoatRobotics/service"
)

type ChatHandler struct {
	Clients *service.Clients
}

// Constructor for ChatHandler
func NewChatHandler(clients *service.Clients) *ChatHandler {
	return &ChatHandler{Clients: clients}
}

// Ping handler for health check
func Ping(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "Pinged Successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// JoinHandler handles client joining
func (h *ChatHandler) JoinHandler(w http.ResponseWriter, r *http.Request) {
	// Handle timeout via context
	select {
	case <-r.Context().Done():
		http.Error(w, "Request timed out", http.StatusRequestTimeout)
		return
	default:
	}

	clientID := r.URL.Query().Get("id")
	if clientID == "" {
		http.Error(w, "Client ID is required", http.StatusBadRequest)
		return
	}

	// Attempt to join client
	if err := h.Clients.JoinClient(clientID); err != nil {
		http.Error(w, err.Error(), http.StatusConflict) // Handle if the client already exists
		return
	}

	// Successful response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User joined successfully"})
}

// LeaveHandler handles client leaving
func (h *ChatHandler) LeaveHandler(w http.ResponseWriter, r *http.Request) {
	// Handle timeout via context
	select {
	case <-r.Context().Done():
		http.Error(w, "Request timed out", http.StatusRequestTimeout)
		return
	default:
	}

	clientID := r.URL.Query().Get("id")
	if clientID == "" {
		http.Error(w, "Client ID is required", http.StatusBadRequest)
		return
	}

	// Attempt to remove client
	if err := h.Clients.LeaveClient(clientID); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // Handle if the client doesn't exist
		return
	}

	// Successful response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User left successfully"})
}

// SendMessageHandler handles sending messages
func (h *ChatHandler) SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Handle timeout via context
	select {
	case <-r.Context().Done():
		http.Error(w, "Request timed out", http.StatusRequestTimeout)
		return
	default:
	}

	clientID := r.URL.Query().Get("id")
	message := r.URL.Query().Get("message")

	// Validate inputs
	if clientID == "" {
		http.Error(w, "Client ID is required", http.StatusBadRequest)
		return
	}
	if message == "" {
		http.Error(w, "Message content is required", http.StatusBadRequest)
		return
	}

	// Construct message model
	chatMessage := models.Chat{
		ClientID: clientID,
		Message:  message,
	}

	// Attempt to send message
	if err := h.Clients.SendMessage(chatMessage); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // Handle if the client doesn't exist
		return
	}

	// Successful response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Message sent successfully"})
}

func (h *ChatHandler) GetMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Handle timeout via context
	select {
	case <-r.Context().Done():
		http.Error(w, "Request timed out", http.StatusRequestTimeout)
		return
	default:
	}

	clientID := r.URL.Query().Get("id")
	if clientID == "" {
		http.Error(w, "Client ID is required", http.StatusBadRequest)
		return
	}

	// Fetch messages for the client
	response, err := h.Clients.GetMessage(clientID)
	if err != nil {
		// If an error occurred in fetching messages
		http.Error(w, err.Error(), http.StatusNotFound) // Handle case if the client doesn't exist
		return
	}

	// Successful response with the messages
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
