package service

import (
	"time"

	"github.com/google/uuid"
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

func (s *ExpenseService) GetExpenses(limit int, offset int, categoryID uuid.UUID, dateFrom time.Time, dateTo time.Time) ([]model.Expense, error) {

	dateFrom, dateTo = verifyAndAddEndDateInclusive(dateFrom, dateTo)
	return s.repository.GetExpenses(limit, offset, categoryID, dateFrom, dateTo)
}

func (s *ExpenseService) GetTotalExpense(categoryID uuid.UUID, dateFrom time.Time, dateTo time.Time) (int64, error) {

	dateFrom, dateTo = verifyAndAddEndDateInclusive(dateFrom, dateTo)
	return s.repository.GetTotalExpense(categoryID, dateFrom, dateTo)
}

func (s *ExpenseService) UpdateExpense(expense model.Expense) (*model.Expense, error) {
	return s.repository.UpdateExpense(expense)
}

func (s *ExpenseService) DeleteExpense(id uuid.UUID) error {
	return s.repository.DeleteExpense(id)
}

func verifyAndAddEndDateInclusive(dateFrom time.Time, dateTo time.Time) (time.Time, time.Time) {

	if dateFrom.IsZero() || dateTo.IsZero() {
		dateFrom = time.Time{}
		dateTo = time.Time{}
	} else {
		dateTo = exclusiveEndDate(dateTo)
	}

	return dateFrom, dateTo
}

func exclusiveEndDate(dateTo time.Time) time.Time {
	if dateTo.IsZero() {
		return dateTo
	}

	return dateTo.AddDate(0, 0, 1)
}
