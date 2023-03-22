package api

import (
	"vacabulary/pkg/hasher"
	"vacabulary/pkg/s3"
	"vacabulary/pkg/token"
	"vacabulary/pkg/translator"
	"vacabulary/repositories/elastic"
	"vacabulary/repositories/postgres"

	"github.com/gin-gonic/gin"
)

type App struct {
	userRepo       postgres.Users
	wordRepo       elastic.Words
	collectionRepo postgres.Collections

	tokenService      token.TokenService
	translatorManager translator.TranslatorManager
	s3Manager         s3.S3Manager
	hasher            hasher.Hasher
}

func NewApp(userRepo postgres.Users, collectionRepo postgres.Collections, wordRepo elastic.Words, tokenService token.TokenService, translatorManager translator.TranslatorManager, s3Manager s3.S3Manager, hasher hasher.Hasher) App {
	return App{
		userRepo:       userRepo,
		wordRepo:       wordRepo,
		collectionRepo: collectionRepo,

		tokenService:      tokenService,
		translatorManager: translatorManager,
		s3Manager:         s3Manager,
		hasher:            hasher,
	}
}

func (a *App) AttachEndpoints(gr *gin.Engine) {
	a.InjectWords(gr)
	a.InjectUsers(gr)
	a.InjectCollections(gr)
	a.InjectStatistic(gr)
}
