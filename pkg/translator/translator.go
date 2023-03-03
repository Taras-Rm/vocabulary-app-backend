package translator

import (
	"errors"
	"vacabulary/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/translate"
)

type TranslatorManager struct {
	awsSession *session.Session
	config     config.AWSConfig
}

func (tm *TranslatorManager) TranslateWord(origin string, langFrom, langTo string) (string, error) {
	sess := tm._getAwsSession()

	transl := translate.New(sess)

	sourceLangCode, err := tm.languageCodeToAwsTranslationCode(langFrom)
	if err != nil {
		return "", nil
	}

	targetLangCode, err := tm.languageCodeToAwsTranslationCode(langTo)
	if err != nil {
		return "", nil
	}

	response, err := transl.Text(&translate.TextInput{
		SourceLanguageCode: aws.String(sourceLangCode),
		TargetLanguageCode: aws.String(targetLangCode),
		Text:               aws.String(origin),
	})
	if err != nil {
		return "", nil
	}

	return *response.TranslatedText, nil
}

func (tm *TranslatorManager) _getAwsSession() *session.Session {
	if tm.awsSession == nil {
		tm._setAwsSession()
	}

	return tm.awsSession
}

func (tm *TranslatorManager) _setAwsSession() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      &tm.config.Region,
		Credentials: credentials.NewStaticCredentials(tm.config.AccessId, tm.config.Secret, ""),
	}))

	tm.awsSession = sess
}

func NewTranslatorManager(config config.AWSConfig) TranslatorManager {
	return TranslatorManager{
		config: config,
	}
}

func (tm *TranslatorManager) languageCodeToAwsTranslationCode(langCode string) (string, error) {
	switch langCode {
	case "ua":
		return "uk", nil
	case "en":
		return "en", nil
	}

	return "", errors.New("unknown language code")
}
