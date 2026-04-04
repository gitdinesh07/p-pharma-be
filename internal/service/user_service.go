package service

import (
	"errors"
	"time"

	"ppharma/backend/internal/domain/user"
	"ppharma/backend/support-pkg/notification"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo        user.Repository
	emailSender notification.EmailSender
	smsSender   notification.SMSSender
	now         func() time.Time
}

func NewUserService(repo user.Repository, emailSender notification.EmailSender, smsSender notification.SMSSender) *UserService {
	return &UserService{
		repo:        repo,
		emailSender: emailSender,
		smsSender:   smsSender,
		now:         time.Now,
	}
}

func (s *UserService) checkDuplicate(u *user.User) error {
	if u.Email != "" {
		if ex, _ := s.repo.GetByEmail(u.Email); ex != nil && ex.ID != u.ID {
			return errors.New("user already exists with this email")
		}
	}
	if u.Mobile != "" {
		if ex, _ := s.repo.GetByMobile(u.Mobile); ex != nil && ex.ID != u.ID {
			return errors.New("user already exists with this mobile")
		}
	}
	return nil
}

func (s *UserService) CreateUser(u *user.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if u.ID == "" {
		u.ID = uuid.New().String()
	}

	if err := s.checkDuplicate(u); err != nil {
		return err
	}

	if u.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashed)
	}

	u.IsVerified = false

	// if u.Email != "" && s.emailSender != nil {
	// 	go s.emailSender.SendEmail(context.Background(), []string{u.Email}, "Your OTP Code", "Your code is: "+u.OTP, false)
	// }
	// if u.Mobile != "" && s.smsSender != nil {
	// 	go s.smsSender.SendSMS(context.Background(), u.Mobile, "Your code is: "+u.OTP)
	// }

	u.CreatedAt = s.now().UTC()
	u.UpdatedAt = s.now().UTC()

	return s.repo.Create(u)
}

func (s *UserService) VerifyOTP(identifier string, otp string) error {
	u, err := s.repo.GetByEmail(identifier)
	if err != nil {
		u, err = s.repo.GetByMobile(identifier)
		if err != nil {
			return errors.New("user not found")
		}
	}

	if u.IsVerified {
		return errors.New("user already verified")
	}

	u.IsVerified = true
	u.UpdatedAt = s.now().UTC()

	return s.repo.Update(u)
}

func (s *UserService) UpdateUser(u *user.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := s.checkDuplicate(u); err != nil {
		return err
	}

	u.UpdatedAt = s.now().UTC()
	return s.repo.Update(u)
}

func (s *UserService) GetUser(id string) (*user.User, error) {
	return s.repo.GetByID(id)
}
