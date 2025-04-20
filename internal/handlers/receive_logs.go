package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/persistence"
	"monitoring/internal/scripts"
)

// ReceiveLogs godoc
// @Summary      ReceiveLogs
// @Description  ReceiveLogs
// @Accept       json
// @Produce      json
// @Param        body  body    scripts.ReceiveLogsReq    true    "Request"
// @Success      201    {object}    scripts.ReceiveLogsResp
// @Failure      400    {object}    ErrorResp
// @Failure      401    {object}    ErrorResp
// @Failure      500    {object}    ErrorResp
// @Router       /api/v1/apps/logs [post]
func ReceiveLogs(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req scripts.ReceiveLogsReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Message: err.Error()})
			return
		}
		req.AppKey = c.GetHeader("x-app-key")

		script := scripts.NewReceiveLogsScript(persistence.NewLogRepo(db), persistence.NewAppRepo(db))
		resp, err := script.Exec(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResp{Message: err.Error()})
			return
		}
		c.JSON(http.StatusCreated, resp)
	}
}
