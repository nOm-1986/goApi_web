package user

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user *User) error
	Get(id string) (*User, error)
	GetAll() ([]User, error)
	Delete(id string) (*User, error)
}

type repo struct {
	log *log.Logger
	db  *gorm.DB
}

func NewRepo(log *log.Logger, db *gorm.DB) Repository {
	return &repo{
		log: log,
		db:  db,
	}
}

func (repo *repo) Create(user *User) error {
	user.ID = uuid.New().String()

	if err := repo.db.Create(user).Error; err != nil {
		repo.log.Println(err)
		return err
	}
	repo.log.Println("User created with id: ", user.ID)
	return nil
}

func (repo *repo) GetAll() ([]User, error) {
	var u []User
	//La info que nos traiga sea acorde a la estructura usuario
	//El Find lo que hace es poblar la data en la estructura
	result := repo.db.Model(&u).Order("created_at desc").Find(&u)
	if result.Error != nil {
		return nil, result.Error
	}
	return u, nil

}

func (repo *repo) Get(id string) (*User, error) {
	user := User{ID: id}
	result := repo.db.First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repo *repo) Delete(id string) (*User, error) {
	return &User{}, nil
}
