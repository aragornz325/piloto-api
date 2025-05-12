package authService

import (
	"context"
	"fmt"
	"os"
	"time"

	authModel "github.com/aragornz325/piloto-api/internal/auth/model"
	"github.com/aragornz325/piloto-api/internal/user/service"
	 "github.com/aragornz325/piloto-api/pkg/errors"
	"github.com/aragornz325/piloto-api/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtService interface {
	SignToken(SignTokenFuncParams) (string, error)
	ValidateToken(TokenFuncParams) (bool, error)
	ParseToken(TokenFuncParams) (string, error)
	GenerateJWT(GenerateTokenFuncParams) (string, error)
}

type jwtService struct {
	UserService service.UserService
}

func NewJwtService() JwtService {
	return &jwtService{}
}

// SignToken generates a signed JWT token for a user specified by the given options.
// It retrieves the user by ID, determines the user's role (defaulting to "user" if not set),
// constructs the token payload, and signs the JWT. The token is valid for 24 hours.
// Returns the signed JWT token string or an error if the operation fails.
// The function uses the PerformServiceOperation utility to handle errors and logging.
// The token is signed using the HS256 algorithm and a secret key from the environment variable "
// JWT_SECRET".
func (s *jwtService) SignToken(opts SignTokenFuncParams) (string, error) {
	var signedToken string

	err := utils.PerformServiceOperation(utils.PerformServiceOperationFunc{
		Name:        "SignToken",
		ServiceName: "JWT Service",
		Operation: func() error {
			user, err := s.UserService.GetUserById(service.GetUserByIdFuncParams{
				Ctx:    opts.Ctx,
				UserId: opts.UserId,
			})
			if err != nil {
				return errors.NewBadRequest(errors.ErrorFuncOptions{
					Message: "error getting user",
					Err:     err,
				})
			}

			var role string
			if len(user.Role) > 0 {
				role = user.Role[0]
			} else {
				role = "user"
			}

			tokenPayload := authModel.TokenPayload{
				UserId: user.ID,
				Email:  user.Email,
				Role:   role,
				Exp:    time.Now().Add(24 * time.Hour).Unix(),
			}

			signedToken, err = s.GenerateJWT(GenerateTokenFuncParams{
				Ctx:     opts.Ctx,
				Payload: tokenPayload,
			})
			if err != nil {
				return errors.NewBadRequest(errors.ErrorFuncOptions{
					Message: "error generating JWT",
					Err:     err,
				})
			}

			return nil
		},
	})

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// GenerateJWT generates a JSON Web Token (JWT) using the provided options.
// It creates a token with user-specific claims such as user ID, email, role, and expiration time.
// The token is signed using the HS256 algorithm and a secret key from the environment variable "JWT_SECRET".
// Returns the signed JWT as a string, or an error if token generation fails.
func (s *jwtService) GenerateJWT(opts GenerateTokenFuncParams) (string, error) {
	var signedToken string

	err := utils.PerformServiceOperation(utils.PerformServiceOperationFunc{
		Name:        "GenerateJWT",
		ServiceName: "JWT Service",
		Operation: func() error {
			claims := jwt.MapClaims{
				"user_id": opts.Payload.UserId,
				"email":   opts.Payload.Email,
				"role":    opts.Payload.Role,
				"exp":     opts.Payload.Exp,
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			var err error
			signedToken, err = token.SignedString([]byte(os.Getenv("JWT_SECRET")))
			if err != nil {
				return errors.NewInternal(errors.ErrorFuncOptions{
					Message: "error signing token",
					Err:     err,
				})
			}

			return nil
		},
	})

	if err != nil {
		return "", errors.NewBadRequest(errors.ErrorFuncOptions{
			Message: "error generating JWT",
			Err:     err,
		})
	}

	return signedToken, nil
}

// ValidateToken validates a JWT token using the provided TokenFuncParams.
// It parses the token, checks its signing method, and verifies its validity
// using the secret from the environment variable "JWT_SECRET".
// Returns true if the token is valid, otherwise returns false and an error
// describing the validation failure.
func (s *jwtService) ValidateToken(opts TokenFuncParams) (bool, error) {
	var valid bool

	err := utils.PerformServiceOperation(utils.PerformServiceOperationFunc{
		Name:        "ValidateToken",
		ServiceName: "JWT Service",
		Operation: func() error {
			token, err := jwt.Parse(opts.Token, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(os.Getenv("JWT_SECRET")), nil
			})

			if err != nil || !token.Valid {
				valid = false
				return fmt.Errorf("invalid token: %w", err)
			}

			valid = true
			return nil
		},
	})

	if err != nil {
		return false, errors.NewBadRequest(errors.ErrorFuncOptions{
			Message: "error validating token",
			Err:     err,
		})
	}

	return valid, nil
}

// ParseToken parses a JWT token string provided in the TokenFuncParams and extracts the user ID from its claims.
// It validates the token's signing method and checks the token's validity. If successful, it returns the user ID as a string.
// Returns an error if the token is invalid, the signing method is unexpected, or the user ID claim is missing or malformed.
func (s *jwtService) ParseToken(opts TokenFuncParams) (string, error) {
	var userID string

	err := utils.PerformServiceOperation(utils.PerformServiceOperationFunc{
		Name:        "ParseToken",
		ServiceName: "JWT Service",
		Operation: func() error {
			token, err := jwt.Parse(opts.Token, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.NewInternal(errors.ErrorFuncOptions{
						Message: "unexpected signing method",
						Err:     fmt.Errorf("unexpected signing method: %v", token.Header["alg"]),
					})
				}
				return []byte(os.Getenv("JWT_SECRET")), nil
			})

			if err != nil {
				return errors.NewInternal(errors.ErrorFuncOptions{
					Message: "error parsing token",
					Err:     err,
				})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				return errors.NewInternal(errors.ErrorFuncOptions{
					Message: "invalid token claims",
					Err:     fmt.Errorf("invalid token claims"),
				})
			}

			uid, ok := claims["user_id"].(string)
			if !ok {
				return errors.NewInternal(errors.ErrorFuncOptions{
					Message: "invalid user ID in token claims",
					Err:     fmt.Errorf("invalid user ID in token claims"),
				})
			}

			userID = uid
			return nil
		},
	})

	if err != nil {
		return "", errors.NewBadRequest(errors.ErrorFuncOptions{
			Message: "error parsing token",
			Err:     err,
		})
	}

	return userID, nil
}


// -------structs-------//
type SignTokenFuncParams struct {
	Ctx    context.Context
	UserId uuid.UUID
}

type TokenFuncParams struct {
	Ctx   context.Context
	Token string
}

type GenerateTokenFuncParams struct {
	Ctx     context.Context
	Payload authModel.TokenPayload
}
