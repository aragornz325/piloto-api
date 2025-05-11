package service

import (
	"context"
	"fmt"
	"time"

	userModel "github.com/aragornz325/piloto-api/internal/user/model"
	db "github.com/aragornz325/piloto-api/pkg/database"

	"github.com/aragornz325/piloto-api/pkg/utils"
	"github.com/google/uuid"
)

// Package service contiene la lógica de negocio para el manejo de usuarios

// UserService define la interfaz pública del servicio
// que puede ser implementada por diferentes versiones (mock, prod, etc).
type UserService interface {
	CreateUser(CreateUserFuncParams) (*userModel.User, error)
	GetAllUsers(ctx context.Context) ([]*userModel.User, error)
	GetUserById(GetUserByIdFuncParams) (*userModel.User, error)
	UpdateUser(UpdateUserFuncParams) (*userModel.User, error)
	SoftDeleteUser(GetUserByIdFuncParams) (*userModel.User, error)
	GetUserByEmail(GetUserByEmailFuncParams) (*userModel.User, error)
}

type userService struct{}

func NewUserService() UserService {
	return &userService{}
}


// CreateUser creates a new user in the database with the provided options.
// It sets the CreatedAt timestamp to the current time and marks the user as active.
// The operation is performed within the context specified in opts.Ctx.
// Returns the created user and any error encountered during the operation.
// Parameters:
//   - opts: CreateUserFuncParams containing the user data and context.
// Returns:
//   - *userModel.User: Pointer to the created user model.
//   - error: An error if the operation fails, otherwise nil.
// This function uses a service operation wrapper to handle the database interaction
// and error handling. It creates a new user record in the database.
// The user is marked as active and the CreatedAt timestamp is set to the current time.
// The function returns the created user model and any error that occurred during the operation.
// The function uses a service operation wrapper to handle the database interaction
// and error handling. It creates a new user record in the database.
func (s *userService) CreateUser(opts CreateUserFuncParams) (*userModel.User, error) {
	opts.User.CreatedAt = time.Now()
	opts.User.IsActive = true
	 err := utils.PerformServiceOperation(utils.PerformServiceOperationFunc{
		Ctx:  opts.Ctx,
		Name: "CreateUser",
		Operation: func() error {
			if err := db.DB.
			WithContext(opts.Ctx).
			Create(opts.User).Error; err != nil {
				return fmt.Errorf("error creating user: %w", err)
			}
			return nil
		},
	})
	if err != nil {
		return nil, err
	}
	return opts.User, nil
}

