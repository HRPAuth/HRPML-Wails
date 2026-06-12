package services

import (
	"archive/zip"
	"bufio"
	"crypto/md5"
	"encoding/hex"
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

// ==================== Authlib Injector Service ====================

// AuthService handles Minecraft authentication
type AuthService struct {
	client *http.Client
}

// NewAuthService creates a new AuthService
func NewAuthService() *AuthService {
	return &AuthService{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// AuthlibServerMeta represents authlib-injector server metadata
type AuthlibServerMeta struct {
	ServerName            string `json:"serverName"`
	ImplementationName    string `json:"implementationName"`
	ImplementationVersion string `json:"implementationVersion"`
	Links                 struct {
		Homepage string `json:"homepage"`
		Register string `json:"register"`
	} `json:"links"`
	SignaturePublickey string   `json:"signaturePublickey"`
	SkinDomains        []string `json:"skinDomains"`
}

// AuthlibAuthRequest represents authentication request to authlib-injector
type AuthlibAuthRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	ClientToken string `json:"clientToken,omitempty"`
	RequestUser bool   `json:"requestUser"`
}

// AuthlibAuthResponse represents authentication response from authlib-injector
type AuthlibAuthResponse struct {
	AccessToken       string           `json:"accessToken"`
	ClientToken       string           `json:"clientToken"`
	AvailableProfiles []AuthlibProfile `json:"availableProfiles"`
	SelectedProfile   AuthlibProfile   `json:"selectedProfile"`
	User              *AuthlibUser     `json:"user,omitempty"`
}

type AuthlibProfile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AuthlibUser struct {
	ID         string `json:"id"`
	Properties []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"properties"`
}

// GetAuthlibMeta fetches authlib-injector server metadata
func (a *AuthService) GetAuthlibMeta(authServerURL string) *models.ApiResponse {
	url := strings.TrimSuffix(authServerURL, "/") + "/"

	resp, err := a.client.Get(url)
	if err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &models.ApiResponse{Success: false, Error: fmt.Sprintf("HTTP %d", resp.StatusCode)}
	}

	var meta AuthlibServerMeta
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}

	return &models.ApiResponse{Success: true, Data: meta}
}

// AuthenticateWithAuthlib authenticates with an authlib-injector server
func (a *AuthService) AuthenticateWithAuthlib(authServerURL, username, password, clientToken string) *models.ApiResponse {
	url := strings.TrimSuffix(authServerURL, "/") + "/authserver/authenticate"

	reqBody := AuthlibAuthRequest{
		Username:    username,
		Password:    password,
		ClientToken: clientToken,
		RequestUser: true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		if errMsg, ok := errResp["errorMessage"].(string); ok {
			return &models.ApiResponse{Success: false, Error: errMsg}
		}
		return &models.ApiResponse{Success: false, Error: fmt.Sprintf("HTTP %d", resp.StatusCode)}
	}

	var authResp AuthlibAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}

	return &models.ApiResponse{Success: true, Data: authResp}
}

// RefreshAuthlibToken refreshes an authlib-injector access token
func (a *AuthService) RefreshAuthlibToken(authServerURL, accessToken, clientToken string) *models.ApiResponse {
	url := strings.TrimSuffix(authServerURL, "/") + "/authserver/refresh"

	reqBody := map[string]string{
		"accessToken": accessToken,
		"clientToken": clientToken,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &models.ApiResponse{Success: false, Error: fmt.Sprintf("HTTP %d", resp.StatusCode)}
	}

	var authResp AuthlibAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}

	return &models.ApiResponse{Success: true, Data: authResp}
}

// ==================== BMCLAPI Service ====================

// BMCLAPIService handles BMCLAPI downloads
type BMCLAPIService struct {
	client  *http.Client
	mirrors []string
	baseURL string
}

// NewBMCLAPIService creates a new BMCLAPIService
func NewBMCLAPIService() *BMCLAPIService {
	return &BMCLAPIService{
		client:  &http.Client{Timeout: 60 * time.Second},
		mirrors: []string{"https://bmclapi2.bangbang93.com", "https://bmclapi.ggsdream.com"},
		baseURL: "https://bmclapi2.bangbang93.com",
	}
}

// BMCLVersionManifestResponse represents the version manifest response from BMCLAPI
type BMCLVersionManifestResponse struct {
	Latest   map[string]string `json:"latest"`
	Versions []BMCLVersion     `json:"versions"`
}

type BMCLVersion struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	URL             string `json:"url"`
	Time            string `json:"time"`
	ReleaseTime     string `json:"releaseTime"`
	Sha1            string `json:"sha1"`
	ComplianceLevel int    `json:"complianceLevel"`
}

// BMCLVersionJSON represents version.json
type BMCLVersionJSON struct {
	ID                     string                  `json:"id"`
	AssetIndex             BMCLAssetIndex          `json:"assetIndex"`
	Assets                 string                  `json:"assets"`
	ComplianceLevel        int                     `json:"complianceLevel"`
	Downloads              map[string]BMCLDownload `json:"downloads"`
	JavaVersion            BMCLJavaVersion         `json:"javaVersion"`
	Libraries              []BMCLLibrary           `json:"libraries"`
	MainClass              string                  `json:"mainClass"`
	MinecraftArguments     string                  `json:"minecraftArguments,omitempty"`
	Arguments              BMCLArguments           `json:"arguments"`
	MinimumLauncherVersion int                     `json:"minimumLauncherVersion"`
	ReleaseTime            string                  `json:"releaseTime"`
	Time                   string                  `json:"time"`
	Type                   string                  `json:"type"`
}

type BMCLAssetIndex struct {
	ID   string `json:"id"`
	URL  string `json:"url"`
	SHA1 string `json:"sha1"`
	Size int64  `json:"size"`
}

type BMCLDownload struct {
	SHA1 string `json:"sha1"`
	URL  string `json:"url"`
	Size int64  `json:"size"`
}

type BMCLJavaVersion struct {
	Component string `json:"component"`
	SHA1      string `json:"sha1"`
	Size      int64  `json:"size"`
	URL       string `json:"url"`
}

type BMCLLibrary struct {
	Name      string               `json:"name"`
	Downloads BMCLLibraryDownloads `json:"downloads"`
	Rules     []BMCLRule           `json:"rules,omitempty"`
}

type BMCLLibraryDownloads struct {
	Artifact    BMCLArtifact            `json:"artifact"`
	Classifiers map[string]BMCLArtifact `json:"classifiers,omitempty"`
}

type BMCLArtifact struct {
	Path string `json:"path"`
	URL  string `json:"url"`
	SHA1 string `json:"sha1"`
	Size int64  `json:"size"`
}

type BMCLArguments struct {
	Game []interface{} `json:"game"`
	JVM  []interface{} `json:"jvm"`
}

type BMCLRule struct {
	Action string  `json:"action"`
	OS     *BMCLOS `json:"os,omitempty"`
}

type BMCLOS struct {
	Name string `json:"name"`
}

// GetVersionList returns all available Minecraft versions
// BMCLAPI endpoint: /mc/game/version_manifest.json
func (b *BMCLAPIService) GetVersionList() *models.ApiResponse {
	url := b.baseURL + "/mc/game/version_manifest.json"

	resp, err := b.client.Get(url)
	if err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &models.ApiResponse{Success: false, Error: fmt.Sprintf("HTTP %d", resp.StatusCode)}
	}

	var manifest BMCLVersionManifestResponse
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}

	return &models.ApiResponse{Success: true, Data: manifest}
}

// GetVersionManifest returns version JSON for a specific version
// The version URL from manifest should be replaced with BMCLAPI domain
func (b *BMCLAPIService) GetVersionManifest(versionURL string) *models.ApiResponse {
	// Replace official URLs with BMCLAPI mirror
	bmclURL := versionURL
	bmclURL = strings.ReplaceAll(bmclURL, "https://launchermeta.mojang.com/", b.baseURL+"/")
	bmclURL = strings.ReplaceAll(bmclURL, "https://launcher.mojang.com/", b.baseURL+"/")

	resp, err := b.client.Get(bmclURL)
	if err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &models.ApiResponse{Success: false, Error: fmt.Sprintf("HTTP %d", resp.StatusCode)}
	}

	var versionJSON BMCLVersionJSON
	if err := json.NewDecoder(resp.Body).Decode(&versionJSON); err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}

	return &models.ApiResponse{Success: true, Data: versionJSON}
}

// GetVersionJAR returns download URL for version JAR from version JSON downloads
func (b *BMCLAPIService) GetVersionJAR(versionJSON *BMCLVersionJSON) *models.ApiResponse {
	if versionJSON.Downloads == nil {
		return &models.ApiResponse{Success: false, Error: "No downloads found in version JSON"}
	}

	clientDownload, ok := versionJSON.Downloads["client"]
	if !ok {
		return &models.ApiResponse{Success: false, Error: "No client download found"}
	}

	// Replace official URL with BMCLAPI mirror
	url := clientDownload.URL
	url = strings.ReplaceAll(url, "https://launcher.mojang.com/", b.baseURL+"/")
	url = strings.ReplaceAll(url, "https://launchermeta.mojang.com/", b.baseURL+"/")

	return &models.ApiResponse{Success: true, Data: map[string]interface{}{
		"url":  url,
		"sha1": clientDownload.SHA1,
		"size": clientDownload.Size,
	}}
}

// DownloadFile downloads a file with progress callback
func (b *BMCLAPIService) DownloadFile(url, destPath string, expectedSize int64, progressCallback func(downloaded, total int64)) *models.ApiResponse {
	// Create directory if not exists
	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}

	resp, err := b.client.Get(url)
	if err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &models.ApiResponse{Success: false, Error: fmt.Sprintf("HTTP %d", resp.StatusCode)}
	}

	totalSize := expectedSize
	if totalSize == 0 {
		totalSize = resp.ContentLength
	}

	file, err := os.Create(destPath)
	if err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}
	defer file.Close()

	downloaded := int64(0)
	buf := make([]byte, 32*1024)
	reader := bufio.NewReader(resp.Body)

	for {
		n, err := reader.Read(buf)
		if n > 0 {
			file.Write(buf[:n])
			downloaded += int64(n)
			if progressCallback != nil {
				progressCallback(downloaded, totalSize)
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return &models.ApiResponse{Success: false, Error: err.Error()}
		}
	}

	return &models.ApiResponse{Success: true, Data: map[string]interface{}{
		"path": destPath,
		"size": downloaded,
	}}
}

// DownloadLibraries downloads all required libraries
// BMCLAPI: https://libraries.minecraft.net/ -> https://bmclapi2.bangbang93.com/maven
func (b *BMCLAPIService) DownloadLibraries(libraries []BMCLLibrary, libsDir string, progressCallback func(name string, downloaded, total int64)) *models.ApiResponse {
	var lastErr error
	for _, lib := range libraries {
		if !shouldDownloadLibrary(lib) {
			continue
		}

		artifact := lib.Downloads.Artifact
		// Replace official library URL with BMCLAPI maven mirror
		url := artifact.URL
		url = strings.ReplaceAll(url, "https://libraries.minecraft.net/", b.baseURL+"/maven/")

		parts := strings.Split(lib.Name, ":")
		if len(parts) != 3 {
			continue
		}
		group := strings.ReplaceAll(parts[0], ".", "/")
		artifactPath := filepath.Join(libsDir, group, parts[1], parts[2], parts[1]+"-"+parts[2]+".jar")

		if _, err := os.Stat(artifactPath); err == nil {
			continue // Already exists
		}

		result := b.DownloadFile(url, artifactPath, artifact.Size, nil)
		if !result.Success {
			lastErr = fmt.Errorf("failed to download %s: %s", lib.Name, result.Error)
			continue // Try next library instead of failing completely
		}

		if progressCallback != nil {
			progressCallback(lib.Name, artifact.Size, artifact.Size)
		}
	}

	if lastErr != nil {
		return &models.ApiResponse{Success: false, Error: lastErr.Error()}
	}
	return &models.ApiResponse{Success: true}
}

func shouldDownloadLibrary(lib BMCLLibrary) bool {
	if len(lib.Rules) == 0 {
		return true
	}

	allowed := false
	for _, rule := range lib.Rules {
		if rule.Action == "allow" {
			if rule.OS == nil || rule.OS.Name == runtime.GOOS {
				allowed = true
			}
		} else if rule.Action == "disallow" {
			if rule.OS != nil && rule.OS.Name == runtime.GOOS {
				return false
			}
		}
	}

	return allowed
}

// GetFabricVersions returns available Fabric versions
// BMCLAPI: https://meta.fabricmc.net -> https://bmclapi2.bangbang93.com/fabric-meta
func (b *BMCLAPIService) GetFabricVersions(mcVersion string) *models.ApiResponse {
	url := b.baseURL + "/fabric-meta/v2/versions/loader/" + mcVersion

	resp, err := b.client.Get(url)
	if err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &models.ApiResponse{Success: false, Error: fmt.Sprintf("HTTP %d", resp.StatusCode)}
	}

	var versions []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}

	return &models.ApiResponse{Success: true, Data: versions}
}

// GetFabricInstaller downloads Fabric installer
// BMCLAPI: https://meta.fabricmc.net -> https://bmclapi2.bangbang93.com/fabric-meta
func (b *BMCLAPIService) GetFabricInstaller(mcVersion, loaderVersion, targetDir string) *models.ApiResponse {
	url := fmt.Sprintf("%s/fabric-meta/v2/versions/loader/%s/%s/profile/jar", b.baseURL, mcVersion, loaderVersion)

	resp, err := b.client.Get(url)
	if err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &models.ApiResponse{Success: false, Error: fmt.Sprintf("HTTP %d", resp.StatusCode)}
	}

	destPath := filepath.Join(targetDir, "fabric-loader-"+loaderVersion+".jar")
	file, err := os.Create(destPath)
	if err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}

	return &models.ApiResponse{Success: true, Data: map[string]string{"path": destPath}}
}

// GetForgeVersions returns available Forge versions for a Minecraft version
func (b *BMCLAPIService) GetForgeVersions(mcVersion string) *models.ApiResponse {
	url := b.baseURL + "/forge/minecraft/" + mcVersion

	resp, err := b.client.Get(url)
	if err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &models.ApiResponse{Success: false, Error: fmt.Sprintf("HTTP %d", resp.StatusCode)}
	}

	var versions []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}

	return &models.ApiResponse{Success: true, Data: versions}
}

// DownloadForgeInstaller downloads Forge installer
func (b *BMCLAPIService) DownloadForgeInstaller(mcVersion, forgeVersion, targetDir string) *models.ApiResponse {
	url := fmt.Sprintf("%s/forge/minecraft/%s/%s/installer", b.baseURL, mcVersion, forgeVersion)

	resp, err := b.client.Get(url)
	if err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &models.ApiResponse{Success: false, Error: fmt.Sprintf("HTTP %d", resp.StatusCode)}
	}

	destPath := filepath.Join(targetDir, "forge-installer-"+forgeVersion+".jar")
	file, err := os.Create(destPath)
	if err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}

	return &models.ApiResponse{Success: true, Data: map[string]string{"path": destPath}}
}

// DownloadAssets downloads game assets
// BMCLAPI: http://resources.download.minecraft.net -> https://bmclapi2.bangbang93.com/assets
func (b *BMCLAPIService) DownloadAssets(assetIndex *BMCLAssetIndex, assetsDir string, progressCallback func(name string, downloaded, total int64)) *models.ApiResponse {
	// Download asset index JSON
	assetIndexURL := assetIndex.URL
	assetIndexURL = strings.ReplaceAll(assetIndexURL, "https://launchermeta.mojang.com/", b.baseURL+"/")
	assetIndexURL = strings.ReplaceAll(assetIndexURL, "https://launcher.mojang.com/", b.baseURL+"/")

	resp, err := b.client.Get(assetIndexURL)
	if err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &models.ApiResponse{Success: false, Error: fmt.Sprintf("HTTP %d", resp.StatusCode)}
	}

	var assetIndexJSON struct {
		Objects map[string]struct {
			Hash string `json:"hash"`
			Size int64  `json:"size"`
		} `json:"objects"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&assetIndexJSON); err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}

	// Download each asset
	for name, obj := range assetIndexJSON.Objects {
		hash := obj.Hash
		subDir := hash[:2]
		url := fmt.Sprintf("%s/assets/%s/%s", b.baseURL, subDir, hash)

		destPath := filepath.Join(assetsDir, "objects", subDir, hash)
		if _, err := os.Stat(destPath); err == nil {
			continue // Already exists
		}

		result := b.DownloadFile(url, destPath, obj.Size, nil)
		if !result.Success {
			return &models.ApiResponse{Success: false, Error: fmt.Sprintf("failed to download asset %s: %s", name, result.Error)}
		}

		if progressCallback != nil {
			progressCallback(name, obj.Size, obj.Size)
		}
	}

	return &models.ApiResponse{Success: true}
}

