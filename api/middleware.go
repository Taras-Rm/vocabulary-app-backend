package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"vacabulary/models"

	"github.com/gin-gonic/gin"
)

func (a *App) authorizeRequest(ctx *gin.Context) {
	header := strings.Split(ctx.GetHeader("Authorization"), " ")

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

func (a *App) superUser(ctx *gin.Context) {
	user := a.getContextUser(ctx)

	if !user.IsSuper {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("not super user").Error())
		return
	}

	ctx.Next()
}

func (a *App) idParam(param string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		idStr := ctx.Param(param)
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, "id not valid")
			return
		}

		ctx.Set(param, id)

		ctx.Next()
	}
}
