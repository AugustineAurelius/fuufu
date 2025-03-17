package server

import (
	"context"
	"time"

	"github.com/AugustineAurelius/fuufu/api/auth"
	user_repository "github.com/AugustineAurelius/fuufu/internal/repository/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	auth.StrictServerInterface

	Repo      *user_repository.CommandRepository
	Telemetry trace.Tracer

	Secret string
}

// Example protected endpoint
// (GET /protected)
func (h *AuthHandler) GetProtected(ctx context.Context, request auth.GetProtectedRequestObject) (auth.GetProtectedResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Authenticate user
// (POST api/v1/auth/signin)
func (h *AuthHandler) PostApiV1AuthSignin(ctx context.Context, request auth.PostApiV1AuthSigninRequestObject) (auth.PostApiV1AuthSigninResponseObject, error) {
	ctx, span := h.Telemetry.Start(ctx, "AuthHandler.PostApiV1AuthSignin", trace.WithAttributes(
		attribute.String("username", request.Body.Username),
	))
	defer span.End()

	user, err := h.Repo.Get(ctx,
		user_repository.WithName(request.Body.Username),
	)
	if err != nil {
		return auth.PostApiV1AuthSignin500JSONResponse{
			Error: err.Error(),
		}, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(request.Body.Password)); err != nil {
		return auth.PostApiV1AuthSignin500JSONResponse{
			Error: err.Error(),
		}, nil
	}

	expireTime := time.Now().Add(24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Name,
		"exp":      expireTime.Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.Secret))
	if err != nil {
		return auth.PostApiV1AuthSignin500JSONResponse{
			Error: err.Error(),
		}, nil
	}

	return auth.PostApiV1AuthSignin200JSONResponse{
		Expires: &expireTime,
		Token:   &tokenString,
	}, nil
}

// Register new user
// (POST api/v1/auth/signup)
func (h *AuthHandler) PostApiV1AuthSignup(ctx context.Context, request auth.PostApiV1AuthSignupRequestObject) (auth.PostApiV1AuthSignupResponseObject, error) {
	ctx, span := h.Telemetry.Start(ctx, "AuthHandler.PostApiV1AuthSignup", trace.WithAttributes(
		attribute.String("username", request.Body.Username),
		attribute.String("email", string(request.Body.Email)),
	))
	defer span.End()

	id := uuid.New()
	created := time.Now().UTC()

	hashed, err := bcrypt.GenerateFromPassword([]byte(request.Body.Password), 10)
	if err != nil {
		return auth.PostApiV1AuthSignup500JSONResponse{
			Error: err.Error(),
		}, nil
	}

	if err = h.Repo.Create(ctx, &user_repository.User{
		ID:             id,
		Name:           request.Body.Username,
		Email:          string(request.Body.Email),
		HashedPassword: string(hashed),
		CreatedAt:      created,
	}); err != nil {
		return auth.PostApiV1AuthSignup500JSONResponse{
			Error: err.Error(),
		}, nil
	}

	return auth.PostApiV1AuthSignup201JSONResponse{
		Id:        id,
		CreatedAt: created,
		Username:  request.Body.Username,
	}, nil
}
