package api

import (
	"errors"
	"net/http"
	"strings"
	"vacabulary/models"

	"github.com/gin-gonic/gin"
)

func (a *App) authorizeRequest(ctx *gin.Context) {
	h := ctx.GetHeader("Authorization")

	if h == "" {
		h = ctx.GetHeader("authorization")
	}

	header := strings.Split(h, " ")

	if len(header) < 2 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("invalid auth header").Error())
		return
	}

	token := header[1]
	userId, err := a.tokenService.ParseToken(token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		return
	}

	user, err := a.userRepo.GetById(userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		return
	}

	if user == nil {
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("can not find user").Error())
			return
		}
	}

	ctx.Set("userId", userId)

	ctx.Next()
}

func (a *App) getContextUser(ctx *gin.Context) *models.User {
	userId := ctx.GetUint64("userId")

	if userId == 0 {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.New("empty user context").Error())
		return nil
	}

	user, err := a.userRepo.GetById(userId)
	if userId == 0 {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return nil
	}

	return user
}
