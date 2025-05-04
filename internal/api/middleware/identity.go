package middleware

import (
	"fmt"
	"go-messenger/internal/api/errs"

	"github.com/gin-gonic/gin"
)

func IdentityRequired(storage storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		// получение номера телефона из запроса
		phoneRequest := c.Param("phone")

		// получение ID пользователя из БД
		u, err := storage.UserGetPhone(phoneRequest)
		if err != nil {
			c.Error(errs.NotFound(fmt.Sprintf("user with phone %s is not found in database", phoneRequest)))
			c.Abort()
			return
		}

		idRequest := u.ID.String()

		// получение id пользователя из токена
		tokenIDRaw, exists := c.Get("id")
		if !exists {
			c.Error(errs.Unauthorized("not authenticated user"))
			c.Abort()
			return
		}
		tokenID := tokenIDRaw.(string)

		// проверка возможности доступа к данным пользователя
		if idRequest != tokenID {
			c.Error(errs.Forbidden("you cannot access another user's data"))
			c.Abort()
			return
		}

		c.Set("id", u.ID)
		c.Next()
	}
}
