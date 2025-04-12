package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"

	"monitoring/config"
	"monitoring/internal/persistence"
	"monitoring/internal/scripts"
)

// OAuthCallback godoc
// @Summary      OAuthCallback
// @Description  OAuthCallback
// @Accept       json
// @Produce      json
// @Success      200    {object}    scripts.GithubAuthResp
// @Failure      400    {object}    ErrorResp
// @Failure      500    {object}    ErrorResp
// @Router       /api/v1/backoffice/auth/{provider}/callback [get]
func OAuthCallback(
	db *mongo.Database,
	infoExtractor OAuthInfoExtractor,
	oauthCfg *oauth2.Config,
	cfg config.Config,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Query("code")
		state := c.Query("state")

		existsState, err := persistence.ExistOAuthStateByState(c, db, state)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResp{Message: err.Error()})
			return
		}

		if !existsState {
			c.JSON(http.StatusBadRequest, ErrorResp{Message: "invalid state"})
			return
		}

		token, err := oauthCfg.Exchange(c, code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResp{Message: err.Error()})
			return
		}

		username, email, err := infoExtractor(token.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResp{Message: err.Error()})
			return
		}
		fmt.Println(username)

		script := scripts.NewOAuthScript(db, []byte(cfg.JWTSecret))
		resp, err := script.Exec(c, scripts.OAuthReq{
			Username: username,
			Email:    email,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResp{Message: err.Error()})
			return
		}

		isVisitor := "false"
		if resp.User.IsVisitor {
			isVisitor = "true"
		}

		url := fmt.Sprintf("%s/login?token=%s&id=%s&email=%s&username=%s&isVisitor=%s",
			cfg.WebBaseURI,
			resp.Token,
			resp.User.ID,
			resp.User.Email,
			resp.User.Username,
			isVisitor,
		)
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}
