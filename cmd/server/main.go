package main

import (
	"log"
	"net/http"
	"os"

	"personal-manager/internal/handler"
	"personal-manager/internal/service"
	"personal-manager/internal/store"
)

func main() {
	dbPath := os.Getenv("PERSONAL_MANAGER_DB")
	if dbPath == "" {
		dbPath = "personal_manager.db"
	}

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

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
