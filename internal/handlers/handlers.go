package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"HRPMLW/internal/models"
	"HRPMLW/internal/services"
)

// Handler holds all handlers
type Handler struct {
	shellSvc *services.ShellService
	fileSvc  *services.FileService
	sysSvc   *services.SysService
}

// NewHandler creates a new Handler
func NewHandler() *Handler {
	return &Handler{
		shellSvc: services.NewShellService(),
		fileSvc:  services.NewFileService(),
		sysSvc:   services.NewSysService(),
	}
}

// Health handles GET / - health check
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, models.HealthResponse{
		Status:  "ok",
		Message: "Hello, Gin!",
	})
}

// Ping handles GET /ping
func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, models.PingResponse{
		Message: "pong",
	})
}

// Hello handles GET /hello/:name
func (h *Handler) Hello(c *gin.Context) {
	name := c.Param("name")
	c.JSON(http.StatusOK, models.HelloResponse{
		Message: "Hello " + name,
	})
}

// Echo handles POST /echo
func (h *Handler) Echo(c *gin.Context) {
	var requestBody map[string]interface{}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, requestBody)
}

// Shell handles POST /shell
func (h *Handler) Shell(c *gin.Context) {
	var request models.ShellRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	if request.Shell == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "shell command is required"})
		return
	}

	if request.Type == "spawn" {
		h.handleSpawnShell(c, request.Shell)
	} else {
		h.handleSimpleShell(c, request.Shell)
	}
}

func (h *Handler) handleSimpleShell(c *gin.Context, shell string) {
	response := h.shellSvc.ExecuteSimple(shell)
	c.JSON(http.StatusOK, response)
}

func (h *Handler) handleSpawnShell(c *gin.Context, shell string) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.WriteHeader(http.StatusOK)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "streaming not supported"})
		return
	}

	sendLog := func(logType, content string) {
		entry := models.NewLogEntry(logType, content)
		jsonBytes, _ := json.Marshal(entry)
		c.Writer.Write(jsonBytes)
		c.Writer.Write([]byte("\n"))
		flusher.Flush()
	}

	h.shellSvc.ExecuteSpawn(shell, c.Writer, flusher, sendLog)
}

// FileOperation handles POST /file
func (h *Handler) FileOperation(c *gin.Context) {
	var request models.FileOperationRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	if request.Operation == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "operation is required"})
		return
	}
	if request.Path == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "path is required"})
		return
	}

	var response *models.FileOperationResponse

	switch request.Operation {
	case "create":
		response = h.fileSvc.Create(request.Path, request.Content)
	case "delete":
		response = h.fileSvc.Delete(request.Path)
	case "append":
		response = h.fileSvc.Append(request.Path, request.Content)
	case "overwrite":
		response = h.fileSvc.Overwrite(request.Path, request.Content)
	case "add_json_key", "modify_json_value":
		response = h.fileSvc.AddJSONKey(request.Path, request.Key, request.Value)
	case "delete_json_key":
		response = h.fileSvc.DeleteJSONKey(request.Path, request.Key)
	default:
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid operation"})
		return
	}

	if response.Success {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusInternalServerError, response)
	}
}

// SysInfo handles GET /sysinfo
func (h *Handler) SysInfo(c *gin.Context) {
	info := h.sysSvc.GetInfo()
	c.JSON(http.StatusOK, info)
}

// NotFound handles 404
func (h *Handler) NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, models.ErrorResponse{
		Error: fmt.Sprintf("route %s not found", c.Request.URL.Path),
	})
}
