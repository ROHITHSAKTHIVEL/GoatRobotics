package apihandlers

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/ROHITHSAKTHIVEL/GoatRobotics/models"
	"golang.org/x/time/rate"
)

type MiddleWare struct {
	rateLimitter map[string]*rate.Limiter
	mu           sync.Mutex
	limit        rate.Limit
	burst        int
}

func NewMiddleware(limit float64, burst int) *MiddleWare {
	return &MiddleWare{
		rateLimitter: make(map[string]*rate.Limiter),
		limit:        rate.Limit(limit),
		burst:        burst,
	}
}

// getLimiter will either return an existing limiter or create a new one for the given clientID
func (m *MiddleWare) getLimiter(clientID string) *rate.Limiter {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if the client already has a limiter
	limiter, exists := m.rateLimitter[clientID]
	if !exists {
		// If not, create a new limiter for the client
		limiter = rate.NewLimiter(m.limit, m.burst)
		m.rateLimitter[clientID] = limiter

		// Remove the limiter after 1 minute (client's session time)
		go func(clientID string) {
			time.Sleep(1 * time.Minute)
			m.mu.Lock()
			delete(m.rateLimitter, clientID)
			m.mu.Unlock()
		}(clientID)
	}

	return limiter
}

func (m *MiddleWare) MiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload models.Request
		var clientID string

		if r.Method == http.MethodPost {
			if err := json.NewDecoder(r.Body).Decode(&payload); err == nil && payload.ClientID != "" {
				clientID = payload.ClientID
			}
		} else {
			clientID = r.URL.Query().Get("id")
		}

		if clientID == "" {
			http.Error(w, "Client ID is required for rate limiting", http.StatusBadRequest)
			return
		}

		// Retrieve the rate limiter for the client
		limiter := m.getLimiter(clientID)

		// Check if the request is allowed by the rate limiter
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		// Create a context for the request (optional)
		ctx := context.Background()
		ctx1, _ := context.WithCancel(ctx)
		r = r.WithContext(ctx1)

		// Call the next handler in the chain
		next(w, r)
	}
}
