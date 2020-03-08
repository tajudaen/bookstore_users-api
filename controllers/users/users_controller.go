package users

import (
	"net/http"

	"github.com/tajud99n/bookstore_users-api/services"
	"github.com/tajud99n/bookstore_users-api/utils/errors"

	"github.com/tajud99n/bookstore_users-api/domain/users"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	result, saveErr := services.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func GetUser(c *gin.Context) {

}
