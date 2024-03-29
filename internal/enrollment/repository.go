package enrollment

import (
	"log"

	"github.com/nOm-1986/goApi_web/internal/domain"
	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(enroll *domain.Enrollment) error
	}

	repo struct {
		db  *gorm.DB
		log *log.Logger
	}
)

func NewRepo(l *log.Logger, db *gorm.DB) Repository {
	return &repo{
		db:  db,
		log: l,
	}
}

func (r *repo) Create(enroll *domain.Enrollment) error {

	if err := r.db.Create(enroll).Error; err != nil {
		r.log.Printf("error: %v", err)
		return err
	}

	r.log.Println("enrollment created with id: ", enroll.ID)
	return nil
}
