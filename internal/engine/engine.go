// SHIFT ::: Runtime
// Lightweight adaptive middleware for API compatibility
// (c) 2026 ShiftAdaptive

package engine

import (
	"errors"
	"fmt"
	"shift/internal/models"
	"shift/internal/utils"

	"github.com/go-resty/resty/v2"
)

func ProcessRequest(req models.Request) (string, error) {
	client := resty.New()

	resp, err := client.R().
		SetQueryParams(utils.ConvertParams(req.Params)).
		SetBody(req.Body).
		Execute(req.Method, req.Target)

	if err != nil {
		return "", err
	}

	if resp.StatusCode() >= 400 {
		return handleFailure(req, resp.String())
	}

	return resp.String(), nil
}

func handleFailure(req models.Request, errorBody string) (string, error) {
	fmt.Println("SHIFT detected failure")
	fmt.Println("Error:", errorBody)

	// CALL Python agent

	return "", errors.New("external API request failed")
}