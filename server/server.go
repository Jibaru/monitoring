package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"monitoring/config"
	"monitoring/internal/handlers"
	"monitoring/internal/middlewares"
)

func New(cfg config.Config, db *mongo.Database) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.UseCORS())

	githubOAuthConfig := &oauth2.Config{
		RedirectURL:  fmt.Sprintf("%s/api/v1/backoffice/auth/github/callback", cfg.APIBaseURI),
		ClientID:     cfg.GithubClientID,
		ClientSecret: cfg.GithubClientSecret,
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}

	backoffice := router.Group("/api/v1/backoffice")
	{
		backoffice.POST("/register", handlers.Register(db, cfg))
		backoffice.POST("/login", handlers.Login(db, cfg.JWTSecret))
		backoffice.POST("/visitor-login", handlers.VisitorLogin(db, cfg.JWTSecret))
		backoffice.POST("/users/:userID/validate", handlers.ValidateUser(db))
		backoffice.GET("/auth/github", handlers.GithubAuth(db, githubOAuthConfig))
		backoffice.GET("/auth/github/callback", handlers.GithubAuthCallback(db, githubOAuthConfig, cfg.JWTSecret))

		backoffice.Use(middlewares.HasAuthorization(cfg.JWTSecret))
		{
			backoffice.GET("/apps", handlers.ListApps(db))
			backoffice.POST("/apps", handlers.CreateApp(db))
			backoffice.PATCH("/apps/:appID", handlers.UpdateApp(db))
			backoffice.DELETE("/apps/:appID", handlers.DeleteApp(db))
			backoffice.GET("/logs", handlers.SearchLogs(db))
			backoffice.GET("/dashboard/overview", handlers.GetDashboardOverview(db))
		}
	}

	appsGroup := router.Group("/api/v1/apps")
	{
		appsGroup.POST("/logs", handlers.ReceiveLogs(db))
	}

	router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
