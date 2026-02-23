package controller

import (
	"net/http"
	"restfulapi/models"
	"restfulapi/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
	GetAllUsers(ctx *gin.Context)
	GetUserById(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type authController struct {
	authservice services.AuthService
}

func NewAuthController(authService services.AuthService) AuthController {
	return &authController{authservice: authService}
}

func (c *authController) Login(ctx *gin.Context) {
	var req models.UserLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := c.authservice.Login(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "berhasil login",
		"token":   token,
	})
}

func (c *authController) Register(ctx *gin.Context) {
	var req models.UserRegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// serahkan ke service
	user, err := c.authservice.Register(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// krim response
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "register behasil",
		"user":    user,
	})
}

func (c *authController) GetAllUsers(ctx *gin.Context) {
	users, err := c.authservice.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "berhasil",
		"Users":   users,
	})
}

func (c *authController) GetUserById(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// serahkan ke service

	user, err := c.authservice.GetUserById(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "sukes",
		"user":    user,
	})

}

func (c *authController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//baca request dulu dari json body yang dikirim
	var req models.UserRegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// serahkan ke srvice
	user, err := c.authservice.UpdateUser(uint(id), req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "update berhasil",
		"user":    user,
	})
}

func (c *authController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//serahkan ke service

	if err := c.authservice.DeleteUser(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "user berhasil dihapus",
	})
}
