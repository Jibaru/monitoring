package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/persistence"
	"monitoring/internal/scripts"
)

// ValidateUser godoc
// @Summary      ValidateUser
// @Description  ValidateUser
// @Accept       json
// @Produce      json
// @Param        body  body    scripts.ValidateUserReq    true    "Request"
// @Success      201    {object}    scripts.ValidateUserResp
// @Failure      400    {object}    ErrorResp
// @Router       /api/v1/backoffice/users/:userID/validate [post]
func ValidateUser(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req scripts.ValidateUserReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Message: err.Error()})
			return
		}

		req.UserID = c.Param("userID")

		resp, err := scripts.NewValidateUserScript(persistence.NewUserRepo(db)).Exec(c, req)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
