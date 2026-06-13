package repository

import (
	"database/sql"
	"fmt"
	"time"

	"HRPMLW/internal/database"
	"HRPMLW/internal/models"
)

// UserRepository handles user database operations
type UserRepository struct{}

// NewUserRepository creates a new UserRepository
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	query := `INSERT INTO users (username, email, password_hash, role) VALUES (?, ?, ?, ?)`
	result, err := database.GetDB().Exec(query, user.Username, user.Email, user.PasswordHash, user.Role)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	user.ID, _ = result.LastInsertId()
	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, role, created_at, updated_at FROM users WHERE id = ?`
	err := database.GetDB().QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password_hash, role, created_at, updated_at FROM users WHERE username = ?`
	err := database.GetDB().QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

// GetAll retrieves all users
func (r *UserRepository) GetAll() ([]models.User, error) {
	query := `SELECT id, username, email, role, created_at, updated_at FROM users ORDER BY created_at DESC`
	rows, err := database.GetDB().Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	query := `UPDATE users SET username = ?, email = ?, role = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := database.GetDB().Exec(query, user.Username, user.Email, user.Role, user.ID)
	return err
}

// Delete deletes a user
func (r *UserRepository) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := database.GetDB().Exec(query, id)
	return err
}

// SettingRepository handles setting database operations
type SettingRepository struct{}

// NewSettingRepository creates a new SettingRepository
func NewSettingRepository() *SettingRepository {
	return &SettingRepository{}
}

// Set creates or updates a setting
func (r *SettingRepository) Set(key, value, settingType, description string) error {
	query := `INSERT INTO settings (key, value, type, description) VALUES (?, ?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET value = ?, type = ?, description = ?, updated_at = CURRENT_TIMESTAMP`
	_, err := database.GetDB().Exec(query, key, value, settingType, description, value, settingType, description)
	return err
}

// Get retrieves a setting by key
func (r *SettingRepository) Get(key string) (*models.Setting, error) {
	setting := &models.Setting{}
	query := `SELECT id, key, value, type, description, created_at, updated_at FROM settings WHERE key = ?`
	err := database.GetDB().QueryRow(query, key).Scan(
		&setting.ID, &setting.Key, &setting.Value, &setting.Type, &setting.Description, &setting.CreatedAt, &setting.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("setting not found: %w", err)
	}
	return setting, nil
}

// GetAll retrieves all settings
func (r *SettingRepository) GetAll() ([]models.Setting, error) {
	query := `SELECT id, key, value, type, description, created_at, updated_at FROM settings ORDER BY key`
	rows, err := database.GetDB().Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}
	defer rows.Close()

	var settings []models.Setting
	for rows.Next() {
		var s models.Setting
		if err := rows.Scan(&s.ID, &s.Key, &s.Value, &s.Type, &s.Description, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		settings = append(settings, s)
	}
	return settings, nil
}

// Delete deletes a setting
func (r *SettingRepository) Delete(key string) error {
	query := `DELETE FROM settings WHERE key = ?`
	_, err := database.GetDB().Exec(query, key)
	return err
}

// LogRepository handles log database operations
type LogRepository struct{}

// NewLogRepository creates a new LogRepository
func NewLogRepository() *LogRepository {
	return &LogRepository{}
}

// Create creates a new log entry
func (r *LogRepository) Create(log *models.Log) error {
	query := `INSERT INTO logs (level, message, source, metadata) VALUES (?, ?, ?, ?)`
	result, err := database.GetDB().Exec(query, log.Level, log.Message, log.Source, log.Metadata)
	if err != nil {
		return fmt.Errorf("failed to create log: %w", err)
	}
	log.ID, _ = result.LastInsertId()
	return nil
}

// GetAll retrieves logs with optional filtering
func (r *LogRepository) GetAll(level string, limit, offset int) ([]models.Log, error) {
	var query string
	var args []interface{}

	if level != "" {
		query = `SELECT id, level, message, source, metadata, created_at FROM logs WHERE level = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
		args = []interface{}{level, limit, offset}
	} else {
		query = `SELECT id, level, message, source, metadata, created_at FROM logs ORDER BY created_at DESC LIMIT ? OFFSET ?`
		args = []interface{}{limit, offset}
	}

	rows, err := database.GetDB().Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs: %w", err)
	}
	defer rows.Close()

	var logs []models.Log
	for rows.Next() {
		var log models.Log
		var source, metadata sql.NullString
		if err := rows.Scan(&log.ID, &log.Level, &log.Message, &source, &metadata, &log.CreatedAt); err != nil {
			return nil, err
		}
		log.Source = source.String
		log.Metadata = metadata.String
		logs = append(logs, log)
	}
	return logs, nil
}

