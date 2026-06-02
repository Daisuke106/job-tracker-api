package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"job-tracker-api/internal/db"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	database, err := db.NewPostgresDB()
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	defer database.Close()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("server started on :" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}