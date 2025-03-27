package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/scripts"
)

// Login godoc
// @Summary      Login
// @Description  Login
// @Accept       json
// @Produce      json
// @Param        body  body    scripts.LoginReq    true    "Request"
// @Success      200    {object}    scripts.LoginResp
// @Failure      401    {object}    ErrorResp
// @Failure      400    {object}    ErrorResp
// @Router       /api/v1/backoffice/login [post]
func Login(db *mongo.Database, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req scripts.LoginReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Message: err.Error()})
			return
		}
		resp, err := scripts.NewLoginScript(db, jwtSecret).Exec(c, req)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResp{Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
