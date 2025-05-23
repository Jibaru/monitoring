package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/config"
	"monitoring/internal/mail"
	"monitoring/internal/persistence"
	"monitoring/internal/scripts"
)

// CreateNoRootUser godoc
// @Summary      CreateNoRootUser
// @Description  CreateNoRootUser
// @Accept       json
// @Produce      json
// @Param        body  body    scripts.CreateUserReq    true    "Request"
// @Success      201    {object}    scripts.CreateUserResp
// @Failure      400    {object}    ErrorResp
// @Router       /api/v1/backoffice/users [post]
func CreateNoRootUser(db *mongo.Database, cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req scripts.CreateUserReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Message: err.Error()})
			return
		}

		req.RootID = c.GetString("user_id")

		mailSender := mail.NewMailSender(cfg.MailFromEmail, cfg.MailAppPassword)
		resp, err := scripts.NewCreateUserScript(persistence.NewUserRepo(db), mailSender, cfg.WebBaseURI).Exec(c, req)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Message: err.Error()})
			return
		}
		c.JSON(http.StatusCreated, resp)
	}
}
