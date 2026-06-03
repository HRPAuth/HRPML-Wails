package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ApiResponse represents the standard API response structure
type ApiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error,omitempty"`
	Time    int64       `json:"time,omitempty"`
}

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// TestHealthCheck calls the backend health check API
func (a *App) TestHealthCheck() ApiResponse {
	start := time.Now()
	resp, err := http.Get("http://localhost:34501")
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
			Time:    time.Since(start).Milliseconds(),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
			Time:    time.Since(start).Milliseconds(),
		}
	}

	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
			Time:    time.Since(start).Milliseconds(),
		}
	}

	return ApiResponse{
		Success: resp.StatusCode >= 200 && resp.StatusCode < 300,
		Data:    data,
		Time:    time.Since(start).Milliseconds(),
	}
}

// TestPing calls the backend ping API
func (a *App) TestPing() ApiResponse {
	start := time.Now()
	resp, err := http.Get("http://localhost:34501/ping")
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
			Time:    time.Since(start).Milliseconds(),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
			Time:    time.Since(start).Milliseconds(),
		}
	}

	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
			Time:    time.Since(start).Milliseconds(),
		}
	}

	return ApiResponse{
		Success: resp.StatusCode >= 200 && resp.StatusCode < 300,
		Data:    data,
		Time:    time.Since(start).Milliseconds(),
	}
}
