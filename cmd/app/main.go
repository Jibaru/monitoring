package main

import (
	"context"
	"monitoring/config"
	"monitoring/db"
	"monitoring/server"

	_ "monitoring/docs"
)

// @title           Monitoring API
// @version         1.0
// @description     This is the monitoring API.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	cfg := config.Load()
	db, client := db.New(cfg)
	defer client.Disconnect(context.Background())
	router := server.New(cfg, db)
	router.Run(":" + cfg.APIPort)
}
