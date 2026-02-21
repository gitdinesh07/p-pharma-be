package order

import "testing"

func TestDeriveStatus(t *testing.T) {
	d := DefaultStatusDeriver{}
	tests := []struct {
		name string
		in   []OrderItem
		want OrderStatus
	}{
		{name: "all pending", in: []OrderItem{{Status: ItemStatusPending}, {Status: ItemStatusPending}}, want: OrderStatusPending},
		{name: "partial shipped", in: []OrderItem{{Status: ItemStatusShipped}, {Status: ItemStatusPending}}, want: OrderStatusPartialShipped},
		{name: "partial delivered", in: []OrderItem{{Status: ItemStatusDelivered}, {Status: ItemStatusShipped}}, want: OrderStatusPartialDelivered},
		{name: "completed", in: []OrderItem{{Status: ItemStatusDelivered}, {Status: ItemStatusDelivered}}, want: OrderStatusCompleted},
		{name: "all cancelled", in: []OrderItem{{Status: ItemStatusCancelled}, {Status: ItemStatusCancelled}}, want: OrderStatusCancelled},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := d.Derive(tc.in)
			if got != tc.want {
				t.Fatalf("got %s want %s", got, tc.want)
			}
		})
	}
}
