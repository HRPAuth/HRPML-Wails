package handlers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"HRPMLW/internal/models"
	"HRPMLW/internal/services"
)

// MCHandler holds Minecraft-related handlers
type MCHandler struct {
	authService     *services.AuthService
	bmclService     *services.BMCLAPIService
	launcherService *services.LauncherService
}

// NewMCHandler creates a new MCHandler
func NewMCHandler() *MCHandler {
	return &MCHandler{
		authService:     services.NewAuthService(),
		bmclService:     services.NewBMCLAPIService(),
		launcherService: services.NewLauncherService(),
	}
}

// ==================== Authlib Injector Handlers ====================

// GetAuthlibMeta handles GET /api/auth/meta
func (h *MCHandler) GetAuthlibMeta(c *gin.Context) {
	authServer := c.Query("server")
	if authServer == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "server parameter is required"})
		return
	}

	result := h.authService.GetAuthlibMeta(authServer)
	if result.Success {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusBadRequest, result)
	}
}

// AuthlibLogin handles POST /api/auth/login
type AuthlibLoginRequest struct {
	Server      string `json:"server" binding:"required"`
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	ClientToken string `json:"client_token"`
}

func (h *MCHandler) AuthlibLogin(c *gin.Context) {
	var req AuthlibLoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	if req.ClientToken == "" {
		req.ClientToken = generateClientToken()
	}

	result := h.authService.AuthenticateWithAuthlib(req.Server, req.Username, req.Password, req.ClientToken)
	if result.Success {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusUnauthorized, result)
	}
}

// AuthlibRefresh handles POST /api/auth/refresh
type AuthlibRefreshRequest struct {
	Server      string `json:"server" binding:"required"`
	AccessToken string `json:"access_token" binding:"required"`
	ClientToken string `json:"client_token" binding:"required"`
}

func (h *MCHandler) AuthlibRefresh(c *gin.Context) {
	var req AuthlibRefreshRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	result := h.authService.RefreshAuthlibToken(req.Server, req.AccessToken, req.ClientToken)
	if result.Success {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusUnauthorized, result)
	}
}

// OfflineLogin handles POST /api/auth/offline
type OfflineLoginRequest struct {
	Username string `json:"username" binding:"required"`
}

func (h *MCHandler) OfflineLogin(c *gin.Context) {
	var req OfflineLoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Generate offline UUID based on username
	uuid := generateOfflineUUID(req.Username)

	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Data: map[string]interface{}{
			"accessToken": uuid,
			"clientToken": generateClientToken(),
			"selectedProfile": map[string]string{
				"id":   uuid,
				"name": req.Username,
			},
		},
	})
}

// ==================== BMCLAPI Handlers ====================

// GetVersions handles GET /api/versions
func (h *MCHandler) GetVersions(c *gin.Context) {
	result := h.bmclService.GetVersionList()
	if result.Success {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusBadRequest, result)
	}
}

// GetVersionManifest handles GET /api/versions/:id
// First gets version list to find the version URL, then fetches the version JSON
func (h *MCHandler) GetVersionManifest(c *gin.Context) {
	versionID := c.Param("id")
	if versionID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "version id is required"})
		return
	}

	// First get the version list to find the version URL
	listResult := h.bmclService.GetVersionList()
	if !listResult.Success {
		c.JSON(http.StatusBadRequest, listResult)
		return
	}

	manifest, ok := listResult.Data.(services.BMCLVersionManifestResponse)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to parse version list"})
		return
	}

	// Find the version URL
	var versionURL string
	for _, v := range manifest.Versions {
		if v.ID == versionID {
			versionURL = v.URL
			break
		}
	}

	if versionURL == "" {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "version not found"})
		return
	}

	// Get the version JSON using the URL
	result := h.bmclService.GetVersionManifest(versionURL)
	if result.Success {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusBadRequest, result)
	}
}

// GetFabricVersions handles GET /api/fabric/versions
func (h *MCHandler) GetFabricVersions(c *gin.Context) {
	mcVersion := c.Query("mc")
	if mcVersion == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "mc version is required"})
		return
	}

	result := h.bmclService.GetFabricVersions(mcVersion)
	if result.Success {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusBadRequest, result)
	}
}

