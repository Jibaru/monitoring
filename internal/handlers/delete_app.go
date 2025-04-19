package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/persistence"
	"monitoring/internal/scripts"
)

// DeleteApp godoc
// @Summary      DeleteApp
// @Description  DeleteApp
// @Accept       json
// @Produce      json
// @Param        body  body    scripts.DeleteAppReq    true    "Request"
// @Success      204
// @Failure      500    {object}    ErrorResp
// @Router       /api/v1/backoffice/apps/{appID} [delete]
func DeleteApp(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		appID := c.Param("appID")
		req := scripts.DeleteAppReq{AppID: appID}
		script := scripts.NewDeleteAppScript(persistence.NewAppRepo(db))
		err := script.Exec(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResp{Message: err.Error()})
			return
		}
		c.JSON(http.StatusNoContent, nil)
	}
}
