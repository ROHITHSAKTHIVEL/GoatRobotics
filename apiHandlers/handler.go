package apihandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ROHITHSAKTHIVEL/GoatRobotics/logs"
	"github.com/ROHITHSAKTHIVEL/GoatRobotics/models"
	"github.com/ROHITHSAKTHIVEL/GoatRobotics/service"
)

type ChatHandler struct {
	Clients *service.Clients
}


func NewChatHandler(clients *service.Clients) *ChatHandler {
	return &ChatHandler{Clients: clients}
}


func Ping(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "Pinged Successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func Logs(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "Pinged Successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// JoinHandler handles client joining
func (h *ChatHandler) JoinHandler(w http.ResponseWriter, r *http.Request) {
	
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

	
	if err := h.Clients.JoinClient(clientID); err != nil {
		http.Error(w, err.Error(), http.StatusConflict) 
		return
	}

	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User joined successfully"})
}


func (h *ChatHandler) LeaveHandler(w http.ResponseWriter, r *http.Request) {
	
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

	
	if err := h.Clients.LeaveClient(clientID); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) 
		return
	}

	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User left successfully"})
}


func (h *ChatHandler) SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	
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

	
	chatMessage := models.Chat{
		ClientID: clientID,
		Message:  message,
	}

	
	if err := h.Clients.SendMessage(chatMessage); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Message sent successfully"})
}

func (h *ChatHandler) GetMessageHandler(w http.ResponseWriter, r *http.Request) {
	
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

	
	response, err := h.Clients.GetMessage(clientID)
	if err != nil {
		
		http.Error(w, err.Error(), http.StatusNotFound) 
		return
	}

	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetLogHandler(w http.ResponseWriter, r *http.Request) {
	
	clientID := r.URL.Query().Get("clientID")
	
	startTimeStr := r.URL.Query().Get("startTime")

	// If neither clientID nor startTime is provided, return all logs
	if clientID == "" && startTimeStr == "" {
		logs, err := logs.GetAllLogs() // Fetch all logs
		if err != nil {
			http.Error(w, "Error fetching logs: "+err.Error(), http.StatusInternalServerError)
			return
		}

		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(logs); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding logs: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// If a startTime is provided, parse it
	var startTime time.Time
	if startTimeStr != "" {
		var err error
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid startTime format: %v", err), http.StatusBadRequest)
			return
		}
	}

	// Call the GetLog function to retrieve logs based on clientID and startTime
	logs, err := logs.GetLog(clientID, startTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(logs); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding logs: %v", err), http.StatusInternalServerError)
	}
}
