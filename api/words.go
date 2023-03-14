package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"
	"vacabulary/models"
	"vacabulary/repositories/elastic"

	"github.com/gin-gonic/gin"
)

func (a *App) InjectWords(gr *gin.Engine) {
	words := gr.Group("/word", a.authorizeRequest)
	words.POST("", a.createWord)
	words.POST("/bulk", a.createWords)

	words.GET(":id/collection/:collectionId", a.getWord)
	words.DELETE(":id/collection/:collectionId", a.deleteWord)
	words.GET("/collection/:collectionId", a.getAllWords)
	words.PUT(":id/collection/:collectionId", a.updateWord)

	words.POST("/translate", a.translateWord)
}

type createWordInp struct {
	Word         string `json:"word"`
	Translation  string `json:"translation"`
	PartOfSpeech string `json:"partOfSpeech"`
	Scentance    string `json:"scentance"`
	CollectionId uint64 `json:"collectionId"`
}

func (a *App) createWord(ctx *gin.Context) {
	var input createWordInp
	err := ctx.BindJSON(&input)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user := a.getContextUser(ctx)

	// check: such origin already esists in selected collection or not
	word, err := a.wordRepo.Get(input.Word, elastic.CollectionWordsOperationCtx{CollectionId: input.CollectionId, UserId: user.Id})
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if word != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, errors.New("such word already esists").Error())
		return
	}

	// create word in elastic too
	err = a.wordRepo.Create(models.Word{
		Word:         input.Word,
		Translation:  input.Translation,
		PartOfSpeech: input.PartOfSpeech,
		Scentance:    input.Scentance,
		CollectionId: input.CollectionId,
	}, elastic.CollectionWordsOperationCtx{CollectionId: input.CollectionId, UserId: user.Id})
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
	})
}

type getAllWordsResponse struct {
	Words      []models.Word `json:"words"`
	TotalWords uint64        `json:"totalWords"`
}

func (a *App) getAllWords(ctx *gin.Context) {
	collectionIdStr := ctx.Param("collectionId")
	if collectionIdStr == "" {
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("can not get collection id").Error())
		return
	}

	collectionId, err := strconv.Atoi(collectionIdStr)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// get pagination params
	size, page, err := getPaginationParams(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user := a.getContextUser(ctx)

	words, totalWords, err := a.wordRepo.GetAll(uint64(size), uint64(page), elastic.CollectionWordsOperationCtx{CollectionId: uint64(collectionId), UserId: user.Id})
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, getAllWordsResponse{
		Words:      words,
		TotalWords: totalWords,
	})

}

func getPaginationParams(ctx *gin.Context) (uint64, uint64, error) {
	// get pagination params
	sizeStr := ctx.Query("size")
	pageStr := ctx.Query("page")

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return 0, 0, err
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, 0, err
	}

	return uint64(size), uint64(page), nil
}

type getWordResponse struct {
	Word models.Word `json:"word"`
}

func (a *App) getWord(ctx *gin.Context) {

	id := ctx.Param("id")
	if id == "" {
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("can not get id").Error())
		return
	}

	collectionIdStr := ctx.Param("collectionId")
	if id == "" {
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("can not get collection id").Error())
		return
	}

	collectionId, err := strconv.Atoi(collectionIdStr)
	if id == "" {
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("can not convert collection id").Error())
		return
	}

	user := a.getContextUser(ctx)

	word, err := a.wordRepo.GetById(id, elastic.CollectionWordsOperationCtx{CollectionId: uint64(collectionId), UserId: user.Id})
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if word == nil {
		newErrorResponse(ctx, http.StatusInternalServerError, errors.New("word with such origin not found").Error())
		return
	}

	ctx.JSON(http.StatusOK, getWordResponse{
		Word: *word,
	})

}

type updateWordInp struct {
	Id           string `json:"id"`
	Word         string `json:"word"`
	Translation  string `json:"translation"`
	PartOfSpeech string `json:"partOfSpeech"`
	Scentance    string `json:"scentance"`
	CollectionId uint64 `json:"collectionId"`
}

