package helpers

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

type Flags struct {
	Port       int
	Directory  string
	EnableAuth bool
	Username   string
	Password   string
	AbsolutePath string
}

func GetFlags() *Flags {
	port := flag.Int("port", 8080, "Port to serve on")
	directory := flag.String("dir", ".", "Directory to serve files from")
	enableAuth := flag.Bool("auth", false, "Enable basic authentication")
	username := flag.String("user", "", "Username for basic authentication")
	password := flag.String("pass", "", "Password for basic authentication")
	flag.Parse()

	if *enableAuth && (*username == "" || *password == "") {
		log.Fatal("'-auth' enabled, '-user' and '-pass' must be set")
	}

	absPath, err := filepath.Abs(*directory)
	if err != nil {
		log.Fatalf("Error getting absolute path: %v", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		log.Fatalf("Directory does not exist: %s", absPath)
	}

	return &Flags{
		Port:       *port,
		Directory:  *directory,
		EnableAuth: *enableAuth,
		Username:   *username,
		Password:   *password,
		AbsolutePath: absPath,
	}
}