// GetForgeVersions handles GET /api/forge/versions
func (h *MCHandler) GetForgeVersions(c *gin.Context) {
	mcVersion := c.Query("mc")
	if mcVersion == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "mc version is required"})
		return
	}

	result := h.bmclService.GetForgeVersions(mcVersion)
	if result.Success {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusBadRequest, result)
	}
}

// ==================== Download Handlers ====================

// DownloadVersion handles POST /api/download/version
//
// Performs a full install of a Minecraft version per the launcher standard:
//   1. writes <gameDir>/versions/<id>/<id>.json (required by BuildLaunchCommand)
//   2. downloads the version JAR
//   3. downloads every library listed in the version JSON
//   4. extracts natives to <gameDir>/versions/<id>/natives
type DownloadVersionRequest struct {
	VersionID string `json:"version_id" binding:"required"`
	GameDir   string `json:"game_dir"`
}

func (h *MCHandler) DownloadVersion(c *gin.Context) {
	var req DownloadVersionRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	gameDir := req.GameDir
	if gameDir == "" {
		gameDir = h.launcherService.GetMinecraftDir()
	}

	// First get the version list to find the version URL
	listResult := h.bmclService.GetVersionList()
	if !listResult.Success {
		c.JSON(http.StatusBadRequest, listResult)
		return
	}

	manifest, ok := listResult.Data.(services.BMCLVersionManifestResponse)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to parse version list"})
		return
	}

	// Find the version URL
	var versionURL string
	for _, v := range manifest.Versions {
		if v.ID == req.VersionID {
			versionURL = v.URL
			break
		}
	}

	if versionURL == "" {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "version not found"})
		return
	}

	// Get version manifest
	manifestResult := h.bmclService.GetVersionManifest(versionURL)
	if !manifestResult.Success {
		c.JSON(http.StatusBadRequest, manifestResult)
		return
	}

	versionJSON, ok := manifestResult.Data.(services.BMCLVersionJSON)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to parse manifest"})
		return
	}

	versionDir := gameDir + "/versions/" + req.VersionID
	jarPath := versionDir + "/" + req.VersionID + ".jar"
	versionJSONPath := versionDir + "/" + req.VersionID + ".json"
	nativesDir := versionDir + "/natives"
	librariesDir := gameDir + "/libraries"

	for _, dir := range []string{versionDir, librariesDir, nativesDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: fmt.Sprintf("create %s: %v", dir, err)})
			return
		}
	}

	// 1. Persist the version JSON. This is the file BuildLaunchCommand reads to
	//    obtain mainClass, assetIndex, libraries, etc. Without it the game
	//    cannot be launched.
	jsonData, err := json.MarshalIndent(versionJSON, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "marshal version json: " + err.Error()})
		return
	}
	if err := os.WriteFile(versionJSONPath, jsonData, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "write version json: " + err.Error()})
		return
	}

	// 2. Download the version JAR
	jarResult := h.bmclService.GetVersionJAR(&versionJSON)
	if !jarResult.Success {
		c.JSON(http.StatusBadRequest, jarResult)
		return
	}

	jarData, ok := jarResult.Data.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to parse JAR data"})
		return
	}

	jarURL := jarData["url"].(string)
	jarSize := jarData["size"].(int64)

	jarDownload := h.bmclService.DownloadFile(jarURL, jarPath, jarSize, nil)
	if !jarDownload.Success {
		c.JSON(http.StatusBadRequest, jarDownload)
		return
	}

	// 3. Download every library. Required for the -cp classpath the launcher
	//    builds in BuildLaunchCommand.
	libDownload := h.bmclService.DownloadLibraries(versionJSON.Libraries, librariesDir, nil)
	if !libDownload.Success {
		c.JSON(http.StatusBadRequest, libDownload)
		return
	}

	// 4. Extract natives (used by -Djava.library.path).
	nativesDownload := h.bmclService.DownloadNatives(versionJSON.Libraries, nativesDir, nil)
	if !nativesDownload.Success {
		c.JSON(http.StatusBadRequest, nativesDownload)
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Data: map[string]interface{}{
			"version_id":    req.VersionID,
			"version_dir":   versionDir,
			"version_json":  versionJSONPath,
			"version_jar":   jarPath,
			"libraries_dir": librariesDir,
			"natives_dir":   nativesDir,
		},
	})
}

// ==================== Launcher Handlers ====================

