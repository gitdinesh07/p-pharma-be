package order_test

import (
	"errors"
	"testing"
	"time"

	"ppharma/backend/internal/domain/order"
)

type repoStub struct {
	orders map[string]*order.Order
}

func (r *repoStub) GetByID(orderID string) (*order.Order, error) {
	o, ok := r.orders[orderID]
	if !ok {
		return nil, order.ErrOrderNotFound
	}
	copy := *o
	copy.Items = append([]order.OrderItem(nil), o.Items...)
	return &copy, nil
}

func (r *repoStub) GetByIDForCustomer(orderID, customerID string) (*order.Order, error) {
	o, err := r.GetByID(orderID)
	if err != nil {
		return nil, err
	}
	if o.CustomerID != customerID {
		return nil, order.ErrOrderNotFound
	}
	return o, nil
}

func (r *repoStub) Save(o *order.Order) error {
	copy := *o
	copy.Items = append([]order.OrderItem(nil), o.Items...)
	r.orders[o.OrderID] = &copy
	return nil
}

func TestTransitionValidation(t *testing.T) {
	if err := order.ValidateTransition(order.ItemStatusPending, order.ItemStatusDelivered); err == nil {
		t.Fatal("expected invalid transition error")
	}
	if err := order.ValidateTransition(order.ItemStatusPending, order.ItemStatusConfirmed); err != nil {
		t.Fatalf("expected valid transition, got %v", err)
	}
}

func TestUpdateItemStatusDerived(t *testing.T) {
	repo := &repoStub{orders: map[string]*order.Order{
		"o1": {
			OrderID:    "o1",
			CustomerID: "c1",
			Items: []order.OrderItem{
				{ItemID: "i1", Qty: 1, Status: order.ItemStatusShipped},
				{ItemID: "i2", Qty: 1, Status: order.ItemStatusPending},
			},
			DerivedStatus: order.OrderStatusPartialShipped,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}}
	svc := order.NewService(repo, order.DefaultStatusDeriver{})
	ord, err := svc.UpdateItemStatus("o1", "i1", order.ItemStatusDelivered, "done", "admin1")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if ord.DerivedStatus != order.OrderStatusPartialDelivered {
		t.Fatalf("got %s", ord.DerivedStatus)
	}
}

func TestReturnQtyRule(t *testing.T) {
	repo := &repoStub{orders: map[string]*order.Order{
		"o1": {
			OrderID:    "o1",
			CustomerID: "c1",
			Items:      []order.OrderItem{{ItemID: "i1", Qty: 1, ReturnedQty: 1, Status: order.ItemStatusReturned}},
		},
	}}
	svc := order.NewService(repo, order.DefaultStatusDeriver{})
	_, err := svc.UpdateItemStatus("o1", "i1", order.ItemStatusReturned, "again", "admin1")
	if err == nil {
		t.Fatal("expected err")
	}
	if !errors.Is(err, order.ErrInvalidReturnedQty) && !errors.Is(err, order.ErrInvalidTransition) {
		t.Fatalf("unexpected err: %v", err)
	}
}
