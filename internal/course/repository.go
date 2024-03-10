package course

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(course *Course) error
		GetAll(filters Filters, offset, limit int) ([]Course, error)
		Get(id string) (*Course, error)
		Update(id string, name *string, startDate, endDate *time.Time) error
		Count(filters Filters) (int, error)
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

func (repo *repo) GetAll(filters Filters, offset, limit int) ([]Course, error) {
	var courses []Course
	tx := repo.db.Model(&courses)
	tx = applyFilters(tx, filters)
	//Gorm ya trae para implementar el offset
	tx = tx.Limit(limit).Offset(offset)
	result := tx.Order("created_at desc").Find(&courses)
	if result.Error != nil {
		return nil, result.Error
	}
	return courses, nil
}

func (repo *repo) Get(id string) (*Course, error) {
	course := Course{ID: id}
	if err := repo.db.First(&course).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

func (repo *repo) Update(id string, name *string, startDate, endDate *time.Time) error {
	values := make(map[string]interface{})
	if name != nil {
		values["name"] = *name
	}
	if startDate != nil {
		values["start_date"] = *startDate
	}
	if endDate != nil {
		values["end_date"] = *endDate
	}
	if err := repo.db.Model(&Course{}).Where("id = ?", id).Updates(values); err.Error != nil {
		return err.Error
	}

	return nil
}

func (repo *repo) Count(filters Filters) (int, error) {
	var count int64
	tx := repo.db.Model(&Course{})
	tx = applyFilters(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

// Función para implementar filtros
func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.Name != "" {
		filters.Name = fmt.Sprintf("%%%s%%", strings.ToLower(filters.Name))
		tx = tx.Where("lower(name) like ?", filters.Name)
	}
	/* Realizar filtro por fecha. Ojo between and
	if filters.StartDate != "" {
		filters.StartDate = fmt.Sprintf("%%%s%%", strings.ToLower(filters.StartDate))
		tx = tx.Where("lower(start_date) like ?", filters.StartDate)
	}

	if filters.StartDate != "" {
		filters.StartDate = fmt.Sprintf("%%%s%%", strings.ToLower(filters.StartDate))
		tx = tx.Where("lower(start_date) like ?", filters.StartDate)
	}*/
	return tx
}
