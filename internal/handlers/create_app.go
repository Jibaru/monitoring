package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/scripts"
)

// CreateApp godoc
// @Summary      CreateApp
// @Description  CreateApp
// @Accept       json
// @Produce      json
// @Param        body  body    scripts.CreateAppReq    true    "Request"
// @Success      201    {object}    scripts.CreateAppResp
// @Failure      401    {object}    ErrorResp
// @Failure      500    {object}    ErrorResp
// @Router       /api/v1/backoffice/register [post]
func CreateApp(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req scripts.CreateAppReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Message: err.Error()})
			return
		}
		req.UserID = c.GetString("user_id")

		script := scripts.NewCreateAppScript(db)
		resp, err := script.Exec(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResp{Message: err.Error()})
			return
		}
		c.JSON(http.StatusCreated, resp)
	}
}
