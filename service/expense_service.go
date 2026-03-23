package service

import (
	"github.com/snck/book-keeper-api/model"
	"github.com/snck/book-keeper-api/repository"
)

type ExpenseService struct {
	repository *repository.ExpenseRepository
}

func NewExpenseService(repository *repository.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repository: repository}
}

func (s *ExpenseService) CreateExpense(expense model.Expense) (model.Expense, error) {
	return s.repository.CreateExpense(expense)
}
