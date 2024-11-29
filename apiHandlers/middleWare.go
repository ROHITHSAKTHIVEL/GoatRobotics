package apihandlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/ROHITHSAKTHIVEL/GoatRobotics/logs"
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

// Middleware function for rate-limiting and logging
func (m *MiddleWare) MiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload models.Request
		var clientID string

		// Log start time
		startTime := time.Now()

		// Create an intermediate response recorder to capture the response details
		rec := &ResponseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		defer func() {
			// Log entry creation happens after the response is written
			endTime := time.Now()
			duration := endTime.Sub(startTime)

			// Prepare log entry
			logEntry := models.Logs{
				ClientID:         clientID,
				MethodName:       r.Method,
				RequestURL:       r.URL.String(),
				RequestBody:      getRequestBody(r, &payload),
				RequestHeaders:   formatHeaders(r.Header),
				QueryParameters:  formatQueryParameters(r.URL.Query()),
				ResponseDuration: duration,
				ResponseBody:     string(rec.Body),            // Log the full response body
				ResponseHeaders:  formatHeaders(rec.Header()), // Log response headers
				StatusCode:       rec.statusCode,
				StartTime:        startTime,
				EndTime:          endTime,
			}

			// Log the details
			logs.AddLog(logEntry)
		}()

		// Handle rate limiting and clientID extraction
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

		// Call the next handler (response is captured by rec)
		next(rec, r)
	}
}

// formatHeaders converts a map[string][]string (http.Header) into a string representation
func formatHeaders(headers http.Header) string {
	var result string
	for key, values := range headers {
		for _, value := range values {
			result += key + ": " + value + "\n"
		}
	}
	return result
}

// ResponseRecorder is a custom wrapper for capturing response details
type ResponseRecorder struct {
	http.ResponseWriter
	statusCode int
	Body       []byte
}

func (rec *ResponseRecorder) Write(p []byte) (n int, err error) {
	rec.Body = append(rec.Body, p...) // Append to Body instead of replacing
	return rec.ResponseWriter.Write(p)
}

func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.statusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

// getRequestBody retrieves the request body as a string (if POST request)
func getRequestBody(r *http.Request, payload *models.Request) string {
	var requestBody string
	if r.Method == http.MethodPost {
		// Read the request body
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return ""
		}
		// Store it for logging, and reset the body so it can be used later
		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		if err := json.NewDecoder(r.Body).Decode(payload); err == nil {
			body, err := json.Marshal(payload)
			if err == nil {
				requestBody = string(body)
			}
		}
		// Restore the body for subsequent use
		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}
	return requestBody
}

// formatQueryParameters converts a url.Values to a query string
func formatQueryParameters(values url.Values) string {
	return values.Encode() // Encodes the query parameters into a URL-encoded query string
}
