// SHIFT ::: Runtime
// Lightweight adaptive middleware for API compatibility
// (c) 2026 ShiftAdaptive

package engine

import (
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

	return resp.String(), nil
}