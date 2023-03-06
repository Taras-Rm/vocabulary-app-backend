package main

import (
	"net/http"
	"vacabulary/api"
	"vacabulary/config"
	"vacabulary/pkg/s3"
	"vacabulary/pkg/token"
	"vacabulary/pkg/translator"

	elrepositories "vacabulary/repositories/elastic"
	postgresRepo "vacabulary/repositories/postgres"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/olivere/elastic/v7"

	"vacabulary/server"
)

func main() {
	cfg := config.Config

	// elClient := myel.NewElasticClient(config.Elastic.Host, config.Elastic.Port)

	// pgClient := postgres.NewPostgres(config.Postgres)

	// createIndices
	// err := elClient.CreateVocabularyIndices()
	// if err != nil {
	// 	panic("can't create elastic indices")
	// }

	tokenService := token.NewTokenService(cfg.Salt)
	translatorManager := translator.NewTranslatorManager(cfg.AWS)
	s3Manager := s3.NewS3Manager(cfg.AWS)

	// elWordsRepo := elrepositories.NewWordsRepo(elClient.Client)
	elWordsRepo := elrepositories.NewWordsRepo(&elastic.Client{})
	usersRepo := postgresRepo.NewUsersRepo(&pg.DB{})
	collectionsRepo := postgresRepo.NewCollectionsRepo(&pg.DB{})

	router := server.NewServer()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	app := api.NewApp(usersRepo, collectionsRepo, elWordsRepo, *tokenService, translatorManager, s3Manager)

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "pong")
	})

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "hello "+cfg.Salt)
	})

	app.AttachEndpoints(router)

	router.Run()
}
