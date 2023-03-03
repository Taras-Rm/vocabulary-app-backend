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
	CreatedAt time.Time `pg:"created_at"`
}

func (u *UserModel) FromModel() models.User {
	return models.User{
		Id:        u.ID,
		Name:      u.Name,
		Password:  u.Password,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}

func ToUserModel(u models.User) *UserModel {
	return &UserModel{
		ID:        u.Id,
		Name:      u.Name,
		Password:  u.Password,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}

type userRepo struct {
	db *pg.DB
}

type Users interface {
	Create(user models.User) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetById(id uint64) (*models.User, error)
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

	createdUser := userModel.FromModel()
	return &createdUser, nil
}

func (r *userRepo) GetByEmail(email string) (*models.User, error) {
	user := UserModel{}
	err := r.db.Model(&user).Where("email=?", email).First()
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
	err := r.db.Model(&user).Where("id=?", id).First()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	createdUser := user.FromModel()
	return &createdUser, nil
}
