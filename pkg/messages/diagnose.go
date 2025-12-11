package messages

// DiagnoseMessage represents a diagnostic message sent to services
type DiagnoseMessage struct {
	TransactionID string `json:"transaction_id"`
	Operation     string `json:"operation"`
	Message       string `json:"message"`
}
