package customer

import (
	"errors"
	"time"

	"context"

	"ppharma/backend/internal/domain/common"
	"ppharma/backend/support-pkg/notification"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrCustomerNotFound   = errors.New("customer not found")
	LOGIN_OTP_EXPIRE_TIME = 5 * time.Minute
)

type Service struct {
	repo        Repository
	emailSender notification.EmailSender
	smsSender   notification.SMSSender
	now         func() time.Time
}

func NewService(repo Repository, emailSender notification.EmailSender, smsSender notification.SMSSender) *Service {
	return &Service{repo: repo, emailSender: emailSender, smsSender: smsSender, now: time.Now}
}

func (s *Service) CreateCustomer(customer *Customer) error {
	if err := customer.Validate(); err != nil {
		return err
	}

	if customer.ID == "" {
		customer.ID = uuid.New().String()
	}

	if customer.Email != "" {
		if c, _ := s.GetCustomerByIdentifier(customer.Email); c != nil && c.ID != customer.ID {
			return errors.New("customer already exists with this email")
		}
	}
	if customer.Mobile != "" {
		if c, _ := s.GetCustomerByIdentifier(customer.Mobile); c != nil && c.ID != customer.ID {
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

	customer.IsVerified = false
	customer.LoginInfo.OTP = notification.GenerateOTP()
	customer.LoginInfo.OTPExpiry = s.now().UTC().Add(LOGIN_OTP_EXPIRE_TIME)

	customer.CreatedAt = s.now().UTC()
	customer.UpdatedAt = s.now().UTC()
	err := s.repo.Create(customer)
	if err == nil {
		if customer.Mobile != "" && s.smsSender != nil {
			go s.smsSender.SendSMS(context.Background(), customer.Mobile, "Your code is: "+customer.LoginInfo.OTP)
		}

		// if customer.Email != "" && s.emailSender != nil {
		// 	go s.emailSender.SendEmail(context.Background(), []string{customer.Email}, "Your OTP Code", "Your code is: "+customer.LoginInfo.OTP, false)
		// }
	}
	return err
}

func (s *Service) VerifyOTP(identifier string, otp string) (*Customer, error) {
	cust, err := s.GetCustomerByIdentifier(identifier)
	if err != nil {
		return nil, err
	}

	if cust.IsVerified {
		return nil, errors.New("customer already verified")
	}

	if cust.LoginInfo.OTP != otp {
		return nil, errors.New("invalid otp")
	}

	if s.now().UTC().After(cust.LoginInfo.OTPExpiry) {
		return nil, errors.New("otp expired")
	}

	cust.IsVerified = true
	cust.IsActive = true
	cust.LoginInfo.LoginAttempts += 1
	cust.UpdatedAt = s.now().UTC()

	return cust, s.repo.Update(cust)
}

func (s *Service) VerifyPassword(identifier, password string) (*Customer, error) {

	if s.IsTestCustomer(identifier, password) {
		return &Customer{ID: "Test-customer-id", Email: identifier, Mobile: identifier}, nil
	}

	cust, err := s.GetCustomerByIdentifier(identifier)
	if err != nil {
		return nil, err
	}

	if cust.Password != password {
		return nil, errors.New("invalid password")
	}

	return cust, nil
}

func (s *Service) SendCustomerOtp(identifier string) error {
	cust, err := s.GetCustomerByIdentifier(identifier)
	if err != nil {
		return err
	}

	if cust.IsVerified {
		return errors.New("customer already verified")
	}

	cust.LoginInfo.OTP = notification.GenerateOTP()
	cust.LoginInfo.OTPExpiry = s.now().UTC().Add(LOGIN_OTP_EXPIRE_TIME)

	cust.UpdatedAt = s.now().UTC()
	return s.repo.Update(cust)
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

func (s *Service) GetCustomerByIdentifier(identifier string) (*Customer, error) {
	cust, err := s.repo.GetByEmail(identifier)
	if err != nil {
		cust, err = s.repo.GetByMobile(identifier)
		if err != nil {
			return nil, errors.New("customer not found")
		}
	}
	return cust, nil
}

func (s *Service) IsTestCustomer(identifier, password string) bool {
	return (identifier == common.TEST_CUSTOMER_EMAIL || identifier == common.TEST_CUSTOMER_MOBILE) && (password == common.TEST_CUSTOMER_PASSWORD || password == common.TEST_CUSTOMER_OTP)
}