// DeleteOld deletes logs older than specified days
func (r *LogRepository) DeleteOld(days int) error {
	query := `DELETE FROM logs WHERE created_at < datetime('now', '-' || ? || ' days')`
	_, err := database.GetDB().Exec(query, days)
	return err
}

// TaskRepository handles task database operations
type TaskRepository struct{}

// NewTaskRepository creates a new TaskRepository
func NewTaskRepository() *TaskRepository {
	return &TaskRepository{}
}

// Create creates a new task
func (r *TaskRepository) Create(task *models.Task) error {
	query := `INSERT INTO tasks (name, description, status, priority) VALUES (?, ?, ?, ?)`
	result, err := database.GetDB().Exec(query, task.Name, task.Description, task.Status, task.Priority)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}
	task.ID, _ = result.LastInsertId()
	return nil
}

// GetByID retrieves a task by ID
func (r *TaskRepository) GetByID(id int64) (*models.Task, error) {
	task := &models.Task{}
	query := `SELECT id, name, description, status, priority, progress, result, error, created_at, updated_at, completed_at FROM tasks WHERE id = ?`
	var completedAt sql.NullTime
	err := database.GetDB().QueryRow(query, id).Scan(
		&task.ID, &task.Name, &task.Description, &task.Status, &task.Priority,
		&task.Progress, &task.Result, &task.Error, &task.CreatedAt, &task.UpdatedAt, &completedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}
	if completedAt.Valid {
		task.CompletedAt = &completedAt.Time
	}
	return task, nil
}

// GetAll retrieves all tasks
func (r *TaskRepository) GetAll(status string) ([]models.Task, error) {
	var query string
	var args []interface{}

	if status != "" {
		query = `SELECT id, name, description, status, priority, progress, result, error, created_at, updated_at, completed_at FROM tasks WHERE status = ? ORDER BY priority DESC, created_at DESC`
		args = []interface{}{status}
	} else {
		query = `SELECT id, name, description, status, priority, progress, result, error, created_at, updated_at, completed_at FROM tasks ORDER BY priority DESC, created_at DESC`
	}

	rows, err := database.GetDB().Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		var completedAt sql.NullTime
		if err := rows.Scan(
			&task.ID, &task.Name, &task.Description, &task.Status, &task.Priority,
			&task.Progress, &task.Result, &task.Error, &task.CreatedAt, &task.UpdatedAt, &completedAt,
		); err != nil {
			return nil, err
		}
		if completedAt.Valid {
			task.CompletedAt = &completedAt.Time
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// Update updates a task
func (r *TaskRepository) Update(task *models.Task) error {
	query := `UPDATE tasks SET name = ?, description = ?, status = ?, priority = ?, progress = ?, result = ?, error = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := database.GetDB().Exec(query, task.Name, task.Description, task.Status, task.Priority, task.Progress, task.Result, task.Error, task.ID)
	return err
}

// UpdateStatus updates task status
func (r *TaskRepository) UpdateStatus(id int64, status string) error {
	var query string
	var args []interface{}

	if status == models.TaskStatusCompleted || status == models.TaskStatusFailed {
		query = `UPDATE tasks SET status = ?, updated_at = CURRENT_TIMESTAMP, completed_at = CURRENT_TIMESTAMP WHERE id = ?`
	} else {
		query = `UPDATE tasks SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	}
	args = []interface{}{status, id}

	_, err := database.GetDB().Exec(query, args...)
	return err
}

// UpdateProgress updates task progress
func (r *TaskRepository) UpdateProgress(id int64, progress int) error {
	query := `UPDATE tasks SET progress = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := database.GetDB().Exec(query, progress, id)
	return err
}

// Delete deletes a task
func (r *TaskRepository) Delete(id int64) error {
	query := `DELETE FROM tasks WHERE id = ?`
	_, err := database.GetDB().Exec(query, id)
	return err
}

// FileRepository handles file database operations
type FileRepository struct{}

// NewFileRepository creates a new FileRepository
func NewFileRepository() *FileRepository {
	return &FileRepository{}
}

// Create creates a new file record
func (r *FileRepository) Create(file *models.File) error {
	query := `INSERT INTO files (name, path, size, mime_type, checksum, metadata) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := database.GetDB().Exec(query, file.Name, file.Path, file.Size, file.MimeType, file.Checksum, file.Metadata)
	if err != nil {
		return fmt.Errorf("failed to create file record: %w", err)
	}
	file.ID, _ = result.LastInsertId()
	return nil
}

// GetByID retrieves a file by ID
func (r *FileRepository) GetByID(id int64) (*models.File, error) {
	file := &models.File{}
	query := `SELECT id, name, path, size, mime_type, checksum, metadata, created_at, updated_at FROM files WHERE id = ?`
	var mimeType, checksum, metadata sql.NullString
	err := database.GetDB().QueryRow(query, id).Scan(
		&file.ID, &file.Name, &file.Path, &file.Size, &mimeType, &checksum, &metadata, &file.CreatedAt, &file.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}
	file.MimeType = mimeType.String
	file.Checksum = checksum.String
	file.Metadata = metadata.String
	return file, nil
}

// GetByPath retrieves a file by path
func (r *FileRepository) GetByPath(path string) (*models.File, error) {
	file := &models.File{}
	query := `SELECT id, name, path, size, mime_type, checksum, metadata, created_at, updated_at FROM files WHERE path = ?`
	var mimeType, checksum, metadata sql.NullString
	err := database.GetDB().QueryRow(query, path).Scan(
		&file.ID, &file.Name, &file.Path, &file.Size, &mimeType, &checksum, &metadata, &file.CreatedAt, &file.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}
	file.MimeType = mimeType.String
	file.Checksum = checksum.String
	file.Metadata = metadata.String
	return file, nil
}

// GetAll retrieves all files
func (r *FileRepository) GetAll() ([]models.File, error) {
	query := `SELECT id, name, path, size, mime_type, checksum, metadata, created_at, updated_at FROM files ORDER BY created_at DESC`
	rows, err := database.GetDB().Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get files: %w", err)
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var file models.File
		var mimeType, checksum, metadata sql.NullString
		if err := rows.Scan(
			&file.ID, &file.Name, &file.Path, &file.Size, &mimeType, &checksum, &metadata, &file.CreatedAt, &file.UpdatedAt,
		); err != nil {
			return nil, err
		}
		file.MimeType = mimeType.String
		file.Checksum = checksum.String
		file.Metadata = metadata.String
		files = append(files, file)
	}
	return files, nil
}

