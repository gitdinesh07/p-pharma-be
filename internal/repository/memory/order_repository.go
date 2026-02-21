package memory

import (
	"sync"

	"ppharma/backend/internal/domain/order"
)

type OrderRepository struct {
	mu     sync.RWMutex
	orders map[string]*order.Order
}

func NewOrderRepository(seed []*order.Order) *OrderRepository {
	orders := make(map[string]*order.Order, len(seed))
	for _, o := range seed {
		copy := *o
		orders[o.OrderID] = &copy
	}
	return &OrderRepository{orders: orders}
}

func (r *OrderRepository) GetByID(orderID string) (*order.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ord, ok := r.orders[orderID]
	if !ok {
		return nil, order.ErrOrderNotFound
	}
	copy := *ord
	copy.Items = append([]order.OrderItem(nil), ord.Items...)
	return &copy, nil
}

func (r *OrderRepository) GetByIDForCustomer(orderID, customerID string) (*order.Order, error) {
	ord, err := r.GetByID(orderID)
	if err != nil {
		return nil, err
	}
	if ord.CustomerID != customerID {
		return nil, order.ErrOrderNotFound
	}
	return ord, nil
}

func (r *OrderRepository) Save(o *order.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	copy := *o
	copy.Items = append([]order.OrderItem(nil), o.Items...)
	r.orders[o.OrderID] = &copy
	return nil
}
