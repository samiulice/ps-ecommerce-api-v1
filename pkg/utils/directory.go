package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// GetPublicImagesDirectoryPath returns the absolute or relative path to the
// root public images directory and ensures it exists.
// Example: "assets/public/images"
func GetPublicImagesDirectoryPath() string {
	// filepath.Join ensures the path separators are correct for the OS
	path := filepath.Join(".", "assets", "public", "images")

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
	basePath := filepath.Join(".", "assets", "public", "images", "categories")

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

// GetCategoryThumbnailURL constructs the url for a category image.
// It accepts a filename, ext, sanitizes it, and appends it to the category directory.
func GetCategoryThumbnailURL(filename, ext string) string {
	filename = sanitizeFilename(filename)
	return strings.Join([]string{"public", "images", "categories", filename + ext}, "/")

}

// GetBrandFolderPath constructs the full file path for a brand image.
// It accepts a filename, sanitizes it, and appends it to the brand directory.
// If the filename is empty, it returns the directory path itself.
// It also ensures the directory exists.
func GetBrandFolderPath(filename string) string {
	// Define the specific subdirectory for brands
	basePath := filepath.Join(".", "assets", "public", "images", "brands")

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

// GetBrandThumbnailURL constructs the url for a brand image.
// It accepts a filename, ext, sanitizes it, and appends it to the brand directory.
func GetBrandThumbnailURL(filename, ext string) string {
	filename = sanitizeFilename(filename)
	return strings.Join([]string{"public", "images", "brands", filename + ext}, "/")

}

// GetProductFolderPath constructs the full file path for a Product images.
// It accepts a filename, sanitizes it, and appends it to the product directory.
// If the filename is empty, it returns the directory path itself.
// It also ensures the directory exists.
func GetProductFolderPath(filename string) string {
	// Define the specific subdirectory for products
	basePath := filepath.Join(".", "assets", "public", "images", "products")

	// Ensure directory exists
	_ = os.MkdirAll(basePath, os.ModePerm)

	// Security: Clean the filename to remove unsafe characters
	filename = sanitizeFilename(filename)

	// If no filename is provided after sanitization, return the folder path
	if filename == "" {
		return basePath
	}

	// Combine the base path with the specific filename
	return filepath.Join(basePath, filename)
}

// GetProductThumbnailURL constructs the url for a product image.
// It accepts a filename, ext, sanitizes it, and appends it to the product directory.
func GetProductThumbnailURL(filename, ext string) string {
	filename = sanitizeFilename(filename)
	return strings.Join([]string{"public", "images", "products", filename + ext}, "/")
}

// GetCustomerFolderPath constructs the full file path for a customer image.
// It accepts a filename, sanitizes it, and appends it to the customer directory.
// If the filename is empty, it returns the directory path itself.
// It also ensures the directory exists.
func GetCustomerFolderPath(filename string) string {
	// Define the specific subdirectory for customers
	basePath := filepath.Join(".", "assets", "public", "images", "customers")

	// Ensure directory exists
	_ = os.MkdirAll(basePath, os.ModePerm)

	// Security: Clean the filename to remove unsafe characters
	filename = sanitizeFilename(filename)

	// If no filename is provided after sanitization, return the folder path
	if filename == "" {
		return basePath
	}

	// Combine the base path with the specific filename
	return filepath.Join(basePath, filename)
}

// GetCustomerImageURL constructs the url for a customer image.
// It accepts a filename, ext, sanitizes it, and appends it to the customer directory.
func GetCustomerImageURL(filename, ext string) string {
	filename = sanitizeFilename(filename)
	return strings.Join([]string{"public", "images", "customers", filename + ext}, "/")
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
			r == '-' || r == '_' || r == '.' {
			clean = append(clean, r)
		}
	}
	return string(clean)
}

// GetHeroFolderPath constructs the full file path for a hero banner image.
// It accepts a filename, sanitizes it, and appends it to the hero directory.
// If the filename is empty, it returns the directory path itself.
// It also ensures the directory exists.
func GetHeroFolderPath(filename string) string {
	// Define the specific subdirectory for hero banners
	basePath := filepath.Join(".", "assets", "public", "images", "hero")

	// Ensure directory exists
	_ = os.MkdirAll(basePath, os.ModePerm)

	// Security: Clean the filename to remove unsafe characters
	filename = sanitizeFilename(filename)

	// If no filename is provided after sanitization, return the folder path
	if filename == "" {
		return basePath
	}

	// Combine the base path with the specific filename
	return filepath.Join(basePath, filename)
}

// GetHeroBannerURL constructs the url for a hero banner image.
// It accepts a filename, ext, sanitizes it, and appends it to the hero directory.
func GetHeroBannerURL(filename, ext string) string {
	filename = sanitizeFilename(filename)
	return strings.Join([]string{"public", "images", "hero", filename + ext}, "/")
}

// GetPurchaseAttachmentFolderPath constructs the full file path for purchase attachments.
// If filename is empty, it returns the directory path itself.
func GetPurchaseAttachmentFolderPath(filename string) string {
	basePath := filepath.Join(".", "assets", "public", "images", "purchases")

	_ = os.MkdirAll(basePath, os.ModePerm)

	filename = sanitizeFilename(filename)
	if filename == "" {
		return basePath
	}

	return filepath.Join(basePath, filename)
}

// GetPurchaseAttachmentURL constructs the URL for a purchase attachment.
func GetPurchaseAttachmentURL(filename, ext string) string {
	filename = sanitizeFilename(filename)
	return strings.Join([]string{"public", "images", "purchases", filename + ext}, "/")
}
