package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"go-messenger/internal/api/dto"
	"go-messenger/internal/api/errs"
	"go-messenger/internal/models"
)

// создание пользователя
func (s *Server) userSignUp(c *gin.Context) {
	// проверка тела запроса
	body := &dto.UserSignUp{}
	if err := c.BindJSON(body); err != nil {
		c.Error(errs.BadRequest(err.Error()))
		return
	}

	// проверка переданного запроса
	if err := body.ValidateUserSignUp(); err != nil {
		c.Error(errs.BadRequest(err.Error()))
		return
	}

	// сохранение структуры пользователя
	user, err := s.service.UserCreate(body)
	if err != nil {
		c.Error(errs.BadRequest(err.Error()))
		return
	}

	userResponse := newUserResponse(user)
	c.JSON(http.StatusCreated, userResponse)
}

// вход в личный кабинет пользователя
func (s *Server) userSignIn(c *gin.Context) {
	// проверка тела запроса
	body := &dto.UserSignIn{}
	if err := c.BindJSON(body); err != nil {
		c.Error(errs.BadRequest(err.Error()))
		return
	}

	// проверка переданного запроса
	if err := body.ValidateUserSignIn(); err != nil {
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
	accessToken, httpError := s.service.UserGenerateToken(user.ID)
	if httpError != nil {
		c.Error(httpError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"access token": accessToken})
}

// обновление JWT-token
func (s *Server) userRefreshToken(c *gin.Context) {
	// проверка тела запроса
	body := &dto.RefreshTokenUpdate{}
	if err := c.BindJSON(body); err != nil {
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
func (s *Server) userGetPhone(c *gin.Context) {
	// получение номера телефона из запроса
	phoneRequest := c.Param("phone")

	// получение данных о пользователе
	user, err := s.service.UserGetPhone(phoneRequest)
	if err != nil {
		c.Error(errs.InternalServerError(err.Error()))
		return
	}

	// проверка наличия пользователя в бд
	if user == nil {
		c.Error(errs.NotFound(fmt.Sprintf("user with phone %s not found", phoneRequest)))
		return
	}

	userResponse := newUserResponse(user)
	c.JSON(http.StatusOK, userResponse)
}

// удаление данных о пользователе из базы данных
func (s *Server) userDelete(c *gin.Context) {
	// получение номера телефона из запроса
	phoneRequest := c.Param("phone")

	// удаление пользователя
	err := s.service.UserDelete(phoneRequest)
	if err != nil {
		c.Error(errs.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("user with phone %s deleted", phoneRequest)})
}

// обновление данных о пользователе в базе данных
func (s *Server) userUpdate(c *gin.Context) {
	// получение номера телефона из запроса
	phoneRequest := c.Param("phone")

	// проверка тела запроса
	body := &dto.UserUpdate{}
	if err := c.BindJSON(body); err != nil {
		c.Error(errs.BadRequest(err.Error()))
		return
	}

	// проверка переданного запроса
	if err := body.ValidateUserUpdate(); err != nil {
		c.Error(errs.BadRequest(err.Error()))
		return
	}

	// обновление данных о пользователе
	err := s.service.UserUpdate(phoneRequest, body)
	if err != nil {
		c.Error(errs.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("user with phone %s updated", phoneRequest)})
}

// logout пользователя из системы
func (s *Server) userLogout(c *gin.Context) {
	// получение номера телефона из запроса
	phoneRequest := c.Param("phone")

	// получение id пользователя из базы данных
	user, err := s.service.UserGetPhone(phoneRequest)
	if err != nil {
		c.Error(errs.InternalServerError(err.Error()))
		return
	}

	// revoke refresh JWT-token
	httpError := s.service.RefreshTokenRevoke(user.ID)
	if httpError != nil {
		c.Error(httpError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("user with phone %s was logget out", phoneRequest)})
}

// создание структуры для вывода результата
func newUserResponse(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Login:     user.Login,
		Phone:     user.Phone,
		BirthDate: user.BirthDate.Format("2006-01-02"),
		CreatedOn: user.CreatedOn.Format("2006-01-02 15:04:05"),
	}
}
