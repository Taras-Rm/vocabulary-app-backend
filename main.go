package main

import (
	"fmt"
	"net/http"
	"vacabulary/api"
	"vacabulary/config"
	"vacabulary/db/elastic"
	"vacabulary/db/postgres"

	"vacabulary/pkg/hasher"
	"vacabulary/pkg/s3"
	"vacabulary/pkg/token"
	"vacabulary/pkg/translator"

	elrepositories "vacabulary/repositories/elastic"
	postgresRepo "vacabulary/repositories/postgres"

	"github.com/gin-gonic/gin"

	"vacabulary/server"
)

func main() {
	cfg := config.Config

	elClient := elastic.NewElasticClient(cfg.Elastic)
	fmt.Println(cfg.Postgres.Host)

	postgres.MigrateDB()
	pgClient := postgres.NewPostgres(cfg.Postgres)

	tokenService := token.NewTokenService(cfg.Salt)
	translatorManager := translator.NewTranslatorManager(cfg.AWS)
	s3Manager := s3.NewS3Manager(cfg.AWS)
	hasher := hasher.NewHasher(cfg.Hasher.Cost)

	elWordsRepo := elrepositories.NewCollectionWordsRepo(elClient.Client)
	usersRepo := postgresRepo.NewUsersRepo(pgClient)
	collectionsRepo := postgresRepo.NewCollectionsRepo(pgClient)

	router := server.NewServer()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	app := api.NewApp(usersRepo, collectionsRepo, elWordsRepo, *tokenService, translatorManager, s3Manager, hasher)

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "hello from api new")
	})

	app.AttachEndpoints(router)

	router.Run()
}
