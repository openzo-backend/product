package service

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/product/internal/models"
	"github.com/tanush-128/openzo_backend/product/internal/pb"
	"github.com/tanush-128/openzo_backend/product/internal/repository"
	"github.com/tanush-128/openzo_backend/product/internal/utils"
)

type ProductService interface {

	//CRUD
	CreateProduct(ctx *gin.Context, req models.Product) (models.Product, error)
	GetProductByID(ctx *gin.Context, id string) (models.Product, error)
	GetProductsByStoreID(ctx *gin.Context, storeID string) ([]models.Product, error)
	UpdateProduct(ctx *gin.Context, req models.Product) (models.Product, error)
	DeleteProduct(ctx *gin.Context, id string) error
}

type productService struct {
	ProductRepository repository.ProductRepository
	imageClient       pb.ImageServiceClient
}

func NewProductService(ProductRepository repository.ProductRepository,
	imageClient pb.ImageServiceClient,
) ProductService {
	return &productService{ProductRepository: ProductRepository, imageClient: imageClient}
}

func (s *productService) CreateProduct(ctx *gin.Context, req models.Product) (models.Product, error) {
	form, err := ctx.MultipartForm()
	if err != nil {

		return models.Product{}, err
	}
	req.Images = []models.ProductImage{}
	for _, file := range form.File["images"] {
		log.Println(file.Filename)

		imageBytes, err := utils.FileHeaderToBytes(file)

		if err != nil {
			return models.Product{}, err
		}
		imageURL, err := s.imageClient.UploadImage(ctx, &pb.ImageMessage{
			ImageData: imageBytes,
		})
		if err != nil {
			return models.Product{}, err
		}

		req.Images = append(req.Images, models.ProductImage{
			Image: imageURL.Url,
		})

	}

	createdProduct, err := s.ProductRepository.CreateProduct(req)
	if err != nil {
		return models.Product{}, err // Propagate error
	}

	return createdProduct, nil
}

func (s *productService) UpdateProduct(ctx *gin.Context, req models.Product) (models.Product, error) {
	updatedProduct, err := s.ProductRepository.UpdateProduct(req)
	if err != nil {
		return models.Product{}, err
	}

	return updatedProduct, nil
}

func (s *productService) DeleteProduct(ctx *gin.Context, id string) error {
	err := s.ProductRepository.DeleteProduct(id)
	if err != nil {
		return err
	}

	return nil
}
