package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"personal-manager/internal/handler"
	"personal-manager/internal/service"
	"personal-manager/internal/store"
)

const defaultPort = 8080

func main() {
	port, err := parsePort(os.Args[1:], os.Stderr)
	if err != nil {
		log.Fatal(err)
	}

	dbPath := os.Getenv("PERSONAL_MANAGER_DB")
	if dbPath == "" {
		dbPath = "personal_manager.db"
	}

	addr := fmt.Sprintf(":%d", port)

	st, err := store.Open(dbPath)
	if err != nil {
		log.Fatalf("open store: %v", err)
	}
	defer st.Close()

	svc := service.New(st)
	router := handler.New(svc).Routes()

	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}

func parsePort(args []string, output io.Writer) (int, error) {
	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	fs.SetOutput(output)

	port := fs.Int("port", defaultPort, "external service port")
	if err := fs.Parse(args); err != nil {
		return 0, err
	}

	portSet := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == "port" {
			portSet = true
		}
	})

	switch fs.NArg() {
	case 0:
		return validatePort(*port)
	case 1:
		if portSet {
			return 0, fmt.Errorf("port must be provided either as -port or as an argument, not both")
		}
		positionalPort, err := strconv.Atoi(fs.Arg(0))
		if err != nil {
			return 0, fmt.Errorf("port must be a number")
		}
		return validatePort(positionalPort)
	default:
		return 0, fmt.Errorf("unexpected arguments: %s", strings.Join(fs.Args(), " "))
	}
}

func validatePort(port int) (int, error) {
	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("port must be between 1 and 65535")
	}
	return port, nil
}