// Update updates a file record
func (r *FileRepository) Update(file *models.File) error {
	query := `UPDATE files SET name = ?, path = ?, size = ?, mime_type = ?, checksum = ?, metadata = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := database.GetDB().Exec(query, file.Name, file.Path, file.Size, file.MimeType, file.Checksum, file.Metadata, file.ID)
	return err
}

// Delete deletes a file record
func (r *FileRepository) Delete(id int64) error {
	query := `DELETE FROM files WHERE id = ?`
	_, err := database.GetDB().Exec(query, id)
	return err
}

// DeleteOld deletes files older than specified days
func (r *FileRepository) DeleteOld(days int) error {
	query := `DELETE FROM files WHERE created_at < datetime('now', '-' || ? || ' days')`
	_, err := database.GetDB().Exec(query, days)
	return err
}

// JavaRepository handles Java installation database operations
type JavaRepository struct{}

// NewJavaRepository creates a new JavaRepository
func NewJavaRepository() *JavaRepository {
	return &JavaRepository{}
}

// Create creates a new Java installation
func (r *JavaRepository) Create(java *models.JavaInstallation) error {
	query := `INSERT INTO java_installations (path, friendly_name, version, is_default) VALUES (?, ?, ?, ?)`
	result, err := database.GetDB().Exec(query, java.Path, java.FriendlyName, java.Version, java.IsDefault)
	if err != nil {
		return fmt.Errorf("failed to create java installation: %w", err)
	}
	java.ID, _ = result.LastInsertId()
	return nil
}

// GetByID retrieves a Java installation by ID
func (r *JavaRepository) GetByID(id int64) (*models.JavaInstallation, error) {
	java := &models.JavaInstallation{}
	query := `SELECT id, path, friendly_name, version, is_default, created_at, updated_at FROM java_installations WHERE id = ?`
	err := database.GetDB().QueryRow(query, id).Scan(
		&java.ID, &java.Path, &java.FriendlyName, &java.Version, &java.IsDefault, &java.CreatedAt, &java.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("java installation not found: %w", err)
	}
	return java, nil
}

// GetAll retrieves all Java installations
func (r *JavaRepository) GetAll() ([]models.JavaInstallation, error) {
	query := `SELECT id, path, friendly_name, version, is_default, created_at, updated_at FROM java_installations ORDER BY created_at DESC`
	rows, err := database.GetDB().Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get java installations: %w", err)
	}
	defer rows.Close()

	var javas []models.JavaInstallation
	for rows.Next() {
		var java models.JavaInstallation
		if err := rows.Scan(&java.ID, &java.Path, &java.FriendlyName, &java.Version, &java.IsDefault, &java.CreatedAt, &java.UpdatedAt); err != nil {
			return nil, err
		}
		javas = append(javas, java)
	}
	return javas, nil
}

// GetDefault retrieves the default Java installation
func (r *JavaRepository) GetDefault() (*models.JavaInstallation, error) {
	java := &models.JavaInstallation{}
	query := `SELECT id, path, friendly_name, version, is_default, created_at, updated_at FROM java_installations WHERE is_default = 1`
	err := database.GetDB().QueryRow(query).Scan(
		&java.ID, &java.Path, &java.FriendlyName, &java.Version, &java.IsDefault, &java.CreatedAt, &java.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("default java installation not found: %w", err)
	}
	return java, nil
}

// Update updates a Java installation
func (r *JavaRepository) Update(java *models.JavaInstallation) error {
	query := `UPDATE java_installations SET path = ?, friendly_name = ?, version = ?, is_default = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := database.GetDB().Exec(query, java.Path, java.FriendlyName, java.Version, java.IsDefault, java.ID)
	return err
}