// DownloadNatives downloads native libraries for the current platform
func (b *BMCLAPIService) DownloadNatives(libraries []BMCLLibrary, nativesDir string, progressCallback func(name string, downloaded, total int64)) *models.ApiResponse {
	for _, lib := range libraries {
		if len(lib.Downloads.Classifiers) == 0 {
			continue
		}

		// Determine native classifier based on OS
		classifier := ""
		switch runtime.GOOS {
		case "windows":
			classifier = "natives-windows"
		case "darwin":
			classifier = "natives-macos"
		case "linux":
			classifier = "natives-linux"
		}

		if classifier == "" {
			continue
		}

		artifact, ok := lib.Downloads.Classifiers[classifier]
		if !ok {
			continue
		}

		url := artifact.URL
		url = strings.ReplaceAll(url, "https://libraries.minecraft.net/", b.baseURL+"/maven/")

		// Extract the native jar
		resp, err := b.client.Get(url)
		if err != nil {
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			continue
		}

		// Create temp file
		tmpFile, err := os.CreateTemp("", "native-*.jar")
		if err != nil {
			resp.Body.Close()
			continue
		}
		tmpPath := tmpFile.Name()

		_, err = io.Copy(tmpFile, resp.Body)
		tmpFile.Close()
		resp.Body.Close()

		if err != nil {
			os.Remove(tmpPath)
			continue
		}

		// Extract native files
		err = b.extractNatives(tmpPath, nativesDir)
		os.Remove(tmpPath)

		if err != nil {
			continue
		}

		if progressCallback != nil {
			progressCallback(lib.Name, artifact.Size, artifact.Size)
		}
	}

	return &models.ApiResponse{Success: true}
}

