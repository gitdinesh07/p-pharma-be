package api

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Meta struct {
	RequestID string `json:"request_id,omitempty"`
}

type APIResponse[T any] struct {
	Success bool      `json:"success"`
	Message string    `json:"message,omitempty"`
	Data    *T        `json:"data,omitempty"`
	Error   *APIError `json:"error,omitempty"`
	Meta    *Meta     `json:"meta,omitempty"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type PaginatedResponse[T any] struct {
	Items      []T        `json:"items"`
	Pagination Pagination `json:"pagination"`
}

type APIRequest[T any] struct {
	RequestID string `json:"request_id,omitempty"`
	Payload   T      `json:"payload"`
}
