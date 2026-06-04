package models

import "time"

// ApiResponse represents the standard API response structure
type ApiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error,omitempty"`
	Time    int64       `json:"time,omitempty"`
}

// LogEntry represents a log entry for streaming output
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
	Content   string `json:"content"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// PingResponse represents ping response
type PingResponse struct {
	Message string `json:"message"`
}

// HelloResponse represents hello response
type HelloResponse struct {
	Message string `json:"message"`
}

// ShellRequest represents shell command request
type ShellRequest struct {
	Shell string `json:"shell"`
	Type  string `json:"type"`
}

// ShellResponse represents shell command response
type ShellResponse struct {
	Success bool   `json:"success"`
	Output  string `json:"output,omitempty"`
	Error   string `json:"error,omitempty"`
}

// FileOperationRequest represents file operation request
type FileOperationRequest struct {
	Operation string      `json:"operation"`
	Path      string      `json:"path"`
	Content   string      `json:"content"`
	Key       string      `json:"key"`
	Value     interface{} `json:"value"`
}

// FileOperationResponse represents file operation response
type FileOperationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// SysInfoResponse represents system info response
type SysInfoResponse struct {
	Arch string `json:"arch"`
	OS   string `json:"os"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// NewApiResponse creates a new API response
func NewApiResponse(success bool, data interface{}, err string) *ApiResponse {
	return &ApiResponse{
		Success: success,
		Data:    data,
		Error:   err,
		Time:    time.Now().UnixMilli(),
	}
}

// NewLogEntry creates a new log entry
func NewLogEntry(logType, content string) *LogEntry {
	return &LogEntry{
		Timestamp: time.Now().Format("2006-01-02 15:04:05.000"),
		Type:      logType,
		Content:   content,
	}
}
