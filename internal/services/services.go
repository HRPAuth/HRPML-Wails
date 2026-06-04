package services

import (
	"bufio"
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

	"HRPMLW/internal/models"
)

// ShellService handles shell command operations
type ShellService struct{}

// NewShellService creates a new ShellService
func NewShellService() *ShellService {
	return &ShellService{}
}

// ExecuteSimple executes a shell command and returns combined output
func (s *ShellService) ExecuteSimple(shell string) *models.ShellResponse {
	cmd := exec.Command("sh", "-c", shell)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return &models.ShellResponse{
			Success: false,
			Error:   err.Error(),
			Output:  string(output),
		}
	}

	return &models.ShellResponse{
		Success: true,
		Output:  string(output),
	}
}

// ExecuteSpawn executes a shell command with streaming output
func (s *ShellService) ExecuteSpawn(shell string, w http.ResponseWriter, flusher http.Flusher, sendLog func(logType, content string)) {
	cmd := exec.Command("sh", "-c", shell)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		sendLog("error", fmt.Sprintf("stdout pipe error: %v", err))
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		sendLog("error", fmt.Sprintf("stderr pipe error: %v", err))
		return
	}

	if err := cmd.Start(); err != nil {
		sendLog("error", fmt.Sprintf("command start error: %v", err))
		return
	}

	go s.streamOutput(stdout, "stdout", sendLog)
	go s.streamOutput(stderr, "stderr", sendLog)

	cmd.Wait()
	sendLog("system", "Process completed")
}

func (s *ShellService) streamOutput(reader io.Reader, outputType string, sendLog func(logType, content string)) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		sendLog(outputType, scanner.Text())
	}
}

// SanitizeShellCommand sanitizes shell command
func (s *ShellService) SanitizeShellCommand(shell string) string {
	shell = strings.TrimSpace(shell)
	dangerousPatterns := []string{"rm -rf /", "> /dev/null", "2>&1 &"}
	for _, pattern := range dangerousPatterns {
		shell = strings.ReplaceAll(shell, pattern, "")
	}
	return shell
}

// FileService handles file operations
type FileService struct{}

// NewFileService creates a new FileService
func NewFileService() *FileService {
	return &FileService{}
}

// Create creates a new file
func (f *FileService) Create(path, content string) *models.FileOperationResponse {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return &models.FileOperationResponse{
			Success: false,
			Error:   err.Error(),
		}
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return &models.FileOperationResponse{
			Success: false,
			Error:   err.Error(),
		}
	}
	return &models.FileOperationResponse{
		Success: true,
		Message: "file created successfully",
	}
}

// Delete deletes a file
func (f *FileService) Delete(path string) *models.FileOperationResponse {
	if err := os.Remove(path); err != nil {
		return &models.FileOperationResponse{
			Success: false,
			Error:   err.Error(),
		}
	}
	return &models.FileOperationResponse{
		Success: true,
		Message: "file deleted successfully",
	}
}

// Append appends content to a file
func (f *FileService) Append(path, content string) *models.FileOperationResponse {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return &models.FileOperationResponse{
			Success: false,
			Error:   err.Error(),
		}
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return &models.FileOperationResponse{
			Success: false,
			Error:   err.Error(),
		}
	}
	return &models.FileOperationResponse{
		Success: true,
		Message: "file appended successfully",
	}
}

// Overwrite overwrites a file
func (f *FileService) Overwrite(path, content string) *models.FileOperationResponse {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return &models.FileOperationResponse{
			Success: false,
			Error:   err.Error(),
		}
	}
	return &models.FileOperationResponse{
		Success: true,
		Message: "file overwritten successfully",
	}
}

// AddJSONKey adds or modifies a key in a JSON file
func (f *FileService) AddJSONKey(path, key string, value interface{}) *models.FileOperationResponse {
	content, err := os.ReadFile(path)
	if err != nil {
		return &models.FileOperationResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	var data map[string]interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return &models.FileOperationResponse{
			Success: false,
			Error:   "invalid JSON file",
		}
	}

	data[key] = value
	updatedContent, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return &models.FileOperationResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	if err := os.WriteFile(path, updatedContent, 0644); err != nil {
		return &models.FileOperationResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return &models.FileOperationResponse{
		Success: true,
		Message: "JSON key added/modified successfully",
	}
}

// DeleteJSONKey deletes a key from a JSON file
func (f *FileService) DeleteJSONKey(path, key string) *models.FileOperationResponse {
	content, err := os.ReadFile(path)
	if err != nil {
		return &models.FileOperationResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	var data map[string]interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return &models.FileOperationResponse{
			Success: false,
			Error:   "invalid JSON file",
		}
	}

	delete(data, key)
	updatedContent, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return &models.FileOperationResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	if err := os.WriteFile(path, updatedContent, 0644); err != nil {
		return &models.FileOperationResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return &models.FileOperationResponse{
		Success: true,
		Message: "JSON key deleted successfully",
	}
}

// SysService handles system information
type SysService struct{}

// NewSysService creates a new SysService
func NewSysService() *SysService {
	return &SysService{}
}

// GetInfo returns system information
func (s *SysService) GetInfo() *models.SysInfoResponse {
	return &models.SysInfoResponse{
		Arch: runtime.GOARCH,
		OS:   runtime.GOOS,
	}
}

// HTTPService handles HTTP requests
type HTTPService struct {
	baseURL string
}

// NewHTTPService creates a new HTTPService
func NewHTTPService(baseURL string) *HTTPService {
	return &HTTPService{baseURL: baseURL}
}

// Get performs a GET request
func (h *HTTPService) Get(endpoint string) *models.ApiResponse {
	start := time.Now()
	resp, err := http.Get(h.baseURL + endpoint)
	if err != nil {
		return &models.ApiResponse{
			Success: false,
			Error:   err.Error(),
			Time:    time.Since(start).Milliseconds(),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &models.ApiResponse{
			Success: false,
			Error:   err.Error(),
			Time:    time.Since(start).Milliseconds(),
		}
	}

	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return &models.ApiResponse{
			Success: false,
			Error:   err.Error(),
			Time:    time.Since(start).Milliseconds(),
		}
	}

	return &models.ApiResponse{
		Success: resp.StatusCode >= 200 && resp.StatusCode < 300,
		Data:    data,
		Time:    time.Since(start).Milliseconds(),
	}
}
