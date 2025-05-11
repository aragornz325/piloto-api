package profileService

import (
	"context"
	"fmt"
	"time"

	profileModel "github.com/aragornz325/piloto-api/internal/profile/model"
	db "github.com/aragornz325/piloto-api/pkg/database"
	"github.com/google/uuid"
)

type ProfileService interface {
	CreateProfile(CreateProfileFuncParams) (*profileModel.Profile, error)
	GetUserProfile(GeProfileByUserIdFuncParams) (*profileModel.Profile, error)
	UpdateUserProfile(UpdateUserProfileFuncParams) (*profileModel.Profile, error)
	SoftDeleteUserProfile(SoftDeleteUserProfileFuncParams) (*profileModel.Profile, error)
}

type profileService struct{}

// NewUserService devuelve una instancia de UserService
func NewProfileService() ProfileService {
	return &profileService{}
}

// CreateProfile creates a new user profile in the database using the provided options.
// It initializes a Profile model with the data from opts.Profile and saves it to the database
// within the context specified by opts.Ctx. Returns the created Profile and any error encountered.
//
// Parameters:
//   - opts: CreateProfileFuncParams containing the context and profile data.
//
// Returns:
//   - *profileModel.Profile: Pointer to the newly created profile.
//   - error: Error encountered during creation, or nil if successful.
func (s *profileService) CreateProfile(opts CreateProfileFuncParams) (*profileModel.Profile, error) {
	if err := db.DB.WithContext(opts.Ctx).Create(opts.Profile).Error; err != nil {
		return nil, err
	}
	opts.Profile.CreatedAt = time.Now()
	return opts.Profile, nil
}

// GetUserProfile retrieves the user profile(s) associated with the given user ID from the database.
// It accepts a GeProfileByUserIdFuncParams struct containing the context and user ID.
// Returns a pointer to a Profile model and an error if the operation fails.
func (s *profileService) GetUserProfile(opts GeProfileByUserIdFuncParams) (*profileModel.Profile, error) {
	fmt.Println(opts.Ctx, opts.UserId)
	var profile profileModel.Profile
	if err := db.DB.WithContext(opts.Ctx).Where("user_id = ?", opts.UserId).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

// UpdateUserProfile updates the profile information of a user identified by UserId.
// It retrieves the existing profile from the database, updates its fields with the values
// provided in opts.Profile, and saves the changes back to the database.
// Returns the updated profile on success, or an error if the operation fails.
//
// Parameters:
//   - opts: UpdateUserProfileFuncParams containing the context, updated profile data, and user ID.
//
// Returns:
//   - *profileModel.Profile: Pointer to the updated profile.
//   - error: Error if the update fails, otherwise nil.
func (s *profileService) UpdateUserProfile(opts UpdateUserProfileFuncParams) (*profileModel.Profile, error) {
	var profile profileModel.Profile
	if err := db.DB.WithContext(opts.Ctx).First(&profile, opts.UserId).Error; err != nil {
		return nil, err
	}

	if err := db.DB.WithContext(opts.Ctx).Save(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

// SoftDeleteUserProfile performs a soft delete operation on a user's profile identified by the provided parameters.
// It retrieves the user profile using the given context and user ID, then marks the profile as deleted in the database
// without physically removing the record. Returns the deleted profile and any error encountered during the process.
//
// Parameters:
//   - opts: SoftDeleteUserProfileFuncParams containing the context and user ID.
//
// Returns:
//   - *profileModel.Profile: Pointer to the deleted profile.
//   - error: Error encountered during retrieval or deletion, if any.
func (s *profileService) SoftDeleteUserProfile(opts SoftDeleteUserProfileFuncParams) (*profileModel.Profile, error) {
	profilePtr, err := s.GetUserProfile(GeProfileByUserIdFuncParams(opts))
	if err != nil {
		return nil, err
	}
	profilePtr.IsActive = false
	profilePtr.DeletedAt.Time = time.Now()
	if _, err := s.UpdateUserProfile(UpdateUserProfileFuncParams{
		Ctx:     opts.Ctx,
		Profile: profilePtr,
		UserId:  opts.UserId,
	}); err != nil {	
		return nil, err
	}
	return profilePtr, nil		
}

// /------ structs ------///
type CreateProfileFuncParams struct {
	Profile *profileModel.Profile
	Ctx     context.Context
}

type GeProfileByUserIdFuncParams struct {
	UserId uuid.UUID
	Ctx    context.Context
}
type UpdateUserProfileFuncParams struct {
	Profile *profileModel.Profile
	Ctx     context.Context
	UserId  uuid.UUID
}
type SoftDeleteUserProfileFuncParams struct {
	UserId uuid.UUID
	Ctx    context.Context
}