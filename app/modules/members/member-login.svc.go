package members

import (
	"context"
	"os"
	"strings"
	"time"

	"balance/app/utils/authx"
	"balance/app/utils/hashing"
)

const memberTokenTTL = 10 * time.Minute
const memberRefreshTokenTTL = 7 * 24 * time.Hour

type LoginRequestService struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponseService struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int64  `json:"expires_in"`
	RefreshExpiresIn int64  `json:"refresh_expires_in"`
	MemberID         string `json:"member_id"`
	Username         string `json:"username"`
}

type RefreshTokenRequestService struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponseService struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int64  `json:"expires_in"`
	RefreshExpiresIn int64  `json:"refresh_expires_in"`
}

func appSecret() string {
	secret := strings.TrimSpace(os.Getenv("APP_KEY"))
	if secret == "" {
		return "secret"
	}
	return secret
}

func (s *Service) LoginMember(ctx context.Context, req *LoginRequestService) (*LoginResponseService, error) {
	username := strings.TrimSpace(req.Username)
	password := strings.TrimSpace(req.Password)
	if username == "" || password == "" {
		return nil, ErrMemberInvalidCredentials
	}

	accounts, err := s.db.ListMemberAccounts(ctx)
	if err != nil {
		return nil, err
	}

	var foundMemberID string
	var foundUsername string
	var foundHash string
	for _, item := range accounts {
		if strings.EqualFold(strings.TrimSpace(item.Username), username) {
			if item.MemberID == nil {
				return nil, ErrMemberInvalidCredentials
			}
			foundMemberID = item.MemberID.String()
			foundUsername = item.Username
			foundHash = item.Password
			break
		}
	}

	if foundMemberID == "" || foundHash == "" {
		return nil, ErrMemberInvalidCredentials
	}

	if !hashing.CheckPasswordHash([]byte(foundHash), []byte(req.Password)) {
		return nil, ErrMemberInvalidCredentials
	}

	now := time.Now()
	if _, err := s.db.UpdateMember(ctx, foundMemberID, nil, nil, nil, nil, nil, nil, &now); err != nil {
		return nil, err
	}

	accessToken, exp, err := authx.GenerateMemberToken(appSecret(), foundMemberID, foundUsername, memberTokenTTL)
	if err != nil {
		return nil, err
	}
	refreshToken, refreshExp, err := authx.GenerateMemberRefreshToken(appSecret(), foundMemberID, foundUsername, memberRefreshTokenTTL)
	if err != nil {
		return nil, err
	}

	return &LoginResponseService{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		TokenType:        "Bearer",
		ExpiresIn:        exp - time.Now().Unix(),
		RefreshExpiresIn: refreshExp - time.Now().Unix(),
		MemberID:         foundMemberID,
		Username:         foundUsername,
	}, nil
}

func (s *Service) RefreshMemberToken(ctx context.Context, req *RefreshTokenRequestService) (*RefreshTokenResponseService, error) {
	token := strings.TrimSpace(req.RefreshToken)
	if token == "" {
		return nil, ErrMemberUnauthorized
	}

	claims, err := authx.ParseMemberRefreshToken(appSecret(), token)
	if err != nil {
		return nil, ErrMemberUnauthorized
	}

	member, err := s.db.GetMemberByID(ctx, claims.MemberID)
	if err != nil {
		return nil, ErrMemberUnauthorized
	}

	accessToken, exp, err := authx.GenerateMemberToken(appSecret(), claims.MemberID, claims.Username, memberTokenTTL)
	if err != nil {
		return nil, err
	}
	refreshToken, refreshExp, err := authx.GenerateMemberRefreshToken(appSecret(), claims.MemberID, claims.Username, memberRefreshTokenTTL)
	if err != nil {
		return nil, err
	}

	_ = member
	return &RefreshTokenResponseService{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		TokenType:        "Bearer",
		ExpiresIn:        exp - time.Now().Unix(),
		RefreshExpiresIn: refreshExp - time.Now().Unix(),
	}, nil
}
