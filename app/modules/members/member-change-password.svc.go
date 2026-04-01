package members

import (
	"context"
	"strings"

	"balance/app/utils/hashing"

	"github.com/google/uuid"
)

type ChangeMePasswordRequestService struct {
	MemberID        string `json:"member_id"`
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type ChangeMePasswordResponseService struct {
	MemberID string `json:"member_id"`
}

func (s *Service) ChangeMePassword(ctx context.Context, req *ChangeMePasswordRequestService) (*ChangeMePasswordResponseService, error) {
	memberID := strings.TrimSpace(req.MemberID)
	currentPassword := strings.TrimSpace(req.CurrentPassword)
	newPassword := strings.TrimSpace(req.NewPassword)

	if memberID == "" {
		return nil, ErrMemberUnauthorized
	}
	if _, err := uuid.Parse(memberID); err != nil {
		return nil, ErrMemberUnauthorized
	}
	if currentPassword == "" || newPassword == "" {
		return nil, ErrMemberPasswordRequired
	}

	accounts, err := s.db.ListMemberAccounts(ctx)
	if err != nil {
		return nil, err
	}

	var accountID string
	var accountHash string
	for _, item := range accounts {
		if item.MemberID == nil {
			continue
		}
		if item.MemberID.String() != memberID {
			continue
		}
		accountID = item.ID.String()
		accountHash = item.Password
		break
	}

	if accountID == "" || accountHash == "" {
		return nil, ErrMemberAccountNotFound
	}

	if !hashing.CheckPasswordHash([]byte(accountHash), []byte(currentPassword)) {
		return nil, ErrMemberInvalidCredentials
	}

	hashedPassword, err := hashing.HashPassword(newPassword)
	if err != nil {
		return nil, err
	}

	password := string(hashedPassword)
	if _, err := s.db.UpdateMemberAccount(ctx, accountID, nil, nil, &password); err != nil {
		return nil, err
	}

	return &ChangeMePasswordResponseService{MemberID: memberID}, nil
}