// extractNatives extracts native libraries from a jar file
func (b *BMCLAPIService) extractNatives(jarPath, destDir string) error {
	r, err := zip.OpenReader(jarPath)
	if err != nil {
		return err
	}
	defer r.Close()

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	for _, f := range r.File {
		// Skip directories and META-INF
		if f.FileInfo().IsDir() || strings.HasPrefix(f.Name, "META-INF/") {
			continue
		}

		// Only extract native files (.dll, .so, .dylib, .jnilib)
		ext := strings.ToLower(filepath.Ext(f.Name))
		if ext != ".dll" && ext != ".so" && ext != ".dylib" && ext != ".jnilib" {
			continue
		}

		destPath := filepath.Join(destDir, filepath.Base(f.Name))

		rc, err := f.Open()
		if err != nil {
			continue
		}

		destFile, err := os.Create(destPath)
		if err != nil {
			rc.Close()
			continue
		}

		io.Copy(destFile, rc)
		destFile.Close()
		rc.Close()
	}

	return nil
}

// ==================== Minecraft Launcher Service ====================

// LauncherService handles Minecraft game launching
type LauncherService struct {
	authService *AuthService
	bmclService *BMCLAPIService
}

// NewLauncherService creates a new LauncherService
func NewLauncherService() *LauncherService {
	return &LauncherService{
		authService: NewAuthService(),
		bmclService: NewBMCLAPIService(),
	}
}

