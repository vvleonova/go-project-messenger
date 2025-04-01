package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"go-messenger/internal/api/dto"
)

// создание пользователя
func (s Server) userSignUp(c *gin.Context) {
	// проверка тела запроса
	body := new(dto.UserSignUp)
	if err := c.BindJSON(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// проверка переданного запроса
	if err := body.ValidateUserSignUp(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// сохранение структуры пользователя
	user, err := s.service.UserCreate(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// вход в личный кабинет пользователя
func (s Server) userSignIn(c *gin.Context) {
	// проверка тела запроса
	body := new(dto.UserSignIn)
	if err := c.BindJSON(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// проверка переданного запроса
	if err := body.ValidateUserSignIn(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// получение данных о пользователе
	user, err := s.service.UserGet(body.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// генерация токена
	token, err := s.service.TokenGenerate(user.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// получение данных о пользователе из базы данных
func (s Server) userGet(c *gin.Context) {
	// получение данных о пользователе
	phone := c.Param("phone")
	user, err := s.service.UserGet(phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// проверка наличия пользователя в бд (это лишнее, если есть авторизация?)
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("user with phone %s not found", phone)})
		return
	}

	c.JSON(http.StatusOK, user)
}

// удаление данных о пользователе из базы данных
func (s Server) userDelete(c *gin.Context) {
	// удаление пользователя
	phone := c.Param("phone")
	err := s.service.UserDelete(phone)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("user with phone %s not found", phone)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("user with phone %s deleted", phone)})
}

// обновление данных о пользователе в базе данных
func (s Server) userUpdate(c *gin.Context) {
	// получение номера телефона
	phone := c.Param("phone")

	// проверка тела запроса
	body := new(dto.UserUpdate)
	if err := c.BindJSON(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// проверка переданного запроса
	if err := body.ValidateUserUpdate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// обновление данных о пользователе
	err := s.service.UserUpdate(phone, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("user with phone %s updated", phone)})
}
