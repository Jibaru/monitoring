package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"

	"monitoring/config"
	"monitoring/internal/domain/services"
	"monitoring/internal/persistence"
	"monitoring/internal/scripts"
)

// OAuthCallback godoc
// @Summary      OAuthCallback
// @Description  OAuthCallback
// @Accept       json
// @Produce      json
// @Success      307
// @Failure      400    {object}    ErrorResp
// @Router       /api/v1/backoffice/auth/{provider}/callback [get]
func OAuthCallback(
	db *mongo.Database,
	infoExtractor services.OAuthInfoExtractor,
	oauthCfg *oauth2.Config,
	cfg config.Config,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := scripts.NewFinishOAuthScript(
			persistence.NewUserRepo(db),
			persistence.NewOAuthStateRepo(db),
			oauthCfg,
			infoExtractor,
			cfg,
		).Exec(c, scripts.FinishOAuthReq{
			Code:  c.Query("code"),
			State: c.Query("state"),
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Message: err.Error()})
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, resp.URL)
	}
}
