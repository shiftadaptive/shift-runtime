// SHIFT ::: Runtime
// Lightweight adaptive middleware for API compatibility
// (c) 2026 ShiftAdaptive

package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"shift/internal/cache"
	"shift/internal/models"
	"shift/internal/utils"

	"github.com/go-resty/resty/v2"
)

type AgentRequest struct {
	Request   map[string]interface{} `json:"request"`
	Error     string                 `json:"error"`
	RequestID string                 `json:"requestId"`
}

type AgentResponse struct {
	Params map[string]interface{} `json:"params"`
}

func ProcessRequest(req models.Request, requestID string) (string, error) {
	applyCache(&req, requestID)

	client := resty.New()

	resp, err := client.R().
		SetQueryParams(utils.ConvertParams(req.Params)).
		SetBody(req.Body).
		Execute(req.Method, req.Target)

	if err != nil {
		return "", err
	}

	if resp.StatusCode() >= 400 {
		errorMessage := extractErrorMessage(resp.String())
		return handleFailure(req, errorMessage, requestID)
	}

	return resp.String(), nil
}

func handleFailure(req models.Request, errorMsg string, requestID string) (string, error) {
	slog.Warn(fmt.Sprintf("[%s] SHIFT detected failure", requestID), "error", errorMsg)

	// 🔥 Call Python agent
	correction, err := callAgent(req, errorMsg, requestID)
	if err != nil {
		return "", err
	}

	slog.Info(fmt.Sprintf("[%s] Agent response received", requestID), "params", correction.Params)

	storeMapping(req.Params, correction.Params)

	req.Params = correction.Params

	return retryRequest(req, requestID)
}

func callAgent(req models.Request, errorMsg string, requestID string) (*AgentResponse, error) {
	slog.Info(fmt.Sprintf("[%s] Calling agent for correction", requestID), "params", req.Params, "error", errorMsg)

	payload := AgentRequest{
		Request: map[string]interface{}{
			"params": req.Params,
		},
		Error:     errorMsg,
		RequestID: requestID,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(
		"http://localhost:8000/correct",
		"application/json",
		bytes.NewBuffer(body),
	)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result AgentResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func retryRequest(req models.Request, requestID string) (string, error) {
	slog.Info(fmt.Sprintf("[%s] Retrying request with corrected parameters", requestID), "params", req.Params)

	client := resty.New()

	resp, err := client.R().
		SetQueryParams(utils.ConvertParams(req.Params)).
		SetBody(req.Body).
		Execute(req.Method, req.Target)

	if err != nil {
		return "", err
	}

	return resp.String(), nil
}

func extractErrorMessage(body string) string {
	var parsed map[string]interface{}

	err := json.Unmarshal([]byte(body), &parsed)
	if err != nil {
		return body
	}

	if errObj, ok := parsed["error"].(map[string]interface{}); ok {
		if msg, ok := errObj["message"].(string); ok {
			return msg
		}
	}

	return body
}

func applyCache(req *models.Request, requestID string) {
	for key, value := range req.Params {
		if mappedKey, exists := cache.GetMapping(key); exists {
			slog.Info(fmt.Sprintf("[%s] Cache hit", requestID), "original", key, "mapped", mappedKey)
			req.Params[mappedKey] = value
			delete(req.Params, key)
		}
	}
}

func storeMapping(original map[string]interface{}, corrected map[string]interface{}) {
	for oldKey, oldVal := range original {
		for newKey, newVal := range corrected {
			// Only map if keys are different but values are the same (indicates a rename)
			if oldKey != newKey && fmt.Sprintf("%v", oldVal) == fmt.Sprintf("%v", newVal) {
				cache.StoreMapping(oldKey, newKey)
			}
		}
	}
}
