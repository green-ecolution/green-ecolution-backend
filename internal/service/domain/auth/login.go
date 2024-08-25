package auth

import (
	"context"
	"net/url"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities/auth"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/pkg/errors"
)

func (s *AuthService) LoginRequest(_ context.Context, loginRequest *auth.LoginRequest) (*auth.LoginResp, error) {
	loginURL, err := url.Parse(s.cfg.IdentityAuth.KeyCloak.Frontend.AuthURL)
	if err != nil {
		return nil, service.NewError(service.InternalError, errors.Wrap(err, "failed to parse auth url").Error())
	}

	query := loginURL.Query()
	query.Add("client_id", s.cfg.IdentityAuth.KeyCloak.Frontend.ClientID)
	query.Add("response_type", "code")
	query.Add("redirect_uri", loginRequest.RedirectURL.String())

	loginURL.RawQuery = query.Encode()
	resp := &auth.LoginResp{
		LoginURL: loginURL,
	}

	return resp, nil
}

func (s *AuthService) ClientTokenCallback(ctx context.Context, loginCallback *auth.LoginCallback) (*auth.ClientToken, error) {
	if err := s.validator.Struct(loginCallback); err != nil {
		return nil, service.NewError(service.BadRequest, errors.Wrap(err, "validation error").Error())
	}

	token, err := s.authRepository.GetAccessTokenFromClientCode(ctx, loginCallback.Code, loginCallback.RedirectURL.String())
	if err != nil {
		return nil, service.NewError(service.InternalError, errors.Wrap(err, "failed to get access token").Error())
	}

	return token, nil
}
