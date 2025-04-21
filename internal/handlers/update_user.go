package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/persistence"
	"monitoring/internal/scripts"
)

// UpdateUser godoc
// @Summary      UpdateUser
// @Description  UpdateUser
// @Accept       json
// @Produce      json
// @Param        body  body    scripts.UpdateUserReq    true    "Request"
// @Success      201    {object}    scripts.UpdateUserResp
// @Failure      401    {object}    ErrorResp
// @Failure      500    {object}    ErrorResp
// @Router       /api/v1/backoffice/users/me [patch]
func UpdateUser(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req scripts.UpdateUserReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Message: err.Error()})
			return
		}

		req.ID = c.GetString("user_id")

		script := scripts.NewUpdateUserScript(persistence.NewUserRepo(db))
		resp, err := script.Exec(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResp{Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
