package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// SaveMultipartImage saves an uploaded multipart image file to disk.
func SaveMultipartImage(
	file multipart.File,
	header *multipart.FileHeader,
	dir string,
	filename string,
) (string, error) {

	if file == nil || header == nil {
		return "", errors.New("invalid file upload")
	}

	// 1. Get the extension from the ACTUAL uploaded file (Source of Truth)
	ext := strings.ToLower(filepath.Ext(header.Filename))

	// 2. Validate extension
	allowed := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	if !allowed[ext] {
		return "", fmt.Errorf("unsupported image type: %s", ext)
	}

	// 3. Ensure directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload dir: %w", err)
	}

	// 4. Generate Filename
	if filename == "" {
		return "", errors.New("file name required")
	} else {
		// FIX START: Prevent double extensions (e.g., "image.jpg.jpg")

		// Clean the input
		filename = sanitizeFilename(filename)

		// Remove any extension the user might have manually typed
		// e.g., if input is "my-pic.png", this makes it "my-pic"
		filename = strings.TrimSuffix(filename, filepath.Ext(filename))

		// Append the correct extension from the file header
		filename = filename + ext
		// FIX END
	}

	path := filepath.Join(dir, filename)

	// 5. Create destination file
	dst, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// 6. Reset file pointer
	if seeker, ok := file.(io.Seeker); ok {
		if _, err := seeker.Seek(0, 0); err != nil {
			return "", fmt.Errorf("failed to seek file: %w", err)
		}
	}

	// 7. Copy contents
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return path, nil
}

// DeleteFile removes a file from the specified path.
func DeleteFile(path string) error {
	if path == "" {
		return errors.New("file path cannot be empty")
	}

	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", path)
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}
