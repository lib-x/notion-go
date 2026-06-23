package notion

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIError is returned for non-2xx Notion responses.
type APIError struct {
	StatusCode     int
	Status         string
	Code           string
	Message        string
	Object         string
	RequestID      string
	AdditionalData Object
	Header         http.Header
	Body           []byte
}

func (e *APIError) Error() string {
	switch {
	case e.Code != "" && e.Message != "":
		return fmt.Sprintf("notion: %s: %s", e.Code, e.Message)
	case e.Message != "":
		return fmt.Sprintf("notion: %s", e.Message)
	case e.Status != "":
		return fmt.Sprintf("notion: %s", e.Status)
	default:
		return fmt.Sprintf("notion: status %d", e.StatusCode)
	}
}

func decodeAPIError(resp *http.Response, body []byte) *APIError {
	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Header:     resp.Header.Clone(),
		Body:       append([]byte(nil), body...),
	}
	var payload struct {
		Object         string `json:"object"`
		Status         int    `json:"status"`
		Code           string `json:"code"`
		Message        string `json:"message"`
		RequestID      string `json:"request_id"`
		AdditionalData Object `json:"additional_data"`
	}
	if err := json.Unmarshal(body, &payload); err == nil {
		apiErr.Object = payload.Object
		apiErr.Code = payload.Code
		apiErr.Message = payload.Message
		apiErr.RequestID = payload.RequestID
		apiErr.AdditionalData = payload.AdditionalData
	}
	if apiErr.RequestID == "" {
		apiErr.RequestID = resp.Header.Get("X-Request-Id")
	}
	return apiErr
}