func (a *App) updateWord(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("can not get id").Error())
		return
	}

	var input updateWordInp
	err := ctx.BindJSON(&input)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	collectionIdStr := ctx.Param("collectionId")
	if id == "" {
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("can not get collection id").Error())
		return
	}

	collectionId, err := strconv.Atoi(collectionIdStr)
	if id == "" {
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("can not convert collection id").Error())
		return
	}

	user := a.getContextUser(ctx)

	word, err := a.wordRepo.GetById(id, elastic.CollectionWordsOperationCtx{CollectionId: uint64(collectionId), UserId: user.Id})
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// create word in elastic too
	err = a.wordRepo.Update(models.Word{
		Id:           id,
		Word:         input.Word,
		Translation:  input.Translation,
		PartOfSpeech: input.PartOfSpeech,
		Scentance:    input.Scentance,
		CreatedAt:    time.Now(),
		CollectionId: word.CollectionId,
	}, elastic.CollectionWordsOperationCtx{CollectionId: uint64(collectionId), UserId: user.Id})
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "success update",
	})

}

func (a *App) deleteWord(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("can not get id").Error())
		return
	}

	collectionIdStr := ctx.Param("collectionId")
	if id == "" {
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("can not get collection id").Error())
		return
	}

	collectionId, err := strconv.Atoi(collectionIdStr)
	if id == "" {
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("can not convert collection id").Error())
		return
	}

	user := a.getContextUser(ctx)

	err = a.wordRepo.DeleteById(id, elastic.CollectionWordsOperationCtx{CollectionId: uint64(collectionId), UserId: user.Id})
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
	})

}

type createWordsWordInp struct {
	Word         string `json:"word"`
	Translation  string `json:"translation"`
	PartOfSpeech string `json:"partOfSpeech"`
	Scentance    string `json:"scentance"`
}

type createWordsInp struct {
	Words        []createWordsWordInp `json:"words"`
	CollectionId uint64               `json:"collectionId"`
}

func (a *App) createWords(ctx *gin.Context) {
	var input createWordsInp
	err := ctx.BindJSON(&input)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	wordsWord := []string{}
	for _, w := range input.Words {
		wordsWord = append(wordsWord, w.Word)
	}

	user := a.getContextUser(ctx)

	// check: such words already esists or not
	words, err := a.wordRepo.GetByWords(wordsWord, elastic.CollectionWordsOperationCtx{CollectionId: uint64(input.CollectionId), UserId: user.Id})
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if words != nil || len(words) > 0 {
		newErrorResponse(ctx, http.StatusInternalServerError, errors.New("such word already esists").Error())
		return
	}

	words = []models.Word{}
	for _, w := range input.Words {
		words = append(words, models.Word{
			Word:         w.Word,
			Translation:  w.Translation,
			PartOfSpeech: w.PartOfSpeech,
			Scentance:    w.Scentance,
			CollectionId: input.CollectionId,
		})
	}

	// create words in elastic too
	err = a.wordRepo.BulkCreate(words, elastic.CollectionWordsOperationCtx{CollectionId: uint64(input.CollectionId), UserId: user.Id})
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
	})

}

type translateWordInp struct {
	Word  string `json:"word"`
	Index int64  `json:"index"`
}

func (a *App) translateWord(ctx *gin.Context) {
	var input translateWordInp
	err := ctx.BindJSON(&input)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	langFrom, langTo, err := getTranslationParams(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	translatedWord, err := a.translatorManager.TranslateWord(input.Word, langFrom, langTo)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"word":  translatedWord,
		"index": input.Index,
	})
}

func getTranslationParams(ctx *gin.Context) (string, string, error) {
	// get translation params
	langFrom := ctx.Query("langFrom")
	langTo := ctx.Query("langTo")

	if langFrom == "" || langTo == "" {
		return "", "", errors.New("can not get language params")
	}

	return langFrom, langTo, nil
}
