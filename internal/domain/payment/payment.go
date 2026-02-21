package payment

import "time"

type Status string

const (
	StatusPending  Status = "pending"
	StatusPaid     Status = "paid"
	StatusFailed   Status = "failed"
	StatusRefunded Status = "refunded"
)

type ItemAllocation struct {
	OrderItemID string `json:"order_item_id" bson:"order_item_id"`
	Amount      int64  `json:"amount" bson:"amount"`
}

type Payment struct {
	ID              string           `json:"id" bson:"_id"`
	OrderID         string           `json:"order_id" bson:"order_id"`
	CustomerID      string           `json:"customer_id" bson:"customer_id"`
	Method          string           `json:"method" bson:"method"`
	Status          Status           `json:"status" bson:"status"`
	Amount          int64            `json:"amount" bson:"amount"`
	GatewayRef      string           `json:"gateway_ref,omitempty" bson:"gateway_ref,omitempty"`
	ItemAllocations []ItemAllocation `json:"item_allocations,omitempty" bson:"item_allocations,omitempty"`
	CreatedAt       time.Time        `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at" bson:"updated_at"`
}

type Repository interface {
	Create(entry *Payment) error
	ListByOrder(orderID string) ([]Payment, error)
	UpdateStatus(id string, status Status, reason string) error
}
