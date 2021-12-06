package types

import (
	"encoding/json"
)

// APIError is one of the following
// 400 Bad Request: The data of your request is invalid or incomplete.
// 401 Unauthorized: You are not correctly authenticated.
// 403 Forbidden: You are authenticated, but you are not authorized.
// 404 Not Found: The resource you try to access does not exist.
// 415 Unsupported Media Type: Invalid Content-Type or Accept HTTP header.
// 417 Expectation Failed: Your token quota has been exhausted.
// 422 Unprocessable Entity: Validation error (missing required attributes, number not in range etc.)
// 423 Locked: The API is locked for write operations during maintenance.
// 429 Too Many Requests: You have been rate limited.
// 500 Internal Server Error: Something wrong happened in the server
// 502 Bad Gateway, 503 Service Unavailable, 504 Gateway Timeout: Temporary communication failure.
// There may be additional error codes in certain circumstances. Please refer to the HTTP error code documentation for more information.

// DataAPIError is the outer wrapper of an error message returned by the Prisma API
type DataAPIError []struct {
	Code int `json:"code"`
	Data struct {
		Attribute string `json:"attribute"`
	} `json:"data"`
	Description string `json:"description"`
	Subject     string `json:"subject"`
	Title       string `json:"title"`
	Trace       string `json:"trace"`
}

// APIError Is an API Error
type APIError struct {
	Code int `json:"code"`
	Data struct {
		Attribute string `json:"attribute"`
	} `json:"data"`
	Description string `json:"description"`
	Subject     string `json:"subject"`
	Title       string `json:"title"`
	Trace       string `json:"trace"`
}

// IsForbiddenOrUnauthorized returns true if the error is an Forbidden or Unauthorized
func (m *APIError) IsForbiddenOrUnauthorized() bool {
	switch m.Code {
	case 401, 403:
		return true
	}
	return false

}

// IsNotFound returns true if the error is because the entity was not found
func (m *APIError) IsNotFound() bool {
	return m.Code == 404
}

func (m *APIError) Error() string {
	return string(m.Description)
}

// NewAPIError returns a new APIError from a byte slice
func NewAPIError(input []byte) *APIError {
	var raw *DataAPIError
	json.Unmarshal(input, &raw)

	e := &APIError{}

	for _, v := range *raw {
		e.Code = v.Code
		e.Description = v.Description
		e.Subject = v.Subject
		e.Title = v.Title
		e.Trace = v.Trace
	}

	return e
}

// TokenExpiredError is an Token Expired Error
type TokenExpiredError struct {
	Message string
}

func (m *TokenExpiredError) Error() string {
	return m.Message
}

// NewTokenExpiredError returns a new TokenExpiredError
func NewTokenExpiredError() *TokenExpiredError {
	return &TokenExpiredError{
		Message: "Token is expired",
	}
}
