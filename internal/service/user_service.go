package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"ppharma/backend/internal/domain/user"
)

type UserService struct {
	repo user.Repository
	now  func() time.Time
}

func NewUserService(repo user.Repository) *UserService {
	return &UserService{
		repo: repo,
		now:  time.Now,
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

	u.CreatedAt = s.now().UTC()
	u.UpdatedAt = s.now().UTC()

	return s.repo.Create(u)
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
