package order

type DefaultStatusDeriver struct{}

func (d DefaultStatusDeriver) Derive(items []OrderItem) OrderStatus {
	if len(items) == 0 {
		return OrderStatusPending
	}

	counts := map[ItemStatus]int{}
	for _, item := range items {
		counts[item.Status]++
	}

	if counts[ItemStatusCancelled] == len(items) {
		return OrderStatusCancelled
	}
	if counts[ItemStatusDelivered] == len(items) {
		return OrderStatusCompleted
	}
	if counts[ItemStatusReturned] > 0 {
		if counts[ItemStatusReturned] == len(items) {
			return OrderStatusPartiallyReturned
		}
		return OrderStatusMixed
	}
	if counts[ItemStatusDelivered] > 0 {
		return OrderStatusPartialDelivered
	}
	if counts[ItemStatusShipped] > 0 {
		return OrderStatusPartialShipped
	}
	if counts[ItemStatusPending] == len(items) {
		return OrderStatusPending
	}
	return OrderStatusMixed
}
