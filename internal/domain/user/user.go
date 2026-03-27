package user

import "time"

type Role string

const RoleAdmin Role = "admin"

type User struct {
	ID        string    `json:"id" bson:"_id"`
	Name      string    `json:"name" bson:"name"`
	Email     string    `json:"email,omitempty" bson:"email,omitempty"`
	Mobile    string    `json:"mobile,omitempty" bson:"mobile,omitempty"`
	Password  string    `json:"-" bson:"password,omitempty"`
	Role      Role      `json:"role" bson:"role"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type Repository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByMobile(mobile string) (*User, error)
}