// GetLauncherConfig handles GET /api/launcher/config
func (h *MCHandler) GetLauncherConfig(c *gin.Context) {
	config := models.LauncherConfig{
		JavaPath:   h.launcherService.GetJavaPath(),
		MaxMemory:  4096,
		MinMemory:  1024,
		GameWidth:  854,
		GameHeight: 480,
		GameDir:    h.launcherService.GetMinecraftDir(),
	}

	c.JSON(http.StatusOK, models.ApiResponse{Success: true, Data: config})
}

// LaunchGame handles POST /api/launcher/launch
type LaunchGameRequest struct {
	VersionID     string `json:"version_id" binding:"required"`
	VersionType   string `json:"version_type"`
	Username      string `json:"username" binding:"required"`
	UUID          string `json:"uuid" binding:"required"`
	AccessToken   string `json:"access_token" binding:"required"`
	ClientToken   string `json:"client_token" binding:"required"`
	AuthType      string `json:"auth_type"`       // "offline" | "microsoft" | "authlib-injector"
	AuthServer    string `json:"auth_server"`     // required when AuthType == "authlib-injector"
	AuthlibJar    string `json:"authlib_jar"`     // absolute path to authlib-injector.jar (optional)
	GameDir       string `json:"game_dir"`
	JavaPath      string `json:"java_path"`
	MaxMemory     int    `json:"max_memory"`
	MinMemory     int    `json:"min_memory"`
	JavaArgs      string `json:"java_args"`
	GameWidth     int    `json:"game_width"`
	GameHeight    int    `json:"game_height"`
	LoaderType    string `json:"loader_type"`
	LoaderVersion string `json:"loader_version"`
}

func (h *MCHandler) LaunchGame(c *gin.Context) {
	var req LaunchGameRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	gameDir := req.GameDir
	if gameDir == "" {
		gameDir = h.launcherService.GetMinecraftDir()
	}

	config := models.LauncherConfig{
		JavaPath:   req.JavaPath,
		MaxMemory:  req.MaxMemory,
		MinMemory:  req.MinMemory,
		JavaArgs:   req.JavaArgs,
		GameWidth:  req.GameWidth,
		GameHeight: req.GameHeight,
		GameDir:    gameDir,
		AuthlibJar: req.AuthlibJar,
	}

	if config.MaxMemory == 0 {
		config.MaxMemory = 4096
	}
	if config.MinMemory == 0 {
		config.MinMemory = 1024
	}
	if config.GameWidth == 0 {
		config.GameWidth = 854
	}
	if config.GameHeight == 0 {
		config.GameHeight = 480
	}

	version := &models.MinecraftVersion{
		VersionID:     req.VersionID,
		VersionType:   req.VersionType,
		LoaderType:    req.LoaderType,
		LoaderVersion: req.LoaderVersion,
	}

	account := &models.MinecraftAccount{
		Username:    req.Username,
		UUID:        req.UUID,
		AccessToken: req.AccessToken,
		ClientToken: req.ClientToken,
		AuthType:    req.AuthType,
		AuthServer:  req.AuthServer,
	}

	result := h.launcherService.LaunchGame(config, version, account)
	if result.Success {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusInternalServerError, result)
	}
}

// GetMinecraftDir handles GET /api/launcher/minecraft-dir
func (h *MCHandler) GetMinecraftDir(c *gin.Context) {
	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Data: map[string]string{
			"path": h.launcherService.GetMinecraftDir(),
		},
	})
}

// ==================== Utility Functions ====================

func generateClientToken() string {
	// Generate a simple random client token
	return randomString(32)
}

func generateOfflineUUID(username string) string {
	// Generate deterministic UUID from username (offline mode)
	// Uses MD5 hash as per Minecraft offline UUID specification
	data := []byte("OfflinePlayer:" + username)
	hash := md5.Sum(data)

	// Set version to 3 (name-based)
	hash[6] = (hash[6] & 0x0f) | 0x30
	// Set variant to RFC 4122
	hash[8] = (hash[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", hash[0:4], hash[4:6], hash[6:8], hash[8:10], hash[10:])
}

func md5Hash(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

func formatUUID(hash string) string {
	// Format hash as UUID string
	if len(hash) < 32 {
		return "00000000-0000-0000-0000-000000000000"
	}
	return fmt.Sprintf("%s-%s-%s-%s-%s", hash[0:8], hash[8:12], hash[12:16], hash[16:20], hash[20:32])
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[i%len(charset)]
	}
	return string(result)
}
