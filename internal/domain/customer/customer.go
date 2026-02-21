package customer

import "time"

type Customer struct {
	ID        string    `json:"id" bson:"_id"`
	Name      string    `json:"name" bson:"name"`
	Email     string    `json:"email,omitempty" bson:"email,omitempty"`
	Mobile    string    `json:"mobile,omitempty" bson:"mobile,omitempty"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type Repository interface {
	Create(customer *Customer) error
	GetByID(id string) (*Customer, error)
	GetByEmail(email string) (*Customer, error)
	GetByMobile(mobile string) (*Customer, error)
}
