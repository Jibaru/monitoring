package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"

	"monitoring/internal/persistence"
	"monitoring/internal/scripts"
)

// StartOAuth godoc
// @Summary      StartOAuth
// @Description  StartOAuth
// @Accept       json
// @Produce      json
// @Success      307
// @Failure      500    {object}    ErrorResp
// @Router       /api/v1/backoffice/auth/{provider} [get]
func StartOAuth(db *mongo.Database, cfg *oauth2.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := scripts.NewStartOAuthScript(persistence.NewOAuthStateRepo(db), cfg).Exec(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResp{Message: err.Error()})
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, resp.URL)
	}
}
