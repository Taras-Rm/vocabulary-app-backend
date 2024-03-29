package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"vacabulary/config"
	el "vacabulary/db/elastic"
	"vacabulary/models"
	pdfgenerator "vacabulary/pkg/pdfGenerator"
	"vacabulary/repositories/elastic"

	"github.com/gin-gonic/gin"
)

var (
	errEmptyQueryText = errors.New("empty query text")
)

func (a *App) InjectCollections(gr *gin.Engine) {
	collections := gr.Group("/collection", a.authorizeRequest)

	collections.POST("", a.createCollection)                                  // OK
	collections.GET("/all", a.getAllCollections)                              // OK
	collections.GET(":id", a.idParam("id"), a.getCollection)                  // OK
	collections.PUT(":id", a.idParam("id"), a.updateCollection)               // OK
	collections.DELETE(":id", a.idParam("id"), a.deleteCollection)            // OK
	collections.GET(":id/search", a.idParam("id"), a.searchWordsInCollection) // OK

	collections.POST(":id/generatePdf", a.idParam("id"), a.generatePdfCollection)
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

	elClient := el.NewElasticClient(config.Config.Elastic)
	err = elClient.CreateCollectionAliases(user.Id, collection.Id)
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
		words, _, err := a.wordRepo.GetAll(0, 0, elastic.CollectionWordsOperationCtx{UserId: user.Id, CollectionId: c.Id})
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
	id := ctx.GetUint64("id")
	if id == 0 {
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("can not get id").Error())
		return
	}

	user := a.getContextUser(ctx)

	collection, err := a.collectionRepo.GetById(id)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	words, _, err := a.wordRepo.GetAll(0, 0, elastic.CollectionWordsOperationCtx{UserId: user.Id, CollectionId: uint64(id)})
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

	file, err := pdfgenerator.GenerateCollectionPdf(contents, collection.Name)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s.pdf", collection.Name))
	ctx.Writer.Header().Add("Content-type", "application/pdf")
	ctx.Writer.Write(file)
}

type getCollectionResponse struct {
	Collection *models.Collection `json:"collection"`
}

func (a *App) getCollection(ctx *gin.Context) {
	id := ctx.GetUint64("id")
	if id == 0 {
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("can not get id").Error())
		return
	}

	collection, err := a.collectionRepo.GetById(id)
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
	id := ctx.GetUint64("id")
	if id == 0 {
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("can not get id").Error())
		return
	}

	var input updateCollectionInp
	err := ctx.BindJSON(&input)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	_, err = a.collectionRepo.Update(&models.Collection{
		Id:   id,
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
	id := ctx.GetUint64("id")
	if id == 0 {
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("can not get id").Error())
		return
	}

	err := a.collectionRepo.DeleteById(id)
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

	id := ctx.GetUint64("id")
	if id == 0 {
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("can not get id").Error())
		return
	}

	user := a.getContextUser(ctx)

	words, err := a.wordRepo.Search(*searchSettings, elastic.CollectionWordsOperationCtx{UserId: user.Id, CollectionId: id})
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
