package profileHandler

import (
	"net/http"

	m "github.com/aragornz325/piloto-api/internal/profile/model"
	profileService "github.com/aragornz325/piloto-api/internal/profile/service"
	"github.com/aragornz325/piloto-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)


type ErrorResponse struct {
	Error string `json:"error"`
}

type ProfileHandler struct {
	ProfileService profileService.ProfileService
}

var validate = validator.New()

func NewProfileHandler(profileService profileService.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		ProfileService: profileService,
	}
}

//----------------------------------------------------

// @Summary Create a new Profile
// @Description Create a user profile
// @Tags profile
// @Accept json
// @Produce json
// @Param input body profileModel.UserProfileDTO true "Data to create a new user profile"
// @Success 201 {object} profileModel.Profile
// @Failure 400 {object} ErrorResponse
// @Router /profile [post]
func (h *ProfileHandler) CreateProfileHandler(c *gin.Context) {
	var payload m.UserProfileDTO
	var profile m.Profile

	// Bind JSON
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	// Validar struct con validator
	if err := validate.Struct(payload); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	// Mapear campos no nulos al modelo
	if err := utils.CopyNonNilFields(&payload, &profile); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// Llamar al servicio
	result, err := h.ProfileService.CreateProfile(profileService.CreateProfileFuncParams{
		Ctx:     c.Request.Context(),
		Profile: &profile,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// @Summary Get profile by ID
// @Description Get a user profile by ID
// @Tags profile
// @Accept json
// @Produce json
// @Param id path string true "Profile ID"
// @Success 200 {object} profileModel.Profile
// @Failure 404 {object} ErrorResponse
// @Router /profile/{id} [get]
func (h *ProfileHandler) GetProfileByIdHandler(c *gin.Context) {
	id := c.Param("id")
	userId, err := uuid.Parse(id)
	// Llamar al servicio
	result, err := h.ProfileService.GetUserProfile(profileService.GeProfileByUserIdFuncParams{
		Ctx:      c.Request.Context(),
		UserId:  userId,
	})

	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Update profile
// @Description Update a user profile
// @Tags profile
// @Accept json
// @Produce json
// @Param id path string true "Profile ID"
// @Param input body profileModel.UserProfileDTO true "Data to update a user profile"
// @Success 200 {object} profileModel.Profile
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /profile/{id} [put]
func (h *ProfileHandler) UpdateProfileHandler(c *gin.Context) {	
	id := c.Param("id")
	userId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	var payload m.UserProfileDTO
	var profile m.Profile

	// Bind JSON
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Validar struct con validator
	if err := validate.Struct(payload); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Mapear campos no nulos al modelo
	if err := utils.CopyNonNilFields(&payload, &profile); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	result, err := h.ProfileService.UpdateUserProfile(profileService.UpdateUserProfileFuncParams{
		Ctx:     c.Request.Context(),
		UserId:  userId,
		Profile: &profile,
	})

	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Soft delete profile
// @Description Soft delete a user profile
// @Tags profile
// @Accept json
// @Produce json
// @Param id path string true "Profile ID"
// @Success 200 {object} profileModel.Profile
// @Failure 404 {object} ErrorResponse
// @Router /profile/{id} [delete]
func (h *ProfileHandler) SoftDeleteProfileHandler(c *gin.Context) {
	id := c.Param("id")
	userId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	result, err := h.ProfileService.SoftDeleteUserProfile(profileService.SoftDeleteUserProfileFuncParams{
		Ctx:     c.Request.Context(),
		UserId:  userId,
	})

	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}