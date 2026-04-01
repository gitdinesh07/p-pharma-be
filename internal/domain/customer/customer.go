package customer

import (
	"errors"
	"regexp"
	"time"
)

var (
	emailRegex  = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	mobileRegex = regexp.MustCompile(`^91\d{10}$`)
)

type GenderEnum int

const (
	GenderMale   GenderEnum = 1
	GenderFemale GenderEnum = 2
	GenderOthers GenderEnum = 3
)

type Customer struct {
	ID        string     `json:"id" bson:"_id"`
	Name      string     `json:"name" bson:"name"`
	Email     string     `json:"email,omitempty" bson:"email,omitempty"`
	Mobile    string     `json:"mobile,omitempty" bson:"mobile,omitempty"`
	Password  string     `json:"-" bson:"password,omitempty"`
	Gender    GenderEnum `json:"gender,omitempty" bson:"gender,omitempty"`
	Age       int        `json:"age,omitempty" bson:"age,omitempty"`
	Address   []Address  `json:"address,omitempty" bson:"address,omitempty"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" bson:"updated_at"`
}

type GeoLocation struct {
	Type        string    `json:"type" bson:"type"`               // e.g., "Point"
	Coordinates []float64 `json:"coordinates" bson:"coordinates"` // User's [longitude, latitude]
}

type Address struct {
	AddressLine1 string      `json:"address_line1" bson:"address_line1"`
	AddressLine2 string      `json:"address_line2,omitempty" bson:"address_line2,omitempty"`
	City         string      `json:"city" bson:"city"`
	State        string      `json:"state" bson:"state"`
	Landmark     string      `json:"landmark,omitempty" bson:"landmark,omitempty"`
	Pincode      string      `json:"pincode" bson:"pincode"`
	Country      string      `json:"country" bson:"country"`
	Location     GeoLocation `json:"location,omitempty" bson:"location,omitempty"`
	IsDefault    bool        `json:"is_default" bson:"is_default"`
	AddressType  string      `json:"address_type" bson:"address_type"`
	CreatedAt    time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at" bson:"updated_at"`
}

func (c *Customer) Validate() error {
	if c.Name == "" {
		return errors.New("name is required")
	}
	if c.Email != "" {
		err := c.ValidateEmail()
		if err != nil {
			return err
		}
	}
	if c.Mobile != "" {
		err := c.ValidateMobile()
		if err != nil {
			return err
		}
	}
	if c.Gender != GenderMale && c.Gender != GenderFemale && c.Gender != GenderOthers {
		return errors.New("gender is required and must be 1 (Male), 2 (Female), or 3 (Others)")
	}
	if c.Age == 0 {
		return errors.New("age is required")
	}
	return nil
}

func (c *Customer) ValidateEmail() error {
	if !emailRegex.MatchString(c.Email) {
		return errors.New("invalid email format")
	}
	return nil
}

func (c *Customer) ValidateMobile() error {
	if !mobileRegex.MatchString(c.Mobile) {
		return errors.New("invalid mobile number format")
	}
	return nil
}

func (a *Address) ValidateAddress() error {
	if a.AddressLine1 == "" {
		return errors.New("address line 1 is required")
	}
	if a.City == "" {
		return errors.New("city is required")
	}
	if a.State == "" {
		return errors.New("state is required")
	}
	if a.Pincode == "" {
		return errors.New("pincode is required")
	}
	if a.Country == "" {
		a.Country = "India"
	}
	if a.AddressType == "" {
		a.AddressType = "home"
	}
	if a.IsDefault == false {
		a.IsDefault = true
	}
	return nil
}

type Repository interface {
	Create(customer *Customer) error
	GetByID(id string) (*Customer, error)
	GetByEmail(email string) (*Customer, error)
	GetByMobile(mobile string) (*Customer, error)
	Update(customer *Customer) error
}
