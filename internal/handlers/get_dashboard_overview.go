package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/persistence"
	"monitoring/internal/scripts"
)

// GetDashboardOverview godoc
// @Summary      GetDashboardOverview
// @Description  GetDashboardOverview
// @Accept       json
// @Produce      json
// @Param        body  body    scripts.GetDashboardOverviewReq    true    "Request"
// @Success      200    {object}    scripts.GetDashboardOverviewResp
// @Failure      500    {object}    ErrorResp
// @Router       /api/v1/dashboard/overview [get]
func GetDashboardOverview(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req scripts.GetDashboardOverviewReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Message: err.Error()})
			return
		}
		req.UserID = c.GetString("user_id")

		script := scripts.NewGetDashboardOverviewScript(persistence.NewDashboardRepo(db))
		resp, err := script.Exec(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResp{Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
