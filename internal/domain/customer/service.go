package customer

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

	if customer.ID == "" {
		customer.ID = uuid.New().String()
	}

	if customer.Email != "" {
		if c, _ := s.repo.GetByEmail(customer.Email); c != nil && c.ID != customer.ID {
			return errors.New("customer already exists with this email")
		}
	}
	if customer.Mobile != "" {
		if c, _ := s.repo.GetByMobile(customer.Mobile); c != nil && c.ID != customer.ID {
			return errors.New("customer already exists with this mobile")
		}
	}

	if customer.Password != "" {
		hp, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		customer.Password = string(hp)
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

	if customer.Email != "" {
		if c, _ := s.repo.GetByEmail(customer.Email); c != nil && c.ID != customer.ID {
			return errors.New("customer already exists with this email")
		}
	}
	if customer.Mobile != "" {
		if c, _ := s.repo.GetByMobile(customer.Mobile); c != nil && c.ID != customer.ID {
			return errors.New("customer already exists with this mobile")
		}
	}

	customer.UpdatedAt = s.now().UTC()
	return s.repo.Update(customer)
}
