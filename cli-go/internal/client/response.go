package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	HTTPResponse *http.Response
	RequestURL   string
}

func (r *Response) JSON(v interface{}) error {
	defer r.HTTPResponse.Body.Close()

	body, err := readAllLimited(r.HTTPResponse.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if r.HTTPResponse.StatusCode >= 400 {
		msg := string(body)
		var apiMsg struct {
			Message string `json:"message"`
		}
		if json.Unmarshal(body, &apiMsg) == nil && apiMsg.Message != "" {
			msg = apiMsg.Message
		}
		return &APIError{
			StatusCode: r.HTTPResponse.StatusCode,
			Message:    msg,
			URL:        r.RequestURL,
		}
	}

	if v != nil && len(body) > 0 {
		if err := json.Unmarshal(body, v); err != nil {
			return fmt.Errorf("unmarshaling response: %w", err)
		}
	}

	return nil
}

func (r *Response) Error() error {
	defer r.HTTPResponse.Body.Close()

	if r.HTTPResponse.StatusCode >= 400 {
		body, err := readAllLimited(r.HTTPResponse.Body)
		if err != nil {
			return &APIError{
				StatusCode: r.HTTPResponse.StatusCode,
				Message:    fmt.Sprintf("failed to read error body: %v", err),
				URL:        r.RequestURL,
			}
		}
		msg := string(body)
		var apiMsg struct {
			Message string `json:"message"`
		}
		if json.Unmarshal(body, &apiMsg) == nil && apiMsg.Message != "" {
			msg = apiMsg.Message
		}
		return &APIError{
			StatusCode: r.HTTPResponse.StatusCode,
			Message:    msg,
			URL:        r.RequestURL,
		}
	}
	return nil
}

func (r *Response) StatusCode() int {
	return r.HTTPResponse.StatusCode
}

func (r *Response) RawBody() ([]byte, error) {
	defer r.HTTPResponse.Body.Close()
	return readAllLimited(r.HTTPResponse.Body)
}
