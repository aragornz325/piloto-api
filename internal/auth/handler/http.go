package authHandler

import (
	"net/http"

	authModel "github.com/aragornz325/piloto-api/internal/auth/model"
	authService "github.com/aragornz325/piloto-api/internal/auth/service"
	userModel "github.com/aragornz325/piloto-api/internal/user/model"
	"github.com/aragornz325/piloto-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type AuthHandler struct {
	AuthService authService.AuthService
}

var validate = validator.New()

func NewAuthHandler(authService authService.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

//--------------------------------------------------//

// @Summary Register a new user
// @Description Register a new user on the system
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body authModel.RegisterDTO true "Data to register a new user"
// @Success 201 {object} string
// @Failure 400 {object} ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) RegisterUser(c *gin.Context) {
	err := utils.PerformHandlerOperation(utils.PerformHandlerOperationFunc{
		Ctx:         c.Request.Context(),
		Name:        "RegisterUser",
		HandlerName: "auth",
		Operation: func() error {
			var payload authModel.RegisterDTO
			var user userModel.User

			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid input"})
				return err
			}

			if err := validate.Struct(payload); err != nil {
				c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
				return err
			}

			if err := utils.CopyNonNilFields(utils.CopyNonNilFieldsFuncParams{
				Source: &payload,
				Dest:   &user,
			}); err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
				return err
			}

			createdUser, err := h.AuthService.RegisterUser(authService.RegisterUserFuncParams{
				User: &user,
				Ctx:  c.Request.Context(),
			})
			if err != nil {
				c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
				return err
			}

			c.JSON(http.StatusCreated, createdUser)
			return nil
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "unexpected error occurred",
		})
	}
}

// @Summary Login a user
// @Description Login a user and return a JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body authModel.LoginDTO true "Data to login a user"
// @Success 200 {object} string
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) LoginUser(c *gin.Context) {
	err := utils.PerformHandlerOperation(utils.PerformHandlerOperationFunc{
		Ctx:         c.Request.Context(),
		Name:        "LoginUser",
		HandlerName: "auth",
		Operation: func() error {
			var payload authModel.LoginDTO

			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid input"})
				return err
			}

			if err := validate.Struct(payload); err != nil {
				c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
				return err
			}

			token, err := h.AuthService.LoginUser(authService.LoginUserFuncParams{
				Ctx:      c.Request.Context(),
				LoginDTO: &payload,
			})
			if err != nil {
				c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
				return err
			}

			c.JSON(http.StatusOK, gin.H{
				"token": token,
			})

			return nil
		},
	})

	if err != nil {
		// Solo logueás o devolvés un error genérico al cliente
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "unexpected error occurred",
		})
	}
}
