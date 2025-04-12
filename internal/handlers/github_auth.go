package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"

	"monitoring/internal/persistence"
)

// GithubAuth godoc
// @Summary      GithubAuth
// @Description  GithubAuth
// @Accept       json
// @Produce      json
// @Success      307
// @Failure      500    {object}    ErrorResp
// @Router       /api/v1/backoffice/auth/github [get]
func GithubAuth(db *mongo.Database, cfg *oauth2.Config) gin.HandlerFunc {
	generateState := func(length int) (string, error) {
		bytes := make([]byte, length)
		_, err := rand.Read(bytes)
		if err != nil {
			return "", err
		}
		return hex.EncodeToString(bytes), nil
	}

	return func(c *gin.Context) {
		state, err := generateState(16)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResp{Message: err.Error()})
			return
		}

		err = persistence.SaveOAuthState(c, db, persistence.OAuthState{
			ID:    primitive.NewObjectID(),
			State: state,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResp{Message: err.Error()})
			return
		}

		url := cfg.AuthCodeURL(state)
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}
