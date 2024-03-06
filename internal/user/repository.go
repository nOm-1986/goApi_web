package user

import (
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user *User) error
	Get(id string) (*User, error)
	GetAll(filters Filters) ([]User, error)
	Delete(id string) error
	// Se pasan por puntero para que lleguen nil y poder identificarlos, si no se agregar el puntero llegan vacíos.
	Update(id string, firstName *string, lastName *string, email *string, phone *string) error
	Count(filters Filters) (int, error)
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

func (repo *repo) GetAll(filters Filters) ([]User, error) {
	var u []User

	//La info que nos traiga sea acorde a la estructura usuario
	//El Find lo que hace es poblar la data en la estructura
	// if err := repo.db.Model(&u).Order("created_at desc").Find(&u).Error; err != nil {
	// 	return nil, err
	// }
	tx := repo.db.Model(&u)
	tx = applyFilters(tx, filters)
	//result := tx.Order("created_at desc").Limit(5).Find(&u) - Para agregar un límite
	result := tx.Order("created_at desc").Find(&u)
	if result.Error != nil {
		return nil, result.Error
	}
	return u, nil
}

func (repo *repo) Get(id string) (*User, error) {
	user := User{ID: id}
	if err := repo.db.First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *repo) Delete(id string) error {
	user := User{ID: id}
	if err := repo.db.Delete(&user).Error; err != nil {
		return err
	}
	return nil
}

// Remember that if the field came nil, it was due it wat not sent
func (repo *repo) Update(id string, firstName *string, lastName *string, email *string, phone *string) error {

	values := make(map[string]interface{})

	if firstName != nil {
		values["first_name"] = *firstName
	}
	if lastName != nil {
		values["last_name"] = *lastName
	}
	if email != nil {
		values["email"] = *email
	}
	if phone != nil {
		values["phone"] = *phone
	}

	if err := repo.db.Model(&User{}).Where("id = ?", id).Updates(values).Error; err != nil {
		return err
	}
	return nil
}

func (repo *repo) Count(filters Filters) (int, error) {
	var count int64
	tx := repo.db.Model(&User{})
	tx = applyFilters(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

// Función para implementar filtros
func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.FirstName != "" {
		filters.FirstName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.FirstName))
		tx = tx.Where("lower(first_name) like ?", filters.FirstName)
	}
	if filters.LastName != "" {
		filters.LastName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.LastName))
		tx = tx.Where("lower(last_name) like ?", filters.LastName)
	}
	return tx
}
