package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	port := flag.Int("port", 8080, "Port to serve on")
	directory := flag.String("dir", ".", "Directory to serve files from")
	flag.Parse()

	absPath, err := filepath.Abs(*directory)
	if err != nil {
		log.Fatalf("Error getting absolute path: %v", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		log.Fatalf("Directory does not exist: %s", absPath)
	}

	fileServer := http.FileServer(http.Dir(absPath))
	http.Handle("/", http.StripPrefix("/", fileServer))
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		log.Println("hello")
	})

	// Get the local IP address
	localIP := GetOutboundIP()

	log.Printf("Starting server at http://%s:%d", localIP, *port)
	log.Printf("Serving files from: %s", absPath)
	log.Printf("Access on your network via: http://%s:%d", localIP, *port)

	addr := fmt.Sprintf(":%d", *port)
	if err := http.ListenAndServe(addr, nil); err != nil {
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