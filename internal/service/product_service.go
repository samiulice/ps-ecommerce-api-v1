package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
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

// Helper to process gallery uploads
func (s *ProductService) processGalleryImages(files []*multipart.FileHeader, productName string) ([]string, error) {
	var imagePaths []string
	for i, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open gallery image %d: %w", i+1, err)
		}
		defer file.Close()

		// Generate unique name for this gallery image
		imageUniqueName := fmt.Sprintf("%s-gallery-%d-%s", productName, i+1, utils.GenerateRandomString(5))

		// Save image with the unique name
		createdPath, saveErr := utils.SaveMultipartImage(file, fileHeader, utils.GetProductFolderPath(""), imageUniqueName)
		if saveErr != nil {
			return nil, fmt.Errorf("failed to save gallery image %d: %w", i+1, saveErr)
		}
		fmt.Println("created gallery image: ", createdPath)

		// Generate URL/Path using the same unique name
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		path := utils.GetProductThumbnailURL(imageUniqueName, ext)
		imagePaths = append(imagePaths, path)
	}
	return imagePaths, nil
}

// Create creates a new product with thumbnail and optional gallery images
func (s *ProductService) Create(ctx context.Context, p *model.Product, thumbFile multipart.File, thumbHeader *multipart.FileHeader, galleryFiles []*multipart.FileHeader) error {
	if p.Name == "" {
		return errors.New("name is required")
	}
	if p.SKU == "" {
		// Generate SKU if missing (basic example)
		p.SKU = strings.ToUpper(strings.ReplaceAll(p.Name, " ", "-") + "-" + utils.GenerateRandomString(4))
	}

	// 1. Handle Thumbnail
	if thumbFile != nil {
		ext := strings.ToLower(filepath.Ext(thumbHeader.Filename))
		thumbnailURL := utils.GetProductThumbnailURL(p.Name, ext)
		p.Thumbnail = thumbnailURL
	}

	// 2. Handle Gallery Images
	if len(galleryFiles) > 0 {
		paths, err := s.processGalleryImages(galleryFiles, p.Name)
		if err != nil {
			return err
		}
		p.GalleryImages = paths
		fmt.Println("paths:", paths)
	} else {
		fmt.Printf("No Gallery Images Found\n")
	}
	// 3. Insert into database
	err := s.repo.Create(ctx, p)

	// 4. Save Thumbnail File on successful DB insert
	if err == nil && thumbFile != nil {
		defer thumbFile.Close()
		_, saveErr := utils.SaveMultipartImage(thumbFile, thumbHeader, utils.GetProductFolderPath(""), p.Name)
		if saveErr != nil {
			return saveErr
		}
	}

	return err
}

// Update modifies an existing product
func (s *ProductService) Update(ctx context.Context, p *model.Product, thumbFile multipart.File, thumbHeader *multipart.FileHeader, galleryFiles []*multipart.FileHeader) error {
	if p.ID == 0 {
		return errors.New("id is required")
	}

	existingProduct, err := s.GetByID(ctx, p.ID)
	if err != nil {
		return err
	}

	// Handle Thumbnail Update
	if thumbFile != nil && thumbHeader != nil {
		defer thumbFile.Close()
		ext := strings.ToLower(filepath.Ext(thumbHeader.Filename))
		thumbnailURL := utils.GetProductThumbnailURL(p.Name, ext)
		p.Thumbnail = thumbnailURL

		// Save thumbnail file first (before DB update)
		_, saveErr := utils.SaveMultipartImage(thumbFile, thumbHeader, utils.GetProductFolderPath(""), p.Name)
		if saveErr != nil {
			return fmt.Errorf("failed to save thumbnail: %w", saveErr)
		}

		// Delete old image if different
		if existingProduct.Thumbnail != "" && existingProduct.Thumbnail != p.Thumbnail {
			utils.DeleteFile(utils.GetProductFolderPath(filepath.Base(existingProduct.Thumbnail)))
		}
	} else {
		p.Thumbnail = existingProduct.Thumbnail // Keep old if not provided
	}

	// Handle Gallery Update (Append or Replace logic - implementing Replace for simplicity)
	if len(galleryFiles) > 0 {
		paths, err := s.processGalleryImages(galleryFiles, p.Name)
		if err != nil {
			return err
		}
		p.GalleryImages = paths
	} else {
		p.GalleryImages = existingProduct.GalleryImages
	}

	return s.repo.Update(ctx, p)
}

// Delete removes a product
func (s *ProductService) Delete(ctx context.Context, id int64) error {
	existingProduct, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.repo.Delete(ctx, id)

	// Cleanup files
	if err == nil {
		if existingProduct.Thumbnail != "" {
			utils.DeleteFile(utils.GetProductFolderPath(filepath.Base(existingProduct.Thumbnail)))
		}
		for _, img := range existingProduct.GalleryImages {
			utils.DeleteFile(utils.GetProductFolderPath(filepath.Base(img)))
		}
	}
	return err
}

// DeleteGalleryImage removes a specific gallery image from product and deletes the file
func (s *ProductService) DeleteGalleryImage(ctx context.Context, productID int64, imagePath string) error {
	product, err := s.GetByID(ctx, productID)
	if err != nil {
		return err
	}

	// Find and remove the image from gallery
	found := false
	newGallery := make([]string, 0, len(product.GalleryImages))
	for _, img := range product.GalleryImages {
		if img == imagePath {
			found = true
			continue
		}
		newGallery = append(newGallery, img)
	}

	if !found {
		return errors.New("image not found in product gallery")
	}

	// Update product with new gallery
	product.GalleryImages = newGallery
	err = s.repo.Update(ctx, product)
	if err != nil {
		return err
	}

	// Delete the actual file
	utils.DeleteFile(utils.GetProductFolderPath(filepath.Base(imagePath)))

	return nil
}

func (s *ProductService) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	return s.repo.GetByID(ctx, id)
}

// GetProductVariationsByProductID calls database func to read the product variations from the database
func (s *ProductService) GetProductVariationsByProductID(ctx context.Context, id int64) ([]*model.ProductVariation, error) {
	return s.repo.GetProductVariationsByProductID(ctx, id)
}

func (s *ProductService) GetProducts(ctx context.Context, filter model.ProductFilter) ([]*model.Product, int64, error) {
	return s.repo.GetProducts(ctx, filter)
}
