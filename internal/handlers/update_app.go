package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/persistence"
	"monitoring/internal/scripts"
)

// UpdateApp godoc
// @Summary      UpdateApp
// @Description  UpdateApp
// @Accept       json
// @Produce      json
// @Param        body  body    scripts.UpdateAppReq    true    "Request"
// @Success      201    {object}    scripts.UpdateAppResp
// @Failure      401    {object}    ErrorResp
// @Failure      500    {object}    ErrorResp
// @Router       /api/v1/backoffice/apps/{appID} [patch]
func UpdateApp(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req scripts.UpdateAppReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Message: err.Error()})
			return
		}

		req.ID = c.Param("appID")

		script := scripts.NewUpdateAppScript(persistence.NewAppRepo(db))
		resp, err := script.Exec(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResp{Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
