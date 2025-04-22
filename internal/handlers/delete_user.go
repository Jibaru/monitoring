package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/persistence"
	"monitoring/internal/scripts"
)

// DeleteUser godoc
// @Summary      DeleteUser
// @Description  DeleteUser
// @Accept       json
// @Produce      json
// @Param        body  body    scripts.DeleteUserReq    true    "Request"
// @Success      204
// @Failure      500    {object}    ErrorResp
// @Router       /api/v1/backoffice/users/{userID} [delete]
func DeleteUser(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userID")
		req := scripts.DeleteUserReq{
			UserID:     userID,
			RootUserID: c.GetString("user_id"),
		}
		script := scripts.NewDeleteUserScript(persistence.NewUserRepo(db))
		err := script.Exec(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResp{Message: err.Error()})
			return
		}
		c.JSON(http.StatusNoContent, nil)
	}
}
