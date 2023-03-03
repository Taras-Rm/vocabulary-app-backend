package api

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"vacabulary/models"
	pdfgenerator "vacabulary/pkg/pdfGenerator"

	"github.com/gin-gonic/gin"
)

var (
	errEmptyQueryText = errors.New("empty query text")
)

func (a *App) InjectCollections(gr *gin.Engine) {
	collections := gr.Group("/collection", a.authorizeRequest)

	collections.POST("", a.createCollection)
	collections.GET("/all", a.getAllCollections)
	collections.GET(":id", a.getCollection)
	collections.PUT(":id", a.updateCollection)
	collections.DELETE(":id", a.deleteCollection)
	collections.GET(":id/search", a.searchWordsInCollection)

	collections.POST(":id/generatePdf", a.generatePdfCollection)
}

type createCollectionInp struct {
	Name     string `json:"name"`
	OwnerId  string `json:"ownerId"`
	LangFrom string `json:"langFrom"`
	LangTo   string `json:"langTo"`
}

func (a *App) createCollection(ctx *gin.Context) {
	var input createCollectionInp
	err := ctx.BindJSON(&input)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	collection, err := a.collectionRepo.GetByName(input.Name)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// check collection name
	if collection != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, errors.New("collection with such name already exists").Error())
		return
	}

	// check languages
	if input.LangFrom == "" || input.LangTo == "" {
		newErrorResponse(ctx, http.StatusInternalServerError, errors.New("lang from and lang to can't be empty").Error())
		return
	}

	// check languages
	if input.LangFrom == input.LangTo {
		newErrorResponse(ctx, http.StatusInternalServerError, errors.New("lang from and lang to can't be same").Error())
		return
	}

	user := a.getContextUser(ctx)

	// create collection
	collection, err = a.collectionRepo.Create(models.Collection{
		Name:      input.Name,
		OwnerId:   user.Id,
		LangFrom:  input.LangFrom,
		LangTo:    input.LangTo,
		CreatedAt: time.Now(),
	})
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message":    "success",
		"collection": collection,
	})
}

type getAllCollectionsResponse struct {
	Collections []models.Collection `json:"collections"`
}

func (a *App) getAllCollections(ctx *gin.Context) {
	user := a.getContextUser(ctx)

	collections, err := a.collectionRepo.GetByOwnerId(user.Id)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	var collectionsWithWords []models.Collection
	for _, c := range collections {
		words, _, err := a.wordRepo.GetAll(c.Id, 0, 0)
		if err != nil {
			continue
		}

		if len(words) == 0 {
			c.Words = []models.Word{}
		} else {
			c.Words = words
		}

		collectionsWithWords = append(collectionsWithWords, c)
	}

	ctx.JSON(http.StatusOK, getAllCollectionsResponse{
		Collections: collectionsWithWords,
	})
}

func (a *App) generatePdfCollection(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("can not get id").Error())
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user := a.getContextUser(ctx)

	collection, err := a.collectionRepo.GetById(uint64(id))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	words, _, err := a.wordRepo.GetAll(uint64(id), 0, 0)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	contents := [][]string{}
	for _, w := range words {
		row := []string{}
		row = append(row, w.Word)
		row = append(row, w.Translation)

		contents = append(contents, row)
	}

	pathToFile := "api/" + collection.Name + ".pdf"

	file, err := pdfgenerator.GenerateCollectionPdf(contents, collection.Name, pathToFile)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	filePdf := bytes.NewReader(file)

	key := fmt.Sprintf("%d/%s.pdf", user.Id, collection.Name)
	link, err := a.s3Manager.StorePdfFile(key, filePdf)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	collection.PdfFileUrl = link

	_, err = a.collectionRepo.Update(collection)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
	})
}

type getCollectionResponse struct {
	Collection *models.Collection `json:"collection"`
}

func (a *App) getCollection(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("can not get id").Error())
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	collection, err := a.collectionRepo.GetById(uint64(id))
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, getCollectionResponse{
		Collection: collection,
	})
}

type updateCollectionInp struct {
	Name string `json:"name"`
}

func (a *App) updateCollection(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("can not get id").Error())
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var input updateCollectionInp
	err = ctx.BindJSON(&input)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	_, err = a.collectionRepo.Update(&models.Collection{
		Id:   uint64(id),
		Name: input.Name,
	})
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "success update",
	})
}

func (a *App) deleteCollection(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("can not get id").Error())
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = a.collectionRepo.DeleteById(uint64(id))
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "success delete",
	})
}

type searchWordsInCollectionResponse struct {
	Words []models.Word `json:"words"`
}

func (a *App) searchWordsInCollection(ctx *gin.Context) {
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

	idStr := ctx.Param("id")
	if idStr == "" {
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("can not get id").Error())
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	words, err := a.wordRepo.Search(uint64(id), *searchSettings)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, searchWordsInCollectionResponse{
		Words: words,
	})
}

func getSearchWordsInCollectionParams(ctx *gin.Context) (*models.SearchSettings, error) {
	searchSettings := models.SearchSettings{}

	// get search words params
	searchBy := ctx.Query("searchBy")
	if searchBy == "" {
		return nil, errors.New("can not get search params")
	}
	searchSettings.SearchBy = searchBy

	partsOfSpeechStr := ctx.Query("partsOfSpeech")
	if partsOfSpeechStr != "" {
		partsOfSppech := strings.Split(partsOfSpeechStr, ",")
		searchSettings.PartsOfSpeech = partsOfSppech
	}

	text := ctx.Query("text")
	if text == "" {
		return nil, errEmptyQueryText
	}
	searchSettings.TextForSearch = text

	return &searchSettings, nil
}
