package postgres

import (
	"time"
	"vacabulary/models"

	"github.com/go-pg/pg/v10"
)

type CollectionModel struct {
	tableName struct{} `pg:"collections"`

	ID        uint64    `pg:"id"`
	Name      string    `pg:"name"`
	OwnerID   uint64    `pg:"owner_id"`
	CreatedAt time.Time `pg:"created_at"`
	LangFrom  string    `pg:"lang_from"`
	LangTo    string    `pg:"lang_to"`
}

func (u *CollectionModel) FromModel() *models.Collection {
	return &models.Collection{
		Id:        u.ID,
		Name:      u.Name,
		OwnerId:   u.OwnerID,
		CreatedAt: u.CreatedAt,
		LangFrom:  u.LangFrom,
		LangTo:    u.LangTo,
	}
}

func ToCollectionModel(u models.Collection) *CollectionModel {
	return &CollectionModel{
		ID:        u.Id,
		Name:      u.Name,
		OwnerID:   u.OwnerId,
		CreatedAt: u.CreatedAt,
		LangFrom:  u.LangFrom,
		LangTo:    u.LangTo,
	}
}

type collectionRepo struct {
	db *pg.DB
}

type Collections interface {
	Create(collection models.Collection) (*models.Collection, error)
	GetByOwnerId(ownerId uint64) ([]models.Collection, error)
	GetById(id uint64) (*models.Collection, error)
	GetByName(name string) (*models.Collection, error)
	Update(collection *models.Collection) (*models.Collection, error)
	DeleteById(id uint64) error
}

func NewCollectionsRepo(db *pg.DB) Collections {
	return &collectionRepo{
		db: db,
	}
}

func (r *collectionRepo) Create(collection models.Collection) (*models.Collection, error) {
	collectionModel := ToCollectionModel(collection)

	_, err := r.db.Model(collectionModel).Insert()
	if err != nil {
		return nil, err
	}

	createdCollection := collectionModel.FromModel()
	return createdCollection, nil
}

func (r *collectionRepo) GetByOwnerId(ownerId uint64) ([]models.Collection, error) {
	var collectionModels []CollectionModel

	err := r.db.Model(&collectionModels).Where("owner_id=?", ownerId).Order("created_at").Select()
	if err != nil {
		return nil, err
	}

	collections := []models.Collection{}
	for _, c := range collectionModels {
		collections = append(collections, *c.FromModel())
	}

	return collections, nil
}

func (r *collectionRepo) GetById(id uint64) (*models.Collection, error) {
	collection := CollectionModel{}
	err := r.db.Model(&collection).Where("id=?", id).First()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	createdCollection := collection.FromModel()
	return createdCollection, nil
}

func (r *collectionRepo) GetByName(name string) (*models.Collection, error) {
	collection := CollectionModel{}
	err := r.db.Model(&collection).Where("name=?", name).First()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	createdCollection := collection.FromModel()
	return createdCollection, nil
}

func (r *collectionRepo) Update(collection *models.Collection) (*models.Collection, error) {
	model := ToCollectionModel(*collection)

	_, err := r.db.Model(model).Where("id=?", model.ID).Column("name").Update()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	updatedCollection := model.FromModel()
	return updatedCollection, nil
}

func (r *collectionRepo) DeleteById(id uint64) error {
	_, err := r.db.Model(&CollectionModel{}).Where("id=?", id).Delete()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil
		}
		return err
	}

	return nil
}
