// routes.go
package routers

import (
	"net/http"

	apiHandlers "github.com/ROHITHSAKTHIVEL/GoatRobotics/apiHandlers"
	"github.com/ROHITHSAKTHIVEL/GoatRobotics/service"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func SetupRoutes(clients *service.Clients) http.Handler {
	// Create a new mux router
	router := mux.NewRouter()

	// Initialize the ChatHandler
	handler := apiHandlers.NewChatHandler(clients)

	// Initialize middleware for rate limiting
	commonMiddleware := apiHandlers.NewMiddleware(10, 20)

	// Define the routes with middlewares
	router.HandleFunc("/join", commonMiddleware.MiddleWare(handler.JoinHandler)).Methods("POST", "GET")
	router.HandleFunc("/leave", commonMiddleware.MiddleWare(handler.LeaveHandler)).Methods("POST", "GET")
	router.HandleFunc("/send", commonMiddleware.MiddleWare(handler.SendMessageHandler)).Methods("POST", "GET")
	router.HandleFunc("/messages", commonMiddleware.MiddleWare(handler.GetMessageHandler)).Methods("GET")
	router.HandleFunc("/ping", apiHandlers.Ping).Methods("GET")

	// Allow CORS for all routes
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},                             // You can restrict origins here if needed
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},  // Allowed methods
		AllowedHeaders:   []string{"Content-Type", "Authorization"}, // Allowed headers
		AllowCredentials: true,
	})

	// Apply CORS middleware
	return corsMiddleware.Handler(router) // Returns http.Handler
}
