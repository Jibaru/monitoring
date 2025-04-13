package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/scripts"
)

// UpdateUserPassword godoc
// @Summary      UpdateUserPassword
// @Description  UpdateUserPassword
// @Accept       json
// @Produce      json
// @Param        body  body    scripts.UpdateUserPasswordReq    true    "Request"
// @Success      204
// @Failure      401    {object}    ErrorResp
// @Failure      500    {object}    ErrorResp
// @Router       /api/v1/backoffice/users/me/password [put]
func UpdateUserPassword(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req scripts.UpdateUserPasswordReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Message: err.Error()})
			return
		}

		req.ID = c.GetString("user_id")

		script := scripts.NewUpdateUserPasswordScript(db)
		err := script.Exec(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResp{Message: err.Error()})
			return
		}
		c.JSON(http.StatusNoContent, nil)
	}
}
