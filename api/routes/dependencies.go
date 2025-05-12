package router

import (
	"github.com/aragornz325/piloto-api/internal/profile/handler"
	"github.com/aragornz325/piloto-api/internal/profile/service"
	"github.com/aragornz325/piloto-api/internal/user/handler"
	"github.com/aragornz325/piloto-api/internal/user/service"
	"github.com/aragornz325/piloto-api/internal/auth/handler"
	"github.com/aragornz325/piloto-api/internal/auth/service"
)

type AppDependencies struct {
	UserHandler *userHandler.UserHandler
	ProfileHandler *profileHandler.ProfileHandler
	AuthHandler *authHandler.AuthHandler
}

func BuildDependencies() *AppDependencies {
	// User
	userService := service.NewUserService()
	userHandler := userHandler.NewUserHandler(userService)
	// Profile
	profileService := profileService.NewProfileService()
	profileHandler := profileHandler.NewProfileHandler(profileService)
	//auth
	authService := authService.NewAuthService(userService)
	authHandler := authHandler.NewAuthHandler(authService)

	return &AppDependencies{
		UserHandler: userHandler,
		ProfileHandler: profileHandler,
		AuthHandler: authHandler,
	}
}
