package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *App) InjectStatistic(gr *gin.Engine) {
	statistic := gr.Group("/statistic")

	statistic.GET("/users", a.authorizeRequest, a.superUser, a.getUsers)
	statistic.GET("/collections", a.authorizeRequest, a.superUser, a.getCollections)
	statistic.GET("/words/count", a.authorizeRequest, a.superUser, a.getAllWordsCount)
	statistic.GET("/words/perTime", a.authorizeRequest, a.superUser, a.getCountOfWordsPerTime)
}

func (a *App) getUsers(ctx *gin.Context) {
	users, err := a.userRepo.GetAll()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"users":   users,
	})
}

func (a *App) getCollections(ctx *gin.Context) {
	collections, err := a.collectionRepo.GetAll()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message":     "success",
		"collections": collections,
	})
}

func (a *App) getAllWordsCount(ctx *gin.Context) {
	users, err := a.userRepo.GetAll()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	var userIds []uint64
	for _, u := range users {
		userIds = append(userIds, u.Id)
	}

	countOfWords, err := a.wordRepo.GetAllWordsCount(userIds)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"count":   countOfWords,
	})
}

func (a *App) getCountOfWordsPerTime(ctx *gin.Context) {
	users, err := a.userRepo.GetAll()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	perTime := ctx.Query("time")

	var userIds []uint64
	for _, u := range users {
		userIds = append(userIds, u.Id)
	}

	countOfWordsPerTime, err := a.wordRepo.GetCountOfWordsPerTime(userIds, perTime)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message":   "success",
		"statistic": countOfWordsPerTime,
	})
}
