package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/scripts"
)

// VisitorLogin godoc
// @Summary      VisitorLogin
// @Description  VisitorLogin
// @Accept       json
// @Produce      json
// @Success      201    {object}    scripts.VisitorLoginResp
// @Failure      401    {object}    ErrorResp
// @Router       /api/v1/backoffice/visitor-login [post]
func VisitorLogin(db *mongo.Database, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := scripts.NewVisitorLoginScript(db, jwtSecret).Exec(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResp{Message: err.Error()})
			return
		}
		c.JSON(http.StatusCreated, resp)
	}
}
