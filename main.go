package main

import (
	"flag"
	"fmt"
	"khalidibnwalid/dirserve/internal/middlewares"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	port := flag.Int("port", 8080, "Port to serve on")
	directory := flag.String("dir", ".", "Directory to serve files from")
	enableAuth := flag.Bool("auth", false, "Enable basic authentication")
	username := flag.String("user", "", "Username for basic authentication")
	password := flag.String("pass", "", "Password for basic authentication")
	flag.Parse()

	// Check if authentication is enabled but credentials are missing
	if *enableAuth && (*username == "" || *password == "") {
		log.Fatal("'-auth' enabled, '-user' and '-pass' must be set")
	}

	absPath, err := filepath.Abs(*directory)
	if err != nil {
		log.Fatalf("Error getting absolute path: %v", err)
	}

	// Verify directory exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		log.Fatalf("Directory does not exist: %s", absPath)
	}

	// Create a safe file server handler
	fileServer := http.FileServer(http.Dir(absPath))

	fileServerWithMiddlewares := middlewares.ApplyMiddlewares(fileServer, middlewares.LoggingMiddleware, middlewares.SecurityHeadersMiddleware)

	// Apply authentication if enabled
	if *enableAuth {
		fileServerWithMiddlewares = middlewares.BasicAuthMiddleware(fileServerWithMiddlewares, *username, *password)
		log.Println("Basic authentication enabled")
	}

	// Register handlers
	http.Handle("/raw/", http.StripPrefix("/raw/", fileServerWithMiddlewares))
	// http.Handle("/web/", http.StripPrefix("/web/", ""))

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "pong")
		log.Println("Ping received from", r.RemoteAddr)
	})
	// Get the local IP address for display purposes
	localIP := GetOutboundIP()

	log.Printf("Starting server at http://%s:%d", localIP, *port)
	log.Printf("Serving files from: %s", absPath)

	if *enableAuth {
		log.Printf("Authentication required: username=%s", *username)
	}

	server := &http.Server{
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// credit https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
