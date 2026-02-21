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

type BrandService struct {
	repo *repository.BrandRepo
}

func NewBrandService(repo *repository.BrandRepo) *BrandService {
	return &BrandService{repo: repo}
}

func (s *BrandService) Create(ctx context.Context, c *model.Brand, file multipart.File, header *multipart.FileHeader) error {
	if c.Name == "" {
		return errors.New("name is required")
	}

	//set thumbnail url if exist
	if file != nil {
		// Get the extension from the ACTUAL uploaded file (Source of Truth)
		ext := strings.ToLower(filepath.Ext(header.Filename))

		c.Thumbnail = utils.GetBrandThumbnailURL(c.Name, ext)
	}

	// insert the data
	err := s.repo.Create(ctx, c)

	//Save the file if exist for successful database insertion
	if err == nil && file != nil {
		defer file.Close()
		_, err := utils.SaveMultipartImage(file, header, utils.GetBrandFolderPath(""), c.Name)
		if err != nil {
			return err
		}
	}

	return err
}
func (s *BrandService) Update(ctx context.Context, c *model.Brand, file multipart.File, header *multipart.FileHeader) error {
	if c.ID == 0 {
		return errors.New("id is required")
	}

	// 1. Fetch Existing Data to preserve old Logo if not updating
	existingBrand, err := s.GetByID(ctx, c.ID)
	if err != nil {
		return err
	}

	//set thumbnail url if exist
	if file != nil {
		// Get the extension from the ACTUAL uploaded file (Source of Truth)
		ext := strings.ToLower(filepath.Ext(header.Filename))

		c.Thumbnail = utils.GetBrandThumbnailURL(c.Name, ext)
	}
	err = s.repo.Update(ctx, c)
	// 4. Handle New Image
	if err == nil && file != nil {
		defer file.Close()
		// Optional: Delete old image here using os.Remove(existingBrand.LogoURL)
		_, err := utils.SaveMultipartImage(file, header, utils.GetBrandFolderPath(""), c.Name)
		if err != nil {
			return err
		}
		if existingBrand.Thumbnail != c.Thumbnail {
			utils.DeleteFile(utils.GetBrandFolderPath(filepath.Base(existingBrand.Thumbnail)))
		}

	}

	return err
}
func (s *BrandService) Delete(ctx context.Context, id int64) error {
	// 1. Fetch Existing Data
	existingBrand, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}
	err = s.repo.Delete(ctx, id)
	if err == nil {
		fmt.Println("image to be deleted: ", utils.GetBrandFolderPath(filepath.Base(existingBrand.Thumbnail)))
		return utils.DeleteFile(utils.GetBrandFolderPath(filepath.Base(existingBrand.Thumbnail)))
	}
	return err
}
func (s *BrandService) GetByID(ctx context.Context, id int64) (*model.Brand, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *BrandService) GetBrands(ctx context.Context, status string) ([]*model.Brand, error) {
	return s.repo.GetBrands(ctx, status)
}
