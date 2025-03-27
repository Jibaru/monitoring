package server

import (
	"context"
	"log"

	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"monitoring/config"
	"monitoring/internal/handlers"
	"monitoring/internal/middlewares"
)

func New(cfg config.Config) *gin.Engine {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database(cfg.DBName)

	router := gin.Default()
	router.Use(middlewares.UseCORS())

	backoffice := router.Group("/api/v1/backoffice")
	{
		backoffice.POST("/register", handlers.Register(db))
		backoffice.POST("/login", handlers.Login(db, cfg.JWTSecret))

		backoffice.Use(middlewares.HasAuthorization(cfg.JWTSecret))
		{
			backoffice.GET("/apps", handlers.ListApps(db))
			backoffice.POST("/apps", handlers.CreateApp(db))
			backoffice.DELETE("/apps/:appID", handlers.DeleteApp(db))
			backoffice.GET("/logs", handlers.SearchLogs(db))
		}
	}

	appsGroup := router.Group("/api/v1/apps")
	{
		appsGroup.POST("/:appID/logs", handlers.ReceiveLogs(db))
	}

	router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