// GetMinecraftDir returns the Minecraft directory path
func (l *LauncherService) GetMinecraftDir() string {
	home, _ := os.UserHomeDir()
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("APPDATA"), ".minecraft")
	} else if runtime.GOOS == "darwin" {
		return filepath.Join(home, "Library", "Application Support", "minecraft")
	}
	return filepath.Join(home, ".minecraft")
}

// GetJavaPath finds Java installation
func (l *LauncherService) GetJavaPath() string {
	javaHome := os.Getenv("JAVA_HOME")
	if javaHome != "" {
		javapath := filepath.Join(javaHome, "bin", "java")
		if _, err := os.Stat(javapath); err == nil {
			return javapath
		}
	}

	// Check common locations
	paths := []string{
		"/usr/bin/java",
		"/usr/local/bin/java",
		"/opt/java/bin/java",
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	return "java" // Fallback to PATH
}

// loadVersionJSON reads the version JSON file from <gameDir>/versions/<id>/<id>.json
func (l *LauncherService) loadVersionJSON(versionID, gameDir string) (*BMCLVersionJSON, error) {
	versionJSONPath := filepath.Join(gameDir, "versions", versionID, versionID+".json")
	data, err := os.ReadFile(versionJSONPath)
	if err != nil {
		return nil, fmt.Errorf("read version json %q: %w", versionJSONPath, err)
	}
	var versionJSON BMCLVersionJSON
	if err := json.Unmarshal(data, &versionJSON); err != nil {
		return nil, fmt.Errorf("parse version json: %w", err)
	}
	return &versionJSON, nil
}

// buildLibraryClasspath assembles the classpath from the version's libraries,
// honouring OS rules and skipping natives (which are loaded via -Djava.library.path).
func buildLibraryClasspath(versionJSON *BMCLVersionJSON, librariesDir string) string {
	sep := string(os.PathListSeparator)
	parts := make([]string, 0, len(versionJSON.Libraries)+1)
	seen := make(map[string]bool)

	for _, lib := range versionJSON.Libraries {
		if !shouldDownloadLibrary(lib) {
			continue
		}
		// Skip natives: they are not jars on the classpath, they're extracted to the natives dir
		if len(lib.Downloads.Classifiers) > 0 {
			continue
		}
		artifact := lib.Downloads.Artifact
		if artifact.Path == "" {
			continue
		}
		if seen[artifact.Path] {
			continue
		}
		seen[artifact.Path] = true
		parts = append(parts, filepath.Join(librariesDir, artifact.Path))
	}

	return strings.Join(parts, sep)
}

// resolveUserType maps the stored account auth type to the value Mojang-style
// game arguments expect (--userType). Per the launcher standard:
//   - offline accounts use "legacy"
//   - Microsoft accounts use "msa"
//   - authlib-injector / Mojang Yggdrasil accounts use "mojang"
func resolveUserType(authType string) string {
	switch strings.ToLower(authType) {
	case "offline":
		return "legacy"
	case "microsoft", "msa":
		return "msa"
	default:
		return "mojang"
	}
}

// resolveVersionType maps the stored version type to the value game arguments
// expect (--versionType). Defaults to "release" for unknown types.
func resolveVersionType(versionType string) string {
	if versionType == "" {
		return "release"
	}
	return versionType
}

// BuildLaunchCommand builds the Minecraft launch command per the standard:
//   - reads <gameDir>/versions/<id>/<id>.json to obtain mainClass, assetIndex and libraries
//   - assembles -cp with the version jar + all matching libraries
//   - sets -Djava.library.path to the natives dir (must be populated beforehand)
//   - prepends -javaagent:/path/authlib-injector.jar=<auth-server> when an
//     authlib-injector account is in use and the jar path is provided
//   - emits the standard game args (--username, --uuid, --accessToken, --userType,
//     --version, --gameDir, --assetsDir, --assetIndex, --width, --height)
//   - injects --tweakClass for Forge LaunchWrapper main classes
func (l *LauncherService) BuildLaunchCommand(config models.LauncherConfig, version *models.MinecraftVersion, account *models.MinecraftAccount) ([]string, error) {
	if version == nil {
		return nil, fmt.Errorf("version is required")
	}
	if account == nil {
		return nil, fmt.Errorf("account is required")
	}

	javaPath := config.JavaPath
	if javaPath == "" {
		javaPath = l.GetJavaPath()
	}

	gameDir := config.GameDir
	if gameDir == "" {
		gameDir = l.GetMinecraftDir()
	}

	versionDir := filepath.Join(gameDir, "versions", version.VersionID)
	jarPath := filepath.Join(versionDir, version.VersionID+".jar")
	nativesDir := filepath.Join(versionDir, "natives")
	librariesDir := filepath.Join(gameDir, "libraries")

	if err := os.MkdirAll(nativesDir, 0755); err != nil {
		return nil, fmt.Errorf("create natives dir: %w", err)
	}

	// Load the version JSON; fall back to safe defaults if it cannot be read.
	versionJSON, err := l.loadVersionJSON(version.VersionID, gameDir)
	if err != nil {
		return nil, err
	}

	mainClass := versionJSON.MainClass
	if mainClass == "" {
		mainClass = "net.minecraft.client.main.Main"
	}

	assetIndex := versionJSON.AssetIndex.ID
	if assetIndex == "" {
		assetIndex = version.VersionID
	}

	// Build the classpath: <version jar> + every library (with rules applied).
	classpathParts := []string{jarPath}
	if libCP := buildLibraryClasspath(versionJSON, librariesDir); libCP != "" {
		classpathParts = append(classpathParts, libCP)
	}
	classpath := strings.Join(classpathParts, string(os.PathListSeparator))

	userType := resolveUserType(account.AuthType)
	versionType := resolveVersionType(version.VersionType)

	// JVM args
	args := []string{javaPath}
	if config.MaxMemory > 0 {
		args = append(args, fmt.Sprintf("-Xmx%dM", config.MaxMemory))
	}
	if config.MinMemory > 0 {
		args = append(args, fmt.Sprintf("-Xms%dM", config.MinMemory))
	}
	if config.JavaArgs != "" {
		args = append(args, strings.Fields(config.JavaArgs)...)
	}
	args = append(args, "-Djava.library.path="+nativesDir)

	// Authlib-Injector: -javaagent:/path/authlib-injector.jar=<auth-server>
	if strings.EqualFold(account.AuthType, "authlib-injector") && config.AuthlibJar != "" && account.AuthServer != "" {
		args = append(args, fmt.Sprintf("-javaagent:%s=%s", config.AuthlibJar, strings.TrimRight(account.AuthServer, "/")))
	}

	args = append(args, "-cp", classpath, mainClass)

	// Game args per the standard
	gameArgs := []string{
		"--username", account.Username,
		"--version", version.VersionID,
		"--gameDir", gameDir,
		"--assetsDir", filepath.Join(gameDir, "assets"),
		"--assetIndex", assetIndex,
		"--uuid", stripDashes(account.UUID),
		"--accessToken", account.AccessToken,
		"--userType", userType,
		"--versionType", versionType,
	}
	if config.GameWidth > 0 {
		gameArgs = append(gameArgs, "--width", fmt.Sprintf("%d", config.GameWidth))
	}
	if config.GameHeight > 0 {
		gameArgs = append(gameArgs, "--height", fmt.Sprintf("%d", config.GameHeight))
	}

	// Forge 1.12.x and older use LaunchWrapper and need --tweakClass
	if strings.HasPrefix(mainClass, "net.minecraft.launchwrapper.Launch") {
		gameArgs = append(gameArgs, "--tweakClass", "net.minecraftforge.fml.common.launcher.FMLTweaker")
	}

	args = append(args, gameArgs...)
	return args, nil
}

// stripDashes removes dashes from a UUID (game expects compact form).
func stripDashes(s string) string {
	return strings.ReplaceAll(s, "-", "")
}

// LaunchGame launches Minecraft with the given configuration.
//
// The process is detached into its own process group so it survives the
// launcher (Wails) process exiting. stdout/stderr are streamed to a log file
// under <gameDir>/logs and the live file path is returned alongside the PID
// so the frontend can tail it.
func (l *LauncherService) LaunchGame(config models.LauncherConfig, version *models.MinecraftVersion, account *models.MinecraftAccount) *models.ApiResponse {
	args, err := l.BuildLaunchCommand(config, version, account)
	if err != nil {
		return &models.ApiResponse{Success: false, Error: err.Error()}
	}

	gameDir := config.GameDir
	if gameDir == "" {
		gameDir = l.GetMinecraftDir()
	}

	logsDir := filepath.Join(gameDir, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return &models.ApiResponse{Success: false, Error: fmt.Sprintf("create logs dir: %v", err)}
	}
	logPath := filepath.Join(logsDir, fmt.Sprintf("launcher-%s.log", time.Now().Format("20060102-150405")))
	logFile, err := os.Create(logPath)
	if err != nil {
		return &models.ApiResponse{Success: false, Error: fmt.Sprintf("create log file: %v", err)}
	}

	// Write the exact command we are about to run, for debugging.
	fmt.Fprintf(logFile, "# Minecraft launch command\n# %s\n\n", time.Now().Format(time.RFC3339))
	for _, a := range args {
		fmt.Fprintf(logFile, "%s ", shellQuote(a))
	}
	fmt.Fprintln(logFile)

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = gameDir

	// Detach the process so it keeps running when the launcher exits.
	// configureDetachedProcess is implemented per-platform in:
	//   - mc_services_unix.go    (Linux / macOS)
	//   - mc_services_windows.go (Windows)
	configureDetachedProcess(cmd)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logFile.Close()
		return &models.ApiResponse{Success: false, Error: fmt.Sprintf("stdout pipe: %v", err)}
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		logFile.Close()
		return &models.ApiResponse{Success: false, Error: fmt.Sprintf("stderr pipe: %v", err)}
	}

	if err := cmd.Start(); err != nil {
		logFile.Close()
		return &models.ApiResponse{Success: false, Error: fmt.Sprintf("start game: %v", err)}
	}

	// Stream game output into the log file in the background.
	go func() {
		defer logFile.Close()
		mw := io.MultiWriter(logFile)
		go io.Copy(mw, stdout)
		io.Copy(mw, stderr)
		cmd.Wait()
		fmt.Fprintf(logFile, "\n# Process exited\n")
	}()

	return &models.ApiResponse{
		Success: true,
		Data: map[string]interface{}{
			"pid":      cmd.Process.Pid,
			"args":     args,
			"log_path": logPath,
		},
	}
}

// shellQuote quotes an argument for the log file so the recorded command can
// be re-run safely. Single-quote and escape embedded single quotes.
func shellQuote(s string) string {
	if s == "" {
		return "''"
	}
	if !strings.ContainsAny(s, " \t\n\"'\\$`&;|*?<>()[]{}#~!") {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

// ==================== File Utility Service ====================

// VerifyMD5 checks file MD5
func VerifyMD5(filePath, expectedMD5 string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return false
	}

	actualMD5 := hex.EncodeToString(hash.Sum(nil))
	return actualMD5 == expectedMD5
}

// Unzip extracts a zip archive
func Unzip(src, dest string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(dest, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, 0755)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		w, err := os.Create(path)
		if err != nil {
			return err
		}
		defer w.Close()

		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		if _, err := io.Copy(w, rc); err != nil {
			return err
		}
	}

	return nil
}
