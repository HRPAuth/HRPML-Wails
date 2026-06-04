package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ApiResponse represents the standard API response structure
type ApiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error,omitempty"`
	Time    int64       `json:"time,omitempty"`
}

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
	Content   string `json:"content"`
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
	// Start the Gin server in a goroutine
	go a.startGinServer()
}

func (a *App) startGinServer() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, Gin!",
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/hello/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello " + name,
		})
	})

	r.POST("/echo", func(c *gin.Context) {
		var requestBody map[string]interface{}
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, requestBody)
	})

	r.POST("/shell", func(c *gin.Context) {
		var request struct {
			Shell string `json:"shell"`
			Type  string `json:"type"`
		}
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if request.Shell == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "shell command is required"})
			return
		}

		if request.Type == "spawn" {
			a.handleSpawnShell(c, request.Shell)
		} else {
			a.handleSimpleShell(c, request.Shell)
		}
	})

	r.POST("/file", a.handleFileOperation)

	r.GET("/sysinfo", a.handleSysinfo)

	r.Run(":34501")
}

func (a *App) handleSimpleShell(c *gin.Context, shell string) {
	cmd := exec.Command("sh", "-c", shell)

	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error":   err.Error(),
			"output":  string(output),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"output":  string(output),
	})
}

func (a *App) handleSpawnShell(c *gin.Context, shell string) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.WriteHeader(http.StatusOK)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "streaming not supported"})
		return
	}

	cmd := exec.Command("sh", "-c", shell)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		a.sendLogEntry(c, flusher, "error", fmt.Sprintf("stdout pipe error: %v", err))
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		a.sendLogEntry(c, flusher, "error", fmt.Sprintf("stderr pipe error: %v", err))
		return
	}

	if err := cmd.Start(); err != nil {
		a.sendLogEntry(c, flusher, "error", fmt.Sprintf("command start error: %v", err))
		return
	}

	go a.streamOutput(c, flusher, stdout, "stdout")
	go a.streamOutput(c, flusher, stderr, "stderr")

	cmd.Wait()

	a.sendLogEntry(c, flusher, "system", "Process completed")
	flusher.Flush()
}

func (a *App) sendLogEntry(c *gin.Context, flusher http.Flusher, logType string, content string) {
	entry := LogEntry{
		Timestamp: time.Now().Format("2006-01-02 15:04:05.000"),
		Type:      logType,
		Content:   content,
	}
	jsonBytes, _ := json.Marshal(entry)
	c.Writer.Write(jsonBytes)
	c.Writer.Write([]byte("\n"))
	flusher.Flush()
}

func (a *App) streamOutput(c *gin.Context, flusher http.Flusher, reader io.Reader, outputType string) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		a.sendLogEntry(c, flusher, outputType, scanner.Text())
	}
}

func (a *App) handleSysinfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"arch": runtime.GOARCH,
		"os":   runtime.GOOS,
	})
}

func (a *App) sanitizeShellCommand(shell string) string {
	shell = strings.TrimSpace(shell)
	dangerousPatterns := []string{"rm -rf /", "> /dev/null", "2>&1 &"}
	for _, pattern := range dangerousPatterns {
		shell = strings.ReplaceAll(shell, pattern, "")
	}
	return shell
}

func (a *App) handleFileOperation(c *gin.Context) {
	var request struct {
		Operation string      `json:"operation"`
		Path      string      `json:"path"`
		Content   string      `json:"content"`
		Key       string      `json:"key"`
		Value     interface{} `json:"value"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Operation == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "operation is required"})
		return
	}
	if request.Path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
		return
	}

	switch request.Operation {
	case "create":
		dir := filepath.Dir(request.Path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
			return
		}
		if err := os.WriteFile(request.Path, []byte(request.Content), 0644); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "file created successfully"})
	case "delete":
		if err := os.Remove(request.Path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "file deleted successfully"})
	case "append":
		file, err := os.OpenFile(request.Path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
			return
		}
		defer file.Close()
		if _, err := file.WriteString(request.Content); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "file appended successfully"})
	case "overwrite":
		if err := os.WriteFile(request.Path, []byte(request.Content), 0644); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "file overwritten successfully"})
	case "add_json_key", "modify_json_value":
		content, err := os.ReadFile(request.Path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
			return
		}
		var data map[string]interface{}
		if err := json.Unmarshal(content, &data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid JSON file"})
			return
		}
		data[request.Key] = request.Value
		updatedContent, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
			return
		}
		if err := os.WriteFile(request.Path, updatedContent, 0644); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "JSON key added/modified successfully"})
	case "delete_json_key":
		content, err := os.ReadFile(request.Path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
			return
		}
		var data map[string]interface{}
		if err := json.Unmarshal(content, &data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid JSON file"})
			return
		}
		delete(data, request.Key)
		updatedContent, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
			return
		}
		if err := os.WriteFile(request.Path, updatedContent, 0644); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "JSON key deleted successfully"})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid operation"})
	}
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
