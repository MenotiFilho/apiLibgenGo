package handlers

import (
	"apiLibgen/services"
	"encoding/json"
	"net/http"
)

// SearchBooks handles the search requests and responds with JSON-encoded books
func SearchBooks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Query parameter is missing", http.StatusBadRequest)
		return
	}

	searchTypeParam := r.URL.Query().Get("type")
	var searchType services.SearchType
	switch searchTypeParam {
	case "title":
		searchType = services.SearchByTitle
	case "author":
		searchType = services.SearchByAuthor
	case "isbn":
		searchType = services.SearchByISBN
	default:
		searchType = services.SearchByTitle // Default to search by title if no type is specified
	}

	books, err := services.SearchBooks(query, searchType)
	if err != nil {
		http.Error(w, "Failed to fetch search results", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}
