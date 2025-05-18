package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/amirzayi/clean_architect/internal/domain"
	"github.com/amirzayi/clean_architect/pkg/auth"
	"github.com/amirzayi/clean_architect/pkg/errs"
	"github.com/amirzayi/clean_architect/pkg/hash"
	"golang.org/x/crypto/bcrypt"
)

type Auth interface {
	Register(ctx context.Context, auth domain.Auth) error
	Login(ctx context.Context, auth domain.Auth) (token string, err error)
}

type authService struct {
	userService User
	hasher      hash.PasswordHasher
	authManager auth.Manager
	logger      *slog.Logger
}

func NewAuthService(userService User, hasher hash.PasswordHasher, authManager auth.Manager) Auth {
	return &authService{
		userService: userService,
		hasher:      hasher,
		authManager: authManager,
	}
}

func (a *authService) Register(ctx context.Context, auth domain.Auth) error {
	pwd, err := a.hasher.Hash(auth.Password)
	if err != nil {
		a.logger.Error("failed to create hashed password", slog.Any("error", err))
		return errs.New(err, errs.CodeInternal)
	}

	user := domain.User{
		Email:       auth.Email,
		PhoneNumber: auth.PhoneNumber,
		Password:    pwd,
		Role:        domain.UserRoleNormal,
	}

	_, err = a.userService.Create(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
func (a *authService) Login(ctx context.Context, auth domain.Auth) (string, error) {
	// todo: add login via phone no
	user, err := a.userService.GetByEmail(ctx, auth.Email)
	if err != nil {
		return "", err
	}

	if user.Status == domain.UserStatusBanned {
		return "", errs.New(errors.New("user banned"), errs.CodeForbiddenAccess)
	}

	if err = a.hasher.Compare(user.Password, auth.Password); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", errs.NotFound("user by given credentials")
		}
		a.logger.Error("failed to compare hashed password", slog.Any("error", err))
		return "", errs.New(err, errs.CodeInternal)
	}

	token, err := a.authManager.CreateToken(user.ID, string(user.Role))
	if err != nil {
		a.logger.Error("failed to create token", slog.Any("error", err))
		return "", errs.New(err, errs.CodeInternal)
	}
	return token, nil
}
