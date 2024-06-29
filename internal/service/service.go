package service

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
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
	ChangeProductQuantity(ctx *gin.Context, id string, quantity int) error
	UpdateProduct(ctx *gin.Context, req models.Product) (models.Product, error)
	UpdateDisplayOrder(ctx *gin.Context, id string, displayOrder int) error
	BatchUpdateDisplayOrder(ctx *gin.Context, updates []models.Product) error
	DeleteProduct(ctx *gin.Context, id string) error
}

type productService struct {
	ProductRepository repository.ProductRepository
	imageClient       pb.ImageServiceClient
	kafkaProducer     *kafka.Producer
}

func NewProductService(ProductRepository repository.ProductRepository,
	imageClient pb.ImageServiceClient, kafkaProducer *kafka.Producer,
) ProductService {
	return &productService{ProductRepository: ProductRepository, imageClient: imageClient, kafkaProducer: kafkaProducer}
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
	go writeProductToKafka(s.kafkaProducer, createdProduct)
	return createdProduct, nil
}

func writeProductToKafka(p *kafka.Producer, product models.Product) {
	// Produce messages to topic (asynchronously)
	topic := "products"

	value, _ := json.Marshal(product)

	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(product.ID),
		Value:          value,
	}, nil)

}

func (s *productService) UpdateProduct(ctx *gin.Context, req models.Product) (models.Product, error) {

	// product, err := s.ProductRepository.GetProductByID(req.ID)
	// if err != nil {
	// 	return models.Product{}, err
	// }
	log.Printf("Product Images: %+v", req.Images)

	form, err := ctx.MultipartForm()
	if err != nil {

		return models.Product{}, err
	}

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

	updatedProduct, err := s.ProductRepository.UpdateProduct(req)
	if err != nil {
		return models.Product{}, err
	}
	// updatedProduct.Images = req.Images

	go writeProductToKafka(s.kafkaProducer, updatedProduct)
	return updatedProduct, nil
}

func (s *productService) UpdateDisplayOrder(ctx *gin.Context, id string, displayOrder int) error {
	err := s.ProductRepository.UpdateDisplayOrder(id, displayOrder)
	if err != nil {
		return err
	}

	// go writeProductToKafka(s.kafkaProducer, updatedProduct)
	return nil
}

func (s *productService) BatchUpdateDisplayOrder(ctx *gin.Context, updates []models.Product) error {
	err := s.ProductRepository.BatchUpdateDisplayOrder(updates)
	if err != nil {
		return err
	}

	return nil
}

func (s *productService) ChangeProductQuantity(ctx *gin.Context, id string, quantity int) error {
	err := s.ProductRepository.ChangeProductQuantity(id, quantity)
	if err != nil {
		return err
	}

	return nil
}

func (s *productService) DeleteProduct(ctx *gin.Context, id string) error {
	err := s.ProductRepository.DeleteProduct(id)
	if err != nil {
		return err
	}

	return nil
}
