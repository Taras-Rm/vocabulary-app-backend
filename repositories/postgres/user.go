package postgres

import (
	"time"
	"vacabulary/models"

	"github.com/go-pg/pg/v10"
)

type UserModel struct {
	tableName struct{} `pg:"users"`

	ID        uint64 `pg:"id"`
	Name      string `pg:"name"`
	Email     string `pg:"email"`
	Password  string
	CreatedAt time.Time          `pg:"created_at"`
	IsSuper   bool               `pg:"is_super"`
	Settings  *UserSettingsModel `pg:"rel:has-one"`
}

type UserSettingsModel struct {
	tableName struct{} `pg:"user_settings"`

	ID       uint64 `pg:"id"`
	UserID   uint64 `pg:"user_id"`
	Language string `pg:"app_language"`
}

func (u *UserModel) FromModel() models.User {
	user := models.User{
		Id:        u.ID,
		Name:      u.Name,
		Password:  u.Password,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		IsSuper:   u.IsSuper,
	}

	if u.Settings != nil {
		s := u.Settings.FromModel()
		user.Settings = &s
	}

	return user
}

func (u *UserSettingsModel) FromModel() models.UserSettings {
	return models.UserSettings{
		Id:       u.ID,
		UserId:   u.UserID,
		Language: u.Language,
	}
}

func ToUserModel(u models.User) *UserModel {
	return &UserModel{
		ID:        u.Id,
		Name:      u.Name,
		Password:  u.Password,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		IsSuper:   u.IsSuper,
	}
}

func ToUserSettingsModel(u models.UserSettings) *UserSettingsModel {
	return &UserSettingsModel{
		ID:       u.Id,
		UserID:   u.UserId,
		Language: u.Language,
	}
}

type userRepo struct {
	db *pg.DB
}

type Users interface {
	Create(user models.User) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetById(id uint64) (*models.User, error)

	UpdateUserLanguage(language string, userId uint64) error
}

func NewUsersRepo(db *pg.DB) Users {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Create(user models.User) (*models.User, error) {
	userModel := ToUserModel(user)

	_, err := r.db.Model(userModel).Insert()
	if err != nil {
		return nil, err
	}

	userSettingsModel := UserSettingsModel{}
	userSettingsModel.UserID = userModel.ID
	_, err = r.db.Model(&userSettingsModel).Insert()
	if err != nil {
		return nil, err
	}

	createdUser := userModel.FromModel()
	return &createdUser, nil
}

func (r *userRepo) GetByEmail(email string) (*models.User, error) {
	user := UserModel{}
	err := r.db.Model(&user).Relation("Settings").Where("email=?", email).First()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	createdUser := user.FromModel()
	return &createdUser, nil
}

func (r *userRepo) GetById(id uint64) (*models.User, error) {
	user := UserModel{}
	err := r.db.Model(&user).Relation("Settings").Where("user_model.id=?", id).First()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	createdUser := user.FromModel()
	return &createdUser, nil
}

func (r *userRepo) UpdateUserLanguage(language string, userId uint64) error {
	settings := UserSettingsModel{Language: language}
	_, err := r.db.Model(&settings).Column("app_language").Where("user_id=?", userId).Update()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil
		}
		return err
	}

	return nil
}
