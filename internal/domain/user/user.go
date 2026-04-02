package user

import (
	"errors"
	"regexp"
	"time"
)

var (
	emailRegex  = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	mobileRegex = regexp.MustCompile(`^91\d{10}$`)
)

type Role string

const RoleAdmin Role = "admin"

type User struct {
	ID        string    `json:"id" bson:"_id"`
	Name      string    `json:"name" bson:"name"`
	Email     string    `json:"email,omitempty" bson:"email,omitempty"`
	Mobile    string    `json:"mobile,omitempty" bson:"mobile,omitempty"`
	Password  string    `json:"-" bson:"password,omitempty"`
	PhotoURL  string    `json:"photo_url,omitempty" bson:"photo_url,omitempty"`
	Role      Role      `json:"role" bson:"role"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type Repository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByMobile(mobile string) (*User, error)
	Update(user *User) error
}

func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("name is required")
	}
	if u.Email != "" {
		if !emailRegex.MatchString(u.Email) {
			return errors.New("invalid email format")
		}
	}
	if u.Mobile != "" {
		if !mobileRegex.MatchString(u.Mobile) {
			return errors.New("invalid mobile format")
		}
	}
	if u.Role == "" {
		return errors.New("role is required")
	}
	return nil
}
