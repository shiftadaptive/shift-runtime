// SHIFT ::: Runtime
// Lightweight adaptive middleware for API compatibility
// (c) 2026 ShiftAdaptive

package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var L *slog.Logger

type BetterStackHandler struct {
	endpoint string
	token    string
}

func NewBetterStackHandler(endpoint, token string) *BetterStackHandler {
	return &BetterStackHandler{
		endpoint: endpoint,
		token:    token,
	}
}

func (h *BetterStackHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (h *BetterStackHandler) Handle(ctx context.Context, r slog.Record) error {
	payload := map[string]interface{}{
		"dt":      r.Time.UTC().Format("2006-01-02 15:04:05 UTC"),
		"message": r.Message,
		"level":   r.Level.String(),
	}

	r.Attrs(func(a slog.Attr) bool {
		payload[a.Key] = a.Value.Any()
		return true
	})

	go h.send(payload)
	return nil
}

func (h *BetterStackHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// For simplicity, we won't implement complex attribute chaining for the remote handler yet
	return h
}

func (h *BetterStackHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h *BetterStackHandler) send(payload map[string]interface{}) {
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", h.endpoint, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.token)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err == nil {
		resp.Body.Close()
	}
}

func Init() {
	endpoint := "https://s2322564.eu-fsn-3.betterstackdata.com"
	token := os.Getenv("BETTERSTACK_TOKEN")

	// Multi-handler: Console (JSON) + Better Stack
	consoleHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	
	if token != "" {
		remoteHandler := NewBetterStackHandler(endpoint, token)
		// We manually call both handlers to keep it simple without external dependencies
		L = slog.New(&multiHandler{handlers: []slog.Handler{consoleHandler, remoteHandler}})
	} else {
		L = slog.New(consoleHandler)
	}

	slog.SetDefault(L)
}

type multiHandler struct {
	handlers []slog.Handler
}

func (m *multiHandler) Enabled(ctx context.Context, l slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, l) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			_ = h.Handle(ctx, r)
		}
	}
	return nil
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{newHandlers}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithGroup(name)
	}
	return &multiHandler{newHandlers}
}
