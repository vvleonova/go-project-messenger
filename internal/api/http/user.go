package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"go-messenger/internal/api/dto"
	"go-messenger/internal/api/errs"
)

var err error

// создание пользователя
func (s Server) userSignUp(c *gin.Context) {
	// проверка тела запроса
	body := new(dto.UserSignUp)
	if err = c.BindJSON(body); err != nil {
		c.Error(errs.BadRequest(err.Error()))
		return
	}

	// проверка переданного запроса
	if err = body.ValidateUserSignUp(); err != nil {
		c.Error(errs.BadRequest(err.Error()))
		return
	}

	// сохранение структуры пользователя
	user, err := s.service.UserCreate(body)
	if err != nil {
		c.Error(errs.BadRequest(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, user)
}

// вход в личный кабинет пользователя
func (s Server) userSignIn(c *gin.Context) {
	// проверка тела запроса
	body := new(dto.UserSignIn)
	if err = c.BindJSON(body); err != nil {
		c.Error(errs.BadRequest(err.Error()))
		return
	}

	// проверка переданного запроса
	if err = body.ValidateUserSignIn(); err != nil {
		c.Error(errs.BadRequest(err.Error()))
		return
	}

	// проверка пользователя
	user, httpError := s.service.UserCheck(body)
	if httpError != nil {
		c.Error(httpError)
		return
	}

	// генерация токена
	accessToken, httpError := s.service.UserGenerateToken(user.Phone, user.ID)
	if httpError != nil {
		c.Error(httpError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"access token": accessToken})
}

// обновление JWT-token
func (s Server) userRefreshToken(c *gin.Context) {
	// проверка тела запроса
	body := new(dto.RefreshTokenUpdate)
	if err = c.BindJSON(body); err != nil {
		c.Error(errs.BadRequest(err.Error()))
		return
	}

	// генерация нового JWT-token
	accessToken, httpError := s.service.RefreshTokenUpdate(body.RefreshToken)
	if httpError != nil {
		c.Error(httpError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"access token": accessToken})
}

// получение данных о пользователе из базы данных
func (s Server) userGetPhone(c *gin.Context) {
	// получение номера телефона из запроса
	phoneRequest := c.Param("phone")

	// получение номера телефона из токена
	tokenPhoneRaw, exists := c.Get("phone")
	if !exists {
		c.Error(errs.Unauthorized("no authenticated user"))
		return
	}
	tokenPhone := tokenPhoneRaw.(string)

	// проверка возможности доступа к данным пользователя
	if phoneRequest != tokenPhone {
		c.Error(errs.Forbidden("you cannot access another user's data"))
		return
	}

	// получение данных о пользователе
	user, err := s.service.UserGetPhone(phoneRequest)
	if err != nil {
		c.Error(errs.InternalServerError(err.Error()))
		return
	}

	// проверка наличия пользователя в бд (это лишнее, если есть авторизация?)
	if user == nil {
		c.Error(errs.NotFound(fmt.Sprintf("user with phone %s not found", phoneRequest)))
		return
	}

	c.JSON(http.StatusOK, user)
}

// удаление данных о пользователе из базы данных
func (s Server) userDelete(c *gin.Context) {
	// получение номера телефона из запроса
	phoneRequest := c.Param("phone")

	// получение номера телефона из токена
	tokenPhoneRaw, exists := c.Get("phone")
	if !exists {
		c.Error(errs.Unauthorized("no authenticated user"))
		return
	}
	tokenPhone := tokenPhoneRaw.(string)

	// проверка возможности доступа к данным пользователя
	if phoneRequest != tokenPhone {
		c.Error(errs.Forbidden("you cannot access another user's data"))
		return
	}

	// удаление пользователя
	err = s.service.UserDelete(phoneRequest)
	if err != nil {
		c.Error(errs.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("user with phone %s deleted", phoneRequest)})
}

// обновление данных о пользователе в базе данных
func (s Server) userUpdate(c *gin.Context) {
	// получение номера телефона из запроса
	phoneRequest := c.Param("phone")

	// получение номера телефона из токена
	tokenPhoneRaw, exists := c.Get("phone")
	if !exists {
		c.Error(errs.Unauthorized("no authenticated user"))
		return
	}
	tokenPhone := tokenPhoneRaw.(string)

	// проверка возможности доступа к данным пользователя
	if phoneRequest != tokenPhone {
		c.Error(errs.Forbidden("you cannot access another user's data"))
		return
	}

	// проверка тела запроса
	body := new(dto.UserUpdate)
	if err = c.BindJSON(body); err != nil {
		c.Error(errs.BadRequest(err.Error()))
		return
	}

	// проверка переданного запроса
	if err = body.ValidateUserUpdate(); err != nil {
		c.Error(errs.BadRequest(err.Error()))
		return
	}

	// обновление данных о пользователе
	err = s.service.UserUpdate(phoneRequest, body)
	if err != nil {
		c.Error(errs.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("user with phone %s updated", phoneRequest)})
}
