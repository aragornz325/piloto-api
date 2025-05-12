package router

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(deps *AppDependencies) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	v1 := r.Group("/api/v1")
	user := v1.Group("/users")
	profile := v1.Group("/profile")
	auth := v1.Group("/auth")
	{
		user.GET("/", deps.UserHandler.GetAllUsersHandler)
		user.POST("/", deps.UserHandler.CreateUserHandler)
		user.GET("/:id", deps.UserHandler.GetUserByIdHandler)
		user.PUT("/:id", deps.UserHandler.UpdateUserHandler)
		user.DELETE("/:id", deps.UserHandler.SoftDeleteUserHandler)
	}
	{
		profile.POST("/", deps.ProfileHandler.CreateProfileHandler)
		profile.GET("/:id", deps.ProfileHandler.GetProfileByIdHandler)
		profile.PUT("/:id", deps.ProfileHandler.UpdateProfileHandler)
		profile.DELETE("/:id", deps.ProfileHandler.SoftDeleteProfileHandler)
	}
	{
		auth.POST("/register", deps.AuthHandler.RegisterUser)
		auth.POST("/login", deps.AuthHandler.LoginUser)
	}

	return r
}
