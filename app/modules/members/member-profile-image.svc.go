package members

import (
	"context"
	"database/sql"
	"errors"
	"mime/multipart"
	"strings"

	"balance/app/modules/storage"

	"github.com/google/uuid"
)

type UploadProfileImageRequestService struct {
	MemberID string
	Image    *multipart.FileHeader
}

func (s *Service) UploadMeProfileImage(ctx context.Context, req *UploadProfileImageRequestService) (*MeResponseService, error) {
	if _, err := uuid.Parse(strings.TrimSpace(req.MemberID)); err != nil {
		return nil, ErrMemberUnauthorized
	}
	if req.Image == nil {
		return nil, storage.ErrImageRequired
	}
	if s.sto == nil || !s.sto.Enabled() {
		return nil, ErrMemberStorageNotConfigured
	}

	imageURL, err := s.sto.UploadProfileImage(ctx, req.MemberID, req.Image)
	if err != nil {
		return nil, err
	}

	if _, err := s.db.UpdateMember(ctx, req.MemberID, nil, nil, nil, nil, nil, nil, nil, &imageURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMemberNotFound
		}
		return nil, err
	}

	return s.InfoMeMember(ctx, &MeRequestService{MemberID: req.MemberID})
}
