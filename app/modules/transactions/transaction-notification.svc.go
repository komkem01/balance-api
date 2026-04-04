package transactions

import (
	"balance/app/modules/entities/ent"
	"context"
	"fmt"
	"strings"
	"time"
)

func (s *Service) notifyForExpenseTransaction(ctx context.Context, item *ent.TransactionEntity) {
	if item == nil || item.Type != ent.TransactionTypeExpense || item.WalletID == nil {
		return
	}

	wallet, err := s.db.GetWalletByID(ctx, item.WalletID.String())
	if err != nil || wallet.MemberID == nil {
		return
	}

	memberID := wallet.MemberID.String()
	member, err := s.db.GetMemberByID(ctx, memberID)
	if err != nil {
		return
	}

	isTH := strings.EqualFold(strings.TrimSpace(member.PreferredLanguage), "TH")
	todayKey := time.Now().UTC().Format("2006-01-02")

	if member.NotifySecurity && item.Amount >= 10000 {
		title := "Unusual transaction detected"
		description := fmt.Sprintf("High expense detected: %.2f", item.Amount)
		if isTH {
			title = "ตรวจพบรายการผิดปกติ"
			description = fmt.Sprintf("พบรายจ่ายมูลค่าสูง: %.2f", item.Amount)
		}

		dedupe := fmt.Sprintf("security:%s", item.ID.String())
		_, _ = s.db.CreateMemberNotification(
			ctx,
			memberID,
			ent.MemberNotificationTypeSecurity,
			ent.MemberNotificationLevelCritical,
			title,
			description,
			&dedupe,
		)
	}

	if !member.NotifyBudget || item.CategoryID == nil {
		return
	}

	categoryID := item.CategoryID.String()
	budgets, err := s.db.ListBudgets(ctx, &memberID, &categoryID, nil)
	if err != nil {
		return
	}

	categoryName := "Category"
	if category, err := s.db.GetCategoryByID(ctx, categoryID); err == nil {
		categoryName = category.Name
	}

	for _, budget := range budgets {
		if budget == nil || budget.Amount <= 0 {
			continue
		}

		usedPercent := int((budget.SpentAmount / budget.Amount) * 100)
		if usedPercent < 80 {
			continue
		}

		level := ent.MemberNotificationLevelWarning
		title := "Budget usage is high"
		if usedPercent >= 100 {
			level = ent.MemberNotificationLevelCritical
			title = "Budget exceeded"
		}
		if isTH {
			if usedPercent >= 100 {
				title = "งบประมาณเกินกำหนด"
			} else {
				title = "การใช้งบประมาณสูง"
			}
		}

		description := fmt.Sprintf("%s: %d%% used", categoryName, usedPercent)
		if isTH {
			description = fmt.Sprintf("%s: ใช้ไปแล้ว %d%%", categoryName, usedPercent)
		}

		dedupe := fmt.Sprintf("budget:%s:%s:%s", budget.ID.String(), level, todayKey)
		_, _ = s.db.CreateMemberNotification(
			ctx,
			memberID,
			ent.MemberNotificationTypeBudget,
			level,
			title,
			description,
			&dedupe,
		)
	}
}
