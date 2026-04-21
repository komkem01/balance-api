package goals

import (
	"context"

	"balance/app/modules/entities/ent"
	"balance/app/utils/base"
)

type ListRequestService struct {
	MemberID *string
	Status   *ent.GoalStatus
	Type     *ent.GoalType
	Page     int
	Size     int
}

func (s *Service) ListGoal(ctx context.Context, req *ListRequestService) ([]*InfoResponseService, *base.ResponsePaginate, error) {
	items, err := s.db.ListGoals(ctx, req.MemberID, req.Status, req.Type)
	if err != nil {
		return nil, nil, err
	}

	res := make([]*InfoResponseService, 0, len(items))
	for _, item := range items {
		goalRes := &InfoResponseService{
			ID:                 item.ID,
			MemberID:           item.MemberID,
			Name:               item.Name,
			Type:               item.Type,
			TargetAmount:       item.TargetAmount,
			StartAmount:        item.StartAmount,
			CurrentAmount:      item.CurrentAmount,
			StartDate:          item.StartDate,
			TargetDate:         item.TargetDate,
			Status:             item.Status,
			AutoTracking:       item.AutoTracking,
			TrackingSourceType: item.TrackingSourceType,
			TrackingSourceID:   item.TrackingSourceID,
			DepositWalletID:    item.DepositWalletID,
			CreatedAt:          item.CreatedAt,
			UpdatedAt:          item.UpdatedAt,
		}

		resolvedCurrent, resolveErr := s.resolveAutoCurrentAmount(ctx, goalRes)
		if resolveErr == nil {
			goalRes.CurrentAmount = resolvedCurrent
		}

		res = append(res, goalRes)
	}

	total := int64(len(items))
	page := int64(req.Page)
	size := int64(req.Size)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	return res, &base.ResponsePaginate{Page: page, Size: size, Total: total}, nil
}
