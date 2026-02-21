package order

import "time"

type Service struct {
	repo    Repository
	deriver StatusDeriver
	now     func() time.Time
}

func NewService(repo Repository, deriver StatusDeriver) *Service {
	return &Service{repo: repo, deriver: deriver, now: time.Now}
}

func (s *Service) GetOrderForCustomer(orderID, customerID string) (*Order, error) {
	return s.repo.GetByIDForCustomer(orderID, customerID)
}

func (s *Service) UpdateItemStatus(orderID, itemID string, to ItemStatus, reason, changedBy string) (*Order, error) {
	ord, err := s.repo.GetByID(orderID)
	if err != nil {
		return nil, err
	}

	itemIdx := -1
	for i := range ord.Items {
		if ord.Items[i].ItemID == itemID {
			itemIdx = i
			break
		}
	}
	if itemIdx == -1 {
		return nil, ErrItemNotFound
	}

	item := &ord.Items[itemIdx]
	if err := ValidateTransition(item.Status, to); err != nil {
		return nil, err
	}
	if to == ItemStatusReturned {
		deliveredOrReturned := item.Status == ItemStatusDelivered || item.Status == ItemStatusReturned
		if !deliveredOrReturned {
			return nil, ErrInvalidTransition
		}
		if item.ReturnedQty >= item.Qty {
			return nil, ErrInvalidReturnedQty
		}
		item.ReturnedQty = item.Qty
	}

	item.Status = to
	item.StatusHistory = append(item.StatusHistory, ItemStatusHistory{
		State:     to,
		Reason:    reason,
		ChangedBy: changedBy,
		ChangedAt: s.now().UTC(),
	})

	ord.DerivedStatus = s.deriver.Derive(ord.Items)
	ord.UpdatedAt = s.now().UTC()

	if err := s.repo.Save(ord); err != nil {
		return nil, err
	}
	return ord, nil
}
