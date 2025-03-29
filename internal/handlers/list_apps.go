package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/scripts"
)

// ListApps godoc
// @Summary      ListApps
// @Description  ListApps
// @Accept       json
// @Produce      json
// @Param        body  body    scripts.ListAppsReq   true    "Request"
// @Param 	     page  query   int                   false   "Page"
// @Param        limit query   int                   false   "Limit"
// @Success      200    {object}    scripts.ListAppsResp
// @Failure      500    {object}    ErrorResp
// @Router       /api/v1/backoffice/apps [get]
func ListApps(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		sortOrder := c.DefaultQuery("sortOrder", "desc")
		searchTerm := c.DefaultQuery("searchTerm", "")

		script := scripts.NewListAppsScript(db)
		resp, err := script.Exec(c, scripts.ListAppsReq{
			UserID:     c.GetString("user_id"),
			Page:       page,
			Limit:      limit,
			SortOrder:  sortOrder,
			SearchTerm: searchTerm,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResp{Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
