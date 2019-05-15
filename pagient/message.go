package pagient

// Message represents a standard response.
type Message struct {
	StatusCode  int    `json:"status"`
	Message     string `json:"message"`
	ErrorText   string `json:"error"`
}
