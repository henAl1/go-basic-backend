package main

import (
	"crypto/tls"
	"log"
	"myapp/config"
	"myapp/internal/handler"
	"myapp/internal/middleware"
	"myapp/internal/repository"
	"myapp/internal/service"
	"myapp/pkg/db"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func setupTLS() (*tls.Config, error) {

	certPath := os.Getenv("TLS_CERT_PATH")
	keyPath := os.Getenv("TLS_KEY_PATH")

	if certPath == "" || keyPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		certPath = filepath.Join(cwd, "cert", "server.crt")
		keyPath = filepath.Join(cwd, "cert", "server.key")
		log.Println("certPath = ", certPath)
		log.Println("keyPath = ", keyPath)
	}

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates:             []tls.Certificate{cert},
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
	}, nil
}

func main() {

	config := config.LoadConfig()
	log.Println(config)

	tlsConfig, err := setupTLS()
	if err != nil {
		log.Fatalf("Failed to setup TLS: %v", err)
	}

	server := &http.Server{
		Addr:         ":8443",
		TLSConfig:    tlsConfig,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		apiKey = "default-api-key" // For development only
		log.Println("Warning: Using default API key")
	}

	// Initialize database connection
	dbConn, err := db.NewPostgresConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Initialize repository
	userRepo := repository.NewUserRepository(dbConn)

	// Initialize service
	userService := service.NewUserService(userRepo)

	// Use the service
	userHandler := handler.NewUserHandler(userService)

	//auth := middleware.NewAPIKeyAuth(apiKey)

	authMiddleware := middleware.NewAuthMiddleware(config)

	//mux := http.NewServeMux()

	http.HandleFunc("/users", authMiddleware.RequireAPIKey(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			log.Println("reading all users.... ")
			userHandler.GetAllUsers(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}))

	http.HandleFunc("/users/", authMiddleware.RequireAPIKey(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.GetUser(w, r)
		case http.MethodPost:
			log.Println("Here.... ")
			userHandler.CreateUser(w, r)
		case http.MethodPut:
			userHandler.UpdateUser(w, r)
		case http.MethodDelete:
			userHandler.DeleteUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start the server
	go func() {
		log.Println("Starting HTTP server on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	// log.Println("Starting HTTPS server on :8443")
	// if err := http.ListenAndServeTLS(":8443", "server.crt", "server.key", nil); err != nil {
	// 	log.Fatalf("HTTPS server error: %v", err)
	// }

	log.Println("Starting HTTPS server on :8443")
	err = server.ListenAndServeTLS("", "") // Empty strings because cert is in TLSConfig
	if err != nil {
		log.Fatalf("HTTPS server error: %v", err)
	}
}
