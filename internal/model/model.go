package model

// WebResponse is a generic response wrapper
type WebResponse[T any] struct {
	Data   T      `json:"data"`
	Errors string `json:"errors,omitempty"`
}
