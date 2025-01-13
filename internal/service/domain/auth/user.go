package auth

import (
	"context"
	"net/url"

	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/pkg/errors"
)

func (s *AuthService) Register(ctx context.Context, user *domain.RegisterUser) (*domain.User, error) {
	if err := s.validator.Struct(user); err != nil {
		return nil, service.NewError(service.BadRequest, errors.Wrap(err, "validation error").Error())
	}

	createdUser, err := s.userRepo.Create(ctx, &user.User, user.Password, user.Roles)
	if err != nil {
		return nil, service.NewError(service.InternalError, errors.Wrap(err, "failed to create user").Error())
	}

	return createdUser, nil
}

func (s *AuthService) LoginRequest(_ context.Context, loginRequest *domain.LoginRequest) (*domain.LoginResp, error) {
	loginURL, err := url.ParseRequestURI(s.cfg.OidcProvider.AuthURL)
	if err != nil {
		return nil, service.NewError(service.InternalError, errors.Wrap(err, "failed to parse auth url in config").Error())
	}

	query := loginURL.Query()
	query.Add("client_id", s.cfg.OidcProvider.Frontend.ClientID)
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

func (s *AuthService) LogoutRequest(ctx context.Context, logoutRequest *domain.Logout) error {
	if err := s.validator.Struct(logoutRequest); err != nil {
		return service.NewError(service.BadRequest, errors.Wrap(err, "validation error").Error())
	}

	err := s.userRepo.RemoveSession(ctx, logoutRequest.RefreshToken)
	if err != nil {
		return service.NewError(service.InternalError, errors.Wrap(err, "failed to remove user session").Error())
	}

	return nil
}

func (s *AuthService) GetAll(ctx context.Context) ([]*domain.User, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, service.NewError(service.InternalError, errors.Wrap(err, "failed to get all users").Error())
	}

	return users, nil
}

func (s *AuthService) GetByIDs(ctx context.Context, ids []string) ([]*domain.User, error) {
	users, err := s.userRepo.GetByIDs(ctx, ids)
	if err != nil {
		return nil, service.NewError(service.InternalError, errors.Wrap(err, "failed to get users by ids").Error())
	}

	return users, nil
}

func (s *AuthService) GetAllByRole(ctx context.Context, role domain.Role) ([]*domain.User, error) {
	users, err := s.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var filteredUsers []*domain.User
	for _, user := range users {
		for _, userRole := range user.Roles {
			if userRole.Name == role.Name {
				filteredUsers = append(filteredUsers, user)
				break
			}
		}
	}

	return filteredUsers, nil
}
