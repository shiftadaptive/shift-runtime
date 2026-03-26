// SHIFT ::: Runtime
// Lightweight adaptive middleware for API compatibility
// (c) 2026 ShiftAdaptive

package main

import (
	"log/slog"
	"net/http"
	"os"

	"shift/internal/handler"
	"shift/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("Error loading .env file, continuing with environment variables")
	}

	logger.Init()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		handler.HandleRequest(w, r)
	})

	slog.Info("SHIFT running on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
