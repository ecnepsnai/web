package api

// Response describes an API response object
type Response struct {
	Code  int         `json:"code,omitempty"`
	Error Error       `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}
