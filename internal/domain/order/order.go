package order

import (
	"errors"
	"fmt"
	"time"
)

type ItemStatus string

type OrderStatus string

const (
	ItemStatusPending   ItemStatus = "pending"
	ItemStatusConfirmed ItemStatus = "confirmed"
	ItemStatusPacked    ItemStatus = "packed"
	ItemStatusShipped   ItemStatus = "shipped"
	ItemStatusDelivered ItemStatus = "delivered"
	ItemStatusReturned  ItemStatus = "returned"
	ItemStatusCancelled ItemStatus = "cancelled"
)

const (
	OrderStatusPending           OrderStatus = "pending"
	OrderStatusPartialShipped    OrderStatus = "partial_shipped"
	OrderStatusPartialDelivered  OrderStatus = "partial_delivered"
	OrderStatusCompleted         OrderStatus = "completed"
	OrderStatusPartiallyReturned OrderStatus = "partially_returned"
	OrderStatusCancelled         OrderStatus = "cancelled"
	OrderStatusMixed             OrderStatus = "mixed"
)

type ProductSnapshot struct {
	Name      string `json:"name" bson:"name"`
	SKU       string `json:"sku" bson:"sku"`
	UnitPrice int64  `json:"unit_price" bson:"unit_price"`
}

type ItemStatusHistory struct {
	State     ItemStatus `json:"state" bson:"state"`
	Reason    string     `json:"reason,omitempty" bson:"reason,omitempty"`
	ChangedBy string     `json:"changed_by" bson:"changed_by"`
	ChangedAt time.Time  `json:"changed_at" bson:"changed_at"`
}

type OrderItem struct {
	ItemID          string              `json:"item_id" bson:"item_id"`
	ProductID       string              `json:"product_id" bson:"product_id"`
	ProductSnapshot ProductSnapshot     `json:"product_snapshot" bson:"product_snapshot"`
	Qty             int                 `json:"qty" bson:"qty"`
	LineTotal       int64               `json:"line_total" bson:"line_total"`
	Status          ItemStatus          `json:"status" bson:"status"`
	ReturnedQty     int                 `json:"returned_qty" bson:"returned_qty"`
	StatusHistory   []ItemStatusHistory `json:"status_history" bson:"status_history"`
}

type Order struct {
	OrderID       string      `json:"order_id" bson:"order_id"`
	CustomerID    string      `json:"customer_id" bson:"customer_id"`
	Currency      string      `json:"currency" bson:"currency"`
	Subtotal      int64       `json:"subtotal" bson:"subtotal"`
	Discount      int64       `json:"discount" bson:"discount"`
	ShippingFee   int64       `json:"shipping_fee" bson:"shipping_fee"`
	GrandTotal    int64       `json:"grand_total" bson:"grand_total"`
	DerivedStatus OrderStatus `json:"derived_status" bson:"derived_status"`
	PaymentStatus string      `json:"payment_status" bson:"payment_status"`
	Items         []OrderItem `json:"items" bson:"items"`
	CreatedAt     time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" bson:"updated_at"`
}

type StatusDeriver interface {
	Derive(items []OrderItem) OrderStatus
}

type Repository interface {
	GetByID(orderID string) (*Order, error)
	GetByIDForCustomer(orderID, customerID string) (*Order, error)
	Save(order *Order) error
}

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrItemNotFound       = errors.New("order item not found")
	ErrInvalidTransition  = errors.New("invalid item status transition")
	ErrInvalidReturnedQty = errors.New("returned quantity exceeds delivered quantity")
)

func ValidateTransition(from, to ItemStatus) error {
	allowed := map[ItemStatus]map[ItemStatus]bool{
		ItemStatusPending:   {ItemStatusConfirmed: true, ItemStatusCancelled: true},
		ItemStatusConfirmed: {ItemStatusPacked: true, ItemStatusCancelled: true},
		ItemStatusPacked:    {ItemStatusShipped: true, ItemStatusCancelled: true},
		ItemStatusShipped:   {ItemStatusDelivered: true, ItemStatusReturned: true},
		ItemStatusDelivered: {ItemStatusReturned: true},
		ItemStatusReturned:  {},
		ItemStatusCancelled: {},
	}
	if from == to {
		return nil
	}
	if !allowed[from][to] {
		return fmt.Errorf("%w: %s -> %s", ErrInvalidTransition, from, to)
	}
	return nil
}

func AllowedStatuses() []ItemStatus {
	return []ItemStatus{
		ItemStatusPending,
		ItemStatusConfirmed,
		ItemStatusPacked,
		ItemStatusShipped,
		ItemStatusDelivered,
		ItemStatusReturned,
		ItemStatusCancelled,
	}
}
