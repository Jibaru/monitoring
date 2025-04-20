package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/persistence"
	"monitoring/internal/scripts"
)

// SearchLogs godoc
// @Summary      SearchLogs
// @Description  SearchLogs
// @Accept       json
// @Produce      json
// @Param        body  body    scripts.SearchLogsReq    true    "Request"
// @Success      200    {object}    scripts.SearchLogsResp
// @Failure      401    {object}    ErrorResp
// @Failure      500    {object}    ErrorResp
// @Router       /api/v1/backoffice/logs [get]
func SearchLogs(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req scripts.SearchLogsReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Message: err.Error()})
			return
		}
		req.UserID = c.GetString("user_id")

		script := scripts.NewSearchLogsScript(persistence.NewAppRepo(db), persistence.NewLogRepo(db))
		resp, err := script.Exec(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResp{Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
