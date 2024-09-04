package auth

import (
	"context"
	"net/url"

	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/pkg/errors"
)

func (s *AuthService) LoginRequest(_ context.Context, loginRequest *domain.LoginRequest) (*domain.LoginResp, error) {
	loginURL, err := url.ParseRequestURI(s.cfg.KeyCloak.Frontend.AuthURL)
	if err != nil {
		return nil, service.NewError(service.InternalError, errors.Wrap(err, "failed to parse auth url in config").Error())
	}

	query := loginURL.Query()
	query.Add("client_id", s.cfg.KeyCloak.Frontend.ClientID)
	query.Add("response_type", "code")
	query.Add("redirect_uri", loginRequest.RedirectURL.String())

	loginURL.RawQuery = query.Encode()
	resp := &domain.LoginResp{
		LoginURL: loginURL,
	}

	return resp, nil
}

func (s *AuthService) ClientTokenCallback(ctx context.Context, loginCallback *domain.LoginCallback) (*domain.ClientToken, error) {
	if err := s.validator.Struct(loginCallback); err != nil {
		return nil, service.NewError(service.BadRequest, errors.Wrap(err, "validation error").Error())
	}

	token, err := s.authRepository.GetAccessTokenFromClientCode(ctx, loginCallback.Code, loginCallback.RedirectURL.String())
	if err != nil {
		return nil, service.NewError(service.InternalError, errors.Wrap(err, "failed to get access token").Error())
	}

	return token, nil
}