package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CopyDBOptions configures the copy behavior
type CopyDBOptions struct {
	// MaxRetries is the number of times to retry opening before copying
	MaxRetries int
	// RetryDelay is the initial delay between retries (doubles each retry)
	RetryDelay time.Duration
}

// DefaultCopyDBOptions returns sensible defaults
func DefaultCopyDBOptions() CopyDBOptions {
	return CopyDBOptions{
		MaxRetries: 3,
		RetryDelay: 100 * time.Millisecond,
	}
}

// CopyDBToTemp copies a database file to a temporary location.
// This is useful when the browser has the database locked.
// Returns the path to the temporary file, which the caller should delete when done.
func CopyDBToTemp(srcPath string) (string, error) {
	// Create temp file with same extension
	ext := filepath.Ext(srcPath)
	base := filepath.Base(srcPath)
	base = strings.TrimSuffix(base, ext)

	tmpFile, err := os.CreateTemp("", fmt.Sprintf("kooky-%s-*%s", base, ext))
	if err != nil {
		return "", fmt.Errorf("creating temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	// Open source file
	srcFile, err := os.Open(srcPath)
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return "", fmt.Errorf("opening source: %w", err)
	}
	defer srcFile.Close()

	// Copy contents
	_, err = io.Copy(tmpFile, srcFile)
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return "", fmt.Errorf("copying database: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("closing temp file: %w", err)
	}

	return tmpPath, nil
}

// IsDBLocked checks if an error indicates a locked database.
// This checks for common SQLite "database is locked" error patterns.
func IsDBLocked(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "database is locked") ||
		strings.Contains(errStr, "SQLITE_BUSY") ||
		strings.Contains(errStr, "locked") ||
		strings.Contains(errStr, "busy")
}
