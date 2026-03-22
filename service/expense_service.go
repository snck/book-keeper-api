package service

import "github.com/snck/book-keeper-api/repository"

type ExpenseService struct {
	repository *repository.ExpenseRepository
}

func NewExpenseService(repository *repository.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repository: repository}
}
