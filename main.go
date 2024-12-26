package main

import (
	"fmt"
	"log"
	"net/http"

	"apiLibgen/handlers"
)

func main() {
	// This handler responds to requests at the root URL ("/")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Server is up and running! Use /search or /download endpoints for more options.")
	})

	http.HandleFunc("/search", handlers.SearchBooks)
	http.HandleFunc("/download", handlers.DownloadBook)

	log.Println("Starting server on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
