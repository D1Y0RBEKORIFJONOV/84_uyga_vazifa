package handler

import (
	_ "cors/internal/app/docs"
	userentity "cors/internal/entity/user"
	userusecase "cors/internal/usecase/user"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type User struct {
	user userusecase.User
}

func NewUser(user *userusecase.User) *User {
	return &User{
		user: *user,
	}
}

// @title Artisan Connect
// @version 1.0
// @description This is a sample server for a restaurant reservation system.
// @host localhost:7777
// @BasePath        /
// @schemes         http
// @securityDefinitions.apiKey ApiKeyAuth
// @in              header
// @name            Authorization

// Register godoc
// @Summary Register
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param body body userentity.CreateUser true "User registration information"
// @Success 201 {object} userentity.Status
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /user/register [post]
func (u *User) Register(c *gin.Context) {
	log.Println("Starting")
	var req *userentity.CreateUser
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Fatal(err, "HANDLERRRR sukaaaaaaaaaaaaaaaaa")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("Req", req)
	message, err := u.user.CreateUser(c.Request.Context(), req)
	if err != nil {
		log.Fatal(err, "HANDLERRRR sukaaaaaaaaaaaaaaaaa")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println(message, "End")
	c.JSON(http.StatusOK, gin.H{"message": message})
}

// VerifyUser godoc
// @Summary VerifyUser
// @Description VerifyUser a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param body body userentity.VerifyRequest true "User registration information"
// @Success 201 {object} userentity.Status
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /user/verify [post]
func (u *User) VerifyUser(c *gin.Context) {
	var req *userentity.VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	message, err := u.user.VeryFyUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": message})
}

// Login godoc
// @Summary Login
// @Description Login  user
// @Tags auth
// @Accept json
// @Produce json
// @Param body body userentity.LoginRequest true "User registration information"
// @Success 201 {object} userentity.Token
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /user/login [post]
func (u *User) Login(c *gin.Context) {
	var req *userentity.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := u.user.LoginUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, token)
}