// GetAllUsers retrieves all active users from the database.
// It executes the operation within the provided context and returns a slice of user models.
// If an error occurs during the database operation, it returns the error.
// Parameters:
//   - ctx: The context for the operation.
// Returns:
//   - []*userModel.User: A slice of pointers to user models.
//   - error: An error if the operation fails, otherwise nil.
// This function uses a service operation wrapper to handle the database interaction
// and error handling. It queries the database for all users where is_active is true.
func (s *userService) GetAllUsers(ctx context.Context) ([]*userModel.User, error) {
	var users []*userModel.User
	err := utils.PerformServiceOperation(utils.PerformServiceOperationFunc{
		Ctx:  ctx,
		Name: "GetAllUsers",
		Operation: func() error {
			if err := db.DB.
			WithContext(ctx).
			Where("is_active = ?", true).
			Find(&users).Error; err != nil {
				return fmt.Errorf("error getting all users: %w", err)
			}
			return nil
		},
	})
	if err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateUser updates the user information in the database based on the provided options.
// It retrieves the user by ID, updates the user's fields and the UpdatedAt timestamp,
// and persists the changes using a service operation wrapper. Returns the updated user
// or an error if the operation fails.
//
// Parameters:
//   - opts: UpdateUserFuncParams containing the context, user ID, and updated user data.
//
// Returns:
//   - *userModel.User: The updated user model.
//   - error: An error if the update operation fails.
func (s *userService) UpdateUser(opts UpdateUserFuncParams) (*userModel.User, error) {
	userDb, err := s.GetUserById(GetUserByIdFuncParams{
		Ctx:    opts.Ctx,
		UserId: opts.UserId,
	})
	if err != nil {
		return nil, err
	}
	opts.User.UpdatedAt = time.Now()
	err = utils.PerformServiceOperation(utils.PerformServiceOperationFunc{
		Ctx:  opts.Ctx,
		Name: "UpdateUser",
		Operation: func() error {
			if err := db.DB.
			WithContext(opts.Ctx).
			Model(userDb).
			Updates(opts.User).Error; err != nil {
				return fmt.Errorf("error updating user: %w", err)
			}
			return nil
		},
	})
	if err != nil {
		return nil, err
	}
	return userDb, nil

}

// GetUserById retrieves a user by their unique ID if the user is active.
// It performs the operation within the provided context and returns the user model
// or an error if the user could not be found or another error occurs during the operation.
//
// Parameters:
//   - opts: GetUserByIdFuncParams containing the context and user ID.
//
// Returns:
//   - *userModel.User: Pointer to the retrieved user model.
//   - error: Error encountered during the operation, or nil if successful.
func (s *userService) GetUserById(opts GetUserByIdFuncParams) (*userModel.User, error) {
	var user userModel.User
	err := utils.PerformServiceOperation(utils.PerformServiceOperationFunc{
		Ctx:  opts.Ctx,
		Name: "GetUserById",
		Operation: func() error {
			if err := db.DB.WithContext(opts.Ctx).
				Preload("Profile").
				Where("is_active = ?", true).
				First(&user, opts.UserId).Error; err != nil {
				return fmt.Errorf("error getting user by ID: %w", err)
			}
			return nil
		},
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// SoftDeleteUser performs a soft delete on a user by setting the DeletedAt timestamp to the current time
// and marking the user as inactive. It retrieves the user by ID, updates the relevant fields, and persists
// the changes using the UpdateUser method. Returns the updated user or an error if the operation fails.
//
// Parameters:
//   - opts: GetUserByIdFuncParams containing the context and user ID.
//
// Returns:
//   - *userModel.User: Pointer to the updated user model.
//   - error: Error if the operation fails, otherwise nil.
func (s *userService) SoftDeleteUser(opts GetUserByIdFuncParams) (*userModel.User, error) {
	now := time.Now()
	userPtr, err := s.GetUserById(GetUserByIdFuncParams{
		Ctx:    opts.Ctx,
		UserId: opts.UserId,
	})
	if err != nil {
		return nil, fmt.Errorf("error obteniendo el usuario: %w", err)
	}

	userPtr.DeletedAt.Scan(now)
	userPtr.IsActive = false

	err = utils.PerformServiceOperation(utils.PerformServiceOperationFunc{
		Ctx:  opts.Ctx,
		Name: "SoftDeleteUser",
		Operation: func() error {
			_, err := s.UpdateUser(UpdateUserFuncParams{
				User:   userPtr,
				Ctx:    opts.Ctx,
				UserId: opts.UserId,
			})
			if err != nil {
				return fmt.Errorf("error soft deleting user: %w", err)
			}
			return nil
		},
	})
	if err != nil {
		return nil, err
	}

	return userPtr, nil
}

// GetUserByEmail retrieves a user by their email.
// Params:
//   - opts: contains context and the email string
//
// Returns:
//   - *User if found
//   - error if the query fails or user is not found
func (s *userService) GetUserByEmail(opts GetUserByEmailFuncParams) (*userModel.User, error) {
	var user userModel.User

	err := utils.PerformServiceOperation(utils.PerformServiceOperationFunc{
		Ctx:  opts.Ctx,
		Name: "GetUserByEmail",
		Operation: func() error {
			if err := db.DB.WithContext(opts.Ctx).Where("email = ?", opts.Email).First(&user).Error; err != nil {
				return fmt.Errorf("error getting user by email: %w", err)
			}
			return nil
		},
	})

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// /-------- structs ---------///
type CreateUserFuncParams struct {
	User *userModel.User
	Ctx  context.Context
}

type UpdateUserFuncParams struct {
	User   *userModel.User
	Ctx    context.Context
	UserId uuid.UUID
}

type GetUserByIdFuncParams struct {
	Ctx    context.Context
	UserId uuid.UUID
}
type GetUserByEmailFuncParams struct {
	Ctx   context.Context
	Email string
}
