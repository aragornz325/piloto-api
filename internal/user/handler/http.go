package userHandler

import (
	"net/http"

	m "github.com/aragornz325/piloto-api/internal/user/model"
	userService "github.com/aragornz325/piloto-api/internal/user/service"
	"github.com/aragornz325/piloto-api/pkg/errors"
	"github.com/aragornz325/piloto-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type UserHandler struct {
	UserService userService.UserService
}

var validate = validator.New()

func NewUserHandler(userService userService.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

//----------------------------------------------------

// @Summary Create a new user
// @Description Create and register a new user on the system
// @Tags users
// @Accept json
// @Produce json
// @Param input body userModel.CreateUserInput true "Data to create a new user"
// @Success 201 {object} userModel.User
// @Failure 400 {object} ErrorResponse
// @Router /users [post]
func (h *UserHandler) CreateUserHandler(c *gin.Context) {
	var payload m.CreateUserInput
	var user m.User

	// Bind JSON
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	// Validate struct with validator
	if err := validate.Struct(payload); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if *payload.Email == "" || *payload.Password == "" {
		c.JSON(http.StatusBadRequest, errors.
			NewBadRequest(errors.ErrorFuncOptions{
				Message: "Email and password are required",
				Err:     nil,
			}))
		return
	}
	// Map non-null fields to the model
	if err := utils.
		CopyNonNilFields(utils.CopyNonNilFieldsFuncParams{
			Source: payload,
			Dest:   &user,
		}); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// create user
	createdUser, err := h.UserService.
		CreateUser(userService.
			CreateUserFuncParams{
			Ctx:  c.Request.Context(),
			User: &user,
		})

	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

// @Summary Get all users
// @Description Get all users in the system
// @Tags users
// @Produce json
// @Success 200 {array} userModel.User
// @Failure 400 {object} ErrorResponse
// @Router /users [get]
func (h *UserHandler) GetAllUsersHandler(c *gin.Context) {
	users, err := h.UserService.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// @Summary Get user by ID
// @Description Get a user by their ID
// @Tags users
// @Param id path string true "User ID"
// @Produce json
// @Success 200 {object} userModel.User
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByIdHandler(c *gin.Context) {
	idParan := c.Param("id")
	userId, err := uuid.Parse(idParan)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	user, err := h.UserService.GetUserById(userService.GetUserByIdFuncParams{
		Ctx:    c.Request.Context(),
		UserId: userId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Update user
// @Description Update an existing user
// @Tags users
// @Accept json
// @Param id path string true "User ID"
// @Param input body userModel.CreateUserInput true "Data to update the user"
// @Produce json
// @Success 200 {object} userModel.User
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUserHandler(c *gin.Context) {
	idParan := c.Param("id")
	userId, err := uuid.Parse(idParan)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}
	var payload m.CreateUserInput
	var user m.User

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
	if err := utils.CopyNonNilFields(utils.CopyNonNilFieldsFuncParams{
		Source: &payload,
		Dest:   &user,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	result, err := h.UserService.UpdateUser(userService.UpdateUserFuncParams{
		Ctx:    c.Request.Context(),
		UserId: userId,
		User:   &user,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Soft delete user
// @Description Soft delete a user by their ID
// @Tags users
// @Accept json
// @Param id path string true "User ID"
// @Produce json
// @Success 200 {object} userModel.User
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) SoftDeleteUserHandler(c *gin.Context) {
	idParan := c.Param("id")
	userId, err := uuid.Parse(idParan)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}
	user, err := h.UserService.SoftDeleteUser(userService.GetUserByIdFuncParams{
		Ctx:    c.Request.Context(),
		UserId: userId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
