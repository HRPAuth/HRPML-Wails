package routes

import (
	"github.com/gin-gonic/gin"

	"HRPMLW/internal/handlers"
	"HRPMLW/internal/middleware"
)

// Router holds the gin engine and handlers
type Router struct {
	Engine  *gin.Engine
	Handler *handlers.Handler
	DB      *handlers.DBHandler
}

// NewRouter creates a new Router
func NewRouter() *Router {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	// Apply global middleware
	engine.Use(middleware.Recovery())
	engine.Use(middleware.Logger())
	engine.Use(middleware.CORS())

	return &Router{
		Engine:  engine,
		Handler: handlers.NewHandler(),
		DB:      handlers.NewDBHandler(),
	}
}

// Setup sets up all routes
func (r *Router) Setup() {
	// API v1 routes
	v1 := r.Engine.Group("/api/v1")
	{
		// Health routes
		v1.GET("/health", r.Handler.Health)
		v1.GET("/ping", r.Handler.Ping)
		v1.GET("/hello/:name", r.Handler.Hello)

		// Utility routes
		v1.POST("/echo", r.Handler.Echo)
		v1.GET("/sysinfo", r.Handler.SysInfo)

		// Shell routes
		shell := v1.Group("/shell")
		{
			shell.POST("/execute", r.Handler.Shell)
		}

		// File routes
		file := v1.Group("/file")
		{
			file.POST("/operation", r.Handler.FileOperation)
		}

		// Database routes
		r.setupDatabaseRoutes(v1)
	}

	// Legacy routes (for backward compatibility)
	r.setupLegacyRoutes()

	// 404 handler
	r.Engine.NoRoute(r.Handler.NotFound)
}

// setupDatabaseRoutes sets up database CRUD routes
func (r *Router) setupDatabaseRoutes(rg *gin.RouterGroup) {
	db := rg.Group("/db")
	{
		// Database stats and cleanup
		db.GET("/stats", r.DB.GetDBStats)
		db.POST("/cleanup", r.DB.CleanupDB)

		// User routes
		users := db.Group("/users")
		{
			users.GET("", r.DB.GetAllUsers)
			users.POST("", r.DB.CreateUser)
			users.GET("/:id", r.DB.GetUser)
			users.PUT("/:id", r.DB.UpdateUser)
			users.DELETE("/:id", r.DB.DeleteUser)
		}

		// Setting routes
		settings := db.Group("/settings")
		{
			settings.GET("", r.DB.GetAllSettings)
			settings.POST("", r.DB.SetSetting)
			settings.GET("/:key", r.DB.GetSetting)
			settings.DELETE("/:key", r.DB.DeleteSetting)
		}

		// Log routes
		logs := db.Group("/logs")
		{
			logs.GET("", r.DB.GetLogs)
			logs.POST("", r.DB.CreateLog)
			logs.DELETE("/old", r.DB.DeleteOldLogs)
		}

		// Task routes
		tasks := db.Group("/tasks")
		{
			tasks.GET("", r.DB.GetTasks)
			tasks.POST("", r.DB.CreateTask)
			tasks.GET("/:id", r.DB.GetTask)
			tasks.PUT("/:id", r.DB.UpdateTask)
			tasks.PATCH("/:id/status", r.DB.UpdateTaskStatus)
			tasks.PATCH("/:id/progress", r.DB.UpdateTaskProgress)
			tasks.DELETE("/:id", r.DB.DeleteTask)
		}

		// File record routes
		files := db.Group("/files")
		{
			files.GET("", r.DB.GetAllFileRecords)
			files.POST("", r.DB.CreateFileRecord)
			files.GET("/:id", r.DB.GetFileRecord)
			files.PUT("/:id", r.DB.UpdateFileRecord)
			files.DELETE("/:id", r.DB.DeleteFileRecord)
		}
	}
}

// setupLegacyRoutes sets up legacy routes for backward compatibility
func (r *Router) setupLegacyRoutes() {
	r.Engine.GET("/", r.Handler.Health)
	r.Engine.GET("/ping", r.Handler.Ping)
	r.Engine.GET("/hello/:name", r.Handler.Hello)
	r.Engine.POST("/echo", r.Handler.Echo)
	r.Engine.POST("/shell", r.Handler.Shell)
	r.Engine.POST("/file", r.Handler.FileOperation)
	r.Engine.GET("/sysinfo", r.Handler.SysInfo)
}

// SetupWithGroup sets up routes with a custom group prefix
func (r *Router) SetupWithGroup(prefix string) {
	group := r.Engine.Group(prefix)
	{
		group.GET("/health", r.Handler.Health)
		group.GET("/ping", r.Handler.Ping)
		group.GET("/hello/:name", r.Handler.Hello)
		group.POST("/echo", r.Handler.Echo)
		group.GET("/sysinfo", r.Handler.SysInfo)
		group.POST("/shell", r.Handler.Shell)
		group.POST("/file", r.Handler.FileOperation)
	}
}

// GetEngine returns the gin engine
func (r *Router) GetEngine() *gin.Engine {
	return r.Engine
}
