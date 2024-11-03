package main

import (
	"log"

	"github.com/zde37/instashop-task/internal/controller"
)

// @title           Instashop API
// @version         1.0
// @description     A REST API for Instashop.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Instashop Support
// @contact.url    https://instashop.com/support
// @contact.email  support@instashop.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	if err := controller.New().Serve(); err != nil {
		log.Fatal(err)
	}
}
