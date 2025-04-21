package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/persistence"
	"monitoring/internal/scripts"
)

// GetLogsSchema godoc
// @Summary      GetLogsSchema
// @Description  GetLogsSchema
// @Accept       json
// @Produce      json
// @Param        body  body    scripts.GetLogsSchemaReq    true    "Request"
// @Success      200    {object}    scripts.GetLogsSchemaResp
// @Failure      500    {object}    ErrorResp
// @Router       /api/v1/logs/schema [get]
func GetLogsSchema(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req scripts.GetLogsSchemaReq
		req.UserID = c.GetString("user_id")

		script := scripts.NewGetLogsSchemaScript(persistence.NewLogSchemaRepo(db))
		resp, err := script.Exec(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResp{Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
