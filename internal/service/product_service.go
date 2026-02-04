package service

import (
	"context"
	"errors"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/projuktisheba/pse-api-v1/internal/model"
	"github.com/projuktisheba/pse-api-v1/internal/repository"
	"github.com/projuktisheba/pse-api-v1/pkg/utils"
)

type ProductService struct {
	repo *repository.ProductRepo
}

func NewProductService(repo *repository.ProductRepo) *ProductService {
	return &ProductService{repo: repo}
}

// generateSlug creates a URL-friendly slug from the product name
func (s *ProductService) generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.TrimSpace(slug)
	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove special characters
	reg := regexp.MustCompile(`[^a-z0-9\-]`)
	slug = reg.ReplaceAllString(slug, "")
	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")
	return slug
}

// Create creates a new product with optional thumbnail
func (s *ProductService) Create(ctx context.Context, p *model.Product, file multipart.File, header *multipart.FileHeader) error {
	if p.Name == nil || *p.Name == "" {
		return errors.New("name is required")
	}

	// Generate slug from name
	slug := s.generateSlug(*p.Name)
	p.Slug = &slug

	// Set thumbnail URL if file exists
	if file != nil {
		ext := strings.ToLower(filepath.Ext(header.Filename))
		thumbnailURL := utils.GetProductThumbnailURL(*p.Name, ext)
		p.Thumbnail = &thumbnailURL
	}

	// Insert into database
	err := s.repo.Create(ctx, p)

	// Save file on successful database insertion
	if err == nil && file != nil {
		defer file.Close()
		_, saveErr := utils.SaveMultipartImage(file, header, utils.GetProductFolderPath(""), *p.Name)
		if saveErr != nil {
			return saveErr
		}
	}

	return err
}

// Update modifies an existing product with optional new thumbnail
func (s *ProductService) Update(ctx context.Context, p *model.Product, file multipart.File, header *multipart.FileHeader) error {
	if p.ID == 0 {
		return errors.New("id is required")
	}

	// Fetch existing product to preserve old thumbnail if not updating
	existingProduct, err := s.GetByID(ctx, p.ID)
	if err != nil {
		return err
	}

	// Regenerate slug if name changed
	if p.Name != nil && *p.Name != "" {
		slug := s.generateSlug(*p.Name)
		p.Slug = &slug
	}

	// Set new thumbnail URL if file exists
	if file != nil {
		ext := strings.ToLower(filepath.Ext(header.Filename))
		thumbnailURL := utils.GetProductThumbnailURL(*p.Name, ext)
		p.Thumbnail = &thumbnailURL
	}

	err = s.repo.Update(ctx, p)

	// Handle new image
	if err == nil && file != nil {
		defer file.Close()
		_, saveErr := utils.SaveMultipartImage(file, header, utils.GetProductFolderPath(""), *p.Name)
		if saveErr != nil {
			return saveErr
		}
		// Delete old image if thumbnail changed
		if existingProduct.Thumbnail != nil && p.Thumbnail != nil && *existingProduct.Thumbnail != *p.Thumbnail {
			utils.DeleteFile(utils.GetProductFolderPath(filepath.Base(*existingProduct.Thumbnail)))
		}
	}

	return err
}

// Delete removes a product by ID
func (s *ProductService) Delete(ctx context.Context, id int64) error {
	// Fetch existing product to get thumbnail path
	existingProduct, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.repo.Delete(ctx, id)
	if err == nil && existingProduct.Thumbnail != nil {
		utils.DeleteFile(utils.GetProductFolderPath(filepath.Base(*existingProduct.Thumbnail)))
	}
	return err
}

// GetByID retrieves a single product by ID
func (s *ProductService) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	return s.repo.GetByID(ctx, id)
}

// GetProducts retrieves products with optional filters and pagination
func (s *ProductService) GetProducts(ctx context.Context, filter model.ProductFilter) ([]*model.Product, int64, error) {
	return s.repo.GetProducts(ctx, filter)
}
