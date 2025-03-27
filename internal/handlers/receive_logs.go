package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

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
// @Router       /api/v1/apps/{appID}/logs [post]
func ReceiveLogs(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		appID := c.Param("appID")
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, ErrorResp{Message: "API Key requerida"})
			return
		}
		apiKey := strings.TrimPrefix(authHeader, "Bearer ")

		var req scripts.ReceiveLogsReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Message: err.Error()})
			return
		}
		req.AppID = appID
		req.APIKey = apiKey

		script := scripts.NewReceiveLogsScript(db)
		resp, err := script.Exec(context.Background(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResp{Message: err.Error()})
			return
		}
		c.JSON(http.StatusCreated, resp)
	}
}
