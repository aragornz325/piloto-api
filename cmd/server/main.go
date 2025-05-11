package main

import (
	"log"
	"os"

	routes "github.com/aragornz325/piloto-api/api/routes"
	_ "github.com/aragornz325/piloto-api/docs"
	db "github.com/aragornz325/piloto-api/pkg/database"
	"github.com/aragornz325/piloto-api/pkg/logger"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Piloto de Tormenta API
//	@version		1.0
//	@description	This is a piloto de tormenta API.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://pilotodetormenta/support
//	@contact.email	rodrigo.m.quintero@gmail.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:3800
//	@BasePath	/api/v1

//	@securityDefinitions.basic	BasicAuth

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
func main() {

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	envFile := "." +env + ".env"
	log.Println("Cargando archivo de entorno: ", envFile)
	if err := godotenv.Load(envFile); err != nil {
		log.Fatal("💥 Error cargando archivo .env")
	}
	logger.Init()
	db.Init()
	db.ExecuteMigrations()
	deps := routes.BuildDependencies()
	r := routes.SetupRoutes(deps)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":3800"); err != nil {
		log.Fatal("💥 Error al iniciar el servidor: ", err)
	}
}
