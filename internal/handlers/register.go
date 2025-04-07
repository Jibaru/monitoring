package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/config"
	"monitoring/internal/mail"
	"monitoring/internal/scripts"
)

// Register godoc
// @Summary      Register
// @Description  Register
// @Accept       json
// @Produce      json
// @Param        body  body    scripts.RegisterReq    true    "Request"
// @Success      201    {object}    scripts.RegisterResp
// @Failure      400    {object}    ErrorResp
// @Router       /api/v1/backoffice/register [post]
func Register(db *mongo.Database, cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req scripts.RegisterReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Message: err.Error()})
			return
		}
		mailSender := mail.NewMailSender(cfg.MailFromEmail, cfg.MailAppPassword)
		resp, err := scripts.NewRegisterScript(db, mailSender, cfg.WebBaseURI).Exec(c, req)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Message: err.Error()})
			return
		}
		c.JSON(http.StatusCreated, resp)
	}
}
