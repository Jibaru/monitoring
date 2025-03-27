package main

import (
	"context"
	"log"
	"os"
	"time"

	"monitoring/internal/handlers"
	"monitoring/internal/middlewares"

	_ "monitoring/docs" // Importa el paquete generado

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
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
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error cargando .env", err)
	}
	mongoURI, ok := os.LookupEnv("MONGODB_URI")
	if !ok {
		log.Fatal("MONGODB_URI no configurada")
	}

	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		log.Fatal("JWT_SECRET no configurada")
	}

	// Conexión a MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database("monitoringapp")

	// Inicializar router con Gin
	router := gin.Default()
	router.Use(middlewares.UseCORS())

	// Rutas para backoffice (protected via JWT)
	backoffice := router.Group("/api/v1/backoffice")
	{
		// Rutas de autenticación
		backoffice.POST("/register", handlers.Register(db))
		backoffice.POST("/login", handlers.Login(db, jwtSecret))
		// Rutas para gestión de Apps y Logs (requieren validación JWT)
		backoffice.Use(middlewares.HasAuthorization(jwtSecret))
		{
			backoffice.GET("/apps", handlers.ListApps(db))
			backoffice.POST("/apps", handlers.CreateApp(db))
			backoffice.DELETE("/apps/:appID", handlers.DeleteApp(db))
			backoffice.GET("/logs", handlers.SearchLogs(db))
		}
	}

	// Rutas para recepción de logs desde las apps
	appsGroup := router.Group("/api/v1/apps")
	{
		// Se espera que la autenticación se realice vía api key bearer
		appsGroup.POST("/:appID/logs", handlers.ReceiveLogs(db))
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
