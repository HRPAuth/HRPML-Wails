package models

import "time"

// User represents a user in the database
type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email,omitempty"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Setting represents a setting in the database
type Setting struct {
	ID          int64     `json:"id"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Type        string    `json:"type"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Log represents a log entry in the database
type Log struct {
	ID        int64     `json:"id"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Source    string    `json:"source,omitempty"`
	Metadata  string    `json:"metadata,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Task represents a task in the database
type Task struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Status      string     `json:"status"`
	Priority    int        `json:"priority"`
	Progress    int        `json:"progress"`
	Result      string     `json:"result,omitempty"`
	Error       string     `json:"error,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// File represents a file record in the database
type File struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	MimeType  string    `json:"mime_type,omitempty"`
	Checksum  string    `json:"checksum,omitempty"`
	Metadata  string    `json:"metadata,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Task status constants
const (
	TaskStatusPending    = "pending"
	TaskStatusRunning    = "running"
	TaskStatusCompleted  = "completed"
	TaskStatusFailed     = "failed"
	TaskStatusCancelled  = "cancelled"
)

// Log level constants
const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

// User role constants
const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)
