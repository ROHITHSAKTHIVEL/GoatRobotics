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

func NewChatHandler(clients *service.Clients) *ChatHandler {
	return &ChatHandler{Clients: clients}
}

func Ping(w http.ResponseWriter, r *http.Request) {
	response := "Pinged Sucessfully"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

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

	h.Clients.JoinClient(clientID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User joined successfully"))
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

	h.Clients.LeaveClient(clientID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User left successfully"))
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

	if clientID == "" || message == "" {
		http.Error(w, "Client ID and message are required", http.StatusBadRequest)
		return
	}
	mes := models.Chat{
		ClientID: clientID,
		Message:  message,
	}
	h.Clients.SendMessage(mes)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message sent successfully"))
}
