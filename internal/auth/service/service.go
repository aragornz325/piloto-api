package authService

import (
	"context"
	"time"

	authModel "github.com/aragornz325/piloto-api/internal/auth/model"
	userModel "github.com/aragornz325/piloto-api/internal/user/model"
	"github.com/aragornz325/piloto-api/internal/user/service"
	"github.com/aragornz325/piloto-api/pkg/errors"
	"github.com/aragornz325/piloto-api/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

// Auth package contains the logic for authentication and authorization
// using JWT tokens. It provides functions to generate, validate, and refresh tokens.

type AuthService interface {
	RegisterUser(RegisterUserFuncParams) (*userModel.User, error)
	LoginUser(LoginUserFuncParams) (*string, error)
}

type authService struct {
	UserService service.UserService
	JwtService jwtService
}

func NewAuthService(userService service.UserService) AuthService {
	return &authService{
		UserService: userService,
	}
}

// RegisterUser registers a new user in the system.
// It hashes the user's password, sets creation time and activation status,
// and delegates user creation to the UserService. If any error occurs during
// the process (password hashing or user creation), it returns a wrapped error.
// On success, it returns the created user and a nil error.
//
// Parameters:
//   - opts: RegisterUserFuncParams containing the user to register and context.
//
// Returns:
//   - *userModel.User: The newly created user.
//   - error: An error if registration fails, otherwise nil.
func (s *authService) RegisterUser(opts RegisterUserFuncParams) (*userModel.User, error) {
	var user *userModel.User
	err := utils.PerformServiceOperation(utils.PerformServiceOperationFunc{
		Ctx:  opts.Ctx,
		Name: "RegisterUser",
		Operation: func() error {
			now := time.Now().UTC()
			user = opts.User
			ctx := opts.Ctx
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				return errors.NewBadRequest(errors.ErrorFuncOptions{
					Message: "error hashing password",
					Err:     err,
				})
			}
			user.Password = string(hashedPassword)
			user.CreatedAt = now
			user.IsActive = true

			createdUser, err := s.
				UserService.
				CreateUser(service.CreateUserFuncParams{
					User: user,
					Ctx:  ctx,
				})
			if err != nil {
				return errors.NewBadRequest(errors.ErrorFuncOptions{
					Message: "error creating user",
					Err:     err,
				})
			}
			user = createdUser

			return nil
		},
	})

	if err != nil {
		return nil, errors.NewInternal(errors.ErrorFuncOptions{
			Message: "error creating user",
			Err:     err,
		})
	}

	return user, nil
}

// LoginUser authenticates a user using the provided login parameters.
// It retrieves the user by email and verifies the password using bcrypt.
// Returns the authenticated user on success, or an error if authentication fails.
//
// Parameters:
//   - opts: LoginUserFuncParams containing the context and login credentials.
//
// Returns:
//   - *userModel.User: The authenticated user object.
//   - error: An error if authentication fails or user retrieval encounters an issue.
func (s *authService) LoginUser(opts LoginUserFuncParams) (*string, error) {
	var token string

	err := utils.PerformServiceOperation(utils.PerformServiceOperationFunc{
		Ctx:  opts.Ctx,
		Name: "LoginUser",
		Operation: func() error {
			userDb, err := s.UserService.GetUserByEmail(service.GetUserByEmailFuncParams{
				Email: *opts.LoginDTO.Email,
				Ctx:   opts.Ctx,
			})
			if err != nil {
				return errors.NewBadRequest(errors.ErrorFuncOptions{
					Message: "error getting user",
					Err:     err,
				})
			}

			if err := bcrypt.CompareHashAndPassword([]byte(userDb.Password), []byte(*opts.LoginDTO.Password)); err != nil {
				return errors.NewUnauthorized(errors.SimpleErrorFuncOptions{
					Message: "invalid credentials",
				})
			}

			signedToken, err := s.JwtService.SignToken(SignTokenFuncParams{
				Ctx:    opts.Ctx,
				UserId: userDb.ID,
			})
			if err != nil {
				return errors.NewInternal(errors.ErrorFuncOptions{
					Message: "error signing token",
					Err:     err,
				})
			}

			token = signedToken
			return nil
		},
	})

	if err != nil {
		return nil, err
	}

	return &token, nil
}



// ---------Structs---------///
type RegisterUserFuncParams struct {
	User *userModel.User
	Ctx  context.Context
}

type LoginUserFuncParams struct {
	Ctx      context.Context
	LoginDTO *authModel.LoginDTO
}
