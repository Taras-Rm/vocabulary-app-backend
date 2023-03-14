package api

import (
	"errors"
	"net/http"
	"time"
	"vacabulary/config"
	"vacabulary/db/elastic"
	"vacabulary/models"

	"github.com/gin-gonic/gin"
)

const (
	expirationTime = 24 * time.Hour
)

func (a *App) InjectUsers(gr *gin.Engine) {
	words := gr.Group("/user")

	words.POST("", a.createUser)
	words.GET("/me", a.authorizeRequest, a.getMe)

	words.POST("/login", a.loginUser)
}

type createUserInp struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *App) createUser(ctx *gin.Context) {
	var input createUserInp
	err := ctx.BindJSON(&input)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, err := a.userRepo.GetByEmail(input.Email)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// check user email
	if user != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, errors.New("user with such email already exists").Error())
		return
	}

	// create user
	user, err = a.userRepo.Create(models.User{
		Name:      input.Name,
		Email:     input.Email,
		Password:  input.Password,
		CreatedAt: time.Now(),
	})
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	elClient := elastic.NewElasticClient(config.Config.Elastic)
	err = elClient.CreateUserWordsIndices(user.Id)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"user":    user,
	})

}

type loginUserInp struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *App) loginUser(ctx *gin.Context) {
	var input loginUserInp
	err := ctx.BindJSON(&input)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, err := a.userRepo.GetByEmail(input.Email)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// get user by email
	if user == nil {
		newErrorResponse(ctx, http.StatusInternalServerError, errors.New("user with such email not founded").Error())
		return
	}

	if user.Password != input.Password {
		newErrorResponse(ctx, http.StatusInternalServerError, errors.New("uncorrect credentials").Error())
		return
	}

	token, err := a.tokenService.GenerateToken(expirationTime, user.Id)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})

}

func (a *App) getMe(ctx *gin.Context) {
	user := a.getContextUser(ctx)

	if user == nil {
		newErrorResponse(ctx, http.StatusInternalServerError, errors.New("cant get user").Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"user": user,
	})
}
