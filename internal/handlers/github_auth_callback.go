package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"

	"monitoring/config"
	"monitoring/internal/persistence"
	"monitoring/internal/scripts"
)

// GithubAuthCallback godoc
// @Summary      GithubAuthCallback
// @Description  GithubAuthCallback
// @Accept       json
// @Produce      json
// @Success      200    {object}    scripts.GithubAuthResp
// @Failure      400    {object}    ErrorResp
// @Failure      500    {object}    ErrorResp
// @Router       /api/v1/backoffice/auth/github/callback [get]
func GithubAuthCallback(db *mongo.Database, cfg *oauth2.Config, appCfg config.Config) gin.HandlerFunc {
	getGithubUserData := func(token string) (string, string, error) {
		req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
		if err != nil {
			return "", "", err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", "", err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return "", "", fmt.Errorf("error getting user info, estado: %s", resp.Status)
		}

		type GithubUser struct {
			Login string `json:"login"`
			Email string `json:"email"`
		}

		var user GithubUser
		if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
			return "", "", err
		}

		if user.Email == "" {
			req2, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
			if err != nil {
				return "", "", err
			}
			req2.Header.Set("Authorization", "Bearer "+token)
			req2.Header.Set("Accept", "application/json")

			resp2, err := http.DefaultClient.Do(req2)
			if err != nil {
				return "", "", err
			}
			defer resp2.Body.Close()
			if resp2.StatusCode != http.StatusOK {
				return "", "", fmt.Errorf("error getting emails, status: %s", resp2.Status)
			}

			var emails []struct {
				Email   string `json:"email"`
				Primary bool   `json:"primary"`
			}
			if err = json.NewDecoder(resp2.Body).Decode(&emails); err != nil {
				return "", "", err
			}

			for _, e := range emails {
				if e.Primary {
					user.Email = e.Email
					break
				}
			}

			if user.Email == "" && len(emails) > 0 {
				user.Email = emails[0].Email
			}
		}

		return user.Login, user.Email, nil
	}

	return func(c *gin.Context) {
		state := c.Query("code")

		existsState, err := persistence.ExistOAuthStateByState(c, db, state)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResp{Message: err.Error()})
			return
		}

		if !existsState {
			c.JSON(http.StatusBadRequest, ErrorResp{Message: "invalid code/state"})
			return
		}

		token, err := cfg.Exchange(c, state)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResp{Message: err.Error()})
			return
		}

		username, email, err := getGithubUserData(token.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResp{Message: err.Error()})
			return
		}
		fmt.Println(username)

		script := scripts.NewGithubAuthScript(db, []byte(appCfg.JWTSecret))
		resp, err := script.Exec(c, scripts.GithubAuthReq{
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

		url := fmt.Sprintf("%s/login?token=%s&id=%s&email=%s&isVisitor=%s",
			appCfg.WebBaseURI,
			resp.Token,
			resp.User.ID,
			resp.User.Email,
			isVisitor,
		)
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}
