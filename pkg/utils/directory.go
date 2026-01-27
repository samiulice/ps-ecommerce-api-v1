package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// GetImagesDirectoryPath returns the absolute or relative path to the
// root images directory and ensures it exists.
// Example: "assets/images"
func GetImagesDirectoryPath() string {
	// filepath.Join ensures the path separators are correct for the OS
	path := filepath.Join(".", "assets", "images")

	// Create directory if it does not exist
	_ = os.MkdirAll(path, os.ModePerm)

	return path
}

// GetCategoryFolderPath constructs the full file path for a category image.
// It accepts a filename, sanitizes it, and appends it to the category directory.
// If the filename is empty, it returns the directory path itself.
// It also ensures the directory exists.
func GetCategoryFolderPath(filename string) string {
	// Define the specific subdirectory for categories
	basePath := filepath.Join(".", "assets", "images", "category")

	// Ensure directory exists
	_ = os.MkdirAll(basePath, os.ModePerm)

	// Security: Clean the filename to remove unsafe characters
	// Note: sanitizeFilename should already exist in utils
	filename = sanitizeFilename(filename)

	// If no filename is provided after sanitization, return the folder path
	if filename == "" {
		return basePath
	}

	// Combine the base path with the specific filename
	return filepath.Join(basePath, filename)
}


// GetProductFolderPath constructs the full file path for a Product images.
// It accepts a filename, sanitizes it, and appends it to the product directory.
// If the filename is empty, it returns the directory path itself.
// It also ensures the directory exists.
func GetProductFolderPath(filename string) string {
	// Define the specific subdirectory for products
	basePath := filepath.Join(".", "assets", "images", "products")

	// Ensure directory exists
	_ = os.MkdirAll(basePath, os.ModePerm)

	// Security: Clean the filename to remove unsafe characters
	// Note: sanitizeFilename should already exist in utils
	filename = sanitizeFilename(filename)

	// If no filename is provided after sanitization, return the folder path
	if filename == "" {
		return basePath
	}

	// Combine the base path with the specific filename
	return filepath.Join(basePath, filename)
}


// sanitizeFilename removes unsafe characters from filenames
func sanitizeFilename(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")

	// Keep only safe characters
	clean := make([]rune, 0, len(name))
	for _, r := range name {
		if (r >= 'a' && r <= 'z') ||
			(r >= '0' && r <= '9') ||
			r == '-' || r == '_' {
			clean = append(clean, r)
		}
	}
	return string(clean)
}
