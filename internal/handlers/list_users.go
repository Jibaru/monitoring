package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/persistence"
	"monitoring/internal/scripts"
)

// ListNoRootUsers godoc
// @Summary      ListNoRootUsers
// @Description  ListNoRootUsers
// @Accept       json
// @Produce      json
// @Param        body  body    scripts.ListUsersReq  true    "Request"
// @Param 	     page  query   int                   false   "Page"
// @Param        limit query   int                   false   "Limit"
// @Success      200    {object}    scripts.ListUsersResp
// @Failure      500    {object}    ErrorResp
// @Router       /api/v1/backoffice/users [get]
func ListNoRootUsers(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		sortOrder := c.DefaultQuery("sortOrder", "desc")
		searchTerm := c.DefaultQuery("searchTerm", "")

		script := scripts.NewListUsersScript(persistence.NewUserRepo(db))
		resp, err := script.Exec(c, scripts.ListUsersReq{
			RootUserID: c.GetString("user_id"),
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
