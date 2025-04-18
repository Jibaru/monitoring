package server

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"

	"monitoring/config"
	"monitoring/internal/handlers"
	"monitoring/internal/middlewares"
)

//go:embed static
var staticFiles embed.FS

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

	googleOAuthConfig := &oauth2.Config{
		RedirectURL:  fmt.Sprintf("%s/api/v1/backoffice/auth/google/callback", cfg.APIBaseURI),
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	backoffice := router.Group("/api/v1/backoffice")
	{
		backoffice.POST("/register", handlers.Register(db, cfg))
		backoffice.POST("/login", handlers.Login(db, cfg.JWTSecret))
		backoffice.POST("/visitor-login", handlers.VisitorLogin(db, cfg.JWTSecret))
		backoffice.POST("/users/:userID/validate", handlers.ValidateUser(db))
		backoffice.GET("/auth/github", handlers.OAuth(db, githubOAuthConfig))
		backoffice.GET("/auth/google", handlers.OAuth(db, googleOAuthConfig))
		backoffice.GET("/auth/github/callback", handlers.OAuthCallback(db, handlers.GithubInfoExtractor, githubOAuthConfig, cfg))
		backoffice.GET("/auth/google/callback", handlers.OAuthCallback(db, handlers.GoogleInfoExtractor, googleOAuthConfig, cfg))

		backoffice.Use(middlewares.HasAuthorization(cfg.JWTSecret))
		{
			backoffice.GET("/apps", handlers.ListApps(db))
			backoffice.POST("/apps", handlers.CreateApp(db))
			backoffice.PATCH("/apps/:appID", handlers.UpdateApp(db))
			backoffice.DELETE("/apps/:appID", handlers.DeleteApp(db))
			backoffice.GET("/logs", handlers.SearchLogs(db))
			backoffice.GET("/dashboard/overview", handlers.GetDashboardOverview(db))
			backoffice.GET("/logs/schema", handlers.GetLogsSchema(db))
			backoffice.PATCH("/users/me", handlers.UpdateUser(db))
			backoffice.PUT("/users/me/password", handlers.UpdateUserPassword(db))
		}
	}

	appsGroup := router.Group("/api/v1/apps")
	{
		appsGroup.POST("/logs", handlers.ReceiveLogs(db))
	}

	subFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}

	router.StaticFS("/api/static", http.FS(subFS))
	router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
