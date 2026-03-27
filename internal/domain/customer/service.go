package customer

import (
	"errors"
	"time"
)

var (
	ErrCustomerNotFound = errors.New("customer not found")
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo, now: time.Now}
}

func (s *Service) CreateCustomer(customer *Customer) error {
	if err := customer.Validate(); err != nil {
		return err
	}
	customer.CreatedAt = s.now().UTC()
	customer.UpdatedAt = s.now().UTC()
	return s.repo.Create(customer)
}

func (s *Service) GetCustomer(id string) (*Customer, error) {
	return s.repo.GetByID(id)
}

func (s *Service) UpdateCustomer(customer *Customer) error {
	if err := customer.Validate(); err != nil {
		return err
	}
	customer.UpdatedAt = s.now().UTC()
	return s.repo.Update(customer)
}
