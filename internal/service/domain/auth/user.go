package auth

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

func (s *AuthService) Register(ctx context.Context, user *domain.RegisterUser) (*domain.User, error) {
	log := logger.GetLogger(ctx)
	if err := s.validator.Struct(user); err != nil {
		log.Debug("failed to validate registered user", "raw_user", fmt.Sprintf("%+v", user), "error", err)
		return nil, service.MapError(ctx, service.ErrValidation, service.ErrorLogValidation)
	}

	createdUser, err := s.userRepo.Create(ctx, &user.User, user.Password, user.Roles)
	if err != nil {
		log.Debug("failed to create user", "error", err, "user_name", user.User.Username)
		return nil, service.MapError(ctx, errors.Join(err, errors.New("failed to create user")), service.ErrorLogAll)
	}

	return createdUser, nil
}

func (s *AuthService) LoginRequest(ctx context.Context, loginRequest *domain.LoginRequest) *domain.LoginResp {
	log := logger.GetLogger(ctx)
	loginURL, err := url.ParseRequestURI(s.cfg.OidcProvider.AuthURL)
	if err != nil {
		log.Error("failed to parse auth url in config", "error", err, "auth_url", s.cfg.OidcProvider.AuthURL)
		panic("failed to parse auth url in config. Pleas check your configuration")
	}

	query := loginURL.Query()
	query.Add("client_id", s.cfg.OidcProvider.Frontend.ClientID)
	query.Add("response_type", "code")
	query.Add("redirect_uri", loginRequest.RedirectURL.String())

	loginURL.RawQuery = query.Encode()
	resp := &domain.LoginResp{
		LoginURL: loginURL,
	}

	return resp
}

func (s *AuthService) ClientTokenCallback(ctx context.Context, loginCallback *domain.LoginCallback) (*domain.ClientToken, error) {
	log := logger.GetLogger(ctx)
	if err := s.validator.Struct(loginCallback); err != nil {
		log.Debug("failed to validate client token callback", "raw_callback", fmt.Sprintf("%+v", loginCallback), "error", err)
		return nil, service.MapError(ctx, service.ErrValidation, service.ErrorLogValidation)
	}

	token, err := s.authRepository.GetAccessTokenFromClientCode(ctx, loginCallback.Code, loginCallback.RedirectURL.String())
	if err != nil {
		log.Debug("failed to get access token from auth flow", "error", err, "code", loginCallback.Code, "redirect_uri", loginCallback.RedirectURL.String())
		return nil, service.MapError(ctx, errors.Join(err, errors.New("failed to get access token")), service.ErrorLogAll)
	}

	return token, nil
}

func (s *AuthService) LogoutRequest(ctx context.Context, logoutRequest *domain.Logout) error {
	log := logger.GetLogger(ctx)
	if err := s.validator.Struct(logoutRequest); err != nil {
		log.Debug("failed to validate logout request", "raw_request", fmt.Sprintf("%+v", logoutRequest), "error", err)
		return service.MapError(ctx, service.ErrValidation, service.ErrorLogValidation)
	}

	err := s.userRepo.RemoveSession(ctx, logoutRequest.RefreshToken)
	if err != nil {
		log.Debug("failed to remove user session", "error", err)
		return service.MapError(ctx, errors.Join(err, errors.New("failed to remove user session")), service.ErrorLogAll)
	}

	return nil
}

func (s *AuthService) GetAll(ctx context.Context) ([]*domain.User, error) {
	log := logger.GetLogger(ctx)
	users, err := s.userRepo.GetAll(ctx) // TODO: Pagination
	if err != nil {
		log.Debug("failed to fetch all user lists", "error", err)
		return nil, service.MapError(ctx, err, service.ErrorLogEntityNotFound)
	}

	return users, nil
}

func (s *AuthService) GetByIDs(ctx context.Context, ids []string) ([]*domain.User, error) {
	log := logger.GetLogger(ctx)
	users, err := s.userRepo.GetByIDs(ctx, ids) // TODO: Pagination
	if err != nil {
		log.Debug("failed to fetch users by ids", "error", err, "user_ids", ids)
		return nil, service.MapError(ctx, err, service.ErrorLogEntityNotFound)
	}

	return users, nil
}

func (s *AuthService) GetAllByRole(ctx context.Context, role domain.UserRole) ([]*domain.User, error) {
	users, err := s.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var filteredUsers []*domain.User
	for _, user := range users {
		for _, userRole := range user.Roles {
			if userRole == role {
				filteredUsers = append(filteredUsers, user)
				break
			}
		}
	} // TODO: Move to repository

	return filteredUsers, nil
}
