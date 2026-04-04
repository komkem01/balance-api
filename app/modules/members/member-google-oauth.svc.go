package members

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"balance/app/utils/authx"
	"balance/app/utils/hashing"
)

const googleAuthBaseURL = "https://accounts.google.com/o/oauth2/v2/auth"
const googleTokenURL = "https://oauth2.googleapis.com/token"
const googleUserInfoURL = "https://openidconnect.googleapis.com/v1/userinfo"

const googleOAuthStateTTL = 10 * time.Minute

type googleOAuthState struct {
	Nonce string `json:"nonce"`
	Exp   int64  `json:"exp"`
	Iat   int64  `json:"iat"`
}

type googleTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

type googleUserInfoResponse struct {
	Sub         string `json:"sub"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	GivenName   string `json:"given_name"`
	FamilyName  string `json:"family_name"`
	Picture     string `json:"picture"`
	VerifiedRaw any    `json:"email_verified"`
}

func (s *Service) googleConfig() Config {
	if s == nil || s.conf == nil || s.conf.Val == nil {
		return Config{}
	}

	return *s.conf.Val
}

func (s *Service) isGoogleOAuthEnabled() bool {
	conf := s.googleConfig()
	return strings.TrimSpace(conf.GoogleClientId) != "" &&
		strings.TrimSpace(conf.GoogleClientSecret) != "" &&
		strings.TrimSpace(conf.GoogleRedirectUrl) != ""
}

func (s *Service) GoogleOAuthStartURL(ctx context.Context) (string, error) {
	_ = ctx
	if !s.isGoogleOAuthEnabled() {
		return "", ErrMemberGoogleOAuthDisabled
	}

	conf := s.googleConfig()
	state, err := s.newGoogleOAuthStateToken()
	if err != nil {
		return "", err
	}

	scopes := strings.TrimSpace(conf.GoogleScopes)
	if scopes == "" {
		scopes = "openid email profile"
	}

	q := url.Values{}
	q.Set("client_id", strings.TrimSpace(conf.GoogleClientId))
	q.Set("redirect_uri", strings.TrimSpace(conf.GoogleRedirectUrl))
	q.Set("response_type", "code")
	q.Set("scope", scopes)
	q.Set("state", state)
	q.Set("prompt", "select_account")

	return fmt.Sprintf("%s?%s", googleAuthBaseURL, q.Encode()), nil
}

func (s *Service) GoogleOAuthFailureRedirectURL(reason string) string {
	conf := s.googleConfig()
	base := strings.TrimSpace(conf.GoogleFrontendFailureUrl)
	if base == "" {
		base = "http://localhost:3000/?oauth=failed"
	}

	u, err := url.Parse(base)
	if err != nil {
		return base
	}

	q := u.Query()
	if strings.TrimSpace(reason) != "" {
		q.Set("reason", strings.TrimSpace(reason))
	}
	u.RawQuery = q.Encode()

	return u.String()
}

func (s *Service) GoogleOAuthCallback(ctx context.Context, code string, state string) (*LoginResponseService, error) {
	if !s.isGoogleOAuthEnabled() {
		return nil, ErrMemberGoogleOAuthDisabled
	}

	if strings.TrimSpace(code) == "" {
		return nil, ErrMemberGoogleAuthFailed
	}

	if err := s.verifyGoogleOAuthStateToken(strings.TrimSpace(state)); err != nil {
		return nil, err
	}

	token, err := s.exchangeGoogleCode(ctx, strings.TrimSpace(code))
	if err != nil {
		return nil, err
	}

	profile, err := s.fetchGoogleUserInfo(ctx, token.AccessToken)
	if err != nil {
		return nil, err
	}

	memberID, username, err := s.ensureGoogleMember(ctx, profile)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	if _, err := s.db.UpdateMember(ctx, memberID, nil, nil, nil, nil, nil, nil, &now, nil); err != nil {
		return nil, err
	}

	accessToken, exp, err := authx.GenerateMemberToken(appSecret(), memberID, username, memberTokenTTL)
	if err != nil {
		return nil, err
	}
	refreshToken, refreshExp, err := authx.GenerateMemberRefreshToken(appSecret(), memberID, username, memberRefreshTokenTTL)
	if err != nil {
		return nil, err
	}

	return &LoginResponseService{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		TokenType:        "Bearer",
		ExpiresIn:        exp - time.Now().Unix(),
		RefreshExpiresIn: refreshExp - time.Now().Unix(),
		MemberID:         memberID,
		Username:         username,
	}, nil
}

func (s *Service) GoogleOAuthSuccessRedirectURL(login *LoginResponseService) string {
	conf := s.googleConfig()
	base := strings.TrimSpace(conf.GoogleFrontendSuccessUrl)
	if base == "" {
		base = "http://localhost:3000/auth/callback/google"
	}

	u, err := url.Parse(base)
	if err != nil {
		return base
	}

	f := url.Values{}
	f.Set("access_token", login.AccessToken)
	f.Set("refresh_token", login.RefreshToken)
	f.Set("token_type", login.TokenType)
	f.Set("expires_in", fmt.Sprintf("%d", login.ExpiresIn))
	f.Set("refresh_expires_in", fmt.Sprintf("%d", login.RefreshExpiresIn))
	f.Set("member_id", login.MemberID)
	f.Set("username", login.Username)
	u.Fragment = f.Encode()

	return u.String()
}

func (s *Service) newGoogleOAuthStateToken() (string, error) {
	rawNonce := make([]byte, 16)
	if _, err := rand.Read(rawNonce); err != nil {
		return "", err
	}

	now := time.Now().Unix()
	payload, err := json.Marshal(googleOAuthState{
		Nonce: hex.EncodeToString(rawNonce),
		Exp:   now + int64(googleOAuthStateTTL.Seconds()),
		Iat:   now,
	})
	if err != nil {
		return "", err
	}

	payloadPart := base64.RawURLEncoding.EncodeToString(payload)
	h := hmac.New(sha256.New, []byte(appSecret()))
	h.Write([]byte(payloadPart))
	signaturePart := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	return payloadPart + "." + signaturePart, nil
}

func (s *Service) verifyGoogleOAuthStateToken(token string) error {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return ErrMemberGoogleInvalidState
	}

	payloadPart := strings.TrimSpace(parts[0])
	signaturePart := strings.TrimSpace(parts[1])
	if payloadPart == "" || signaturePart == "" {
		return ErrMemberGoogleInvalidState
	}

	h := hmac.New(sha256.New, []byte(appSecret()))
	h.Write([]byte(payloadPart))
	expectedSignature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	if !hmac.Equal([]byte(expectedSignature), []byte(signaturePart)) {
		return ErrMemberGoogleInvalidState
	}

	payloadRaw, err := base64.RawURLEncoding.DecodeString(payloadPart)
	if err != nil {
		return ErrMemberGoogleInvalidState
	}

	state := &googleOAuthState{}
	if err := json.Unmarshal(payloadRaw, state); err != nil {
		return ErrMemberGoogleInvalidState
	}

	if strings.TrimSpace(state.Nonce) == "" || state.Exp == 0 || time.Now().Unix() >= state.Exp {
		return ErrMemberGoogleInvalidState
	}

	return nil
}

func (s *Service) exchangeGoogleCode(ctx context.Context, code string) (*googleTokenResponse, error) {
	conf := s.googleConfig()
	form := url.Values{}
	form.Set("code", code)
	form.Set("client_id", strings.TrimSpace(conf.GoogleClientId))
	form.Set("client_secret", strings.TrimSpace(conf.GoogleClientSecret))
	form.Set("redirect_uri", strings.TrimSpace(conf.GoogleRedirectUrl))
	form.Set("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, googleTokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, ErrMemberGoogleAuthFailed
	}

	decoded := &googleTokenResponse{}
	if err := json.NewDecoder(resp.Body).Decode(decoded); err != nil {
		return nil, err
	}

	if strings.TrimSpace(decoded.AccessToken) == "" {
		return nil, ErrMemberGoogleAuthFailed
	}

	return decoded, nil
}

func (s *Service) fetchGoogleUserInfo(ctx context.Context, accessToken string) (*googleUserInfoResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, googleUserInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(accessToken))
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, ErrMemberGoogleAuthFailed
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	profile := &googleUserInfoResponse{}
	if err := json.Unmarshal(body, profile); err != nil {
		return nil, err
	}

	if strings.TrimSpace(profile.Sub) == "" {
		return nil, ErrMemberGoogleAuthFailed
	}

	return profile, nil
}

func (s *Service) ensureGoogleMember(ctx context.Context, profile *googleUserInfoResponse) (string, string, error) {
	email := strings.ToLower(strings.TrimSpace(profile.Email))
	if email == "" {
		return "", "", ErrMemberGoogleAuthFailed
	}

	accounts, err := s.db.ListMemberAccounts(ctx)
	if err != nil {
		return "", "", err
	}

	for _, account := range accounts {
		if !strings.EqualFold(strings.TrimSpace(account.Username), email) {
			continue
		}
		if account.MemberID == nil {
			return "", "", ErrMemberGoogleAuthFailed
		}

		memberID := account.MemberID.String()
		if _, err := s.db.GetMemberByID(ctx, memberID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return "", "", ErrMemberGoogleAuthFailed
			}
			return "", "", err
		}

		return memberID, account.Username, nil
	}

	firstName := strings.TrimSpace(profile.GivenName)
	lastName := strings.TrimSpace(profile.FamilyName)
	if firstName == "" {
		firstName = strings.TrimSpace(profile.Name)
	}
	if firstName == "" {
		firstName = "Google"
	}

	displayName := strings.TrimSpace(profile.Name)
	if displayName == "" {
		displayName = strings.TrimSpace(firstName + " " + lastName)
	}
	if displayName == "" {
		displayName = email
	}

	rawPassword := "oauth-google:" + profile.Sub + ":" + fmt.Sprintf("%d", time.Now().UnixNano())
	hashedPassword, err := hashing.HashPassword(rawPassword)
	if err != nil {
		return "", "", err
	}

	member, err := s.db.CreateMemberWithAccount(
		ctx,
		nil,
		nil,
		firstName,
		lastName,
		displayName,
		"",
		email,
		string(hashedPassword),
	)
	if err != nil {
		return "", "", err
	}

	return member.ID.String(), email, nil
}