// SetDefault sets a Java installation as the default
func (r *JavaRepository) SetDefault(id int64) error {
	tx, err := database.GetDB().Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE java_installations SET is_default = 0`)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`UPDATE java_installations SET is_default = 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Delete deletes a Java installation
func (r *JavaRepository) Delete(id int64) error {
	query := `DELETE FROM java_installations WHERE id = ?`
	_, err := database.GetDB().Exec(query, id)
	return err
}

// Stats returns database statistics
func Stats() (map[string]int64, error) {
	stats := make(map[string]int64)

	queries := map[string]string{
		"users":    "SELECT COUNT(*) FROM users",
		"settings": "SELECT COUNT(*) FROM settings",
		"logs":     "SELECT COUNT(*) FROM logs",
		"tasks":    "SELECT COUNT(*) FROM tasks",
		"files":    "SELECT COUNT(*) FROM files",
	}

	for name, query := range queries {
		var count int64
		err := database.GetDB().QueryRow(query).Scan(&count)
		if err != nil {
			return nil, err
		}
		stats[name] = count
	}

	return stats, nil
}

// Cleanup performs database cleanup (VACUUM, etc.)
func Cleanup() error {
	_, err := database.GetDB().Exec("VACUUM")
	return err
}

// Backup creates a backup of the database
func Backup(backupPath string) error {
	// For SQLite, we can use the backup API or simply copy the file
	// Here we'll use a simple approach by exporting all data
	db := database.GetDB()

	// Create backup database
	backupDB, err := sql.Open("sqlite", backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup database: %w", err)
	}
	defer backupDB.Close()

	// Copy schema and data
	// This is a simplified backup - in production, use SQLite's backup API
	tables := []string{"users", "settings", "logs", "tasks", "files"}

	for _, table := range tables {
		// Get create table statement
		var createStmt string
		err := db.QueryRow(fmt.Sprintf("SELECT sql FROM sqlite_master WHERE type='table' AND name='%s'", table)).Scan(&createStmt)
		if err != nil {
			continue
		}

		// Create table in backup
		_, err = backupDB.Exec(createStmt)
		if err != nil {
			return fmt.Errorf("failed to create table %s in backup: %w", table, err)
		}

		// Copy data
		rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", table))
		if err != nil {
			continue
		}
		defer rows.Close()

		columns, _ := rows.Columns()
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		for rows.Next() {
			rows.Scan(valuePtrs...)
			placeholders := ""
			args := make([]interface{}, len(values))
			for i, v := range values {
				if i > 0 {
					placeholders += ", "
				}
				placeholders += "?"
				args[i] = v
			}
			backupDB.Exec(fmt.Sprintf("INSERT INTO %s VALUES (%s)", table, placeholders), args...)
		}
	}

	return nil
}

// GetDBStats returns database connection statistics
func GetDBStats() sql.DBStats {
	return database.GetDB().Stats()
}

// SetNow sets the current time for testing
var Now = time.Now
