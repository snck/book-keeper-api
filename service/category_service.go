package service

import (
	"github.com/snck/book-keeper-api/model"
	"github.com/snck/book-keeper-api/repository"
)

type CategoryService struct {
	repository *repository.CategoryRepository
}

func NewCategoryService(repository *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repository: repository}
}

func (s *CategoryService) GetCategories() ([]model.Category, error) {
	return s.repository.GetCategories()
}
