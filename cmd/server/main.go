// SHIFT ::: Runtime
// Lightweight adaptive middleware for API compatibility
// (c) 2026 ShiftAdaptive

package main

import (
	"log"
	"net/http"

	"shift/internal/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		handler.HandleRequest(w, r)
	})

	log.Println("SHIFT running on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}