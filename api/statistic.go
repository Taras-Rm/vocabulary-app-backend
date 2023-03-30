package api

import (
	"errors"
	"fmt"
	"net/http"
	"vacabulary/models"

	"github.com/gin-gonic/gin"
)

func (a *App) InjectStatistic(gr *gin.Engine) {
	statistic := gr.Group("/statistic")

	statistic.GET("/users", a.authorizeRequest, a.superUser, a.getUsers)
	statistic.GET("/collections", a.authorizeRequest, a.superUser, a.getCollections)
	statistic.GET("/words/count", a.authorizeRequest, a.superUser, a.getAllWordsCount)
	statistic.GET("/words/perTime", a.authorizeRequest, a.superUser, a.getCountOfWordsPerTime)

	statistic.GET("/words/search", a.authorizeRequest, a.superUser, a.searchWordsInAllCollections)
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

type searchWordsInAllCollectionsResponse struct {
	Words []CreatorWord `json:"words"`
}

type Creator struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}
type CreatorWord struct {
	Word    models.Word `json:"word"`
	Creator Creator     `json:"creator"`
}

func (a *App) searchWordsInAllCollections(ctx *gin.Context) {
	searchSettings, err := getSearchWordsInCollectionParams(ctx)
	if err != nil {
		if errors.Is(err, errEmptyQueryText) {
			ctx.JSON(http.StatusOK, searchWordsInCollectionResponse{
				Words: []models.Word{},
			})
			return
		}
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("can not get id").Error())
		return
	}

	users, err := a.userRepo.GetAll()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	var userIds []uint64
	for _, u := range users {
		userIds = append(userIds, u.Id)
	}

	words, err := a.wordRepo.SearchOnCollections(*searchSettings, userIds)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	var wordsInfo []CreatorWord
	for _, w := range words {
		wordCreator, err := a.userRepo.GetByCollectionId(w.CollectionId)
		if err != nil {
			fmt.Println(err)
		}

		wordsInfo = append(wordsInfo, CreatorWord{
			Word: w,
			Creator: Creator{
				Name:  wordCreator.Name,
				Email: wordCreator.Email,
			},
		})
	}

	ctx.JSON(http.StatusOK, searchWordsInAllCollectionsResponse{
		Words: wordsInfo,
	})
}
