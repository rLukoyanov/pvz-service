package handlers

import (
	"net/http"
	"pvz-service/internal/models"
	"pvz-service/internal/services"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	services *services.Services
}

func NewAuthHandler(services *services.Services) *AuthHandler {
	return &AuthHandler{services: services}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// @Summary Регистрация пользователя
// @Description Создание нового пользователя с email, паролем и ролью
// @Tags auth
// @Accept json
// @Produce json
// @Param request body registerRequest true "User registration data"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req registerRequest
	if err := c.Bind(&req); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid body"})
	}

	if req.Role != "client" && req.Role != "moderator" {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid role"})
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "encryption error"})
	}

	user := models.User{
		Email:    req.Email,
		Password: string(hashed),
		Role:     req.Role,
	}

	logrus.Debug("creating user", user)
	err = h.services.UserService.CreateUser(c.Request().Context(), user)
	if err != nil {
		logrus.Error(err)
		// TODO - добавить разные варианты ошибок
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "can not create user"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"email": user.Email, "role": user.Role})
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Summary Авторизация пользователя
// @Description Аутентификация пользователя и получение JWT токена
// @Tags auth
// @Accept json
// @Produce json
// @Param request body loginRequest true "User login data"
// @Success 200 {object} string "Token"
// @Failure 401 {object} map[string]string
// @Router /login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid body"})
	}

	user, err := h.services.UserService.GetUserByEmail(c.Request().Context(), req.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "invalid credentials"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"role":  user.Role,
	})

	signed, err := token.SignedString([]byte(h.services.Cfg.SECRET))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "token error"})
	}

	return c.JSON(http.StatusOK, echo.Map{"token": signed})
}
