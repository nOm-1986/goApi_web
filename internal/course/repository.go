package course

import (
	"log"

	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(course *Course) error
	}

	repo struct {
		log *log.Logger
		db  *gorm.DB
	}
)

func NewRepo(l *log.Logger, db *gorm.DB) Repository {
	return &repo{
		log: l,
		db:  db,
	}
}

func (repo *repo) Create(course *Course) error {
	if err := repo.db.Create(course).Error; err != nil {
		repo.log.Printf("error: %v", err)
		return err
	}
	repo.log.Println("course created with id: ", course.ID)
	return nil
}
