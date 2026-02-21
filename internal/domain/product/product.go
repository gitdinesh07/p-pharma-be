package product

import "time"

type Product struct {
	ID             string    `json:"id" bson:"_id"`
	SKU            string    `json:"sku" bson:"sku"`
	Name           string    `json:"name" bson:"name"`
	Description    string    `json:"description,omitempty" bson:"description,omitempty"`
	Category       string    `json:"category,omitempty" bson:"category,omitempty"`
	UnitPrice      int64     `json:"unit_price" bson:"unit_price"`
	InventoryCount int       `json:"inventory_count" bson:"inventory_count"`
	IsActive       bool      `json:"is_active" bson:"is_active"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" bson:"updated_at"`
}

type Repository interface {
	BulkUpsert(products []Product) error
	GetByID(id string) (*Product, error)
	GetBySKU(sku string) (*Product, error)
	UpdateInventory(id string, delta int, reason, orderID, itemID string) error
}
