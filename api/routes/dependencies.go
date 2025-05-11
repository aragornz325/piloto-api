package router

import (
	"github.com/aragornz325/piloto-api/internal/profile/handler"
	"github.com/aragornz325/piloto-api/internal/profile/service"
	"github.com/aragornz325/piloto-api/internal/user/handler"
	"github.com/aragornz325/piloto-api/internal/user/service"
)

type AppDependencies struct {
	UserHandler *userHandler.UserHandler
	ProfileHandler *profileHandler.ProfileHandler
}

func BuildDependencies() *AppDependencies {
	userService := service.NewUserService()
	profileService := profileService.NewProfileService()
	userHandler := userHandler.NewUserHandler(userService)
	profileHandler := profileHandler.NewProfileHandler(profileService)

	return &AppDependencies{
		UserHandler: userHandler,
		ProfileHandler: profileHandler,
	}
}
