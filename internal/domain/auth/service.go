package auth

import (
	"errors"
	"strings"
	"time"

	"ppharma/backend/internal/domain/common"
	"ppharma/backend/internal/domain/customer"
	"ppharma/backend/internal/domain/user"
	"ppharma/backend/support-pkg/auth/jwt"

	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrCustomerNotFound   = errors.New("customer not found")
	ErrInternalError      = errors.New("internal error")
)

type AuthService struct {
	customerRepo    customer.Repository
	customerService *customer.Service
	userRepo        user.Repository
	userService     *user.Service
	jwtProvider     *jwt.Provider
}

func NewAuthService(customerRepo customer.Repository, customerService *customer.Service, userRepo user.Repository, userService *user.Service, jwtProvider *jwt.Provider) *AuthService {
	return &AuthService{
		customerRepo: customerRepo,
		userRepo:     userRepo,
		jwtProvider:  jwtProvider,
	}
}

func (s *AuthService) CustomerLogin(identifier, password, otp string) (string, error) {
	if s.customerService.IsTestCustomer(identifier, password) {
		var email, mobile string
		if strings.Contains(identifier, "@") {
			email = identifier
		} else {
			mobile = identifier
		}
		return s.createCustomerToken(&customer.Customer{ID: common.TEST_CUSTOMER_ID, Email: email, Mobile: mobile})
	}

	if s.customerRepo == nil {
		return "", ErrInvalidCredentials
	}

	var cust *customer.Customer
	var err error

	if otp != "" {
		cust, err = s.customerService.VerifyOTP(identifier, otp)
	} else {
		cust, err = s.customerService.GetCustomerByIdentifier(identifier)
	}

	if err != nil || cust == nil {
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(cust.Password), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.createCustomerToken(cust)
	if err != nil {
		return "", ErrInternalError
	}
	return token, nil
}

func (s *AuthService) SendVerificationOtp(identifier string) error {
	if s.customerService.IsTestCustomer(identifier, "") {
		return nil
	}

	if s.customerRepo == nil {
		return ErrInvalidCredentials
	}

	err := s.customerService.SendCustomerOtp(identifier)
	if err != nil {
		return ErrInternalError
	}

	return nil
}

// func (s *AuthService) VerifyCustomerOtpGenerateToken(identifier, otp string) (string, error) {
// 	if s.customerService.IsTestCustomer(identifier, otp) {
// 		return s.createCustomerToken(&customer.Customer{ID: common.TEST_CUSTOMER_ID, Email: identifier, Mobile: identifier})
// 	}

// 	if s.customerRepo == nil {
// 		return "", ErrInvalidCredentials
// 	}

// 	var custService *customer.Service
// 	var err error

// 	cust, err := custService.VerifyOTP(identifier, otp)
// 	if err != nil {
// 		return "", ErrInvalidCredentials
// 	}

// 	token, err := s.createCustomerToken(cust)
// 	if err != nil {
// 		return "", ErrInternalError
// 	}
// 	return token, nil
// }

func (s *AuthService) ResetCustomerPassword(identifier, newPassword string) error {
	if s.customerService.IsTestCustomer(identifier, newPassword) {
		return nil
	}

	if s.customerRepo == nil {
		return ErrCustomerNotFound
	}

	var cust *customer.Customer
	var err error

	if strings.Contains(identifier, "@") {
		cust, err = s.customerRepo.GetByEmail(identifier)
	} else {
		cust, err = s.customerRepo.GetByMobile(identifier)
	}

	if err != nil || cust == nil {
		return ErrCustomerNotFound
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return ErrInternalError
	}

	cust.Password = string(hashed)
	if err := s.customerRepo.Update(cust); err != nil {
		return ErrInternalError
	}

	return nil
}

func (s *AuthService) UserLogin(identifier, otpCode string) (string, error) {
	if s.userService.IsTestUser(identifier, otpCode) {
		var email, mobile string
		if strings.Contains(identifier, "@") {
			email = identifier
		} else {
			mobile = identifier
		}
		return s.createUserAdminToken(&user.User{ID: "mock-user-id", Role: "admin", Email: email, Mobile: mobile})
	}

	if s.userRepo == nil {
		return "", ErrInvalidCredentials
	}

	u, err := s.VerifyAndEnableUserTOTP(identifier, otpCode)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.createUserAdminToken(u)
	if err != nil {
		return "", ErrInternalError
	}
	return token, nil
}

func (s *AuthService) GenerateUserTOTPConfig(email string) (string, string, error) {
	if s.userRepo == nil {
		return "", "", ErrInternalError
	}
	if email == "" {
		return "", "", errors.New("email is required")
	}

	u, err := s.userRepo.GetByEmail(email)
	if err != nil || u == nil {
		return "", "", errors.New("user not found")
	}

	secret, url, err := s.generateTOTP(email)
	if err != nil {
		return "", "", err
	}

	return secret, url, nil
}

func (s *AuthService) VerifyAndEnableUserTOTP(email, otpCode string) (*user.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	if otpCode == "" {
		return nil, errors.New("otp code is required")
	}

	secret, _, err := s.generateTOTP(email)
	if err != nil {
		return nil, err
	}
	if secret == "" {
		return nil, errors.New("totp not configured")
	}

	valid := totp.Validate(otpCode, secret)
	if !valid {
		return nil, ErrInvalidCredentials
	}

	u, err := s.userRepo.GetByEmail(email)
	if err != nil || u == nil {
		return nil, errors.New("user not found")
	}

	u.TOTPEnabled = true
	u.LastLogin = time.Now()
	if err := s.userRepo.Update(u); err != nil {
		return nil, ErrInternalError
	}
	return u, nil
}

func (s *AuthService) createCustomerToken(cust *customer.Customer) (string, error) {
	return s.jwtProvider.GenerateToken(&common.Principal{ID: cust.ID, Role: "customer", Email: cust.Email, Mobile: cust.Mobile}, 30*24*time.Hour)
}

func (s *AuthService) createUserAdminToken(u *user.User) (string, error) {
	return s.jwtProvider.GenerateToken(&common.Principal{ID: u.ID, Role: string(u.Role), Email: u.Email, Mobile: u.Mobile}, 24*time.Hour)
}

func (s *AuthService) generateTOTP(email string) (string, string, error) {
	if email == "" {
		return "", "", errors.New("email is required")
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      common.PROJECT_NAME,
		AccountName: email,
	})
	if err != nil {
		return "", "", ErrInternalError
	}
	return key.Secret(), key.URL(), nil
}
