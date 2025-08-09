package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/khalidibnwalid/DirServe/internal/helpers"
	"github.com/khalidibnwalid/DirServe/internal/middlewares"
	"github.com/khalidibnwalid/DirServe/internal/templviews"
)

func main() {
	flags := helpers.GetFlags()

	fileServer := http.FileServer(http.Dir(flags.AbsolutePath))
	fileServerWithMiddlewares := middlewares.ApplyMiddlewares(fileServer, middlewares.LoggingMiddleware, middlewares.SecurityHeadersMiddleware)
	templWithMiddlewares := middlewares.ApplyMiddlewares(templviews.Handler(), middlewares.LoggingMiddleware)

	if flags.EnableAuth {
		fileServerWithMiddlewares = middlewares.BasicAuthMiddleware(fileServerWithMiddlewares, flags.Username, flags.Password)
		templWithMiddlewares = middlewares.BasicAuthMiddleware(templWithMiddlewares, flags.Username, flags.Password)
		log.Println("Basic authentication enabled")
	}

	// Register handlers
	http.Handle("/raw/", http.StripPrefix("/raw/", fileServerWithMiddlewares))
	http.Handle("/web/", http.StripPrefix("/web/", templWithMiddlewares))

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "pong")
		log.Println("Ping received from", r.RemoteAddr)
	})
	// Get the local IP address for display purposes
	localIP := GetOutboundIP()

	log.Printf("Starting server at http://%s:%d", localIP, flags.Port)
	log.Printf("Serving files from: %s", flags.AbsolutePath)

	if flags.EnableAuth {
		log.Printf("Authentication required: username=%s", flags.Username)
	}

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", flags.Port),
		Handler:           http.DefaultServeMux,
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
